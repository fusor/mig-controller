/*
Copyright 2021 Red Hat Inc.

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

package miganalytic

import (
	"context"
	"fmt"
	"sync"

	"github.com/konveyor/mig-controller/pkg/compat"
	"github.com/konveyor/mig-controller/pkg/pods"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// ResticPodLabelKey is the key of the label used to discover Restic pod
	ResticPodLabelKey = "name"
	// ResticPodLabelValue is the value of the label used to discover Restic pod
	ResticPodLabelValue = "restic"
)

// ResticDFCommandExecutor uses Restic pods to run DF command
type ResticDFCommandExecutor struct {
	// Namespace is the ns in which Restic pods are present
	Namespace string
	// Client to interact with Restic pods
	Client compat.Client
	// ResticPodReferences is a local cache of known Restic pods
	ResticPodReferences map[string]*corev1.Pod
}

// DF given a podRef and a list of volumes, runs df command, returns with structured command context
// any errors running the df command are suppressed here. DFCommand.stdErr field should be used to determine failure
func (r *ResticDFCommandExecutor) DF(podRef *corev1.Pod, persistentVolumes []MigAnalyticPersistentVolumeDetails) DFCommand {
	// TODO: use the appropriate block size based on PVCs
	dfCmd := DFCommand{
		BaseLocation: "/host_pods",
		BlockSize:    DecimalSIMega,
		StdOut:       "",
		StdErr:       "",
	}
	cmdString := dfCmd.PrepareDFCommand(persistentVolumes)
	restCfg := r.Client.RestConfig()
	podCommand := pods.PodCommand{
		Pod:     podRef,
		RestCfg: restCfg,
		Args:    cmdString,
	}
	err := podCommand.Run()
	if err != nil {
		log.Info(
			fmt.Sprintf("df command inside pod %s returned non-zero error code", podRef.Name))
	}
	dfCmd.StdErr = podCommand.Err.String()
	dfCmd.StdOut = podCommand.Out.String()
	return dfCmd
}

// getResticPodForNode lookup Restic Pod ref in local cache
func (r *ResticDFCommandExecutor) getResticPodForNode(nodeName string) *corev1.Pod {
	if podRef, exists := r.ResticPodReferences[nodeName]; exists {
		return podRef
	}
	return nil
}

// loadResticPodReferences load Restic Pod refs in-memory
func (r *ResticDFCommandExecutor) loadResticPodReferences() error {
	if r.ResticPodReferences == nil {
		r.ResticPodReferences = make(map[string]*corev1.Pod)
	}
	resticPodList := corev1.PodList{}
	labelSelector := client.InNamespace(r.Namespace).MatchingLabels(
		map[string]string{
			ResticPodLabelKey: ResticPodLabelValue})
	err := r.Client.List(context.TODO(), labelSelector, &resticPodList)
	if err != nil {
		return err
	}
	for i := range resticPodList.Items {
		if resticPodList.Items[i].Spec.NodeName != "" {
			r.ResticPodReferences[resticPodList.Items[i].Spec.NodeName] = &resticPodList.Items[i]
		}
	}
	return nil
}

// Execute given a map node->[]pvc, runs Df command for each, returns list of structured df output per pvc
func (r *ResticDFCommandExecutor) Execute(pvcNodeMap map[string][]MigAnalyticPersistentVolumeDetails) ([]DFOutput, error) {
	gatheredData := []DFOutput{}
	err := r.loadResticPodReferences()
	if err != nil {
		return gatheredData, err
	}
	// dfOutputs for n nodes
	dfOutputs := make(map[string]DFCommand, len(pvcNodeMap))
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(pvcNodeMap))
	// run a buffered reader to accept DF outputs on channels
	bufferedOutputChannel, finishSignalChannel := r.runBufferedOutputReader(dfOutputs)
	// run df concurrently for 'n' nodes
	for node := range pvcNodeMap {
		resticPodRef := r.getResticPodForNode(node)
		// if no Restic pod is found for this node, all PVCs on this node are skipped
		if resticPodRef == nil {
			for _, pvc := range pvcNodeMap[node] {
				dfOutput := DFOutput{
					IsError:   true,
					Name:      pvc.Name,
					Namespace: pvc.Namespace,
				}
				gatheredData = append(gatheredData, dfOutput)
			}
			waitGroup.Done()
			continue
		}
		go func(n string, podRef *corev1.Pod) {
			defer waitGroup.Done()
			bufferedOutputChannel <- bufferedOutput{
				CmdOutput: r.DF(podRef, pvcNodeMap[n]),
				Node:      n,
			}
		}(node, resticPodRef)
	}
	// wait for all command instances to return
	waitGroup.Wait()
	// exit buffered reader
	finishSignalChannel <- struct{}{}
	for node, cmdOutput := range dfOutputs {
		for _, pvc := range pvcNodeMap[node] {
			pvcDFInfo := cmdOutput.GetDFOutputForPV(pvc.VolumeName, pvc.PodUID)
			pvcDFInfo.Node = node
			pvcDFInfo.Name = pvc.Name
			pvcDFInfo.Namespace = pvc.Namespace
			gatheredData = append(gatheredData, pvcDFInfo)
		}
	}
	return gatheredData, nil
}

// bufferedOutput is used to buffer outputs coming from concurrent DF commands
// to avoid overwhelming the apiserver in case there are too many nodes running
type bufferedOutput struct {
	CmdOutput DFCommand
	Node      string
}

// runBufferedOutputReader runs a thread to receive updates from df command instances
func (r *ResticDFCommandExecutor) runBufferedOutputReader(dfOutputs map[string]DFCommand) (chan bufferedOutput, chan struct{}) {
	// at most 10 instances of sender can run at the same time
	bufferedOutputs := make(chan bufferedOutput, 10)
	finishSignal := make(chan struct{})
	mutex := sync.Mutex{}
	go func() {
		for {
			select {
			case output := <-bufferedOutputs:
				mutex.Lock()
				dfOutputs[output.Node] = output.CmdOutput
				mutex.Unlock()
				break
			case <-finishSignal:
				return
			}
		}
	}()
	return bufferedOutputs, finishSignal
}
