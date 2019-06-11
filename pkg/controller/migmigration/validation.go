package migmigration

import (
	migapi "github.com/fusor/mig-controller/pkg/apis/migration/v1alpha1"
	migref "github.com/fusor/mig-controller/pkg/reference"
)

// Types
const (
	InvalidPlanRef = "InvalidPlanRef"
	PlanNotReady   = "PlanNotReady"
)

// Categories
const (
	Critical = migapi.Critical
)

// Reasons
const (
	NotSet   = "NotSet"
	NotFound = "NotFound"
)

// Statuses
const (
	True  = migapi.True
	False = migapi.False
)

// Messages
const (
	ReadyMessage          = "The migration is ready."
	InvalidPlanRefMessage = "The `migPlanRef` must reference a `migplan`."
	PlanNotReadyMessage   = "The referenced `migPlanRef` does not have a `Ready` condition."
)

// Validate the plan resource.
// Returns error and the total error conditions set.
func (r ReconcileMigMigration) validate(migration *migapi.MigMigration) error {
	// Plan
	err := r.validatePlan(migration)
	if err != nil {
		return err
	}

	return nil
}

// Validate the referenced plan.
// Returns error and the total error conditions set.
func (r ReconcileMigMigration) validatePlan(migration *migapi.MigMigration) error {
	ref := migration.Spec.MigPlanRef

	// NotSet
	if !migref.RefSet(ref) {
		migration.Status.SetCondition(migapi.Condition{
			Type:     InvalidPlanRef,
			Status:   True,
			Reason:   NotSet,
			Category: Critical,
			Message:  InvalidPlanRefMessage,
		})
		return nil
	}

	plan, err := migapi.GetPlan(r, ref)
	if err != nil {
		return err
	}

	// NotFound
	if plan == nil {
		migration.Status.SetCondition(migapi.Condition{
			Type:     InvalidPlanRef,
			Status:   True,
			Reason:   NotFound,
			Category: Critical,
			Message:  InvalidPlanRefMessage,
		})
		return nil
	}

	// NotReady
	if !plan.Status.IsReady() {
		migration.Status.SetCondition(migapi.Condition{
			Type:     PlanNotReady,
			Status:   True,
			Category: Critical,
			Message:  PlanNotReadyMessage,
		})
		return nil
	}

	return nil
}
