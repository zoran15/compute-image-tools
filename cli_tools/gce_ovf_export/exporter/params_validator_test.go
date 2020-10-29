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
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/gce_ovf_export/domain"
	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var validReleaseTracks = []string{"ga", "beta", "alpha"}

func TestInstanceNameAndMachineImageNameProvidedAtTheSameTime(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	params := ovfexportdomain.GetAllInstanceExportParams()
	params.MachineImageName = "machine-image-name"
	assertErrorOnValidate(t, params, createDefaultParamValidator(mockCtrl, false))
}

func TestInstanceExportFlagsInstanceNameNotProvided(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	params := ovfexportdomain.GetAllInstanceExportParams()
	params.InstanceName = ""
	assertErrorOnValidate(t, params, createDefaultParamValidator(mockCtrl, false))
}

func TestInstanceExportFlagsOvfGcsPathFlagKeyNotProvided(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	params := ovfexportdomain.GetAllInstanceExportParams()
	params.DestinationURI = ""
	assertErrorOnValidate(t, params, createDefaultParamValidator(mockCtrl, false))
}

func TestInstanceExportFlagsOvfGcsPathFlagNotValid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	params := ovfexportdomain.GetAllInstanceExportParams()
	params.DestinationURI = "NOT_GCS_PATH"
	assertErrorOnValidate(t, params, createDefaultParamValidator(mockCtrl, false))
}

func TestInstanceExportZoneInvalid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	params := ovfexportdomain.GetAllInstanceExportParams()
	params.Zone = "not-a-zone"

	mockZoneValidator := mocks.NewMockZoneValidatorInterface(mockCtrl)
	zoneError := fmt.Errorf("invalid zone")
	mockZoneValidator.EXPECT().ZoneValid(ovfexportdomain.TestProject, params.Zone).Return(zoneError)
	validator := &ovfExportParamValidatorImpl{
		validReleaseTracks: validReleaseTracks,
		zoneValidator:      mockZoneValidator,
	}

	assert.Equal(t, zoneError, validator.ValidateAndParseParams(params))
}

func TestInstanceExportFlagsClientIdNotProvided(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	params := ovfexportdomain.GetAllInstanceExportParams()
	params.ClientID = ""
	assertErrorOnValidate(t, params, createDefaultParamValidator(mockCtrl, false))
}

func TestInstanceExportFlagsInvalidReleaseTrack(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	params := ovfexportdomain.GetAllInstanceExportParams()
	params.ReleaseTrack = "not-a-release-track"
	assertErrorOnValidate(t, params, createDefaultParamValidator(mockCtrl, false))
}

func TestInstanceExportFlagsAllValid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	validator := createDefaultParamValidator(mockCtrl, true)
	assert.Nil(t, validator.ValidateAndParseParams(ovfexportdomain.GetAllInstanceExportParams()))
}

func TestInstanceExportFlagsAllValidBucketOnlyPathTrailingSlash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	validator := createDefaultParamValidator(mockCtrl, true)

	params := ovfexportdomain.GetAllInstanceExportParams()
	params.DestinationURI = "gs://bucket_name/"
	assert.Nil(t, validator.ValidateAndParseParams(ovfexportdomain.GetAllInstanceExportParams()))
}

func TestInstanceExportFlagsAllValidBucketOnlyPathNoTrailingSlash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	validator := createDefaultParamValidator(mockCtrl, true)

	params := ovfexportdomain.GetAllInstanceExportParams()
	params.DestinationURI = "gs://bucket_name"
	assert.Nil(t, validator.ValidateAndParseParams(ovfexportdomain.GetAllInstanceExportParams()))
}

func TestInstanceExportFlagsInvalidOvfFormat(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	validator := createDefaultParamValidator(mockCtrl, false)
	params := ovfexportdomain.GetAllInstanceExportParams()
	params.OvfFormat = "zip"
	assertErrorOnValidate(t, params, validator)
}

func TestMachineImageExportFlagsAllValid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	validator := createDefaultParamValidator(mockCtrl, true)
	assert.Nil(t, validator.ValidateAndParseParams(ovfexportdomain.GetAllMachineImageExportParams()))
}

func TestMachineImageExportFlagsMachineImageNameNotProvided(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	validator := createDefaultParamValidator(mockCtrl, false)
	params := ovfexportdomain.GetAllMachineImageExportParams()
	params.MachineImageName = ""
	assertErrorOnValidate(t, params, validator)
}

func TestMachineImageExportFlagsOvfGcsPathFlagKeyNotProvided(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	validator := createDefaultParamValidator(mockCtrl, false)
	params := ovfexportdomain.GetAllMachineImageExportParams()
	params.DestinationURI = ""
	assertErrorOnValidate(t, params, validator)
}

func TestMachineImageExportFlagsOvfGcsPathFlagNotValid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	validator := createDefaultParamValidator(mockCtrl, false)
	params := ovfexportdomain.GetAllMachineImageExportParams()
	params.DestinationURI = "NOT_GCS_PATH"
	assertErrorOnValidate(t, params, validator)
}

func TestMachineImageExportFlagsClientIdNotProvided(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	validator := createDefaultParamValidator(mockCtrl, false)
	params := ovfexportdomain.GetAllMachineImageExportParams()
	params.ClientID = ""
	assertErrorOnValidate(t, params, validator)
}

func TestMachineImageExportFlagsAllValidBucketOnlyPathTrailingSlash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	validator := createDefaultParamValidator(mockCtrl, true)
	params := ovfexportdomain.GetAllMachineImageExportParams()
	params.DestinationURI = "gs://bucket_name/"
	assert.Nil(t, validator.ValidateAndParseParams(ovfexportdomain.GetAllMachineImageExportParams()))
}

func TestMachineImageExportFlagsAllValidBucketOnlyPathNoTrailingSlash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	validator := createDefaultParamValidator(mockCtrl, true)
	params := ovfexportdomain.GetAllMachineImageExportParams()
	params.DestinationURI = "gs://bucket_name"
	assert.Nil(t, validator.ValidateAndParseParams(ovfexportdomain.GetAllMachineImageExportParams()))
}

func assertErrorOnValidate(t *testing.T, params *ovfexportdomain.OVFExportParams, validator *ovfExportParamValidatorImpl) {
	assert.NotNil(t, validator.ValidateAndParseParams(params))
}

func TestMachineImageExportFlagsInvalidOvfFormat(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	validator := createDefaultParamValidator(mockCtrl, false)
	params := ovfexportdomain.GetAllMachineImageExportParams()
	params.OvfFormat = "zip"
	assertErrorOnValidate(t, params, validator)
}

func createDefaultParamValidator(mockCtrl *gomock.Controller, validateZone bool) *ovfExportParamValidatorImpl {
	mockZoneValidator := mocks.NewMockZoneValidatorInterface(mockCtrl)
	if validateZone {
		mockZoneValidator.EXPECT().ZoneValid(ovfexportdomain.TestProject, ovfexportdomain.TestZone).Return(nil)
	}
	validator := &ovfExportParamValidatorImpl{
		validReleaseTracks: validReleaseTracks,
		zoneValidator:      mockZoneValidator,
	}
	return validator
}
