package migplan

import (
	migapi "github.com/konveyor/mig-controller/pkg/apis/migration/v1alpha1"
	migctl "github.com/konveyor/mig-controller/pkg/controller/migmigration"
	migref "github.com/konveyor/mig-controller/pkg/reference"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type PlanPredicate struct {
	predicate.Funcs
}

func (r PlanPredicate) Create(e event.CreateEvent) bool {
	plan, cast := e.Object.(*migapi.MigPlan)
	if cast {
		if !plan.InSandbox() {
			return false
		}
		r.mapRefs(plan)
	}
	return true
}

func (r PlanPredicate) Update(e event.UpdateEvent) bool {
	old, cast := e.ObjectOld.(*migapi.MigPlan)
	if !cast {
		return false
	}
	new, cast := e.ObjectNew.(*migapi.MigPlan)
	if !cast {
		return false
	}
	if !old.InSandbox() {
		return false
	}
	changed := !reflect.DeepEqual(old.Spec, new.Spec) ||
		!reflect.DeepEqual(old.DeletionTimestamp, new.DeletionTimestamp)
	if changed {
		r.unmapRefs(old)
		r.mapRefs(new)
	}
	return changed
}

func (r PlanPredicate) Delete(e event.DeleteEvent) bool {
	plan, cast := e.Object.(*migapi.MigPlan)
	if cast {
		if !plan.InSandbox() {
			return false
		}
		r.unmapRefs(plan)
	}
	return true
}

func (r PlanPredicate) Generic(e event.GenericEvent) bool {
	plan, cast := e.Object.(*migapi.MigPlan)
	if cast {
		if !plan.InSandbox() {
			return false
		}
		r.mapRefs(plan)
	}
	return true
}

func (r PlanPredicate) mapRefs(plan *migapi.MigPlan) {
	refMap := migref.GetMap()

	refOwner := migref.RefOwner{
		Kind:      migref.ToKind(plan),
		Namespace: plan.Namespace,
		Name:      plan.Name,
	}

	// source cluster
	ref := plan.Spec.SrcMigClusterRef
	if migref.RefSet(ref) {
		refMap.Add(refOwner, migref.RefTarget{
			Kind:      migref.ToKind(migapi.MigCluster{}),
			Namespace: ref.Namespace,
			Name:      ref.Name,
		})
	}

	// destination cluster
	ref = plan.Spec.DestMigClusterRef
	if migref.RefSet(ref) {
		refMap.Add(refOwner, migref.RefTarget{
			Kind:      migref.ToKind(migapi.MigCluster{}),
			Namespace: ref.Namespace,
			Name:      ref.Name,
		})
	}

	// storage
	ref = plan.Spec.MigStorageRef
	if migref.RefSet(ref) {
		refMap.Add(refOwner, migref.RefTarget{
			Kind:      migref.ToKind(migapi.MigStorage{}),
			Namespace: ref.Namespace,
			Name:      ref.Name,
		})
	}
}

func (r PlanPredicate) unmapRefs(plan *migapi.MigPlan) {
	refMap := migref.GetMap()

	refOwner := migref.RefOwner{
		Kind:      migref.ToKind(plan),
		Namespace: plan.Namespace,
		Name:      plan.Name,
	}

	// source cluster
	ref := plan.Spec.SrcMigClusterRef
	if migref.RefSet(ref) {
		refMap.Delete(refOwner, migref.RefTarget{
			Kind:      migref.ToKind(migapi.MigCluster{}),
			Namespace: ref.Namespace,
			Name:      ref.Name,
		})
	}

	// destination cluster
	ref = plan.Spec.DestMigClusterRef
	if migref.RefSet(ref) {
		refMap.Delete(refOwner, migref.RefTarget{
			Kind:      migref.ToKind(migapi.MigCluster{}),
			Namespace: ref.Namespace,
			Name:      ref.Name,
		})
	}

	// storage
	ref = plan.Spec.MigStorageRef
	if migref.RefSet(ref) {
		refMap.Delete(refOwner, migref.RefTarget{
			Kind:      migref.ToKind(migapi.MigStorage{}),
			Namespace: ref.Namespace,
			Name:      ref.Name,
		})
	}
}

type ClusterPredicate struct {
	predicate.Funcs
}

func (r ClusterPredicate) Create(e event.CreateEvent) bool {
	return false
}

func (r ClusterPredicate) Update(e event.UpdateEvent) bool {
	new, cast := e.ObjectNew.(*migapi.MigCluster)
	if cast && new.InPrivileged() {
		// Reconciled by the controller.
		return new.HasReconciled()
	}
	return false
}

func (r ClusterPredicate) Delete(e event.DeleteEvent) bool {
	cluster, cast := e.Object.(*migapi.MigCluster)
	return cast && cluster.InPrivileged()
}

func (r ClusterPredicate) Generic(e event.GenericEvent) bool {
	cluster, cast := e.Object.(*migapi.MigCluster)
	return cast && cluster.InPrivileged()
}

type StoragePredicate struct {
	predicate.Funcs
}

func (r StoragePredicate) Create(e event.CreateEvent) bool {
	return false
}

func (r StoragePredicate) Update(e event.UpdateEvent) bool {
	new, cast := e.ObjectNew.(*migapi.MigStorage)
	if cast && new.InSandbox() {
		// Reconciled by the controller.
		return new.HasReconciled()
	}
	return false
}

func (r StoragePredicate) Delete(e event.DeleteEvent) bool {
	storage, cast := e.Object.(*migapi.MigStorage)
	return cast && storage.InSandbox()
}

func (r StoragePredicate) Generic(e event.GenericEvent) bool {
	storage, cast := e.Object.(*migapi.MigStorage)
	return cast && storage.InSandbox()
}

type MigrationPredicate struct {
	predicate.Funcs
}

func (r MigrationPredicate) Create(e event.CreateEvent) bool {
	return false
}

func (r MigrationPredicate) Update(e event.UpdateEvent) bool {
	old, cast := e.ObjectOld.(*migapi.MigMigration)
	if !cast {
		return false
	}
	new, cast := e.ObjectNew.(*migapi.MigMigration)
	if !cast {
		return false
	}
	if !old.InSandbox() {
		return false
	}
	started := !old.Status.HasCondition(migctl.Running) &&
		new.Status.HasCondition(migctl.Running)
	stopped := old.Status.HasCondition(migctl.Running) &&
		!new.Status.HasCondition(migctl.Running)
	if started || stopped {
		return true
	}

	return false
}

func (r MigrationPredicate) Delete(e event.DeleteEvent) bool {
	migration, cast := e.Object.(*migapi.MigMigration)
	return cast && migration.InSandbox()
}

func (r MigrationPredicate) Generic(e event.GenericEvent) bool {
	migration, cast := e.Object.(*migapi.MigMigration)
	return cast && migration.InSandbox()
}
