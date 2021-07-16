package tests

import (
	"context"
	migapi "github.com/konveyor/mig-controller/pkg/apis/migration/v1alpha1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"strconv"
	"strings"
)

// common constants
const (
	E2ETestObjectName  = "e2e-object"
	MigrationNamespace = "openshift-migration"
	TestSecretType     = "Opaque"
)

// cluster constants
const (
	TestDestinationCluster = "host"
	TestClusterSecret      = "e2eclustersecret"
)

// storage constants
const (
	TestStorageSecret = "e2estoragesecret"
	ConfigNamespace   = "openshift-config"
)

// migrationcontroller constants
const (
	MigrationController = "migration-controller"
)

// environment variables
const (
	EXPOSEDREGISTRYPATH       = "EXPOSEDREGISTRYPATH"
	SOURCEURL                 = "SOURCEURL"
	SOURCECONFIG              = "SOURCECONFIG"
	HOSTCONFIG                = "KUBECONFIG"
	VELERO_PLUGIN_FQIN        = "VELERO_PLUGIN_FQIN"
	MIG_CONTROLLER_IMAGE_FQIN = "MIG_CONTROLLER_IMAGE_FQIN"
	RSYNC_TRANSFER_IMAGE_FQIN = "RSYNC_TRANSFER_IMAGE_FQIN"
	MIG_NAMESPACE_LIMIT       = "MIG_NAMESPACE_LIMIT"
	MIG_POD_LIMIT             = "MIG_POD_LIMIT"
	CLUSTER_NAME              = "CLUSTER_NAME"
	RESTIC_TIMEOUT            = "RESTIC_TIMEOUT"
	MIGRATION_VELERO          = "MIGRATION_VELERO"
	MIG_PV_LIMIT              = "MIG_PV_LIMIT"
	VERSION                   = "VERSION"
	AWSBUCKETNAME             = "AWSBUCKETNAME"
	AWSSECRETKEY              = "AWSSECRETKEY"
	AWSACCESSKEY              = "AWSACCESSKEY"
	BACKUPSTORAGEPROVIDER     = "BACKUPSTORAGEPROVIDER"
)

func NewMigMigration(name string, planRef string, quiesce bool, stage bool) *migapi.MigMigration {
	return &migapi.MigMigration{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: MigrationNamespace,
		},
		Spec: migapi.MigMigrationSpec{
			MigPlanRef: &core.ObjectReference{
				Namespace: MigrationNamespace,
				Name:      planRef,
			},
			Stage:       stage,
			QuiescePods: quiesce,
		},
	}
}

func NewMigPlan(namespaces []string, name string) *migapi.MigPlan {
	return &migapi.MigPlan{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: MigrationNamespace,
		},
		Spec: migapi.MigPlanSpec{
			Namespaces: namespaces,
			SrcMigClusterRef: &core.ObjectReference{
				Namespace: MigrationNamespace,
				Name:      E2ETestObjectName,
			},
			DestMigClusterRef: &core.ObjectReference{
				Namespace: MigrationNamespace,
				Name:      TestDestinationCluster,
			},
			MigStorageRef: &core.ObjectReference{
				Name:      E2ETestObjectName,
				Namespace: MigrationNamespace,
			},
		},
	}
}

func NewMigStorage() (*migapi.MigStorage, *core.Secret) {
	return &migapi.MigStorage{
			ObjectMeta: metav1.ObjectMeta{
				Name:      E2ETestObjectName,
				Namespace: MigrationNamespace,
			},
			Spec: migapi.MigStorageSpec{
				BackupStorageConfig: migapi.BackupStorageConfig{
					CredsSecretRef: &core.ObjectReference{
						Namespace: ConfigNamespace,
						Name:      TestStorageSecret,
					},
					AwsBucketName: os.Getenv(AWSBUCKETNAME),
				},
				BackupStorageProvider: os.Getenv(BACKUPSTORAGEPROVIDER),
				// TODO: define env variable. can these both be different?
				VolumeSnapshotProvider: os.Getenv(BACKUPSTORAGEPROVIDER),
				VolumeSnapshotConfig: migapi.VolumeSnapshotConfig{
					CredsSecretRef: &core.ObjectReference{
						Namespace: ConfigNamespace,
						Name:      TestStorageSecret,
					},
				},
			},
		}, &core.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      TestStorageSecret,
				Namespace: ConfigNamespace,
			},
			Type: TestSecretType,
			Data: map[string][]byte{
				"aws-access-key-id":     []byte(os.Getenv(AWSACCESSKEY)),
				"aws-secret-access-key": []byte(os.Getenv(AWSSECRETKEY)),
			},
		}
}

func NewMigCluster(saToken []byte) (*migapi.MigCluster, *core.Secret) {
	return &migapi.MigCluster{
			ObjectMeta: metav1.ObjectMeta{
				Name:      E2ETestObjectName,
				Namespace: MigrationNamespace,
			},
			Spec: migapi.MigClusterSpec{
				IsHostCluster: false,
				URL:           os.Getenv(SOURCEURL),
				ServiceAccountSecretRef: &core.ObjectReference{
					Namespace: ConfigNamespace,
					Name:      TestClusterSecret,
				},
				Insecure:            true,
				ExposedRegistryPath: os.Getenv(EXPOSEDREGISTRYPATH),
			},
		}, &core.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      TestClusterSecret,
				Namespace: ConfigNamespace,
			},
			Data: map[string][]byte{
				"saToken": saToken,
			},
			Type: TestSecretType,
		}
}

// We are assuming that the controller CR is created and controller is running
func NewMigrationController() *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "migration.openshift.io/v1alpha1",
			"kind":       "MigrationController",
			"metadata": map[string]interface{}{
				"name":      MigrationController,
				"namespace": MigrationNamespace,
			},
			// TODO: take variables from env
			"spec": map[string]interface{}{
				"velero_plugin_fqin":        getEnvStr(VELERO_PLUGIN_FQIN, "quay.io/konveyor/openshift-velero-plugin:latest"),
				"mig_controller_image_fqin": getEnvStr(MIG_CONTROLLER_IMAGE_FQIN, "quay.io/konveyor/mig-controller:latest"),
				"mig_namespace_limit":       getEnvInt(MIG_NAMESPACE_LIMIT, 10),
				"migration_ui":              true,
				"mig_pod_limit":             getEnvInt(MIG_POD_LIMIT, 100),
				"migration_controller":      true,
				"migration_log_reader":      true,
				"olm_managed":               true,
				"cluster_name":              getEnvStr(CLUSTER_NAME, "host"),
				"restic_timeout":            getEnvStr(RESTIC_TIMEOUT, "1h"),
				"migration_velero":          getEnvBool(MIGRATION_VELERO, true),
				"rsync_transfer_image_fqin": getEnvStr(RSYNC_TRANSFER_IMAGE_FQIN, "quay.io/konveyor/rsync-transfer:latest"),
				"mig_pv_limit":              getEnvInt(MIG_PV_LIMIT, 100),
				"version":                   getEnvStr(VERSION, "latest"),
				"azure_resource_group":      "",
			},
		},
	}
}

func NewMigrationNS(ns string) *core.Namespace {
	return &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
	}
}

func GetMigSaToken(sourceClient *kubernetes.Clientset) []byte {
	ctx := context.TODO()
	sa, err := sourceClient.CoreV1().ServiceAccounts(MigrationNamespace).Get(ctx, MigrationController, metav1.GetOptions{})
	if err != nil {
		log.Print(err)
	}
	for _, s := range sa.Secrets {
		if strings.Contains(s.Name, "token") {
			secret, err := sourceClient.CoreV1().Secrets(MigrationNamespace).Get(ctx, s.Name, metav1.GetOptions{})
			if err != nil {
				log.Println(err)
			}
			return secret.Data["token"]
		}
	}
	return []byte{}
}

func getEnvInt(name string, val int) int {
	t, err := strconv.Atoi(os.Getenv(name))
	if err != nil {
		return val
	}
	return t
}

func getEnvStr(name string, val string) string {
	t := os.Getenv(name)
	if t == "" {
		return val
	}
	return t
}

func getEnvBool(name string, val bool) bool {
	t, err := strconv.ParseBool(os.Getenv(name))
	if err != nil {
		return val
	}
	return t
}
