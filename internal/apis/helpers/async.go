package helpers

import (
	"encoding/json"
	"fmt"
	"time"

	"resty.dev/v3"
)

// EndpointAsyncTask is the polling endpoint exposed by the Infomaniak API for
// resolving operations whose initial response carried `result: "delayed"`.
const EndpointAsyncTask = "/1/async/tasks/{task_uuid}"

// AsyncTask is the payload of GET /1/async/tasks/{uuid}. When status reaches
// "executed", Response carries the final result of the original POST/PATCH/…
type AsyncTask struct {
	UUID       string          `json:"uuid"`
	Status     string          `json:"status"`
	CreatedAt  int64           `json:"created_at"`
	UpdatedAt  int64           `json:"updated_at"`
	ExecutedAt int64           `json:"executed_at"`
	FinishedAt int64           `json:"finished_at"`
	Response   AsyncTaskResult `json:"response"`
}

type AsyncTaskResult struct {
	Result string          `json:"result"`
	Data   json.RawMessage `json:"data"`
	Error  *ApiError       `json:"error"`
}

// ResolveAsync handles the three Infomaniak result envelopes:
//
//   - `result: "success"` (or "asynchronous" without polling) -> returns data as-is.
//   - `result: "delayed"` (or "asynchronous" with a task uuid) -> polls
//     /1/async/tasks/{uuid} until executed, then returns the inner data.
//   - `result: "error"` -> returns the embedded ApiError.
//
// The caller is responsible for json.Unmarshaling the returned raw bytes into
// the expected concrete type (typically int64 for create endpoints, bool for
// update/delete, or a domain struct).
func ResolveAsync(client *resty.Client, resp *resty.Response, raw NormalizedApiResponse[json.RawMessage], timeout time.Duration) (json.RawMessage, error) {
	if resp.IsError() {
		if raw.Error != nil {
			return nil, raw.Error
		}
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode())
	}

	switch raw.Result {
	case "success":
		return raw.Data, nil
	case "delayed", "asynchronous":
		var taskUUID string
		if err := json.Unmarshal(raw.Data, &taskUUID); err != nil {
			return nil, fmt.Errorf("expected task uuid in %q response, got: %s", raw.Result, raw.Data)
		}
		return PollAsyncTask(client, taskUUID, timeout)
	case "error":
		if raw.Error != nil {
			return nil, raw.Error
		}
		return nil, fmt.Errorf("API returned result=error without details")
	default:
		return nil, fmt.Errorf("unexpected API result envelope: %q", raw.Result)
	}
}

// PollAsyncTask blocks until the task is executed (and returns its inner data),
// fails, or the timeout fires. The poll interval ramps from 1s up to 5s.
func PollAsyncTask(client *resty.Client, taskUUID string, timeout time.Duration) (json.RawMessage, error) {
	deadline := time.Now().Add(timeout)
	interval := 1 * time.Second

	for {
		var env NormalizedApiResponse[*AsyncTask]
		resp, err := client.R().
			SetPathParam("task_uuid", taskUUID).
			SetResult(&env).
			SetError(&env).
			Get(EndpointAsyncTask)
		if err != nil {
			return nil, fmt.Errorf("poll task %s: %w", taskUUID, err)
		}
		if resp.IsError() {
			if env.Error != nil {
				return nil, env.Error
			}
			return nil, fmt.Errorf("poll task %s: HTTP %d", taskUUID, resp.StatusCode())
		}
		if env.Data == nil {
			return nil, fmt.Errorf("poll task %s: empty data", taskUUID)
		}

		switch env.Data.Status {
		case "executed":
			inner := env.Data.Response
			if inner.Result == "error" || inner.Error != nil {
				if inner.Error != nil {
					return nil, inner.Error
				}
				return nil, fmt.Errorf("task %s finished with error", taskUUID)
			}
			return inner.Data, nil
		case "failed", "error":
			return nil, fmt.Errorf("task %s failed", taskUUID)
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for task %s (last status=%q)", taskUUID, env.Data.Status)
		}
		time.Sleep(interval)
		if interval < 5*time.Second {
			interval += 1 * time.Second
		}
	}
}
