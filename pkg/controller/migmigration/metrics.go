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
	"time"

	migapi "github.com/konveyor/mig-controller/pkg/apis/migration/v1alpha1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// 'status' - [ idle, running, completed, error ]
	// 'type'   - [ stage, final ]
	migrationGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cam_app_workload_migrations",
		Help: "Count of MigMigrations sorted by status and type",
	},
		[]string{"type", "status"},
	)
)

func recordMetrics(client client.Client) {
	const (
		// Metrics const values
		//   Separate from mig-controller consts to keep a stable interface for metrics systems
		//   configured to pull from static metrics endpoints.

		// Migration Type
		stage = "stage"
		final = "final"

		// Migration Status
		idle      = "idle"
		running   = "running"
		completed = "completed"
		failed    = "failed"
	)

	go func() {
		for {
			time.Sleep(10 * time.Second)

			// get all migmigration objects
			migrations, err := migapi.ListMigrations(client)

			// if error occurs, retry 10 seconds later
			if err != nil {
				continue
			}

			// Holding counter vars used to make gauge update "atomic"
			var stageIdle, stageRunning, stageCompleted, stageFailed float64
			var finalIdle, finalRunning, finalCompleted, finalFailed float64

			// for all migmigrations, count # in idle, running, completed, failed
			for _, m := range migrations {
				// Stage
				if m.Spec.Stage && m.Status.HasCondition(migapi.Running) {
					stageRunning++
					continue
				}
				if m.Spec.Stage && m.Status.HasCondition(migapi.Succeeded) {
					stageCompleted++
					continue
				}
				if m.Spec.Stage && m.Status.HasCondition(migapi.Failed) {
					stageFailed++
					continue
				}
				if m.Spec.Stage {
					stageIdle++
					continue
				}

				// Final
				if !m.Spec.Stage && m.Status.HasCondition(migapi.Running) {
					finalRunning++
					continue
				}
				if !m.Spec.Stage && m.Status.HasCondition(migapi.Succeeded) {
					finalCompleted++
					continue
				}
				if !m.Spec.Stage && m.Status.HasCondition(migapi.Failed) {
					finalFailed++
					continue
				}
				if !m.Spec.Stage {
					finalIdle++
					continue
				}
			}

			// Stage
			if stageIdle > 0 {
				migrationGauge.With(
					prometheus.Labels{"type": stage, "status": idle}).Set(stageIdle)
			}
			if stageRunning > 0 {

				migrationGauge.With(
					prometheus.Labels{"type": stage, "status": running}).Set(stageRunning)
			}
			// We probably care if someone has Mig installed but no failed or completed migrations
			migrationGauge.With(
				prometheus.Labels{"type": stage, "status": completed}).Set(stageCompleted)
			migrationGauge.With(
				prometheus.Labels{"type": stage, "status": failed}).Set(stageFailed)

			// Final
			if finalIdle > 0 {
				migrationGauge.With(
					prometheus.Labels{"type": final, "status": idle}).Set(finalIdle)
			}
			if finalRunning > 0 {
				migrationGauge.With(
					prometheus.Labels{"type": final, "status": running}).Set(finalRunning)
			}
			migrationGauge.With(
				prometheus.Labels{"type": final, "status": completed}).Set(finalCompleted)
			migrationGauge.With(
				prometheus.Labels{"type": final, "status": failed}).Set(finalFailed)
		}
	}()
}
