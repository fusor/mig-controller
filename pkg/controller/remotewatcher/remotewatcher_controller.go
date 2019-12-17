/*
Copyright 2019 Red Hat Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package remotewatcher

import (
	"github.com/fusor/mig-controller/pkg/logging"
	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	kapi "k8s.io/api/core/v1"
	storageapi "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logging.WithName("remote-watch")

// Add creates a new RemoteWatcher Controller with a forwardChannel
func Add(mgr manager.Manager, forwardChannel chan event.GenericEvent, fowardEvent event.GenericEvent) error {
	return add(mgr, newReconciler(mgr, forwardChannel, fowardEvent))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(
	mgr manager.Manager,
	forwardChannel chan event.GenericEvent,
	forwardEvent event.GenericEvent) reconcile.Reconciler {
	return &ReconcileRemoteWatcher{
		Client:         mgr.GetClient(),
		scheme:         mgr.GetScheme(),
		ForwardChannel: forwardChannel,
		ForwardEvent:   forwardEvent,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("remotewatcher-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		log.Trace(err)
		return err
	}

	// Backup
	err = c.Watch(
		&source.Kind{
			Type: &velerov1.Backup{},
		},
		&handler.EnqueueRequestForObject{},
		&BackupPredicate{})
	if err != nil {
		log.Trace(err)
		return err
	}

	// Restore
	err = c.Watch(
		&source.Kind{
			Type: &velerov1.Restore{},
		},
		&handler.EnqueueRequestForObject{},
		&RestorePredicate{})
	if err != nil {
		log.Trace(err)
		return err
	}

	// BSL
	err = c.Watch(
		&source.Kind{
			Type: &velerov1.BackupStorageLocation{},
		},
		&handler.EnqueueRequestForObject{},
		&BSLPredicate{})
	if err != nil {
		log.Trace(err)
		return err
	}

	// VSL
	err = c.Watch(
		&source.Kind{
			Type: &velerov1.VolumeSnapshotLocation{},
		},
		&handler.EnqueueRequestForObject{},
		&VSLPredicate{})
	if err != nil {
		log.Trace(err)
		return err
	}

	// Secret
	err = c.Watch(
		&source.Kind{
			Type: &kapi.Secret{},
		},
		&handler.EnqueueRequestForObject{},
		&SecretPredicate{})
	if err != nil {
		log.Trace(err)
		return err
	}

	// Pod
	err = c.Watch(
		&source.Kind{
			Type: &kapi.Pod{},
		},
		&handler.EnqueueRequestForObject{},
		&PodPredicate{})
	if err != nil {
		log.Trace(err)
		return err
	}

	// PV
	err = c.Watch(
		&source.Kind{
			Type: &kapi.PersistentVolume{},
		},
		&handler.EnqueueRequestForObject{})
	if err != nil {
		log.Trace(err)
		return err
	}

	// PVC
	err = c.Watch(
		&source.Kind{
			Type: &kapi.PersistentVolumeClaim{},
		},
		&handler.EnqueueRequestForObject{})
	if err != nil {
		log.Trace(err)
		return err
	}

	// Namespaces
	err = c.Watch(
		&source.Kind{
			Type: &kapi.Namespace{},
		},
		&handler.EnqueueRequestForObject{})
	if err != nil {
		log.Trace(err)
		return err
	}

	// StorageClass
	err = c.Watch(
		&source.Kind{
			Type: &storageapi.StorageClass{},
		},
		&handler.EnqueueRequestForObject{})
	if err != nil {
		log.Trace(err)
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileRemoteWatcher{}

// ReconcileRemoteWatcher reconciles a RemoteWatcher object
type ReconcileRemoteWatcher struct {
	client.Client
	scheme *runtime.Scheme
	// channel to forward GenericEvents to
	ForwardChannel chan event.GenericEvent
	// Event to forward when this controller gets event
	ForwardEvent event.GenericEvent
}

// Reconcile reads that state of the cluster for a RemoteWatcher object and makes changes
func (r *ReconcileRemoteWatcher) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Reset()
	// Forward a known Event back to the parent controller
	r.ForwardChannel <- r.ForwardEvent
	return reconcile.Result{}, nil
}
