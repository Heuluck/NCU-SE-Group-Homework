//go:build integration

package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"gophertodo/backend/internal/domain"
	"gophertodo/backend/internal/repository"
	"gophertodo/backend/internal/service"
)

// startIntegrationServer starts a real HTTP server on a random port
// and returns the base URL and a cleanup function.
func startIntegrationServer(t *testing.T) (baseURL string, cleanup func()) {
	t.Helper()

	repo, err := repository.NewJSONTaskRepository(filepath.Join(t.TempDir(), "tasks.json"))
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	svc := service.NewTaskAppService(repo)
	handler := NewServer(svc).Routes()

	// httptest.NewServer starts a real HTTP server on a random port
	ts := httptest.NewServer(handler)
	return ts.URL, ts.Close
}

func TestIntegration_TaskLifecycle(t *testing.T) {
	baseURL, cleanup := startIntegrationServer(t)
	defer cleanup()
	client := &http.Client{}

	// === 1. Create a task ===
	createResp, err := client.Post(baseURL+"/tasks", "application/json",
		bytes.NewBufferString(`{"content":"integration test task"}`))
	if err != nil {
		t.Fatalf("POST /tasks failed: %v", err)
	}
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("POST /tasks status = %d, want 201", createResp.StatusCode)
	}

	var created domain.Task
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatalf("decode response failed: %v", err)
	}
	createResp.Body.Close()

	if created.Content != "integration test task" || created.Status != domain.StatusPending {
		t.Fatalf("unexpected task: %+v", created)
	}
	if created.ID <= 0 {
		t.Fatalf("expected positive ID, got %d", created.ID)
	}

	// === 2. List tasks ===
	listResp, err := client.Get(baseURL + "/tasks")
	if err != nil {
		t.Fatalf("GET /tasks failed: %v", err)
	}
	if listResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /tasks status = %d, want 200", listResp.StatusCode)
	}

	var tasks []domain.Task
	if err := json.NewDecoder(listResp.Body).Decode(&tasks); err != nil {
		t.Fatalf("decode response failed: %v", err)
	}
	listResp.Body.Close()

	if len(tasks) < 1 {
		t.Fatal("expected at least 1 task")
	}

	// === 3. Get task by ID ===
	getResp, err := client.Get(fmt.Sprintf("%s/tasks/%d", baseURL, created.ID))
	if err != nil {
		t.Fatalf("GET /tasks/%d failed: %v", created.ID, err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /tasks/%d status = %d, want 200", created.ID, getResp.StatusCode)
	}

	var fetched domain.Task
	if err := json.NewDecoder(getResp.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode response failed: %v", err)
	}
	getResp.Body.Close()

	if fetched.ID != created.ID {
		t.Fatalf("task ID mismatch: %d vs %d", fetched.ID, created.ID)
	}

	// === 4. Complete task ===
	completeResp, err := client.Post(fmt.Sprintf("%s/tasks/%d/complete", baseURL, created.ID),
		"application/json", nil)
	if err != nil {
		t.Fatalf("POST /tasks/%d/complete failed: %v", created.ID, err)
	}
	if completeResp.StatusCode != http.StatusOK {
		t.Fatalf("POST /tasks/%d/complete status = %d, want 200", created.ID, completeResp.StatusCode)
	}

	var completed domain.Task
	if err := json.NewDecoder(completeResp.Body).Decode(&completed); err != nil {
		t.Fatalf("decode response failed: %v", err)
	}
	completeResp.Body.Close()

	if completed.Status != domain.StatusCompleted {
		t.Fatalf("expected status 'completed', got '%s'", completed.Status)
	}
	if completed.CompletedAt == nil {
		t.Fatal("expected CompletedAt to be set")
	}

	// === 5. Delete task ===
	deleteReq, err := http.NewRequest(http.MethodDelete,
		fmt.Sprintf("%s/tasks/%d", baseURL, created.ID), nil)
	if err != nil {
		t.Fatalf("create DELETE request failed: %v", err)
	}
	deleteResp, err := client.Do(deleteReq)
	if err != nil {
		t.Fatalf("DELETE /tasks/%d failed: %v", created.ID, err)
	}
	if deleteResp.StatusCode != http.StatusNoContent {
		t.Fatalf("DELETE /tasks/%d status = %d, want 204", created.ID, deleteResp.StatusCode)
	}
	deleteResp.Body.Close()

	// === 6. Verify deletion ===
	getDeletedResp, err := client.Get(fmt.Sprintf("%s/tasks/%d", baseURL, created.ID))
	if err != nil {
		t.Fatalf("GET /tasks/%d after delete failed: %v", created.ID, err)
	}
	if getDeletedResp.StatusCode != http.StatusNotFound {
		t.Fatalf("GET /tasks/%d after delete status = %d, want 404", created.ID, getDeletedResp.StatusCode)
	}
	io.Copy(io.Discard, getDeletedResp.Body)
	getDeletedResp.Body.Close()
}

func TestIntegration_CORSHeaders(t *testing.T) {
	baseURL, cleanup := startIntegrationServer(t)
	defer cleanup()
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodOptions, baseURL+"/tasks", nil)
	if err != nil {
		t.Fatalf("create OPTIONS request failed: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("OPTIONS /tasks failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("OPTIONS status = %d, want 204", resp.StatusCode)
	}

	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("missing Access-Control-Allow-Origin header")
	}
	if resp.Header.Get("Access-Control-Allow-Methods") == "" {
		t.Fatalf("missing Access-Control-Allow-Methods header")
	}
}

func TestIntegration_Healthz(t *testing.T) {
	baseURL, cleanup := startIntegrationServer(t)
	defer cleanup()
	client := &http.Client{}

	resp, err := client.Get(baseURL + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /healthz status = %d, want 200", resp.StatusCode)
	}
}
