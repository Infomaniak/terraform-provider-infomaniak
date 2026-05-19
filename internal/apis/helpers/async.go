package helpers

import (
	"fmt"
	"time"

	"resty.dev/v3"
)

// EndpointAsyncTask is the polling endpoint exposed by the Infomaniak API for
// resolving operations whose initial response carried `result: "delayed"`.
const EndpointAsyncTask = "/1/async/tasks/{task_uuid}"

// AsyncResponse is the wire envelope returned by an endpoint known to be
// asynchronous: `data` is always a task UUID. Use it as the SetResult/SetError
// target on the initiating request, then pass it to ResolveAsync.
type AsyncResponse = NormalizedApiResponse[string]

// AsyncTask is the payload of GET /1/async/tasks/{uuid}. When status reaches
// "executed", Response carries the final result of the original POST/PATCH/…
// typed as R.
type AsyncTask[R any] struct {
	UUID       string                   `json:"uuid"`
	Status     string                   `json:"status"`
	CreatedAt  int64                    `json:"created_at"`
	UpdatedAt  int64                    `json:"updated_at"`
	ExecutedAt int64                    `json:"executed_at"`
	FinishedAt int64                    `json:"finished_at"`
	Response   NormalizedApiResponse[R] `json:"response"`
}

// ResolveAsync handles the result envelope of an endpoint known to be
// asynchronous and returns its final, typed result R.
//
//   - `result: "delayed"` / "asynchronous" -> polls /1/async/tasks/{uuid}
//     until executed and unmarshals the inner data into R.
//   - `result: "error"` -> returns the embedded ApiError.
//   - any other envelope (including "success") -> returns an error: an
//     async endpoint must not return a synchronous payload.
func ResolveAsync[R any](client *resty.Client, resp *resty.Response, raw AsyncResponse, timeout time.Duration) (R, error) {
	var zero R
	if resp.IsError() {
		if raw.Error != nil {
			return zero, raw.Error
		}
		return zero, fmt.Errorf("HTTP %d", resp.StatusCode())
	}

	switch raw.Result {
	case "delayed", "asynchronous":
		return PollAsyncTask[R](client, raw.Data, timeout)
	case "error":
		if raw.Error != nil {
			return zero, raw.Error
		}
		return zero, fmt.Errorf("API returned result=error without details")
	case "success":
		return zero, fmt.Errorf("expected task uuid in async response, got synchronous %q envelope", raw.Result)
	default:
		return zero, fmt.Errorf("unexpected API result envelope: %q", raw.Result)
	}
}

// pollInitialInterval is the first wait between polls; it ramps up to pollMaxInterval.
const (
	pollInitialInterval = 1 * time.Second
	pollMaxInterval     = 5 * time.Second
)

// PollAsyncTask blocks until the task is executed (and returns its typed inner
// data), fails, or the timeout fires. The poll interval ramps from 1s up to 5s.
func PollAsyncTask[R any](client *resty.Client, taskUUID string, timeout time.Duration) (R, error) {
	var zero R
	deadline := time.Now().Add(timeout)
	interval := pollInitialInterval

	for {
		task, err := fetchAsyncTask[R](client, taskUUID)
		if err != nil {
			return zero, err
		}
		if result, done, err := interpretAsyncTask(task, taskUUID); done {
			return result, err
		}
		if time.Now().After(deadline) {
			return zero, fmt.Errorf("timeout waiting for task %s (last status=%q)", taskUUID, task.Status)
		}
		time.Sleep(interval)
		if interval < pollMaxInterval {
			interval += 1 * time.Second
		}
	}
}

// fetchAsyncTask issues a single GET on the polling endpoint and returns the
// decoded task, normalising HTTP- and envelope-level errors into Go errors.
func fetchAsyncTask[R any](client *resty.Client, taskUUID string) (*AsyncTask[R], error) {
	var env NormalizedApiResponse[*AsyncTask[R]]
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
	return env.Data, nil
}

// interpretAsyncTask examines a polled task's status. The returned `done` flag
// signals to the caller that polling should stop; when true, (result, err)
// carries the final outcome. When false, polling should continue.
func interpretAsyncTask[R any](task *AsyncTask[R], taskUUID string) (R, bool, error) {
	var zero R
	switch task.Status {
	case "executed":
		inner := task.Response
		if inner.Error != nil {
			return zero, true, inner.Error
		}
		if inner.Result == "error" {
			return zero, true, fmt.Errorf("task %s finished with error", taskUUID)
		}
		return inner.Data, true, nil
	case "failed", "error":
		return zero, true, fmt.Errorf("task %s failed", taskUUID)
	default:
		return zero, false, nil
	}
}
