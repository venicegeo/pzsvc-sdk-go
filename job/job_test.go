/*
Copyright 2015-2016, RadiantBlue Technologies, Inc.

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

package job_test

import (
	"encoding/json"
	"testing"

	"github.com/venicegeo/pzsvc-sdk-go/job"
)

func TestInputMsg(t *testing.T) {
	in := `
    {
      "source": {
        "bucket": "Foo",
        "key": "Bar"
      },
      "function": "Baz",
			"destination": {
				"bucket": "Out",
				"key": "File"
			}
    }`

	b := []byte(in)

	var msg job.InputMsg
	if err := json.Unmarshal(b, &msg); err != nil {
		t.Error("Error parsing InputMsg")
	}
	if msg.Source.Bucket != "Foo" {
		t.Error(msg.Source.Bucket, "!= `Foo`")
	}
	if msg.Source.Key != "Bar" {
		t.Error(msg.Source.Key, "!= `Bar`")
	}
	if *msg.Function != "Baz" {
		t.Error(msg.Function, "!= `Baz`")
	}
	if msg.Destination.Bucket != "Out" {
		t.Error(msg.Destination.Bucket, "!= `Out`")
	}
	if msg.Destination.Key != "File" {
		t.Error(msg.Destination.Key, "!= `File`")
	}
}

func TestOutputMsg(t *testing.T) {
	out := `
    {
      "input": {
        "source": {
          "bucket": "Foo",
          "key": "Bar"
        },
        "function": "Baz"
      },
      "started_at": "2015-12-11T01:31:26.784569058Z",
      "finished_at": "2015-12-11T01:31:26.784569058Z",
      "status": "submitted",
      "response": {"filename":"download_file.laz","pdal_version":"1.1.0 (git-version: 0c36aa)"}
    }`

	b := []byte(out)

	var msg job.OutputMsg
	if err := json.Unmarshal(b, &msg); err != nil {
		t.Error("Error parsing OutputMsg")
	}
	if msg.Input.Source.Bucket != "Foo" {
		t.Error(msg.Input.Source.Bucket, "!= `Foo`")
	}
	if msg.Input.Source.Key != "Bar" {
		t.Error(msg.Input.Source.Key, "!= `Bar`")
	}
	if *msg.Input.Function != "Baz" {
		t.Error(msg.Input.Function, "!= `Baz`")
	}
	// if msg.StartedAt != "2015-12-11T01:31:26.784569058Z" {
	// 	t.Error(msg.StartedAt, "!= `2015-12-11T01:31:26.784569058Z`")
	// }
	// if msg.FinishedAt != "2015-12-11T01:31:26.784569058Z" {
	// 	t.Error(msg.FinishedAt, "!= `2015-12-11T01:31:26.784569058Z`")
	// }
	// if msg.Status != "submitted" {
	// 	t.Error(msg.Status, "!= `submitted`")
	// }
	// if msg.Response != `{"filename":"download_file.laz","pdal_version":"1.1.0 (git-version: 0c36aa)"}` {
	// 	t.Error(msg.Response, `!= {"filename":"download_file.laz","pdal_version":"1.1.0 (git-version: 0c36aa)"}`)
	// }
}

func TestStatusTypes(t *testing.T) {
	if job.Submitted.String() != "submitted" {
		t.Error(job.Submitted.String(), "!= `submitted`")
	}
	if job.Running.String() != "running" {
		t.Error(job.Running.String(), "!= `running`")
	}
	if job.Success.String() != "success" {
		t.Error(job.Success.String(), "!= `success`")
	}
	if job.Cancelled.String() != "cancelled" {
		t.Error(job.Cancelled.String(), "!= `cancelled`")
	}
	if job.Error.String() != "error" {
		t.Error(job.Error.String(), "!= `error`")
	}
	if job.Fail.String() != "fail" {
		t.Error(job.Fail.String(), "!= `fail`")
	}
}
