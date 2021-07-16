package tests

import (
	"context"
	migapi "github.com/konveyor/mig-controller/pkg/apis/migration/v1alpha1"
	"github.com/konveyor/mig-controller/pkg/pods"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"log"
	"path"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
	"time"
)

var ctx = context.Background()

var _ = Describe("Running migmigration when storage, cluster and plan are ready", func() {
	var (
		migPlan     *migapi.MigPlan
		migPlanName string
		namespaces  []string
	)

	JustBeforeEach(func() {
		Expect(hostClient.Create(ctx, migPlan)).Should(Succeed())
		Eventually(func() bool {
			log.Println("waiting for migPlan to be ready")
			hostClient.Get(ctx, client.ObjectKey{Name: migPlanName, Namespace: MigrationNamespace}, migPlan)
			return migPlan.Status.IsReady()
		}, time.Minute*5, time.Second).Should(Equal(true))
	})

	AfterEach(func() {
		err = hostClient.Delete(ctx, &migapi.MigPlan{
			ObjectMeta: metav1.ObjectMeta{
				Name:      migPlanName,
				Namespace: MigrationNamespace,
			},
		})
		if err != nil {
			log.Println(err)
		}
	})

	Context("Testing BZ #1965421, PVC name is more than 63 char long", func() {

		BeforeEach(func() {
			migPlanName = "e2e-1965421"
			namespaces = []string{"ocp-41583-longpvcname"}
			migPlan = NewMigPlan(namespaces, migPlanName)
			_, err = sourceClient.CoreV1().Namespaces().Create(ctx, NewMigrationNS("ocp-41583-longpvcname"), metav1.CreateOptions{})
			if err != nil {
				log.Println(err)
			}
			storage, err := resource.ParseQuantity("2Gi")
			if err != nil {
				log.Println(err)
			}
			pvc := &v1.PersistentVolumeClaim{

				ObjectMeta: metav1.ObjectMeta{
					Name:      "long-name-123456789011121314151617181920212223242526272829303132",
					Namespace: "ocp-41583-longpvcname",
				},
				Spec: v1.PersistentVolumeClaimSpec{
					AccessModes: []v1.PersistentVolumeAccessMode{"ReadWriteOnce"},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							"storage": storage,
						},
					},
				},
			}
			_, err = sourceClient.CoreV1().PersistentVolumeClaims("ocp-41583-longpvcname").Create(ctx, pvc, metav1.CreateOptions{})
			if err != nil {
				log.Println(err)
			}
			res := int32(1)
			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "longpvc-test",
					Namespace: "ocp-41583-longpvcname",
					Labels: map[string]string{
						"app": "longpvc-test",
					},
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &res,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "longpvc-test",
						},
					},
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "longpvc-test",
							},
						},
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								v1.Container{
									Name:            "pod-test",
									ImagePullPolicy: "Always",
									Image:           "quay.io/openshifttest/alpine",
									Command:         []string{"/bin/sh", "-c", "--"},
									Args:            []string{"while true; do sleep 30; done;"},
									VolumeMounts: []v1.VolumeMount{
										v1.VolumeMount{
											Name:      "testvolume",
											MountPath: "/data/test",
										},
									},
								},
							},
							Volumes: []v1.Volume{
								v1.Volume{
									Name: "testvolume",
									VolumeSource: v1.VolumeSource{
										PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
											ClaimName: "long-name-123456789011121314151617181920212223242526272829303132",
										},
									},
								},
							},
						},
					},
				},
			}
			_, err = sourceClient.AppsV1().Deployments("ocp-41583-longpvcname").Create(ctx, deployment, metav1.CreateOptions{})
			if err != nil {
				log.Println(err)
			}
			if err != nil {
				log.Println(err)
			}
			Eventually(func() string {
				podList, err := sourceClient.CoreV1().Pods("ocp-41583-longpvcname").List(ctx, metav1.ListOptions{})
				if err != nil {
					log.Println(err)
				}
				for _, p := range podList.Items {
					if p.Status.Phase == "Running" {
						podCmd := pods.PodCommand{
							Pod:     &p,
							RestCfg: sourceCfg,
							Args:    []string{"sh", "-c", "dd if=/dev/zero of=/data/test/file_1 bs=60M count=1"},
						}
						err = podCmd.Run()
						if err != nil {
							log.Println(err, "Failed running ls command inside destination Pod",
								"pod", path.Join(p.Namespace, p.Name),
								"command", "dd if=/dev/zero of=/data/test/file_1 bs=60M count=1;")
						}
					}
					return string(p.Status.Phase)
				}
				return ""
			}, time.Minute*5, time.Second).Should(Equal("Running"))
		})

		AfterEach(func() {
			for _, ns := range namespaces {
				err = sourceClient.CoreV1().Namespaces().Delete(ctx, ns, metav1.DeleteOptions{})
				if err != nil {
					log.Println(err)
				}
				err = hostClient.Delete(context.TODO(), &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}})
				if err != nil {
					log.Println(err)
				}
			}
		})

		It("Should create a new migration", func() {
			By("Creating a new migmigration")

			migrationName := "foo"
			migration := NewMigMigration(migrationName, migPlanName, false, false)
			Expect(hostClient.Create(ctx, migration)).Should(Succeed())

			fooMigration := &migapi.MigMigration{}

			Eventually(func() string {
				dvms := &migapi.DirectVolumeMigrationList{}
				err := hostClient.List(ctx, dvms, client.InNamespace(MigrationNamespace))
				if err != nil {
					log.Println(err)
				}
				for _, dvm := range dvms.Items {
					if dvm.GenerateName == migrationName+"-" {
						if dvm.Status.Itinerary == "VolumeMigration" {
							dvmps := &migapi.DirectVolumeMigrationProgressList{}
							err = hostClient.List(ctx, dvmps, client.MatchingLabels{"directvolumemigration": string(dvm.UID)})
							if err != nil {
								log.Println(err)
							}
							for _, dvmp := range dvmps.Items {
								if dvmp.Status.TotalProgressPercentage == "100%" {
									return dvm.Status.Phase
								}
							}
						} else {
							log.Println("Direct volume migration Failed")
						}
					}
				}
				return ""
			}, time.Minute*5, time.Second).Should(Equal("Completed"))
			Eventually(func() string {
				hostClient.Get(ctx, client.ObjectKey{Name: migrationName, Namespace: MigrationNamespace}, fooMigration)
				if fooMigration.Status.Itinerary == "Final" {
					return fooMigration.Status.Phase
				}
				return ""
			}, time.Minute*5, time.Second).Should(Equal("Completed"))

			// TODO check the destination pvc for data
			Eventually(func() bool {
				podList := &v1.PodList{}
				hostClient.List(ctx, podList, &client.ListOptions{
					LabelSelector: labels.SelectorFromSet(map[string]string{
						"app": "longpvc-test",
					},
					),
					Namespace: "ocp-41583-longpvcname",
				})
				destinationFile := []string{}
				for _, p := range podList.Items {
					podCmd := pods.PodCommand{
						Pod:     &p,
						RestCfg: hostCfg,
						Args:    []string{"sh", "-c", "ls -l data/test"},
					}
					err = podCmd.Run()
					if err != nil {
						log.Println(err, "Failed running ls command inside destination Pod",
							"pod", path.Join(p.Namespace, p.Name),
							"command", "ls -l data/test")
					}
					destinationFile = strings.Split(podCmd.Out.String(), " ")
				}
				podList, err = sourceClient.CoreV1().Pods("ocp-41583-longpvcname").List(ctx, metav1.ListOptions{
					LabelSelector: labels.SelectorFromSet(map[string]string{
						"app": "longpvc-test",
					},
					).String(),
				})
				for _, p := range podList.Items {
					podCmd := pods.PodCommand{
						Pod:     &p,
						RestCfg: sourceCfg,
						Args:    []string{"sh", "-c", "ls -l data/test"},
					}
					err = podCmd.Run()
					if err != nil {
						log.Println(err, "Failed running ls command inside source Pod",
							"pod", path.Join(p.Namespace, p.Name),
							"command", "ls data/test")
					}
					sourceFile := strings.Split(podCmd.Out.String(), " ")
					if len(sourceFile) > 0 && len(destinationFile) > 0 {
						if sourceFile[0] == destinationFile[0] && sourceFile[len(sourceFile)-1] == destinationFile[len(destinationFile)-1] {
							return true
						}
					}
				}
				return false
			}, time.Minute*5, time.Second).Should(Equal(true))
		})
	})
})
