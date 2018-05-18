/**
How this program works:
Given this payload

payload:
   recordings: ["1","2"]
   token: ""
   engineId: ""
   from_status: ""
   to_status: ""

for each recording id --> get the tasks --> that match the engine id
   look for `from_status`
   and change those tasks to `to_status`

Make sure token can be used to modify job!
*/
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/iron-io/iron_go3/worker"
	vt "github.com/veritone/go-veritone-api"
	"log"
)

type ConductorTaskPayload struct {
	EngineId           string   `json:"engineId"`
	Recordings         []string `json:"recordings"`
	VeritoneApiBaseUrl string   `json:"veritoneApiBaseUrl"`
	Token              string   `json:"token, omitempty"`
	FromStatus         string   `json:"fromStatus"`
	ToStatus           string   `json:"toStatus"`
}

func (p ConductorTaskPayload) String() string {
	b, _ := json.MarshalIndent(p, "", "  ")
	return string(b)
}

func toString(o interface{}) string {
	if o == nil {
		return "nil"
	}
	b, _ := json.MarshalIndent(o, "", "  ")
	return string(b)
}

// Error reporting and also set the status to failing
func updateTask(vtClient *vt.APIClient, jobId string, taskId string, status vt.TaskStatus, taskOutput interface{}) {
	fmt.Printf("updateTask jobId=%v, taskId=%v, status=%v, output=%v\n", jobId, taskId, status, toString(taskOutput))

	err := vtClient.UpdateTaskStatus(context.Background(), jobId, taskId, status, taskOutput)
	if err != nil {
		fmt.Printf("Failed to update taskStatus to %v: %s\n", status, err)
	}

	if status == vt.TaskStatusFailed {
		log.Printf("Updating fail...%v", taskOutput)
	}
}

type EngineExecutionHistory struct {
	JobId    string
	TaskId   string
	EngineId string
	Status   vt.TaskStatus
}

func (e EngineExecutionHistory) String() string {
	p, _ := json.MarshalIndent(e, "", "  ")
	return string(p)
}
func getExecutedEngineArraysOnARecording(vtClient *vt.APIClient, recordingId string) (engineHistory []EngineExecutionHistory) {
	fmt.Printf("Checking past tasks for recording %s ....\n", recordingId)
	engineHistory = make([]EngineExecutionHistory, 0, 10)
	recordingJobResponse, err := vtClient.GetJobsForRecording(context.Background(), recordingId)
	if err != nil {
		fmt.Printf("Just track here for now.. error getting jobs on recording=%v\n", err)
		return engineHistory
	}
	if recordingJobResponse.TotalRecords == 0 {
		return engineHistory
	}

	/**
	iterate over the tasks, collect the engine information incldue the task payload, language etc.
	*/
	for _, job := range recordingJobResponse.Records {
		for _, task := range job.Tasks {
			AnEngine := EngineExecutionHistory{
				JobId:    job.JobID,
				TaskId:   task.TaskID,
				EngineId: task.EngineID,
				Status:   task.Status,
			}
			engineHistory = append(engineHistory, AnEngine)

			// Trace the engine history here
			fmt.Println(AnEngine)
		}
	}
	return engineHistory
}
func toMsgMap(s string) map[string]string {
	return map[string]string{
		"msg": s,
	}
}
func main() {
	worker.ParseFlags()
	payload := &ConductorTaskPayload{}
	if err := worker.PayloadFromJSON(payload); err != nil {
		log.Fatalf("Failed to parse JSON-encoded task payload: %s", err.Error())
	}
	if payload.EngineId == "" || payload.Recordings == nil || payload.FromStatus == "" || payload.ToStatus == "" || payload.Token == "" {
		log.Fatalln("Missing required fields in task payload")
	}

	vtAPIConfig := vt.APIConfig{
		BaseURI: payload.VeritoneApiBaseUrl,
		Token:   payload.Token,
	}

	vtClient, err := vt.New(vtAPIConfig)
	if err != nil {
		log.Fatalf("Failed to create new Veritone API client: %s", err.Error())
	}

	for _, recordingId := range payload.Recordings {
		fmt.Println("--------------------------------------------------")
		recordingEngineHistory := getExecutedEngineArraysOnARecording(vtClient, recordingId)
		// now for each engine in the history, check against the from status and set to toStatus
		for _, e := range recordingEngineHistory {
			if e.EngineId == payload.EngineId && e.Status == vt.TaskStatus(payload.FromStatus) {
				fmt.Println("Recording Id=%s, Found Task for Engine! TaskId=%s, EngineId=%s, TaskStatus=%v\n",
					recordingId, e.TaskId, e.EngineId, e.Status)
				// update status accordingly
				output := toMsgMap(fmt.Sprintf("Auto setting status from %v to %v", e.Status, payload.ToStatus))
				updateTask(vtClient, e.JobId, e.TaskId, vt.TaskStatus(payload.ToStatus), output)
				if err != nil {
					fmt.Printf("Failed to update taskStatus from %v to %v: %s\n", e.Status, payload.ToStatus, err)
				}
			}
		}
	}
}
