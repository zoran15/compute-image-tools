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

package ovfmodel

// Descriptor is OVF descriptor
type Descriptor struct {
	References    *ReferencesSection `xml:"References"`
	Annotation    *AnnotationSection `xml:"AnnotationSection"`
	VirtualSystem *VirtualSystem     `xml:"VirtualSystem"`
	Disk          *DiskSection       `xml:"DiskSection"`
}

// Section is OVF section
type Section struct {
	Required *bool  `xml:"required,attr"`
	Info     string `xml:"Info"`
}

// Content is generic OVF content
type Content struct {
	ID   string `xml:"id,attr"`
	Info string `xml:"Info"`
	Name string `xml:"Name"`
}

// ReferencesSection is OVF references section
type ReferencesSection struct {
	Section
	Files []File `xml:"File"`
}

// VirtualSystem is OVF virtual system
type VirtualSystem struct {
	Content
	VirtualHardware []VirtualHardwareSection `xml:"VirtualHardwareSection"`
}

// VirtualHardwareSection is OVF virtual hardware section
type VirtualHardwareSection struct {
	Item []ResourceAllocationSettingData `xml:"Item"`
}

// File is OVF file reference
type File struct {
	ID   string `xml:"id,attr"`
	Href string `xml:"href,attr"`
	Size uint   `xml:"size,attr"`
}

// AnnotationSection is OVF annotation section
type AnnotationSection struct {
	Section
	Annotation string `xml:"Annotation"`
}

// DiskSection is OVF disk section
type DiskSection struct {
	Section
	Disks []VirtualDisk `xml:"Disk"`
}

// VirtualDisk is OVF virtual disk
type VirtualDisk struct {
	DiskID                  string  `xml:"diskId,attr"`
	FileRef                 *string `xml:"fileRef,attr"`
	Capacity                string  `xml:"capacity,attr"`
	CapacityAllocationUnits *string `xml:"capacityAllocationUnits,attr"`
}

// ResourceAllocationSettingData is CIM Resource Allocation Setting Data
type ResourceAllocationSettingData struct {
	ResourceType         *uint16
	Parent               *string  `xml:"Parent"`
	InstanceID           string   `xml:"InstanceID"`
	AddressOnParent      *string  `xml:"AddressOnParent"`
	HostResource         []string `xml:"HostResource"`
	VirtualQuantity      *uint    `xml:"VirtualQuantity"`
	VirtualQuantityUnits *string  `xml:"VirtualQuantityUnits"`
	AllocationUnits      *string  `xml:"AllocationUnits"`
	ElementName          string   `xml:"ElementName"`
}
