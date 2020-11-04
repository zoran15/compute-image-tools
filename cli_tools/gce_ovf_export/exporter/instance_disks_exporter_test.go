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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/compute-image-tools/cli_tools/mocks"
	"github.com/GoogleCloudPlatform/compute-image-tools/daisy"
	"github.com/golang/mock/gomock"
	"google.golang.org/api/option"
)

func Test(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockComputeClient := mocks.NewMockClient(mockCtrl)
	mockStorageClient, err := newTestGCSClient()
	if err != nil {
		t.Fail()
	}

	cb := func(w *daisy.Workflow) {
		w.ComputeClient = mockComputeClient
		w.StorageClient = mockStorageClient
	}

}

var (
	testGCSObjs   []string
	testGCSObjsMx = sync.Mutex{}
)

func newTestGCSClient() (*storage.Client, error) {
	nameRgx := regexp.MustCompile(`"name":"([^"].*)"`)
	rewriteRgx := regexp.MustCompile(`/b/([^/]+)/o/([^/]+)/rewriteTo/b/([^/]+)/o/([^?]+)`)
	uploadRgx := regexp.MustCompile(`/b/([^/]+)/o?.*uploadType=multipart.*`)
	getObjRgx := regexp.MustCompile(`/b/.+/o/.+alt=json&projection=full`)
	getBktRgx := regexp.MustCompile(`/b/.+alt=json&prettyPrint=false&projection=full`)
	deleteObjRgx := regexp.MustCompile(`/b/.+/o/.+alt=json`)
	listObjsRgx := regexp.MustCompile(`/b/.+/o\?alt=json&delimiter=&pageToken=&prefix=.+&projection=full&versions=false`)
	listObjsNoPrefixRgx := regexp.MustCompile(`/b/.+/o\?alt=json&delimiter=&pageToken=&prefix=&projection=full&versions=false`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		u := r.URL.String()
		m := r.Method

		if match := uploadRgx.FindStringSubmatch(u); m == "POST" && match != nil {
			body, _ := ioutil.ReadAll(r.Body)
			n := nameRgx.FindStringSubmatch(string(body))[1]
			addGCSObj(n)
			fmt.Fprintf(w, `{"kind":"storage#object","bucket":"%s","name":"%s"}`, match[1], n)
		} else if match := rewriteRgx.FindStringSubmatch(u); m == "POST" && match != nil {
			if strings.Contains(match[1], "dne") || strings.Contains(match[2], "dne") {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, storage.ErrObjectNotExist)
				return
			}
			path, err := url.PathUnescape(match[4])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, err)
				return
			}
			addGCSObj(path)
			o := fmt.Sprintf(`{"bucket":"%s","name":"%s"}`, match[3], match[4])
			fmt.Fprintf(w, `{"kind": "storage#rewriteResponse", "done": true, "objectSize": "1", "totalBytesRewritten": "1", "resource": %s}`, o)
		} else if match := getObjRgx.FindStringSubmatch(u); m == "GET" && match != nil {
			// Return StatusNotFound for objects that do not exist.
			if strings.Contains(match[0], "dne") {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			// Yes this object exists, we don't need to fill out the values, just return something.
			fmt.Fprint(w, "{}")
		} else if match := getBktRgx.FindStringSubmatch(u); m == "GET" && match != nil {
			// Return StatusNotFound for objects that do not exist.
			if strings.Contains(match[0], "dne") {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			// Yes this object exists, we don't need to fill out the values, just return something.
			fmt.Fprint(w, "{}")
		} else if match := deleteObjRgx.FindStringSubmatch(u); m == "DELETE" && match != nil {
			// Return StatusNotFound for objects that do not exist.
			if strings.Contains(match[0], "dne") {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			// Yes this object exists, we don't need to fill out the values, just return something.
			fmt.Fprint(w, "{}")
		} else if match := listObjsRgx.FindStringSubmatch(u); m == "GET" && match != nil {
			// Return StatusNotFound for objects that do not exist.
			if strings.Contains(match[0], "dne") {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			// Return 2 objects for testing recursiveGCS.
			fmt.Fprint(w, `{"kind": "storage#objects", "items": [{"kind": "storage#object", "name": "folder/object", "size": "1"},{"kind": "storage#object", "name": "folder/folder/object", "size": "1"}]}`)
		} else if match := listObjsNoPrefixRgx.FindStringSubmatch(u); m == "GET" && match != nil {
			// Return 2 objects for testing recursiveGCS.
			fmt.Fprint(w, `{"kind": "storage#objects", "items": [{"kind": "storage#object", "name": "object", "size": "1"},{"kind": "storage#object", "name": "folder/object", "size": "1"}]}`)
		} else if m == "PUT" && u == "/b/bucket/o/object/acl/allUsers?alt=json&prettyPrint=false" {
			fmt.Fprint(w, `{}`)
		} else if m == "PUT" && u == "/b/bucket/o/object%2Ffolder%2Ffolder%2Fobject/acl/allUsers?alt=json&prettyPrint=false" {
			fmt.Fprint(w, `{}`)
		} else if m == "PUT" && u == "/b/bucket/o/object%2Ffolder%2Fobject/acl/allUsers?alt=json&prettyPrint=false" {
			fmt.Fprint(w, `{}`)
		} else if m == "GET" && u == "/b?alt=json&pageToken=&prefix=&prettyPrint=false&project=foo-project&projection=full" {
			fmt.Fprint(w, `{}`)
		} else if m == "GET" && u == "/b?alt=json&pageToken=&prefix=&prettyPrint=false&project=bar-project&projection=full" {
			fmt.Fprint(w, `{"items": [{"name": "bar-project-daisy-bkt"}]}`)
		} else if m == "POST" && u == "/b?alt=json&prettyPrint=false&project=foo-project" {
			fmt.Fprint(w, `{}`)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "testGCSClient unknown request: %+v\n", r)
		}
	}))

	return storage.NewClient(context.Background(), option.WithEndpoint(ts.URL), option.WithHTTPClient(http.DefaultClient))
}

func addGCSObj(o string) {
	testGCSObjsMx.Lock()
	defer testGCSObjsMx.Unlock()
	testGCSObjs = append(testGCSObjs, o)
}