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

	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/common/domain"
	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/gce_ovf_import/ovf_model"
)

// OvfValidator is responsible for validating OVF packages
type OvfValidator struct {
	storageClient commondomain.StorageClientInterface
}

// NewOvfValidator creates a new OvfValidator
func NewOvfValidator(
	storageClient commondomain.StorageClientInterface) *OvfValidator {
	return &OvfValidator{storageClient: storageClient}
}

// ValidateOvfPackage validates OVF package. This includes checking that references to resources in GCS exist.
func (v *OvfValidator) ValidateOvfPackage(
	ovfDescriptor *ovfmodel.Descriptor, ovfGcsPath string) (*ovfmodel.Descriptor, error) {
	if ovfDescriptor == nil {
		return nil, fmt.Errorf("OVF descriptor cannot be nil")
	}

	if err := v.validateReferencesExistInGcs(ovfDescriptor.References.Files, ovfGcsPath); err != nil {
		return nil, err
	}

	return ovfDescriptor, nil
}

func (v *OvfValidator) validateReferencesExistInGcs(
	references []ovfmodel.File, ovfGcsPath string) error {
	if references == nil {
		return nil
	}

	for _, reference := range references {
		referenceGcsHandle, err := v.storageClient.FindGcsFile(ovfGcsPath, reference.Href)
		if referenceGcsHandle == nil || err != nil {
			return fmt.Errorf("OVF reference %v not found in OVF package in %v", reference.Href, ovfGcsPath)
		}
	}
	return nil
}
