package jobmanager

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/venicegeo/pzsvc-sdk-go/s3helper"
)

// StatusType is a string describing the state of the job.
type StatusType int

// Enumerate valid StatusType values.
const (
	Submitted StatusType = iota
	Running
	Success
	Cancelled
	Error
	Fail
)

var statuses = [...]string{"submitted", "running", "success", "cancelled", "error", "fail"}

func (status StatusType) String() string {
	return statuses[status]
}

// JobInput defines the expected input JSON structure.
// We currently support S3 input (bucket/key), though provider-specific (e.g.,
// GRiD) may be legitimate.
type JobInput struct {
	Source      s3helper.S3Bucket `json:"source"`
	Function    *string           `json:"function"`
	Options     *json.RawMessage  `json:"options"`
	Destination s3helper.S3Bucket `json:"destination"`
}

// JobOutput defines the expected output JSON structure.
type JobOutput struct {
	Input      JobInput                    `json:"input"`
	StartedAt  time.Time                   `json:"started_at"`
	FinishedAt time.Time                   `json:"finished_at"`
	Code       int                         `json:"code"`
	Message    string                      `json:"message"`
	Response   map[string]*json.RawMessage `json:"response"`
}

// UpdateMsg defines the expected output JSON structure for updating the JobManager.
type UpdateMsg struct {
	Status string `json:"status"`
}

// Update handles PDAL status updates.
func Update(t StatusType, r *http.Request) {
	log.Println("Setting job status as \"", t.String(), "\"")
	// var res UpdateMsg
	// res.Status = t.String()
	// //	url := "http://192.168.99.100:8080/manager"
	// url := r.URL.Path + `/manager`
	//
	// jsonStr, err := json.Marshal(res)
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("Content-Type", "application/json")
	//
	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
}

/*
BadRequest handles bad requests.

All bad requests result in a failure in the eyes of the JobManager. The ResponseWriter echos some key aspects of the Request (e.g., input, start time) and appends StatusBadRequest (400) as well as a message to the JobOutput, which is returned as JSON.
*/
func BadRequest(w http.ResponseWriter, r *http.Request, res JobOutput, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	res.Code = http.StatusBadRequest
	res.Message = message
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Fatal(err)
	}
	Update(Fail, r)
}

/*
InternalError handles internal server errors.

All internal server errors result in an error in the eyes of the JobManager. The ResponseWriter echos some key aspects of the Request (e.g., input, start time) and appends StatusInternalServerError (500) as well as a message to the JobOutput, which is returned as JSON.
*/
func InternalError(w http.ResponseWriter, r *http.Request, res *JobOutput, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	res.Code = http.StatusInternalServerError
	res.Message = message
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Fatal(err)
	}
	Update(Error, r)
}

/*
Okay handles successful calls.

All successful calls result in sucess in the eyes of the JobManager. The ResponseWriter echos some key aspects of the Request (e.g., input, start time) and appends StatusOK (200) as well as a message to the JobOutput, which is returned as JSON.
*/
func Okay(w http.ResponseWriter, r *http.Request, res JobOutput, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	res.Code = http.StatusOK
	res.Message = message
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Fatal(err)
	}
	Update(Success, r)
}

// GetJobInput provides a common means of parsing the JobInput JSON.
func GetJobInput(w http.ResponseWriter, r *http.Request, res JobOutput) JobInput {
	var msg JobInput

	// There should always be a body, else how are we to know what to do? Throw
	// 400 if missing.
	if r.Body == nil {
		BadRequest(w, r, res, "No JSON")
	}

	// Throw 500 if we cannot read the body.
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		InternalError(w, r, &res, err.Error())
	}

	// Throw 400 if we cannot unmarshal the body as a valid JobInput.
	if err := json.Unmarshal(b, &msg); err != nil {
		BadRequest(w, r, res, err.Error())
	}

	return msg
}
