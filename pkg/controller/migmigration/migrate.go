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

package migmigration

import (
	migapi "github.com/fusor/mig-controller/pkg/apis/migration/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// Annotations
const (
	pvAnnotationKey             = "openshift.io/migrate-type"
	migrateAnnotationValue      = "final"
	migrateAnnotationKey        = "openshift.io/migrate-copy-phase"
	stageAnnotationValue        = "stage"
	stageAnnotationKey          = "openshift.io/migrate-copy-phase"
	resticPvBackupAnnotationKey = "backup.velero.io/backup-volumes"
)

// Labels
const (
	pvBackupLabelKey   = "openshift.io/pv-backup"
	pvBackupLabelValue = "true"
)

// Backup resources.
var stagingResources = []string{
	"persistentvolumes",
	"persistentvolumeclaims",
	"imagestreams",
	"imagestreamtags",
	"secrets",
	"configmaps",
	"pods",
}

// Perform the migration.
func (r *ReconcileMigMigration) migrate(migration *migapi.MigMigration) (int, error) {
	if migration.Status.HasAnyCondition(Succeeded, Failed) {
		return 0, nil
	}

	// Ready
	plan, err := migration.GetPlan(r)
	if err != nil {
		return 0, err
	}
	if !plan.Status.IsReady() {
		log.Info("Plan not ready.", "name", migration.Name)
		return 0, err
	}

	// Resources
	planResources, err := plan.GetRefResources(r)
	if err != nil {
		log.Trace(err)
		return 0, err
	}

	// Started
	if migration.Status.StartTimestamp == nil {
		migration.Status.StartTimestamp = &metav1.Time{Time: time.Now()}
	}

	// Run
	task := Task{
		Log:             log,
		Client:          r,
		Owner:           migration,
		PlanResources:   planResources,
		Phase:           Phase{Name: migration.Status.Phase},
		Annotations:     r.getAnnotations(migration),
		BackupResources: r.getBackupResources(migration),
	}
	err = task.Run()
	if err != nil {
		log.Trace(err)
		return 0, err
	}

	migration.Status.SetCondition(migapi.Condition{
		Type:     Running,
		Status:   True,
		Reason:   task.Phase.Name,
		Category: Advisory,
		Message:  RunningMessage,
	})

	// TODO: SYNC WITH JEFF TO GET OPINION ON THIS HACK
	// Setting this annotation to not create stage pods after copy restore has
	// run to completion
	if task.Phase.Equals(DeleteStagePodsStarted) {
		if migration.Annotations == nil {
			migration.Annotations = make(map[string]string)
		}
		migration.Annotations["openshift.io/stage-completed"] = "true"
	}

	// Result
	migration.Status.Phase = task.Phase.Name
	if task.Phase.Equals(WaitOnResticRestart) {
		return 10, nil
	}
	if task.Phase.Final() {
		migration.Status.DeleteCondition(Running)
		migration.Status.SetCondition(migapi.Condition{
			Type:     Succeeded,
			Status:   True,
			Reason:   task.Phase.Name,
			Category: Advisory,
			Message:  SucceededMessage,
		})
	}
	if task.Phase.Failed() {
		migration.AddErrors(task.Errors)
		migration.Status.DeleteCondition(Running)
		migration.Status.SetCondition(migapi.Condition{
			Type:     Failed,
			Status:   True,
			Reason:   task.Phase.Name,
			Category: Critical,
			Message:  FailedMessage,
		})
	}

	return 0, nil
}

// Get annotations.
// TODO: Revisit this. We are hardcoding this for now until 2 things occur.
// 1. We are properly setting this annotation from user input to the UI
// 2. We fix the plugin to operate migration specific behavior on the
// migrateAnnnotationKey
func (r *ReconcileMigMigration) getAnnotations(migration *migapi.MigMigration) map[string]string {
	annotations := make(map[string]string)
	if migration.Spec.Stage {
		annotations[stageAnnotationKey] = stageAnnotationValue
	} else {
		annotations[migrateAnnotationKey] = migrateAnnotationValue
	}
	return annotations
}

// Get the resources (kinds) to be included in the backup.
func (r *ReconcileMigMigration) getBackupResources(migration *migapi.MigMigration) []string {
	if migration.Spec.Stage {
		return stagingResources
	}

	return []string{}
}
