package migplan

import (
	"context"
	"errors"
	migapi "github.com/fusor/mig-controller/pkg/apis/migration/v1alpha1"
	velero "github.com/heptio/velero/pkg/apis/velero/v1"
	kapi "k8s.io/api/core/v1"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// Create the velero BackupStorageLocation(s) and VolumeSnapshotLocation(s)
// have been created on both the source and destination clusters associated
// with the migration plan.
func (r ReconcileMigPlan) ensureStorage(plan *migapi.MigPlan) error {
	var client k8sclient.Client
	nEnsured := 0

	if plan.Status.HasCriticalCondition() || plan.Status.HasAnyCondition(Suspended) {
		plan.Status.StageCondition(StorageEnsured)
		return nil
	}
	storage, err := plan.GetStorage(r)
	if err != nil {
		log.Trace(err)
		return err
	}
	if storage == nil || !storage.Status.IsReady() {
		return nil
	}
	clusters, err := r.planClusters(plan)
	if err != nil {
		log.Trace(err)
		return err
	}

	for _, cluster := range clusters {
		if !cluster.Status.IsReady() {
			continue
		}
		client, err = cluster.GetClient(r)
		if err != nil {
			log.Trace(err)
			return err
		}
		pl := PlanStorage{
			Client:        r,
			targetCluster: &cluster,
			targetClient:  client,
			storage:       storage,
			plan:          plan,
		}

		// BSL
		err := pl.ensureBSL()
		if err != nil {
			log.Trace(err)
			return err
		}

		// VSL
		err = pl.ensureVSL()
		if err != nil {
			log.Trace(err)
			return err
		}

		// Cloud Secret
		err = pl.ensureCloudSecret()
		if err != nil {
			log.Trace(err)
			return err
		}

		nEnsured++
	}

	// Condition
	ensured := nEnsured == 2 // Both clusters.
	if ensured {
		plan.Status.SetCondition(migapi.Condition{
			Type:     StorageEnsured,
			Status:   True,
			Category: migapi.Required,
			Message:  StorageEnsuredMessage,
		})
	}

	return err
}

// Get clusters referenced by the plan.
func (r ReconcileMigPlan) planClusters(plan *migapi.MigPlan) ([]migapi.MigCluster, error) {
	list := []migapi.MigCluster{}
	// Source
	cluster, err := plan.GetSourceCluster(r)
	if err != nil {
		return nil, err
	}
	if cluster != nil {
		list = append(list, *cluster)
	}
	// Destination
	cluster, err = plan.GetDestinationCluster(r)
	if err != nil {
		return nil, err
	}
	if cluster != nil {
		list = append(list, *cluster)
	}
	return list, nil
}

//
// PlanStorage
// Client: The controller client.
// targetClient: A client for a cluster.
// plan: A plan resource.
// storage: A storage resource.
//
type PlanStorage struct {
	k8sclient.Client
	targetCluster *migapi.MigCluster
	targetClient  k8sclient.Client
	plan          *migapi.MigPlan
	storage       *migapi.MigStorage
}

// Create the velero BackupStorageLocation has been created.
func (r PlanStorage) ensureBSL() error {
	newBSL := r.BuildBSL()
	newBSL.Labels = r.plan.GetCorrelationLabels()
	foundBSL, err := r.plan.GetBSL(r.targetClient)
	if err != nil {
		log.Trace(err)
		return err
	}
	if foundBSL == nil {
		err = r.targetClient.Create(context.TODO(), newBSL)
		if err != nil {
			log.Trace(err)
			return err
		}
		return nil
	}
	if r.storage.EqualsBSL(foundBSL, newBSL) {
		return nil
	}
	r.UpdateBSL(foundBSL)
	err = r.targetClient.Update(context.TODO(), foundBSL)
	if err != nil {
		log.Trace(err)
		return err
	}

	return nil
}

// Create the velero VolumeSnapshotLocation has been created.
func (r PlanStorage) ensureVSL() error {
	newVSL := r.BuildVSL()
	newVSL.Labels = r.plan.GetCorrelationLabels()
	foundVSL, err := r.plan.GetVSL(r.targetClient)
	if err != nil {
		log.Trace(err)
		return err
	}
	if foundVSL == nil {
		err = r.targetClient.Create(context.TODO(), newVSL)
		if err != nil {
			log.Trace(err)
			return err
		}
		return nil
	}
	if r.storage.EqualsVSL(foundVSL, newVSL) {
		return nil
	}
	r.UpdateVSL(foundVSL)
	err = r.targetClient.Update(context.TODO(), foundVSL)
	if err != nil {
		log.Trace(err)
		return err
	}

	return nil
}

// Create the velero BSL cloud secret has been created.
func (r PlanStorage) ensureCloudSecret() error {
	newSecret := r.storage.BuildBSLCloudSecret()
	newSecret.Labels = r.plan.GetCorrelationLabels()
	foundSecret, err := r.plan.GetCloudSecret(r.targetClient)
	if err != nil {
		log.Trace(err)
		return err
	}
	if foundSecret == nil {
		err = r.targetClient.Create(context.TODO(), newSecret)
		if err != nil {
			log.Trace(err)
			return err
		}
		return nil
	}
	if r.storage.EqualsCloudSecret(foundSecret, newSecret) {
		return nil
	}
	r.UpdateBSLCloudSecret(foundSecret)
	err = r.targetClient.Update(context.TODO(), foundSecret)
	if err != nil {
		log.Trace(err)
		return err
	}

	return nil
}

// Build BSL.
func (r *PlanStorage) BuildBSL() *velero.BackupStorageLocation {
	bsl := r.storage.BuildBSL()
	r.UpdateBSL(bsl)
	return bsl
}

// Update BSL.
func (r *PlanStorage) UpdateBSL(bsl *velero.BackupStorageLocation) {
	provider := r.storage.GetBackupStorageProvider()
	r.targetCluster.UpdateProvider(provider)
	provider.UpdateBSL(bsl)
}

// Build VSL.
func (r *PlanStorage) BuildVSL() *velero.VolumeSnapshotLocation {
	vsl := r.storage.BuildVSL(string(r.plan.UID))
	r.UpdateVSL(vsl)
	return vsl
}

// Update VSL.
func (r *PlanStorage) UpdateVSL(vsl *velero.VolumeSnapshotLocation) {
	provider := r.storage.GetVolumeSnapshotProvider()
	r.targetCluster.UpdateProvider(provider)
	provider.UpdateVSL(vsl)
}

// Build BSL cloud secret.
func (r *PlanStorage) BuildBSLCloudSecret() (*kapi.Secret, error) {
	secret := r.storage.BuildVSLCloudSecret()
	err := r.UpdateBSLCloudSecret(secret)
	if err != nil {
		log.Trace(err)
		return nil, err
	}

	return secret, nil
}

// Update backup cloud-secret.
func (r *PlanStorage) UpdateBSLCloudSecret(cloudSecret *kapi.Secret) error {
	secret, err := r.storage.GetBackupStorageCredSecret(r.Client)
	if err != nil {
		log.Trace(err)
		return err
	}
	if secret == nil {
		return errors.New("Credentials secret not found.")
	}
	provider := r.storage.GetBackupStorageProvider()
	r.targetCluster.UpdateProvider(provider)
	provider.UpdateCloudSecret(secret, cloudSecret)
	return nil
}

// Build VSL cloud secret.
func (r *PlanStorage) BuildVSLCloudSecret() (*kapi.Secret, error) {
	secret := r.storage.BuildVSLCloudSecret()
	err := r.UpdateVSLCloudSecret(secret)
	if err != nil {
		log.Trace(err)
		return nil, err
	}

	return secret, nil
}

// Update snapshot cloud-secret.
func (r *PlanStorage) UpdateVSLCloudSecret(cloudSecret *kapi.Secret) error {
	secret, err := r.storage.GetVolumeSnapshotCredSecret(r.Client)
	if err != nil {
		log.Trace(err)
		return err
	}
	if secret == nil {
		return errors.New("Credentials secret not found.")
	}
	provider := r.storage.GetVolumeSnapshotProvider()
	r.targetCluster.UpdateProvider(provider)
	provider.UpdateCloudSecret(secret, cloudSecret)
	return nil
}
