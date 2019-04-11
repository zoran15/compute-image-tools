//  Copyright 2019 Google Inc. All Rights Reserved.
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

package ovfutils

import (
	"fmt"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/gce_ovf_import/ovf_model"
	"github.com/GoogleCloudPlatform/compute-image-tools/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	ovfPath         = "gs://bucket/folder/descriptor.ovf"
	ovfObjectHandle = &storage.ObjectHandle{}

	infoStr             = "INFO_STR"
	annotationStr       = "ANNOTATION_STR"
	infoSectionRequired = false

	ovfDescriptorStr = fmt.Sprintf(
		"<Descriptor><AnnotationSection ovf:required='false'><Info>%v</Info><Annotation>%v</Annotation></AnnotationSection></Descriptor>",
		infoStr, annotationStr)
	ovfDescriptor = &ovfmodel.Descriptor{
		References: nil,
		Annotation: &ovfmodel.AnnotationSection{
			Section: ovfmodel.Section{
				Required: &infoSectionRequired,
				Info:     infoStr,
			}, Annotation: annotationStr,
		},
	}
)

func TestOvfDescriptorLoader(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStorageClient := mocks.NewMockStorageClientInterface(mockCtrl)
	mockStorageClient.EXPECT().FindGcsFile(ovfPath, ".ovf").Return(ovfObjectHandle, nil).Times(1)
	mockStorageClient.EXPECT().GetGcsFileContent(ovfObjectHandle).Return([]byte(ovfDescriptorStr), nil).Times(1)

	mockOvfDescriptorValidator := mocks.NewMockOvfDescriptorValidatorInterface(mockCtrl)
	mockOvfDescriptorValidator.EXPECT().ValidateOvfPackage(ovfDescriptor, ovfPath).Return(ovfDescriptor, nil).Times(1)

	l := OvfDescriptorLoader{storageClient: mockStorageClient, validator: mockOvfDescriptorValidator}
	result, resultError := l.Load(ovfPath)

	assert.Equal(t, result, ovfDescriptor)
	assert.Nil(t, resultError)
}

func TestOvfDescriptorLoaderNoDescriptorInGcs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	err := fmt.Errorf("no OVF file")
	mockStorageClient := mocks.NewMockStorageClientInterface(mockCtrl)
	mockStorageClient.EXPECT().FindGcsFile(ovfPath, ".ovf").Return(nil, err).Times(1)
	mockOvfDescriptorValidator := mocks.NewMockOvfDescriptorValidatorInterface(mockCtrl)

	l := OvfDescriptorLoader{storageClient: mockStorageClient, validator: mockOvfDescriptorValidator}
	result, resultError := l.Load(ovfPath)

	assert.Equal(t, err, resultError)
	assert.Nil(t, result)
}

func TestOvfDescriptorLoaderErrorLoadingDescriptor(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	err := fmt.Errorf("error loading descriptor")
	mockStorageClient := mocks.NewMockStorageClientInterface(mockCtrl)
	mockStorageClient.EXPECT().FindGcsFile(ovfPath, ".ovf").Return(ovfObjectHandle, nil).Times(1)
	mockStorageClient.EXPECT().GetGcsFileContent(ovfObjectHandle).Return(nil, err).Times(1)
	mockOvfDescriptorValidator := mocks.NewMockOvfDescriptorValidatorInterface(mockCtrl)

	l := OvfDescriptorLoader{storageClient: mockStorageClient, validator: mockOvfDescriptorValidator}
	result, resultError := l.Load(ovfPath)

	assert.Equal(t, err, resultError)
	assert.Nil(t, result)
}

func TestOvfDescriptorLoaderErrorValidatingDescriptor(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	err := fmt.Errorf("error validating descriptor")

	mockStorageClient := mocks.NewMockStorageClientInterface(mockCtrl)
	mockStorageClient.EXPECT().FindGcsFile(ovfPath, ".ovf").Return(ovfObjectHandle, nil).Times(1)
	mockStorageClient.EXPECT().GetGcsFileContent(ovfObjectHandle).Return([]byte(ovfDescriptorStr), nil).Times(1)

	mockOvfDescriptorValidator := mocks.NewMockOvfDescriptorValidatorInterface(mockCtrl)
	mockOvfDescriptorValidator.EXPECT().ValidateOvfPackage(ovfDescriptor, ovfPath).Return(nil, err).Times(1)

	l := OvfDescriptorLoader{storageClient: mockStorageClient, validator: mockOvfDescriptorValidator}
	result, resultError := l.Load(ovfPath)

	assert.Equal(t, err, resultError)
	assert.Nil(t, result)
}
