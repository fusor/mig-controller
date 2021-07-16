package tests

import (
	"context"
	migapi "github.com/konveyor/mig-controller/pkg/apis/migration/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"

	// +kubebuilder:scaffold:imports
)

func TestMigmigration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}

var sourceClient *kubernetes.Clientset
var hostClient client.Client
var hostCfg *rest.Config
var sourceCfg *rest.Config
var dynamicClient dynamic.Interface
var controllerCR schema.GroupVersionResource
var err error
var shouldCreate bool

var _ = BeforeSuite(func() {
	hostCfg, err = clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv(HOSTCONFIG)))
	if err != nil {
		log.Println(err)
	}

	err = migapi.AddToScheme(scheme.Scheme)
	hostClient, err = client.New(hostCfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		log.Println(err)
	}
	sourceCfg, err = clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv(SOURCECONFIG)))
	if err != nil {
		log.Println(err)
	}
	sourceClient, err = kubernetes.NewForConfig(sourceCfg)
	if err != nil {
		log.Println(err)
	}

	dynamicClient, err = dynamic.NewForConfig(hostCfg)
	if err != nil {
		log.Println(err)
	}
	controllerCR = schema.GroupVersionResource{
		Group:    "migration.openshift.io",
		Version:  "v1alpha1",
		Resource: "migrationcontrollers",
	}

	controller, err := dynamicClient.Resource(controllerCR).Namespace(MigrationNamespace).Get(context.TODO(), MigrationController, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
	}

	shouldCreate := true
	if controller != nil {
		shouldCreate = false
	}
	if shouldCreate {
		controller, err = dynamicClient.Resource(controllerCR).Namespace(MigrationNamespace).Create(ctx, NewMigrationController(), metav1.CreateOptions{})
		if err != nil {
			log.Println(err)
		}
	}
	Eventually(func() string {
		controller, err := dynamicClient.Resource(controllerCR).Namespace(MigrationNamespace).Get(context.TODO(), MigrationController, metav1.GetOptions{})
		if err != nil {
			log.Println(err)
		}
		phase, found, err := unstructured.NestedString(controller.Object, "status", "phase")
		if err != nil || !found {
			log.Println("phase not found or error in status", err)
		}
		return phase
	}, time.Minute*5, time.Second).Should(Equal("Reconciled"))

	migCluster, secret := NewMigCluster(GetMigSaToken(sourceClient))
	Expect(hostClient.Create(ctx, secret)).Should(Succeed())
	Expect(hostClient.Create(ctx, migCluster)).Should(Succeed())
	Eventually(func() bool {
		hostClient.Get(ctx, client.ObjectKey{Name: E2ETestObjectName, Namespace: MigrationNamespace}, migCluster)
		return migCluster.Status.IsReady()
	}, time.Minute*5, time.Second).Should(Equal(true))

	migStorage, secret := NewMigStorage()
	Expect(hostClient.Create(ctx, secret)).Should(Succeed())
	Expect(hostClient.Create(ctx, migStorage)).Should(Succeed())
	Eventually(func() bool {
		hostClient.Get(ctx, client.ObjectKey{Name: E2ETestObjectName, Namespace: MigrationNamespace}, migStorage)
		return migStorage.Status.IsReady()
	}, time.Minute*5, time.Second).Should(Equal(true))

}, 60)

var _ = AfterSuite(func() {

	ctx := context.TODO()
	err = hostClient.Delete(ctx, &migapi.MigStorage{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestObjectName,
			Namespace: MigrationNamespace,
		},
	})
	if err != nil {
		log.Println(err)
	}

	err = hostClient.Delete(ctx, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TestStorageSecret,
			Namespace: ConfigNamespace,
		},
	})
	if err != nil {
		log.Println(err)
	}

	err = hostClient.Delete(ctx, &migapi.MigCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      E2ETestObjectName,
			Namespace: MigrationNamespace,
		},
	})
	if err != nil {
		log.Println(err)
	}

	err = hostClient.Delete(ctx, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TestClusterSecret,
			Namespace: ConfigNamespace,
		},
	})
	if err != nil {
		log.Println(err)
	}

	if shouldCreate {
		err = dynamicClient.Resource(controllerCR).Namespace(MigrationNamespace).Delete(ctx, MigrationController, metav1.DeleteOptions{})
		if err != nil {
			log.Println(err)
		}
	}
})
