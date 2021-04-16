package container

import (
	"context"
	"time"

	migapi "github.com/konveyor/mig-controller/pkg/apis/migration/v1alpha1"
	"github.com/konveyor/mig-controller/pkg/controller/discovery/model"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

//
// A collection of k8s MigHook resources.
type MigHook struct {
	// Base
	BaseCollection
}

func (r *MigHook) AddWatch(dsController controller.Controller) error {
	err := dsController.Watch(
		&source.Kind{
			Type: &migapi.MigHook{},
		},
		&handler.EnqueueRequestForObject{},
		r)
	if err != nil {
		Log.Trace(err)
		return err
	}

	return nil
}

func (r *MigHook) Reconcile() error {
	mark := time.Now()
	sr := SimpleReconciler{
		Db: r.ds.Container.Db,
	}
	err := sr.Reconcile(r)
	if err != nil {
		Log.Trace(err)
		return err
	}
	r.hasReconciled = true
	Log.Info(
		"MigHook (collection) reconciled.",
		"ns",
		r.ds.Cluster.Namespace,
		"name",
		r.ds.Cluster.Name,
		"duration",
		time.Since(mark))

	return nil
}

func (r *MigHook) GetDiscovered() ([]model.Model, error) {
	models := []model.Model{}
	onCluster := migapi.MigHookList{}
	err := r.ds.Client.List(context.TODO(), nil, &onCluster)
	if err != nil {
		Log.Trace(err)
		return nil, err
	}
	for _, discovered := range onCluster.Items {
		dim := &model.Hook{}
		dim.With(&discovered)
		models = append(models, dim)
	}

	return models, nil
}

func (r *MigHook) GetStored() ([]model.Model, error) {
	models := []model.Model{}
	list, err := model.Hook{}.List(
		r.ds.Container.Db,
		model.ListOptions{})
	if err != nil {
		Log.Trace(err)
		return nil, err
	}
	for _, dim := range list {
		models = append(models, dim)
	}

	return models, nil
}

//
// Predicate methods.
//

func (r *MigHook) Create(e event.CreateEvent) bool {
	Log.Reset()
	object, cast := e.Object.(*migapi.MigHook)
	if !cast {
		return false
	}
	dim := model.Hook{}
	dim.With(object)
	r.ds.Create(&dim)

	return false
}

func (r *MigHook) Update(e event.UpdateEvent) bool {
	Log.Reset()
	object, cast := e.ObjectNew.(*migapi.MigHook)
	if !cast {
		return false
	}
	restore := model.Hook{}
	restore.With(object)
	r.ds.Update(&restore)

	return false
}

func (r *MigHook) Delete(e event.DeleteEvent) bool {
	Log.Reset()
	object, cast := e.Object.(*migapi.MigHook)
	if !cast {
		return false
	}
	dim := model.Hook{}
	dim.With(object)
	r.ds.Delete(&dim)

	return false
}

func (r *MigHook) Generic(e event.GenericEvent) bool {
	return false
}
