//  Copyright 2020 Google Inc. All Rights Reserved.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package ovfexporter

import (
	"context"
	"fmt"
	"strings"

	daisyutils "github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/utils/daisy"
	ovfexportdomain "github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/gce_ovf_export/domain"
	"github.com/GoogleCloudPlatform/compute-image-tools/daisy"
	"google.golang.org/api/compute/v1"
)

type instanceExportCleanerImpl struct {
	wf              *daisy.Workflow
	attachDiskWfs   []*daisy.Workflow
	startInstanceWf *daisy.Workflow
	serialLogs      []string
}

// NewInstanceExportCleaner creates a new instance export cleaner
func NewInstanceExportCleaner() ovfexportdomain.InstanceExportCleaner {
	return &instanceExportCleanerImpl{}
}

func (iec *instanceExportCleanerImpl) init(instance *compute.Instance, params *ovfexportdomain.OVFExportParams) error {
	// don't use default timeout as it might not be long enough for cleanup,
	// e.g. if it's very short (e.g. 10s)
	wfTimeout := "10m"
	attachedDisks := instance.Disks

	for _, attachedDisk := range attachedDisks {
		attachDiskWf, err := generateWorkflowWithSteps(fmt.Sprintf("ovf-export-clean-attach-disk-%v", attachedDisk.DeviceName),
			wfTimeout,
			func(w *daisy.Workflow) error {
				iec.addAttachDiskStepStep(w, instance, params, attachedDisk)
				return nil
			}, params)
		if err != nil {
			return err
		}
		iec.attachDiskWfs = append(iec.attachDiskWfs, attachDiskWf)
	}

	wasInstanceRunningBeforeExport := isInstanceRunning(instance)
	if wasInstanceRunningBeforeExport {
		var err error
		iec.startInstanceWf, err = generateWorkflowWithSteps("ovf-export-clean-start-instance", wfTimeout,
			func(w *daisy.Workflow) error {
				if wasInstanceRunningBeforeExport {
					iec.addStartInstanceStep(w, instance, params)
				}
				return nil
			}, params)
		if err != nil {
			return err
		}
	}

	return nil
}

func (iec *instanceExportCleanerImpl) addStartInstanceStep(w *daisy.Workflow,
	instance *compute.Instance, params *ovfexportdomain.OVFExportParams) {
	stepName := "start-instance"
	startInstanceStep := daisy.NewStepDefaultTimeout(stepName, w)
	startInstanceStep.StartInstances = &daisy.StartInstances{Instances: []string{getInstancePath(instance, *params.Project)}}
	w.Steps[stepName] = startInstanceStep
}

func (iec *instanceExportCleanerImpl) addAttachDiskStepStep(w *daisy.Workflow,
	instance *compute.Instance, params *ovfexportdomain.OVFExportParams, attachedDisk *compute.AttachedDisk) {
	diskPath := attachedDisk.Source[strings.Index(attachedDisk.Source, "projects/"):]
	attachDiskStepName := "attach-disk"
	attachDiskStep := daisy.NewStepDefaultTimeout(attachDiskStepName, w)
	attachDiskStep.AttachDisks = &daisy.AttachDisks{
		{
			Instance: getInstancePath(instance, *params.Project),
			AttachedDisk: compute.AttachedDisk{
				Mode:       attachedDisk.Mode,
				Source:     diskPath,
				Boot:       attachedDisk.Boot,
				DeviceName: attachedDisk.DeviceName,
			},
		},
	}
	w.Steps[attachDiskStepName] = attachDiskStep
}

func (iec *instanceExportCleanerImpl) Clean(instance *compute.Instance, params *ovfexportdomain.OVFExportParams) error {
	err := iec.init(instance, params)
	if err != nil {
		return err
	}
	for _, attachDiskWf := range iec.attachDiskWfs {
		// ignore errors as these will be due to instance being already started or disks already attached
		_ = daisyutils.RunWorkflowWithCancelSignal(context.Background(), attachDiskWf)
		if attachDiskWf.Logger != nil {
			iec.serialLogs = append(iec.serialLogs, attachDiskWf.Logger.ReadSerialPortLogs()...)
		}
	}
	if iec.startInstanceWf != nil {
		err = daisyutils.RunWorkflowWithCancelSignal(context.Background(), iec.startInstanceWf)
		if iec.startInstanceWf.Logger != nil {
			iec.serialLogs = append(iec.serialLogs, iec.startInstanceWf.Logger.ReadSerialPortLogs()...)
		}
	}
	return err
}

func (iec *instanceExportCleanerImpl) TraceLogs() []string {
	return iec.serialLogs
}

func (iec *instanceExportCleanerImpl) Cancel(reason string) bool {
	// cleaner is not cancelable
	return false
}