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

package directvolumemigrationprogress

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	migapi "github.com/konveyor/mig-controller/pkg/apis/migration/v1alpha1"
	"github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var c client.Client

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo", Namespace: "default"}}

//var depKey = types.NamespacedName{Name: "foo-deployment", Namespace: "default"}

const timeout = time.Second * 5

func TestReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	instance := &migapi.DirectVolumeMigrationProgress{ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "default"}}

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c = mgr.GetClient()

	recFn, requests := SetupTestReconcile(newReconciler(mgr))
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Create the DirectVolumeMigrationProgress object and expect the Reconcile
	err = c.Create(context.TODO(), instance)
	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())
	defer c.Delete(context.TODO(), instance)
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))
}

func Test_parseLogs(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{reader: bytes.NewBufferString(`
          1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)2020/11/03 23:16:34 [1] <f+++++++++ file76
          1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)2020/11/03 23:16:34 [1] <f+++++++++ file76
          1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)2020/11/03 23:16:34 [1] <f+++++++++ file76
          1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)2020/11/03 23:16:34 [1] <f+++++++++ file76
          1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)2020/11/03 23:16:34 [1] <f+++++++++ file76
          1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)2020/11/03 23:16:34 [1] <f+++++++++ file76
          1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)2020/11/03 23:16:34 [1] <f+++++++++ file76
          1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)2020/11/03 23:16:34 [1] <f+++++++++ file76
          1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)2020/11/03 23:16:34 [1] <f+++++++++ file76
          1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)2020/11/03 23:16:34 [1] <f+++++++++ file76
          1.57G  89%   40.90MB/s    0:00:36 (xfr#120, to-chk=42/163)2020/11/04 01:49:41 [1] <f+++++++++ file57
          1.57G  90%   40.91MB/s    0:00:36 (xfr#121, to-chk=41/163)2020/11/04 01:49:41 [1] <f+++++++++ file58
          1.58G  90%   40.91MB/s    0:00:36 (xfr#122, to-chk=40/163)2020/11/04 01:49:41 [1] <f+++++++++ file59
          1.58G  90%   40.91MB/s    0:00:36 (xfr#123, to-chk=39/163)2020/11/04 01:49:41 [1] <f+++++++++ file60
          1.59G  90%   40.91MB/s    0:00:36 (xfr#124, to-chk=38/163)2020/11/04 01:49:42 [1] <f+++++++++ file61
          1.59G  91%   40.91MB/s    0:00:37 (xfr#125, to-chk=37/163)2020/11/04 01:49:42 [1] <f+++++++++ file62
          1.59G  91%   40.91MB/s    0:00:37 (xfr#126, to-chk=36/163)2020/11/04 01:49:42 [1] <f+++++++++ file63
          1.60G  91%   40.91MB/s    0:00:37 (xfr#127, to-chk=35/163)2020/11/04 01:49:42 [1] <f+++++++++ file64
          1.60G  91%   40.91MB/s    0:00:37 (xfr#128, to-chk=34/163)2020/11/04 01:49:42 [1] <f+++++++++ file65
          1.61G  92%   40.91MB/s    0:00:37 (xfr#129, to-chk=33/163)2020/11/04 01:49:42 [1] <f+++++++++ file66`)},
			want: strings.Join([]string{
				"1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)20",
				"1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)20",
				"1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)20",
				"1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)20",
				"1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)20",
				"1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)20",
				"1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)20",
				"1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)20",
				"1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)20",
				"1.65G  94%   40.95MB/s    0:00:38 (xfr#139, to-chk=23/163)20",
				"1.57G  89%   40.90MB/s    0:00:36 (xfr#120, to-chk=42/163)20",
				"1.57G  90%   40.91MB/s    0:00:36 (xfr#121, to-chk=41/163)20",
				"1.58G  90%   40.91MB/s    0:00:36 (xfr#122, to-chk=40/163)20",
				"1.58G  90%   40.91MB/s    0:00:36 (xfr#123, to-chk=39/163)20",
				"1.59G  90%   40.91MB/s    0:00:36 (xfr#124, to-chk=38/163)20",
				"1.59G  91%   40.91MB/s    0:00:37 (xfr#125, to-chk=37/163)20",
				"1.59G  91%   40.91MB/s    0:00:37 (xfr#126, to-chk=36/163)20",
				"1.60G  91%   40.91MB/s    0:00:37 (xfr#127, to-chk=35/163)20",
				"1.60G  91%   40.91MB/s    0:00:37 (xfr#128, to-chk=34/163)20",
				"1.61G  92%   40.91MB/s    0:00:37 (xfr#129, to-chk=33/163)20",
			}, "\n"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLogs(tt.args.reader)
			if tt.wantErr && err == nil {
				t.Errorf("parseLogs() = want error, got a nil error")
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLogs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReconcileDirectVolumeMigrationProgress_getCumulativeProgressPercentage(t *testing.T) {
	type args struct {
		pvProgress *migapi.DirectVolumeMigrationProgress
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "given all valid distinct rsync progress percentages, should return addition of all percentages",
			args: args{
				pvProgress: &migapi.DirectVolumeMigrationProgress{
					Status: migapi.DirectVolumeMigrationProgressStatus{
						RsyncPodStatus: migapi.RsyncPodStatus{PodName: "pod-4", LastObservedProgressPercent: "10%"},
						RsyncPodStatuses: []migapi.RsyncPodStatus{
							{PodName: "pod-1", LastObservedProgressPercent: "0%"},
							{PodName: "pod-2", LastObservedProgressPercent: "10%"},
							{PodName: "pod-3", LastObservedProgressPercent: "20%"},
						},
					},
				},
			},
			want: "40%",
		},
		{
			name: "given all valid, some indistinct rsync progress percentages, should return addition of all percentages, should not count the same pod twice",
			args: args{
				pvProgress: &migapi.DirectVolumeMigrationProgress{
					Status: migapi.DirectVolumeMigrationProgressStatus{
						RsyncPodStatus: migapi.RsyncPodStatus{PodName: "pod-3", LastObservedProgressPercent: "10%"},
						RsyncPodStatuses: []migapi.RsyncPodStatus{
							{PodName: "pod-1", LastObservedProgressPercent: "0%"},
							{PodName: "pod-2", LastObservedProgressPercent: "10%"},
							{PodName: "pod-3", LastObservedProgressPercent: "10%"},
						},
					},
				},
			},
			want: "20%",
		},
		{
			name: "given some valid and some invalid, some distinct and some indistinct rsync progress percentages, should return addition of all percentages, should not count the same pod twice, should not include the invalid percentage",
			args: args{
				pvProgress: &migapi.DirectVolumeMigrationProgress{
					Status: migapi.DirectVolumeMigrationProgressStatus{
						RsyncPodStatus: migapi.RsyncPodStatus{PodName: "pod-4", LastObservedProgressPercent: "10%"},
						RsyncPodStatuses: []migapi.RsyncPodStatus{
							{PodName: "pod-1", LastObservedProgressPercent: "20%"},
							{PodName: "pod-2", LastObservedProgressPercent: "A%"},
							{PodName: "pod-3"},
							{PodName: "pod-4", LastObservedProgressPercent: "10%"},
						},
					},
				},
			},
			want: "30%",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RsyncPodProgressTask{}
			if got := r.getCumulativeProgressPercentage(tt.args.pvProgress); got != tt.want {
				t.Errorf("ReconcileDirectVolumeMigrationProgress.getCumulativeProgressPercentage() = %v, want %v", got, tt.want)
			}
		})
	}
}
