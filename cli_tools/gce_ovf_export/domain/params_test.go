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

package ovfexportdomain

import (
	"testing"

	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/daisycommon"
	"github.com/stretchr/testify/assert"
)

func TestIsInstanceExport(t *testing.T) {
	assert.True(t, GetAllInstanceExportParams().IsInstanceExport())
	assert.False(t, GetAllMachineImageExportParams().IsInstanceExport())
}

func TestIsMachineImageExport(t *testing.T) {
	assert.False(t, GetAllInstanceExportParams().IsMachineImageExport())
	assert.True(t, GetAllMachineImageExportParams().IsMachineImageExport())
}

func TestDaisyAttrs(t *testing.T) {
	params := GetAllInstanceExportParams()
	assert.Equal(t,
		daisycommon.WorkflowAttributes{
			Project: *params.Project, Zone: params.Zone, GCSPath: params.ScratchBucketGcsPath,
			OAuth: params.Oauth, Timeout: params.Timeout.String(), ComputeEndpoint: params.Ce,
			WorkflowDirectory: params.WorkflowDir, DisableGCSLogs: params.GcsLogsDisabled,
			DisableCloudLogs: params.CloudLogsDisabled, DisableStdoutLogs: params.StdoutLogsDisabled,
			NoExternalIP: params.NoExternalIP,
		},
		params.DaisyAttrs())
}