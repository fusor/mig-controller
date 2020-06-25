package migmigration

import (
	"time"

	"github.com/go-logr/logr"
	migapi "github.com/konveyor/mig-controller/pkg/apis/migration/v1alpha1"
	"github.com/konveyor/mig-controller/pkg/compat"
	"github.com/konveyor/mig-controller/pkg/settings"
	"github.com/pkg/errors"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// Application settings.
var Settings = &settings.Settings

// Requeue
var FastReQ = time.Duration(time.Millisecond * 100)
var PollReQ = time.Duration(time.Second * 3)
var NoReQ = time.Duration(0)

// Phases
const (
	Created                         = ""
	Started                         = "Started"
	Prepare                         = "Prepare"
	EnsureCloudSecretPropagated     = "EnsureCloudSecretPropagated"
	PreBackupHooks                  = "PreBackupHooks"
	PostBackupHooks                 = "PostBackupHooks"
	PreRestoreHooks                 = "PreRestoreHooks"
	PostRestoreHooks                = "PostRestoreHooks"
	PreBackupHooksFailed            = "PreBackupHooksFailed"
	PostBackupHooksFailed           = "PostBackupHooksFailed"
	PreRestoreHooksFailed           = "PreRestoreHooksFailed"
	PostRestoreHooksFailed          = "PostRestoreHooksFailed"
	EnsureInitialBackup             = "EnsureInitialBackup"
	InitialBackupCreated            = "InitialBackupCreated"
	InitialBackupFailed             = "InitialBackupFailed"
	AnnotateResources               = "AnnotateResources"
	EnsureStagePodsFromRunning      = "EnsureStagePodsFromRunning"
	EnsureStagePodsFromTemplates    = "EnsureStagePodsFromTemplates"
	EnsureStagePodsFromOrphanedPVCs = "EnsureStagePodsFromOrphanedPVCs"
	StagePodsCreated                = "StagePodsCreated"
	RestartRestic                   = "RestartRestic"
	ResticRestarted                 = "ResticRestarted"
	QuiesceApplications             = "QuiesceApplications"
	EnsureQuiesced                  = "EnsureQuiesced"
	UnQuiesceApplications           = "UnQuiesceApplications"
	EnsureStageBackup               = "EnsureStageBackup"
	StageBackupCreated              = "StageBackupCreated"
	StageBackupFailed               = "StageBackupFailed"
	EnsureInitialBackupReplicated   = "EnsureInitialBackupReplicated"
	EnsureStageBackupReplicated     = "EnsureStageBackupReplicated"
	EnsureStageRestore              = "EnsureStageRestore"
	StageRestoreCreated             = "StageRestoreCreated"
	StageRestoreFailed              = "StageRestoreFailed"
	EnsureFinalRestore              = "EnsureFinalRestore"
	FinalRestoreCreated             = "FinalRestoreCreated"
	FinalRestoreFailed              = "FinalRestoreFailed"
	Verification                    = "Verification"
	EnsureStagePodsDeleted          = "EnsureStagePodsDeleted"
	EnsureStagePodsTerminated       = "EnsureStagePodsTerminated"
	EnsureAnnotationsDeleted        = "EnsureAnnotationsDeleted"
	EnsureLabelsDeleted             = "EnsureLabelsDeleted"
	EnsureMigratedDeleted           = "EnsureMigratedDeleted"
	DeleteMigrated                  = "DeleteMigrated"
	DeleteBackups                   = "DeleteBackups"
	DeleteRestores                  = "DeleteRestores"
	MigrationFailed                 = "MigrationFailed"
	Canceling                       = "Canceling"
	Canceled                        = "Canceled"
	Completed                       = "Completed"
)

// Flags
const (
	Quiesce      = 0x01 // Only when QuiescePods (true).
	HasStagePods = 0x02 // Only when stage pods created.
	HasPVs       = 0x04 // Only when PVs migrated.
	HasVerify    = 0x08 // Only when the plan has enabled verification
)

type Itinerary struct {
	Name  string
	Steps []Step
}

var StageItinerary = Itinerary{
	Name: "Stage",
	Steps: []Step{
		{phase: Created},
		{phase: Started},
		{phase: Prepare},
		{phase: EnsureCloudSecretPropagated},
		{phase: EnsureStagePodsFromRunning, all: HasPVs},
		{phase: EnsureStagePodsFromTemplates, all: HasPVs},
		{phase: EnsureStagePodsFromOrphanedPVCs, all: HasPVs},
		{phase: StagePodsCreated, all: HasStagePods},
		{phase: AnnotateResources, all: HasPVs},
		{phase: RestartRestic, all: HasStagePods},
		{phase: ResticRestarted, all: HasStagePods},
		{phase: QuiesceApplications, all: Quiesce},
		{phase: EnsureQuiesced, all: Quiesce},
		{phase: EnsureStageBackup, all: HasPVs},
		{phase: StageBackupCreated, all: HasPVs},
		{phase: EnsureStageBackupReplicated, all: HasPVs},
		{phase: EnsureStageRestore, all: HasPVs},
		{phase: StageRestoreCreated, all: HasPVs},
		{phase: EnsureStagePodsDeleted, all: HasStagePods},
		{phase: EnsureStagePodsTerminated, all: HasStagePods},
		{phase: EnsureAnnotationsDeleted, all: HasPVs},
		{phase: EnsureLabelsDeleted},
		{phase: Completed},
	},
}

var FinalItinerary = Itinerary{
	Name: "Final",
	Steps: []Step{
		{phase: Created},
		{phase: Started},
		{phase: Prepare},
		{phase: EnsureCloudSecretPropagated},
		{phase: PreBackupHooks},
		{phase: EnsureInitialBackup},
		{phase: InitialBackupCreated},
		{phase: EnsureStagePodsFromRunning, all: HasPVs},
		{phase: EnsureStagePodsFromTemplates, all: HasPVs},
		{phase: EnsureStagePodsFromOrphanedPVCs, all: HasPVs},
		{phase: StagePodsCreated, all: HasStagePods},
		{phase: AnnotateResources, all: HasPVs},
		{phase: RestartRestic, all: HasStagePods},
		{phase: ResticRestarted, all: HasStagePods},
		{phase: QuiesceApplications, all: Quiesce},
		{phase: EnsureQuiesced, all: Quiesce},
		{phase: EnsureStageBackup, all: HasPVs},
		{phase: StageBackupCreated, all: HasPVs},
		{phase: EnsureStageBackupReplicated, all: HasPVs},
		{phase: EnsureStageRestore, all: HasPVs},
		{phase: StageRestoreCreated, all: HasPVs},
		{phase: EnsureStagePodsDeleted, all: HasStagePods},
		{phase: EnsureStagePodsTerminated, all: HasStagePods},
		{phase: EnsureAnnotationsDeleted, all: HasPVs},
		{phase: EnsureInitialBackupReplicated},
		{phase: PostBackupHooks},
		{phase: PreRestoreHooks},
		{phase: EnsureFinalRestore},
		{phase: FinalRestoreCreated},
		{phase: EnsureLabelsDeleted},
		{phase: PostRestoreHooks},
		{phase: Verification, all: HasVerify},
		{phase: Completed},
	},
}

var FinalItineraryNoPVs = Itinerary{
	Name: "FinalNoPVs",
	Steps: []Step{
		{phase: Created},
		{phase: Started},
		{phase: Prepare},
		{phase: EnsureCloudSecretPropagated},
		{phase: PreBackupHooks},
		{phase: EnsureInitialBackup},
		{phase: InitialBackupCreated},
		{phase: QuiesceApplications, all: Quiesce},
		{phase: EnsureQuiesced, all: Quiesce},
		{phase: EnsureInitialBackupReplicated},
		{phase: PostBackupHooks},
		{phase: PreRestoreHooks},
		{phase: EnsureFinalRestore},
		{phase: FinalRestoreCreated},
		{phase: EnsureLabelsDeleted},
		{phase: PostRestoreHooks},
		{phase: Verification, all: HasVerify},
		{phase: Completed},
	},
}

var CancelItinerary = Itinerary{
	Name: "Cancel",
	Steps: []Step{
		{phase: Canceling},
		{phase: DeleteBackups},
		{phase: DeleteRestores},
		{phase: EnsureStagePodsDeleted, all: HasStagePods},
		{phase: EnsureAnnotationsDeleted, all: HasPVs},
		{phase: DeleteMigrated},
		{phase: EnsureMigratedDeleted},
		{phase: UnQuiesceApplications, all: Quiesce},
		{phase: Canceled},
		{phase: Completed},
	},
}

var FailedItinerary = Itinerary{
	Name: "Failed",
	Steps: []Step{
		{phase: MigrationFailed},
		{phase: EnsureStagePodsDeleted, all: HasStagePods},
		{phase: EnsureAnnotationsDeleted, all: HasPVs},
		{phase: DeleteMigrated},
		{phase: EnsureMigratedDeleted},
		{phase: UnQuiesceApplications, all: Quiesce},
		{phase: Completed},
	},
}

// Step
type Step struct {
	// A phase name.
	phase string
	// Step included when ALL flags evaluate true.
	all uint8
	// Step included when ANY flag evaluates true.
	any uint8
}

// Get a progress report.
// Returns: phase, n, total.
func (r Itinerary) progressReport(phase string) (string, int, int) {
	n := 0
	total := len(r.Steps)
	for i, step := range r.Steps {
		if step.phase == phase {
			n = i + 1
			break
		}
	}

	return phase, n, total
}

// A Velero task that provides the complete backup & restore workflow.
// Log - A controller's logger.
// Client - A controller's (local) client.
// Owner - A MigMigration resource.
// PlanResources - A PlanRefResources.
// Annotations - Map of annotations to applied to the backup & restore
// BackupResources - Resource types to be included in the backup.
// Phase - The task phase.
// Requeue - The requeueAfter duration. 0 indicates no requeue.
// Itinerary - The phase itinerary.
// Errors - Migration errors.
// Failed - Task phase has failed.
type Task struct {
	Log               logr.Logger
	Client            k8sclient.Client
	Owner             *migapi.MigMigration
	PlanResources     *migapi.PlanResources
	Annotations       map[string]string
	BackupResources   []string
	ExcludedResources []string
	Phase             string
	Requeue           time.Duration
	Itinerary         Itinerary
	Errors            []string
}

// Run the task.
// Each call will:
//   1. Run the current phase.
//   2. Update the phase to the next phase.
//   3. Set the Requeue (as appropriate).
//   4. Return.
func (t *Task) Run() error {
	t.Log.Info("[RUN]", "stage", t.stage(), "phase", t.Phase)

	t.init()

	// Run the current phase.
	switch t.Phase {
	case Created, Started:
		t.next()
	case Prepare:
		err := t.ensureStagePodsDeleted()
		if err != nil {
			log.Trace(err)
			return err
		}
		err = t.deleteAnnotations()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.next()
	case EnsureCloudSecretPropagated:
		count := 0
		for _, cluster := range t.getBothClusters() {
			propagated, err := t.veleroPodCredSecretPropagated(cluster)
			if err != nil {
				log.Trace(err)
				return err
			}
			if propagated {
				count++
			} else {
				break
			}
		}
		if count == 2 {
			t.next()
		} else {
			t.Requeue = PollReQ
		}
	case PreBackupHooks:
		status, err := t.runHooks(migapi.PreBackupHookPhase)
		if err != nil {
			log.Trace(err)
			t.fail(PreBackupHooksFailed, []string{err.Error()})
			return err
		}
		if status {
			t.next()
		} else {
			t.Requeue = NoReQ
		}
	case EnsureInitialBackup:
		_, err := t.ensureInitialBackup()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.Requeue = NoReQ
		t.next()
	case InitialBackupCreated:
		backup, err := t.getInitialBackup()
		if err != nil {
			log.Trace(err)
			return err
		}
		if backup == nil {
			return errors.New("Backup not found")
		}
		completed, reasons := t.hasBackupCompleted(backup)
		if completed {
			if len(reasons) > 0 {
				t.fail(InitialBackupFailed, reasons)
			} else {
				t.next()
			}
		} else {
			t.Requeue = NoReQ
		}
	case AnnotateResources:
		err := t.annotateStageResources()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.next()
	case EnsureStagePodsFromRunning:
		err := t.ensureStagePodsFromRunning()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.Requeue = NoReQ
		t.next()
	case EnsureStagePodsFromTemplates:
		err := t.ensureStagePodsFromTemplates()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.Requeue = NoReQ
		t.next()
	case EnsureStagePodsFromOrphanedPVCs:
		err := t.ensureStagePodsFromOrphanedPVCs()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.Requeue = NoReQ
		t.next()
	case StagePodsCreated:
		started, err := t.ensureStagePodsStarted()
		if err != nil {
			log.Trace(err)
			return err
		}
		if started {
			t.next()
		} else {
			t.Requeue = NoReQ
		}
	case RestartRestic:
		err := t.restartResticPods()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.Requeue = PollReQ
		t.next()
	case ResticRestarted:
		started, err := t.haveResticPodsStarted()
		if err != nil {
			log.Trace(err)
			return err
		}
		if started {
			t.next()
		} else {
			t.Requeue = PollReQ
		}
	case QuiesceApplications:
		err := t.quiesceApplications()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.next()
	case EnsureQuiesced:
		quiesced, err := t.ensureQuiescedPodsTerminated()
		if err != nil {
			log.Trace(err)
			return err
		}
		if quiesced {
			t.next()
		} else {
			t.Requeue = PollReQ
		}
	case UnQuiesceApplications:
		err := t.unQuiesceApplications()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.next()
	case EnsureStageBackup:
		_, err := t.ensureStageBackup()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.Requeue = NoReQ
		t.next()
	case StageBackupCreated:
		backup, err := t.getStageBackup()
		if err != nil {
			log.Trace(err)
			return err
		}
		if backup == nil {
			return errors.New("Backup not found")
		}
		completed, reasons := t.hasBackupCompleted(backup)
		if completed {
			if len(reasons) > 0 {
				t.fail(StageBackupFailed, reasons)
			} else {
				t.next()
			}
		} else {
			t.Requeue = NoReQ
		}
	case EnsureStageBackupReplicated:
		backup, err := t.getStageBackup()
		if err != nil {
			log.Trace(err)
			return err
		}
		if backup == nil {
			return errors.New("Backup not found")
		}
		replicated, err := t.isBackupReplicated(backup)
		if err != nil {
			log.Trace(err)
			return err
		}
		if replicated {
			t.next()
		} else {
			t.Requeue = NoReQ
		}
	case PostBackupHooks:
		status, err := t.runHooks(migapi.PostBackupHookPhase)
		if err != nil {
			log.Trace(err)
			t.fail(PostBackupHooksFailed, []string{err.Error()})
			return err
		}
		if status {
			t.next()
		} else {
			t.Requeue = NoReQ
		}
	case PreRestoreHooks:
		status, err := t.runHooks(migapi.PreRestoreHookPhase)
		if err != nil {
			log.Trace(err)
			t.fail(PreRestoreHooksFailed, []string{err.Error()})
			return err
		}
		if status {
			t.next()
		} else {
			t.Requeue = NoReQ
		}
	case EnsureStageRestore:
		_, err := t.ensureStageRestore()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.Requeue = NoReQ
		t.next()
	case StageRestoreCreated:
		restore, err := t.getStageRestore()
		if err != nil {
			log.Trace(err)
			return err
		}
		if restore == nil {
			return errors.New("Restore not found")
		}
		completed, reasons := t.hasRestoreCompleted(restore)
		if completed {
			t.setResticConditions(restore)
			if len(reasons) > 0 {
				t.fail(StageRestoreFailed, reasons)
			} else {
				t.next()
			}
		} else {
			t.Requeue = NoReQ
		}
	case EnsureStagePodsDeleted:
		err := t.ensureStagePodsDeleted()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.next()
	case EnsureStagePodsTerminated:
		terminated, err := t.ensureStagePodsTerminated()
		if err != nil {
			log.Trace(err)
			return err
		}
		if terminated {
			t.next()
		} else {
			t.Requeue = PollReQ
		}
	case EnsureAnnotationsDeleted:
		if !t.keepAnnotations() {
			err := t.deleteAnnotations()
			if err != nil {
				log.Trace(err)
				return err
			}
		}
		t.next()
	case EnsureLabelsDeleted:
		if !t.keepAnnotations() {
			err := t.deleteLabels()
			if err != nil {
				log.Trace(err)
				return err
			}
		}
		t.next()
	case EnsureInitialBackupReplicated:
		backup, err := t.getInitialBackup()
		if err != nil {
			log.Trace(err)
			return err
		}
		if backup == nil {
			return errors.New("Backup not found")
		}
		replicated, err := t.isBackupReplicated(backup)
		if err != nil {
			log.Trace(err)
			return err
		}
		if replicated {
			t.next()
		} else {
			t.Requeue = NoReQ
		}
	case EnsureFinalRestore:
		backup, err := t.getInitialBackup()
		if err != nil {
			log.Trace(err)
			return err
		}
		if backup == nil {
			return errors.New("Backup not found")
		}
		_, err = t.ensureFinalRestore()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.Requeue = NoReQ
		t.next()
	case FinalRestoreCreated:
		restore, err := t.getFinalRestore()
		if err != nil {
			log.Trace(err)
			return err
		}
		if restore == nil {
			return errors.New("Restore not found")
		}
		completed, reasons := t.hasRestoreCompleted(restore)
		if completed {
			if len(reasons) > 0 {
				t.fail(FinalRestoreFailed, reasons)
			} else {
				t.next()
			}
		} else {
			t.Requeue = NoReQ
		}
	case PostRestoreHooks:
		status, err := t.runHooks(migapi.PostRestoreHookPhase)
		if err != nil {
			log.Trace(err)
			t.fail(PostRestoreHooksFailed, []string{err.Error()})
			return err
		}
		if status {
			t.next()
		} else {
			t.Requeue = NoReQ
		}
	case Verification:
		completed, err := t.VerificationCompleted()
		if err != nil {
			log.Trace(err)
			return err
		}
		if completed {
			t.next()
		} else {
			t.Requeue = PollReQ
		}
	case Canceling:
		t.Owner.Status.SetCondition(migapi.Condition{
			Type:     Canceling,
			Status:   True,
			Reason:   Cancel,
			Category: Advisory,
			Message:  CancelInProgressMessage,
			Durable:  true,
		})
		t.next()

	case MigrationFailed:
		if Settings.Migration.FailureRollback {
			t.next()
		} else {
			t.Phase = Completed
		}

	case DeleteMigrated:
		err := t.deleteMigrated()
		if err != nil {
			log.Trace(err)
			return err
		}
		t.next()
	case EnsureMigratedDeleted:
		deleted, err := t.ensureMigratedResourcesDeleted()
		if err != nil {
			log.Trace(err)
			return err
		}
		if deleted {
			t.next()
		} else {
			t.Requeue = PollReQ
		}
	case DeleteBackups:
		if err := t.deleteBackups(); err != nil {
			log.Trace(err)
			return err
		}
		t.next()
	case DeleteRestores:
		if err := t.deleteRestores(); err != nil {
			log.Trace(err)
			return err
		}
		t.next()
	case Canceled:
		t.Owner.Status.DeleteCondition(Canceling)
		t.Owner.Status.SetCondition(migapi.Condition{
			Type:     Canceled,
			Status:   True,
			Reason:   Cancel,
			Category: Advisory,
			Message:  CanceledMessage,
			Durable:  true,
		})
		t.next()
	// Out of tree states - needs to be triggered manually with t.fail(...)
	case InitialBackupFailed, FinalRestoreFailed, StageBackupFailed, StageRestoreFailed:
		t.Requeue = NoReQ
		t.next()
	case Completed:
	}

	if t.Phase == Completed {
		t.Requeue = NoReQ
		t.Log.Info("[COMPLETED]")
	}

	return nil
}

// Initialize.
func (t *Task) init() {
	t.Requeue = FastReQ
	if t.failed() {
		t.Itinerary = FailedItinerary
	} else if t.canceled() {
		t.Itinerary = CancelItinerary
	} else if t.stage() {
		t.Itinerary = StageItinerary
	} else if t.PlanResources.MigPlan.MigratePVsInFinal() {
		t.Itinerary = FinalItinerary
	} else {
		t.Itinerary = FinalItineraryNoPVs
	}
	if t.Owner.Status.Itenerary != t.Itinerary.Name {
		t.Phase = t.Itinerary.Steps[0].phase
	}
	if t.stage() && !t.hasPVs() {
		t.Owner.Status.SetCondition(migapi.Condition{
			Type:     StageNoOp,
			Status:   True,
			Category: migapi.Warn,
			Message:  StageNoOpMessage,
			Durable:  true,
		})
	}
}

// Advance the task to the next phase.
func (t *Task) next() {
	current := -1
	for i, step := range t.Itinerary.Steps {
		if step.phase != t.Phase {
			continue
		}
		current = i
		break
	}
	if current == -1 {
		t.Phase = Completed
		return
	}
	for n := current + 1; n < len(t.Itinerary.Steps); n++ {
		next := t.Itinerary.Steps[n]
		if !t.allFlags(next) {
			continue
		}
		if !t.anyFlags(next) {
			continue
		}
		t.Phase = next.phase
		return
	}
	t.Phase = Completed
}

// Evaluate `all` flags.
func (t *Task) allFlags(step Step) bool {
	if step.all&HasPVs != 0 && !t.hasPVs() {
		return false
	}
	if step.all&HasStagePods != 0 && !t.Owner.Status.HasCondition(StagePodsCreated) {
		return false
	}
	if step.all&Quiesce != 0 && !t.quiesce() {
		return false
	}
	if step.all&HasVerify != 0 && !t.hasVerify() {
		return false
	}

	return true
}

// Evaluate `any` flags.
func (t *Task) anyFlags(step Step) bool {
	if step.any&HasPVs != 0 && t.hasPVs() {
		return true
	}
	if step.any&HasStagePods != 0 && t.Owner.Status.HasCondition(StagePodsCreated) {
		return true
	}
	if step.any&Quiesce != 0 && t.quiesce() {
		return true
	}
	if step.any&HasVerify != 0 && t.hasVerify() {
		return true
	}

	return step.any == uint8(0)
}

// Phase fail.
func (t *Task) fail(nextPhase string, reasons []string) {
	t.addErrors(reasons)
	t.Owner.AddErrors(t.Errors)
	t.Owner.Status.SetCondition(migapi.Condition{
		Type:     Failed,
		Status:   True,
		Reason:   t.Phase,
		Category: Advisory,
		Message:  FailedMessage,
		Durable:  true,
	})
	t.Phase = nextPhase
}

// Add errors.
func (t *Task) addErrors(errors []string) {
	for _, error := range errors {
		t.Errors = append(t.Errors, error)
	}
}

// Migration UID.
func (t *Task) UID() string {
	return string(t.Owner.UID)
}

// Get whether the migration has failed
func (t *Task) failed() bool {
	return t.Owner.HasErrors() || t.Owner.Status.HasCondition(Failed)
}

// Get whether the migration is cancelled.
func (t *Task) canceled() bool {
	return t.Owner.Spec.Canceled || t.Owner.Status.HasAnyCondition(Canceled, Canceling)
}

// Get whether the migration is stage.
func (t *Task) stage() bool {
	return t.Owner.Spec.Stage
}

// Get the migration namespaces with mapping.
func (t *Task) namespaces() []string {
	return t.PlanResources.MigPlan.Spec.Namespaces
}

// Get the migration source namespaces without mapping.
func (t *Task) sourceNamespaces() []string {
	return t.PlanResources.MigPlan.GetSourceNamespaces()
}

// Get the migration source namespaces without mapping.
func (t *Task) destinationNamespaces() []string {
	return t.PlanResources.MigPlan.GetDestinationNamespaces()
}

// Get whether to quiesce pods.
func (t *Task) quiesce() bool {
	return t.Owner.Spec.QuiescePods
}

// Get whether to retain annotations
func (t *Task) keepAnnotations() bool {
	return t.Owner.Spec.KeepAnnotations
}

// Get a client for the source cluster.
func (t *Task) getSourceClient() (compat.Client, error) {
	return t.PlanResources.SrcMigCluster.GetClient(t.Client)
}

// Get a client for the destination cluster.
func (t *Task) getDestinationClient() (compat.Client, error) {
	return t.PlanResources.DestMigCluster.GetClient(t.Client)
}

// Get the persistent volumes included in the plan which are not skipped.
func (t *Task) getPVs() migapi.PersistentVolumes {
	volumes := []migapi.PV{}
	for _, pv := range t.PlanResources.MigPlan.Spec.PersistentVolumes.List {
		if pv.Selection.Action != migapi.PvSkipAction {
			volumes = append(volumes, pv)
		}
	}
	pvList := t.PlanResources.MigPlan.Spec.PersistentVolumes.DeepCopy()
	pvList.List = volumes
	return *pvList
}

// Get the persistentVolumeClaims / action mapping included in the plan which are not skipped.
func (t *Task) getPVCs() map[k8sclient.ObjectKey]migapi.PV {
	claims := map[k8sclient.ObjectKey]migapi.PV{}
	for _, pv := range t.getPVs().List {
		claimKey := k8sclient.ObjectKey{
			Name:      pv.PVC.Name,
			Namespace: pv.PVC.Namespace,
		}
		claims[claimKey] = pv
	}
	return claims
}

// Get whether the associated plan lists not skipped PVs.
func (t *Task) hasPVs() bool {
	for _, pv := range t.PlanResources.MigPlan.Spec.PersistentVolumes.List {
		if pv.Selection.Action != migapi.PvSkipAction {
			return true
		}
	}
	return false
}

// Get whether the verification is desired
func (t *Task) hasVerify() bool {
	return t.Owner.Spec.Verify
}

// Get both source and destination clusters.
func (t *Task) getBothClusters() []*migapi.MigCluster {
	return []*migapi.MigCluster{
		t.PlanResources.SrcMigCluster,
		t.PlanResources.DestMigCluster}
}

// Get both source and destination clients.
func (t *Task) getBothClients() ([]k8sclient.Client, error) {
	list := []k8sclient.Client{}
	for _, cluster := range t.getBothClusters() {
		client, err := cluster.GetClient(t.Client)
		if err != nil {
			log.Trace(err)
			return nil, err
		}
		list = append(list, client)
	}

	return list, nil
}

// Get both source and destination clients with associated namespaces.
func (t *Task) getBothClientsWithNamespaces() ([]k8sclient.Client, [][]string, error) {
	clientList, err := t.getBothClients()
	if err != nil {
		log.Trace(err)
		return nil, nil, err
	}
	namespaceList := [][]string{t.sourceNamespaces(), t.destinationNamespaces()}

	return clientList, namespaceList, nil
}
