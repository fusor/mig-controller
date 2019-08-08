package migplan

import (
	migapi "github.com/fusor/mig-controller/pkg/apis/migration/v1alpha1"
	migref "github.com/fusor/mig-controller/pkg/reference"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func MigrationRequests(a handler.MapObject) []reconcile.Request {
	requests := []reconcile.Request{}
	migration, cast := a.Object.(*migapi.MigMigration)
	if !cast {
		return requests
	}
	ref := migration.Spec.MigPlanRef
	if migref.RefSet(ref) {
		requests = append(
			requests, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: ref.Namespace,
					Name:      ref.Name,
				},
			})
	}

	return requests
}
