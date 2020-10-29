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

// Package mocks is a generated GoMock package.
package mocks

import (
	domain "github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/gce_ovf_export/domain"
	gomock "github.com/golang/mock/gomock"
	v1 "google.golang.org/api/compute/v1"
	reflect "reflect"
)

// MockInstanceExportPreparer is a mock of InstanceExportPreparer interface
type MockInstanceExportPreparer struct {
	ctrl     *gomock.Controller
	recorder *MockInstanceExportPreparerMockRecorder
}

// MockInstanceExportPreparerMockRecorder is the mock recorder for MockInstanceExportPreparer
type MockInstanceExportPreparerMockRecorder struct {
	mock *MockInstanceExportPreparer
}

// NewMockInstanceExportPreparer creates a new mock instance
func NewMockInstanceExportPreparer(ctrl *gomock.Controller) *MockInstanceExportPreparer {
	mock := &MockInstanceExportPreparer{ctrl: ctrl}
	mock.recorder = &MockInstanceExportPreparerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockInstanceExportPreparer) EXPECT() *MockInstanceExportPreparerMockRecorder {
	return m.recorder
}

// Cancel mocks base method
func (m *MockInstanceExportPreparer) Cancel(arg0 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cancel", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Cancel indicates an expected call of Cancel
func (mr *MockInstanceExportPreparerMockRecorder) Cancel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cancel", reflect.TypeOf((*MockInstanceExportPreparer)(nil).Cancel), arg0)
}

// Prepare mocks base method
func (m *MockInstanceExportPreparer) Prepare(arg0 *v1.Instance, arg1 *domain.OVFExportParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Prepare", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Prepare indicates an expected call of Prepare
func (mr *MockInstanceExportPreparerMockRecorder) Prepare(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Prepare", reflect.TypeOf((*MockInstanceExportPreparer)(nil).Prepare), arg0, arg1)
}

// TraceLogs mocks base method
func (m *MockInstanceExportPreparer) TraceLogs() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TraceLogs")
	ret0, _ := ret[0].([]string)
	return ret0
}

// TraceLogs indicates an expected call of TraceLogs
func (mr *MockInstanceExportPreparerMockRecorder) TraceLogs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TraceLogs", reflect.TypeOf((*MockInstanceExportPreparer)(nil).TraceLogs))
}
