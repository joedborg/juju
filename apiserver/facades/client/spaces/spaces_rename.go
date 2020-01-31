// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package spaces

import (
	"github.com/juju/errors"
	"gopkg.in/juju/names.v3"
	"gopkg.in/mgo.v2/txn"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/params"
	jujucontroller "github.com/juju/juju/controller"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/permission"
	"github.com/juju/juju/core/settings"
	"github.com/juju/juju/state"
)

// RenameSpaceModelOp describes a model operation for renaming a space.
type RenameSpaceModelOp interface {
	state.ModelOperation
}

// RenameSpace describes a space that can be renamed.
type RenameSpace interface {
	Refresh() error
	Id() string
	Name() string
	RenameSpaceCompleteOps(toName string) ([]txn.Op, error)
}

// RenameSpaceState describes state operations required
// to execute the renameSpace operation.
// * This allows us to indirect state at the operation level instead of the
// * whole API level as currently done in interface.go
type RenameSpaceState interface {
	// ControllerConfig returns current ControllerConfig.
	ControllerConfig() (jujucontroller.Config, error)

	// ConstraintsBySpace returns current constraints using the given spaceName.
	ConstraintsBySpaceName(spaceName string) (map[string]constraints.Value, error)

	// ControllerSettingsGlobalKey returns the global controller settings key..
	ControllerSettingsGlobalKey() string

	// GetConstraintsOps gets the database transaction operations for the given constraints.
	// Cons is a map keyed by the DocID.
	GetConstraintsOps(cons map[string]constraints.Value) ([]txn.Op, error)
}

// Settings describes methods for interacting with settings to apply
// space-based configuration deltas.
type Settings interface {
	DeltaOps(key string, delta settings.ItemChanges) ([]txn.Op, error)
}

// Model describes methods for interacting with Model to
// check whether the current model is a controllerModel.
type Model interface {
	IsControllerModel() bool
}

type spaceRenameModelOp struct {
	st       RenameSpaceState
	space    RenameSpace
	settings Settings
	model    Model
	toName   string
}

func (o *spaceRenameModelOp) Done(err error) error {
	return err
}

func NewRenameSpaceModelOp(model Model, settings Settings, st RenameSpaceState, space RenameSpace, toName string) *spaceRenameModelOp {
	return &spaceRenameModelOp{
		st:       st,
		settings: settings,
		space:    space,
		model:    model,
		toName:   toName,
	}
}

type renameSpaceStateShim struct {
	*state.State
}

// Build (state.ModelOperation) creates and returns a slice of transaction
// operations necessary to rename a space.
func (o *spaceRenameModelOp) Build(attempt int) ([]txn.Op, error) {
	if attempt > 0 {
		if err := o.space.Refresh(); err != nil {
			return nil, errors.Trace(err)
		}
	}

	var totalOps []txn.Op

	settingsDelta, err := o.getSettingsChanges(o.space.Name(), o.toName)
	if err != nil {
		newErr := errors.Annotatef(err, "retrieving setting changes")
		return nil, errors.Trace(newErr)
	}
	newConstraints, err := o.getConstraintsChanges(o.space.Name(), o.toName)
	if err != nil {
		newErr := errors.Annotatef(err, "retrieving constraint changes")
		return nil, errors.Trace(newErr)
	}

	newConstraintsOps, err := o.st.GetConstraintsOps(newConstraints)
	if err != nil {
		return nil, errors.Trace(err)
	}
	newSettingsOps, err := o.settings.DeltaOps(o.st.ControllerSettingsGlobalKey(), settingsDelta)
	if err != nil {
		return nil, errors.Trace(err)
	}

	completeOps, err := o.space.RenameSpaceCompleteOps(o.toName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	totalOps = append(totalOps, completeOps...)
	totalOps = append(totalOps, newConstraintsOps...)
	totalOps = append(totalOps, newSettingsOps...)

	return totalOps, nil
}

// getConstraintsChanges gets new constraints after applying the new space name.
func (o *spaceRenameModelOp) getConstraintsChanges(fromSpaceName, toName string) (map[string]constraints.Value, error) {
	currentConstraints, err := o.st.ConstraintsBySpaceName(fromSpaceName)
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, constraint := range currentConstraints {
		spaces := constraint.Spaces
		if spaces == nil {
			continue
		}
		for i, space := range *spaces {
			if space == fromSpaceName {
				(*spaces)[i] = toName
				break
			}
		}
	}
	return currentConstraints, nil
}

// getSettingsChanges get's skipped and returns nil if we are not in the controllerModel
func (o *spaceRenameModelOp) getSettingsChanges(fromSpaceName, toName string) (settings.ItemChanges, error) {
	if !o.model.IsControllerModel() {
		return nil, nil
	}
	currentControllerConfig, err := o.st.ControllerConfig()
	if err != nil {
		return nil, errors.Trace(err)
	}

	var deltas settings.ItemChanges

	if mgmtSpace := currentControllerConfig.JujuManagementSpace(); mgmtSpace == fromSpaceName {
		change := settings.MakeModification(jujucontroller.JujuManagementSpace, fromSpaceName, toName)
		deltas = append(deltas, change)
	}
	if haSpace := currentControllerConfig.JujuHASpace(); haSpace == fromSpaceName {
		change := settings.MakeModification(jujucontroller.JujuHASpace, fromSpaceName, toName)
		deltas = append(deltas, change)
	}
	return deltas, nil
}

// RenameSpace renames a space.
func (api *API) RenameSpace(args params.RenameSpacesParams) (params.ErrorResults, error) {
	isAdmin, err := api.auth.HasPermission(permission.AdminAccess, api.backing.ModelTag())
	if err != nil && !errors.IsNotFound(err) {
		return params.ErrorResults{}, errors.Trace(err)
	}
	if !isAdmin {
		return params.ErrorResults{}, common.ServerError(common.ErrPerm)
	}
	if err := api.check.ChangeAllowed(); err != nil {
		return params.ErrorResults{}, errors.Trace(err)
	}
	if err = api.checkSupportsProviderSpaces(); err != nil {
		return params.ErrorResults{}, common.ServerError(errors.Trace(err))
	}
	results := params.ErrorResults{
		Results: make([]params.ErrorResult, len(args.SpacesRenames)),
	}

	for i, spaceRename := range args.SpacesRenames {
		fromTag, err := names.ParseSpaceTag(spaceRename.FromSpaceTag)
		if err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
			continue
		}
		toTag, err := names.ParseSpaceTag(spaceRename.ToSpaceTag)
		if err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
			continue
		}
		toSpace, err := api.backing.SpaceByName(toTag.Id())
		if err != nil && !errors.IsNotFound(err) {
			newErr := errors.Annotatef(err, "retrieving space: %q unexpected error, besides not found", toTag.Id())
			results.Results[i].Error = common.ServerError(errors.Trace(newErr))
			continue
		}
		if toSpace != nil {
			newErr := errors.AlreadyExistsf("space: %q", toTag.Id())
			results.Results[i].Error = common.ServerError(errors.Trace(newErr))
			continue
		}
		operation, err := api.opFactory.NewRenameSpaceModelOp(fromTag.Id(), toTag.Id())
		if err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
			continue
		}
		if err = api.backing.ApplyOperation(operation); err != nil {
			results.Results[i].Error = common.ServerError(errors.Trace(err))
			continue
		}
	}
	return results, nil
}
