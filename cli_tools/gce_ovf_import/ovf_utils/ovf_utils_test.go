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

	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/gce_ovf_import/ovf_model"
	"github.com/GoogleCloudPlatform/compute-image-tools/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	diskCapacityAllocationUnits = "byte * 2^30"

	fileRef1     = "file1"
	fileRef2     = "file2"
	defaultDisks = &ovfmodel.DiskSection{Disks: []ovfmodel.VirtualDisk{
		{Capacity: "20", CapacityAllocationUnits: &diskCapacityAllocationUnits, DiskID: "vmdisk1", FileRef: &fileRef1},
		{Capacity: "1", CapacityAllocationUnits: &diskCapacityAllocationUnits, DiskID: "vmdisk2", FileRef: &fileRef2},
	}}

	defaultFiles = &[]ovfmodel.File{
		{Href: "Ubuntu_for_Horizon71_1_1.0-disk1.vmdk", ID: "file1", Size: 1151322112},
		{Href: "Ubuntu_for_Horizon71_1_1.0-disk2.vmdk", ID: "file2", Size: 68096},
	}

	defaultReferences = &ovfmodel.ReferencesSection{
		Files: *defaultFiles,
	}
)

func TestGetDiskFileInfosDisksOnSingleControllerOutOfOrder(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createControllerItem("3", sataController),
			createControllerItem("4", usbController),
			createControllerItem("5", parallelSCSIController),
			createDiskItem("7", "1", "disk1", "ovf:/disk/vmdisk2", "5"),
			createDiskItem("6", "0", "disk0", "ovf:/disk/vmdisk1", "5"),
		},
	}
	doTestGetDiskFileInfosSuccess(t, virtualHardware)
}

func TestGetDiskFileInfosAllocationUnitExtraSpace(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createControllerItem("3", sataController),
			createControllerItem("4", usbController),
			createControllerItem("5", parallelSCSIController),
			createDiskItem("7", "1", "disk1", "ovf:/disk/vmdisk2", "5"),
			createDiskItem("6", "0", "disk0", "ovf:/disk/vmdisk1", "5"),
		},
	}
	extraSpaceDiskCapacityAllocationUnits := "byte * 2^ 30   "
	disks := &ovfmodel.DiskSection{Disks: []ovfmodel.VirtualDisk{
		{Capacity: "11", CapacityAllocationUnits: &extraSpaceDiskCapacityAllocationUnits, DiskID: "vmdisk1", FileRef: &fileRef1},
		{Capacity: "12", CapacityAllocationUnits: &extraSpaceDiskCapacityAllocationUnits, DiskID: "vmdisk2", FileRef: &fileRef2},
	}}

	diskInfos, err := GetDiskInfos(virtualHardware, disks, defaultFiles)

	assert.Nil(t, err)
	assert.NotNil(t, diskInfos)
	assert.Equal(t, 2, len(diskInfos))
	assert.Equal(t, "Ubuntu_for_Horizon71_1_1.0-disk1.vmdk", diskInfos[0].FilePath)
	assert.Equal(t, "Ubuntu_for_Horizon71_1_1.0-disk2.vmdk", diskInfos[1].FilePath)
	assert.Equal(t, 11, diskInfos[0].SizeInGB)
	assert.Equal(t, 12, diskInfos[1].SizeInGB)
}

func TestGetDiskFileInfosDisksOnSeparateControllersOutOfOrder(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createControllerItem("3", sataController),
			createControllerItem("4", usbController),
			createControllerItem("5", parallelSCSIController),
			createDiskItem("7", "0", "disk1", "ovf:/disk/vmdisk2", "5"),
			createDiskItem("6", "0", "disk0", "ovf:/disk/vmdisk1", "3"),
		},
	}

	doTestGetDiskFileInfosSuccess(t, virtualHardware)
}

func TestGetDiskFileInfosInvalidDiskReferenceFormat(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createControllerItem("3", sataController),
			createControllerItem("4", usbController),
			createControllerItem("5", parallelSCSIController),
			createDiskItem("7", "0", "disk1", "ovf:/disk/vmdisk2", "5"),
			createDiskItem("6", "0", "disk0", "INVALID_DISK_HOST_RESOURCE", "3"),
		},
	}

	_, err := GetDiskInfos(virtualHardware, defaultDisks, defaultFiles)
	assert.NotNil(t, err)
}

func TestGetDiskFileInfosMissingDiskReference(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createControllerItem("3", sataController),
			createControllerItem("4", usbController),
			createControllerItem("5", parallelSCSIController),
			createDiskItem("7", "0", "disk1", "ovf:/disk/vmdisk_DOESNT_EXIST", "5"),
			createDiskItem("6", "0", "disk0", "ovf:/disk/vmdisk1", "3"),
		},
	}

	_, err := GetDiskInfos(virtualHardware, defaultDisks, defaultFiles)
	assert.NotNil(t, err)
}

func TestGetDiskFileInfosMissingFileReference(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createControllerItem("3", sataController),
			createControllerItem("4", usbController),
			createControllerItem("5", parallelSCSIController),
			createDiskItem("7", "0", "disk1", "ovf:/disk/vmdisk2", "5"),
			createDiskItem("6", "0", "disk0", "ovf:/disk/vmdisk1", "3"),
		},
	}

	_, err := GetDiskInfos(virtualHardware, defaultDisks, &[]ovfmodel.File{
		{Href: "Ubuntu_for_Horizon71_1_1.0-disk1.vmdk", ID: "file1", Size: 1151322112},
	})
	assert.NotNil(t, err)
}

func TestGetDiskFileInfosDiskWithoutParentController(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createControllerItem("3", sataController),
			createControllerItem("4", usbController),
			createControllerItem("5", parallelSCSIController),
			createDiskItem("7", "0", "disk1", "ovf:/disk/vmdisk2", "5"),
			createDiskItem("6", "0", "disk0", "ovf:/disk/vmdisk1", "3"),
			createDiskItem("8", "0", "disk2", "ovf:/disk/vmdisk3", "123"),
		},
	}

	doTestGetDiskFileInfosSuccess(t, virtualHardware)
}

func TestGetDiskFileInfosNoControllers(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createDiskItem("7", "0", "disk1", "ovf:/disk/vmdisk2", "5"),
			createDiskItem("6", "0", "disk0", "ovf:/disk/vmdisk1", "3"),
			createDiskItem("8", "0", "disk2", "ovf:/disk/vmdisk3", "123"),
		},
	}
	_, err := GetDiskInfos(virtualHardware, defaultDisks, defaultFiles)
	assert.NotNil(t, err)
}

func TestGetDiskFileInfosNilFileReferences(t *testing.T) {
	_, err := GetDiskInfos(&ovfmodel.VirtualHardwareSection{}, defaultDisks, nil)
	assert.NotNil(t, err)
}

func TestGetDiskFileInfosNilDiskSection(t *testing.T) {
	_, err := GetDiskInfos(&ovfmodel.VirtualHardwareSection{}, nil, defaultFiles)
	assert.NotNil(t, err)
}

func TestGetDiskFileInfosNilDisks(t *testing.T) {
	_, err := GetDiskInfos(&ovfmodel.VirtualHardwareSection{}, &ovfmodel.DiskSection{}, defaultFiles)
	assert.NotNil(t, err)
}

func TestGetDiskFileInfosEmptyDisks(t *testing.T) {
	_, err := GetDiskInfos(&ovfmodel.VirtualHardwareSection{},
		&ovfmodel.DiskSection{Disks: []ovfmodel.VirtualDisk{}}, defaultFiles)
	assert.NotNil(t, err)
}

func TestGetDiskFileInfosNilVirtualHardware(t *testing.T) {
	_, err := GetDiskInfos(nil, defaultDisks, defaultFiles)
	assert.NotNil(t, err)
}

func doTestGetDiskFileInfosSuccess(t *testing.T, virtualHardware *ovfmodel.VirtualHardwareSection) {
	diskInfos, err := GetDiskInfos(virtualHardware, defaultDisks, defaultFiles)

	assert.Nil(t, err)
	assert.NotNil(t, diskInfos)
	assert.Equal(t, 2, len(diskInfos))
	assert.Equal(t, "Ubuntu_for_Horizon71_1_1.0-disk1.vmdk", diskInfos[0].FilePath)
	assert.Equal(t, "Ubuntu_for_Horizon71_1_1.0-disk2.vmdk", diskInfos[1].FilePath)
	assert.Equal(t, 20, diskInfos[0].SizeInGB)
	assert.Equal(t, 1, diskInfos[1].SizeInGB)
}

func TestGetVirtualHardwareSection(t *testing.T) {
	expected := ovfmodel.VirtualHardwareSection{}
	virtualSystem := &ovfmodel.VirtualSystem{VirtualHardware: []ovfmodel.VirtualHardwareSection{expected}}

	virtualHardware, err := GetVirtualHardwareSection(virtualSystem)
	assert.Equal(t, &expected, virtualHardware)
	assert.Nil(t, err)
}

func TestGetVirtualHardwareSectionWhenVirtualSystemNil(t *testing.T) {
	_, err := GetVirtualHardwareSection(nil)
	assert.NotNil(t, err)
}

func TestGetVirtualHardwareSectionWhenVirtualHardwareNil(t *testing.T) {
	virtualSystem := &ovfmodel.VirtualSystem{VirtualHardware: nil}
	_, err := GetVirtualHardwareSection(virtualSystem)
	assert.NotNil(t, err)
}

func TestGetVirtualHardwareSectionWhenVirtualHardwareEmpty(t *testing.T) {
	virtualSystem := &ovfmodel.VirtualSystem{VirtualHardware: []ovfmodel.VirtualHardwareSection{}}
	_, err := GetVirtualHardwareSection(virtualSystem)
	assert.NotNil(t, err)
}

func TestGetVirtualSystem(t *testing.T) {
	expected := &ovfmodel.VirtualSystem{}
	ovfDescriptor := &ovfmodel.Descriptor{VirtualSystem: expected}
	virtualSystem, err := GetVirtualSystem(ovfDescriptor)

	assert.Equal(t, expected, virtualSystem)
	assert.Nil(t, err)
}

func TestGetVirtualSystemNilOvfDescriptor(t *testing.T) {
	_, err := GetVirtualSystem(nil)
	assert.NotNil(t, err)
}

func TestGetVirtualSystemNilVirtualSystem(t *testing.T) {
	ovfDescriptor := &ovfmodel.Descriptor{}
	_, err := GetVirtualSystem(ovfDescriptor)
	assert.NotNil(t, err)
}

func TestGetVirtualHardwareSectionFromDescriptor(t *testing.T) {
	expected := ovfmodel.VirtualHardwareSection{}
	virtualSystem := &ovfmodel.VirtualSystem{VirtualHardware: []ovfmodel.VirtualHardwareSection{expected}}
	ovfDescriptor := &ovfmodel.Descriptor{VirtualSystem: virtualSystem}

	virtualHardware, err := GetVirtualHardwareSectionFromDescriptor(ovfDescriptor)
	assert.Equal(t, &expected, virtualHardware)
	assert.Nil(t, err)
}

func TestGetVirtualHardwareSectionFromDescriptorWhenNilVirtualHardware(t *testing.T) {
	virtualSystem := &ovfmodel.VirtualSystem{VirtualHardware: nil}
	ovfDescriptor := &ovfmodel.Descriptor{VirtualSystem: virtualSystem}

	_, err := GetVirtualHardwareSectionFromDescriptor(ovfDescriptor)
	assert.NotNil(t, err)
}

func TestGetVirtualHardwareSectionFromDescriptorWhenNilVirtualSystem(t *testing.T) {
	ovfDescriptor := &ovfmodel.Descriptor{VirtualSystem: nil}

	_, err := GetVirtualHardwareSectionFromDescriptor(ovfDescriptor)
	assert.NotNil(t, err)
}

func TestGetCapacityInGB(t *testing.T) {
	//in GB
	doTestGetCapacityInGB(t, 20, "20", "byte * 2^30")
	doTestGetCapacityInGB(t, 10, "10", "byte * 2^30")
	doTestGetCapacityInGB(t, 1, "1", "byte * 2^30")
	doTestGetCapacityInGB(t, 1024, "1024", "byte * 2^30")
	doTestGetCapacityInGB(t, 5242880, "5242880", "byte * 2^30")

	//in MB
	doTestGetCapacityInGB(t, 1, "1", "byte * 2^20")
	doTestGetCapacityInGB(t, 1, "1024", "byte * 2^20")
	doTestGetCapacityInGB(t, 5*1024, "5242880", "byte * 2^20")

	//in TB
	doTestGetCapacityInGB(t, 1024, "1", "byte * 2^40")
	doTestGetCapacityInGB(t, 5242880*1024, "5242880", "byte * 2^40")
}

func TestGetNumberOfCPUs(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createCPUItem("1", 3),
		},
	}

	result, err := GetNumberOfCPUs(virtualHardware)
	assert.Nil(t, err)
	assert.Equal(t, int64(3), result)
}

func TestGetNumberOfCPUsPicksFirst(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createCPUItem("1", 11),
			createCPUItem("2", 2),
			createCPUItem("3", 4),
		},
	}

	result, err := GetNumberOfCPUs(virtualHardware)
	assert.Nil(t, err)
	assert.Equal(t, int64(11), result)
}

func TestGetNumberOfCPUsErrorWhenVirtualHardwareNil(t *testing.T) {
	_, err := GetNumberOfCPUs(nil)
	assert.NotNil(t, err)
}

func TestGetNumberOfCPUsErrorWhenNoCPUs(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createControllerItem("4", usbController),
			createControllerItem("5", parallelSCSIController),
			createDiskItem("7", "0", "disk1", "ovf:/disk/vmdisk2", "5"),
		},
	}

	_, err := GetNumberOfCPUs(virtualHardware)
	assert.NotNil(t, err)
}

func TestGetMemoryInMB(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createMemoryItem("1", 16),
		},
	}

	result, err := GetMemoryInMB(virtualHardware)
	assert.Nil(t, err)
	assert.Equal(t, int64(16), result)
}

func TestGetMemoryInMBReturnsFirstMemorySpec(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createMemoryItem("1", 33),
			createMemoryItem("1", 16),
			createMemoryItem("1", 1),
		},
	}

	result, err := GetMemoryInMB(virtualHardware)
	assert.Nil(t, err)
	assert.Equal(t, int64(33), result)
}

func TestGetMemoryInMBSpecInGB(t *testing.T) {
	virtualHardware := createVirtualHardwareSectionWithMemoryItem(7, "byte * 2^30")
	result, err := GetMemoryInMB(virtualHardware)
	assert.Nil(t, err)
	assert.Equal(t, int64(7*1024), result)
}

func TestGetMemoryInMBSpecInGBSpacesAroundPower(t *testing.T) {
	virtualHardware := createVirtualHardwareSectionWithMemoryItem(3, "byte * 2^ 30   ")
	result, err := GetMemoryInMB(virtualHardware)
	assert.Nil(t, err)
	assert.Equal(t, int64(3*1024), result)
}

func TestGetMemoryInMBSpecInTB(t *testing.T) {
	virtualHardware := createVirtualHardwareSectionWithMemoryItem(5, "byte * 2^40")
	result, err := GetMemoryInMB(virtualHardware)
	assert.Nil(t, err)
	assert.Equal(t, int64(5*1024*1024), result)
}

func TestGetMemoryInMBInvalidAllocationUnit(t *testing.T) {
	virtualHardware := createVirtualHardwareSectionWithMemoryItem(5, "NOT_VALID_ALLOCATION_UNIT")
	_, err := GetMemoryInMB(virtualHardware)
	assert.NotNil(t, err)
}

func TestGetMemoryInMBEmptyAllocationUnit(t *testing.T) {
	virtualHardware := createVirtualHardwareSectionWithMemoryItem(5, "")
	_, err := GetMemoryInMB(virtualHardware)
	assert.NotNil(t, err)
}

func TestGetMemoryInMBNilAllocationUnit(t *testing.T) {
	memoryItem := createMemoryItem("1", 33)
	memoryItem.AllocationUnits = nil
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			memoryItem,
		},
	}
	_, err := GetMemoryInMB(virtualHardware)
	assert.NotNil(t, err)
}

func TestGetMemoryInMBReturnsErrorWhenVirtualHardwareNil(t *testing.T) {
	_, err := GetMemoryInMB(nil)
	assert.NotNil(t, err)
}

func TestGetMemoryInMBErrorWhenNoMemory(t *testing.T) {
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createControllerItem("4", usbController),
			createControllerItem("5", parallelSCSIController),
			createDiskItem("7", "0", "disk1",
				"ovf:/disk/vmdisk2", "5"),
		},
	}

	_, err := GetMemoryInMB(virtualHardware)
	assert.NotNil(t, err)
}

func TestGetOVFDescriptorAndDiskPaths(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ovfPackagePath := "gs://abucket/apath/"

	virtualHardware := ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			createControllerItem("3", sataController),
			createControllerItem("5", parallelSCSIController),
			createDiskItem("7", "1", "disk1",
				"ovf:/disk/vmdisk2", "5"),
			createDiskItem("6", "0", "disk0",
				"ovf:/disk/vmdisk1", "5"),
		},
	}
	ovfDescriptor := &ovfmodel.Descriptor{
		Disk:       defaultDisks,
		References: defaultReferences,
		VirtualSystem: &ovfmodel.VirtualSystem{
			VirtualHardware: []ovfmodel.VirtualHardwareSection{virtualHardware},
		},
	}

	mockOvfDescriptorLoader := mocks.NewMockOvfDescriptorLoaderInterface(mockCtrl)
	mockOvfDescriptorLoader.EXPECT().Load(ovfPackagePath).Return(ovfDescriptor, nil)

	ovfDescriptorResult, diskPaths, err := GetOVFDescriptorAndDiskPaths(
		mockOvfDescriptorLoader, ovfPackagePath)
	assert.NotNil(t, ovfDescriptorResult)
	assert.NotNil(t, diskPaths)
	assert.Nil(t, err)

	assert.Equal(t, []DiskInfo{
		{"gs://abucket/apath/Ubuntu_for_Horizon71_1_1.0-disk1.vmdk", 20},
		{"gs://abucket/apath/Ubuntu_for_Horizon71_1_1.0-disk2.vmdk", 1},
	}, diskPaths)
	assert.Equal(t, ovfDescriptor, ovfDescriptorResult)
}

func TestGetOVFDescriptorAndDiskPathsErrorWhenLoadingDescriptor(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOvfDescriptorLoader := mocks.NewMockOvfDescriptorLoaderInterface(mockCtrl)
	mockOvfDescriptorLoader.EXPECT().Load(
		"gs://abucket/apath/").Return(nil, fmt.Errorf("error loading descriptor"))

	ovfDescriptorResult, diskPaths, err := GetOVFDescriptorAndDiskPaths(
		mockOvfDescriptorLoader, "gs://abucket/apath/")
	assert.Nil(t, ovfDescriptorResult)
	assert.Nil(t, diskPaths)
	assert.NotNil(t, err)
}

func TestGetOVFDescriptorAndDiskPathsErrorWhenNoVirtualSystem(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOvfDescriptorLoader := mocks.NewMockOvfDescriptorLoaderInterface(mockCtrl)
	mockOvfDescriptorLoader.EXPECT().Load("gs://abucket/apath/").Return(
		&ovfmodel.Descriptor{
			References: defaultReferences,
			Disk:       defaultDisks,
		}, nil)

	ovfDescriptorResult, diskPaths, err := GetOVFDescriptorAndDiskPaths(
		mockOvfDescriptorLoader, "gs://abucket/apath/")
	assert.Nil(t, ovfDescriptorResult)
	assert.Nil(t, diskPaths)
	assert.NotNil(t, err)
}

func TestGetOVFDescriptorAndDiskPathsErrorWhenNoVirtualHardware(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOvfDescriptorLoader := mocks.NewMockOvfDescriptorLoaderInterface(mockCtrl)
	mockOvfDescriptorLoader.EXPECT().Load("gs://abucket/apath/").Return(
		&ovfmodel.Descriptor{
			VirtualSystem: &ovfmodel.VirtualSystem{},
			References:    defaultReferences,
			Disk:          defaultDisks,
		}, nil)

	ovfDescriptorResult, diskPaths, err := GetOVFDescriptorAndDiskPaths(
		mockOvfDescriptorLoader, "gs://abucket/apath/")
	assert.Nil(t, ovfDescriptorResult)
	assert.Nil(t, diskPaths)
	assert.NotNil(t, err)
}

func TestGetOVFDescriptorAndDiskPathsErrorWhenNoDisks(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOvfDescriptorLoader := mocks.NewMockOvfDescriptorLoaderInterface(mockCtrl)
	mockOvfDescriptorLoader.EXPECT().Load("gs://abucket/apath/").Return(
		&ovfmodel.Descriptor{
			VirtualSystem: &ovfmodel.VirtualSystem{VirtualHardware: []ovfmodel.VirtualHardwareSection{
				{Item: []ovfmodel.ResourceAllocationSettingData{
					createControllerItem("3", sataController)},
				},
			}},
			References: defaultReferences,
		}, nil)

	ovfDescriptorResult, diskPaths, err := GetOVFDescriptorAndDiskPaths(
		mockOvfDescriptorLoader, "gs://abucket/apath/")
	assert.Nil(t, ovfDescriptorResult)
	assert.Nil(t, diskPaths)
	assert.NotNil(t, err)
}

func TestGetOVFDescriptorAndDiskPathsErrorWhenNoReferences(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOvfDescriptorLoader := mocks.NewMockOvfDescriptorLoaderInterface(mockCtrl)
	mockOvfDescriptorLoader.EXPECT().Load("gs://abucket/apath/").Return(
		&ovfmodel.Descriptor{
			VirtualSystem: &ovfmodel.VirtualSystem{VirtualHardware: []ovfmodel.VirtualHardwareSection{
				{Item: []ovfmodel.ResourceAllocationSettingData{createControllerItem("3", sataController)}},
			}},
			Disk: defaultDisks,
		}, nil)

	ovfDescriptorResult, diskPaths, err := GetOVFDescriptorAndDiskPaths(
		mockOvfDescriptorLoader, "gs://abucket/apath/")
	assert.Nil(t, ovfDescriptorResult)
	assert.Nil(t, diskPaths)
	assert.NotNil(t, err)
}

func createVirtualHardwareSectionWithMemoryItem(quantity uint, allocationUnit string) *ovfmodel.VirtualHardwareSection {
	memoryItem := createMemoryItem("1", quantity)
	memoryItem.AllocationUnits = &allocationUnit
	virtualHardware := &ovfmodel.VirtualHardwareSection{
		Item: []ovfmodel.ResourceAllocationSettingData{
			memoryItem,
		},
	}
	return virtualHardware
}

func doTestGetCapacityInGB(t *testing.T, expected int, capacity string, allocationUnits string) {
	capacityInGB, err := getCapacityInGB(capacity, allocationUnits)
	assert.Nil(t, err)
	assert.Equal(t, expected, capacityInGB)
}

func createControllerItem(instanceID string, resourceType uint16) ovfmodel.ResourceAllocationSettingData {
	return ovfmodel.ResourceAllocationSettingData{
		InstanceID:   instanceID,
		ResourceType: &resourceType,
	}
}

func createDiskItem(instanceID string, addressOnParent string,
	elementName string, hostResource string, parent string) ovfmodel.ResourceAllocationSettingData {
	diskType := disk
	return ovfmodel.ResourceAllocationSettingData{
		InstanceID:      instanceID,
		ResourceType:    &diskType,
		AddressOnParent: &addressOnParent,
		ElementName:     elementName,
		HostResource:    []string{hostResource},
		Parent:          &parent,
	}
}

func createCPUItem(instanceID string, quantity uint) ovfmodel.ResourceAllocationSettingData {
	resourceType := cpu
	mhz := "hertz * 10^6"
	return ovfmodel.ResourceAllocationSettingData{
		InstanceID:      instanceID,
		ResourceType:    &resourceType,
		VirtualQuantity: &quantity,
		AllocationUnits: &mhz,
	}
}

func createMemoryItem(instanceID string, quantity uint) ovfmodel.ResourceAllocationSettingData {
	resourceType := memory
	mb := "byte * 2^20"

	return ovfmodel.ResourceAllocationSettingData{
		InstanceID:      instanceID,
		ResourceType:    &resourceType,
		VirtualQuantity: &quantity,
		AllocationUnits: &mb,
	}
}
