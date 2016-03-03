// Copyright 2016, RadiantBlue Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sdk

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
)

// RequestDecorator decorates an http.Request
type RequestDecorator interface {
	Decorate(*http.Request) error
}

// RequestFactory creates http.Request decorated as needed
type RequestFactory struct {
	BaseURL    string
	decorators []RequestDecorator
}

var rf *RequestFactory

// GetRequestFactory returns the one RequestFactory
func GetRequestFactory() *RequestFactory {
	if rf == nil {
		rf = new(RequestFactory)
	}
	return rf
}

// AddDecorator adds a RequestDecorator to the RequestFactory
func (rf *RequestFactory) AddDecorator(rd RequestDecorator) {
	rf.decorators = append(rf.decorators, rd)
}

// NewRequest creates a new http.Request, decorating it as needed
func (rf *RequestFactory) NewRequest(method, relativeURL string) *http.Request {
	parsedRelativeURL, _ := url.Parse(relativeURL)
	request, _ := http.NewRequest(method, parsedRelativeURL.String(), nil)
	for inx := 0; inx < len(rf.decorators); inx++ {
		rf.decorators[inx].Decorate(request)
	}
	return request
}

// ConfigBasicAuthDecorator adds basic authentication to a request
// based on the contents of a config file
type ConfigBasicAuthDecorator struct {
	Project string
	auth    string
}

// Decorate is the decorator as per RequestDecorator
func (bad ConfigBasicAuthDecorator) Decorate(request *http.Request) error {
	var err error
	if bad.auth == "" {
		bytes, err := GetConfig(bad.Project)
		if err == nil {
			type config struct {
				Auth string `json:"auth"`
			}
			var unmarshal config
			json.Unmarshal(bytes, &unmarshal)
			bad.auth = unmarshal.Auth
		}
	}
	if err == nil {
		request.Header.Set("Authorization", "Basic "+bad.auth)
	}
	return err
}

// StaticBaseURLDecorator sets the base URL based on a static string
type StaticBaseURLDecorator struct {
	BaseURL string
}

// Decorate is the decorator as per RequestDecorator
func (sbud StaticBaseURLDecorator) Decorate(request *http.Request) error {

	baseURL, _ := url.Parse(sbud.BaseURL)
	resolvedURL := baseURL.ResolveReference(request.URL)
	parsedURL, err := url.Parse(resolvedURL.String())
	request.URL = parsedURL
	return err
}

// LogDecorator decorates the request by logging it
type LogDecorator struct {
}

// Decorate is the decorator as per RequestDecorator
func (ld LogDecorator) Decorate(request *http.Request) error {
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	request.Write(writer)
	writer.Flush()
	log.Printf(buffer.String())
	return nil
}

// HTTPError represents any HTTP error
type HTTPError struct {
	Status  int
	Message string
}

func (err HTTPError) Error() string {
	return fmt.Sprintf("%d: %v", err.Status, err.Message)
}

// DoRequestCallback is a callback function for a generic request
type DoRequestCallback interface {
	Callback(*http.Response, error) error
}

var client *http.Client

// DoRequest performs the request and forwards the response
// to the callback provided
func DoRequest(request *http.Request, callback DoRequestCallback) error {
	if client == nil {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client = &http.Client{Transport: transport}
	}
	response, err := client.Do(request)
	return callback.Callback(response, err)
}

// DownloadCallback is an example callback object for DoRequest
// that performs a download operation
type DownloadCallback struct {
	FileName string
}

// Callback is the callback function for DownloadCallback
func (dc *DownloadCallback) Callback(response *http.Response, err error) error {
	if err != nil {
		return err
	}
	file, err := os.Create("temp")
	if err != nil {
		return err
	}
	defer file.Close()

	cd := response.Header.Get("Content-Disposition")
	_, params, err := mime.ParseMediaType(cd)
	if err != nil {
		// This generally means a broader error which is hopefully contained in the body
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return &HTTPError{Message: string(body), Status: http.StatusNotAcceptable}
	}
	dc.FileName = params["filename"]
	err = os.Rename(file.Name(), dc.FileName)
	return err
}
