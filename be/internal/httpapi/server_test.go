package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"gophertodo/backend/internal/domain"
	"gophertodo/backend/internal/repository"
	"gophertodo/backend/internal/service"
)

func TestTaskAPIFlow(t *testing.T) {
	handler := newTestServer(t)

	createReq := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{"content":"write backend"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	handler.ServeHTTP(createRes, createReq)
	if createRes.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body = %s", createRes.Code, createRes.Body.String())
	}

	var created domain.Task
	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	if created.ID != 1 || created.Status != domain.StatusPending || created.Content != "write backend" {
		t.Fatalf("unexpected created task: %+v", created)
	}

	completeReq := httptest.NewRequest(http.MethodPost, "/tasks/1/complete", nil)
	completeRes := httptest.NewRecorder()
	handler.ServeHTTP(completeRes, completeReq)
	if completeRes.Code != http.StatusOK {
		t.Fatalf("complete status = %d, body = %s", completeRes.Code, completeRes.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	listRes := httptest.NewRecorder()
	handler.ServeHTTP(listRes, listReq)
	if listRes.Code != http.StatusOK {
		t.Fatalf("list status = %d, body = %s", listRes.Code, listRes.Body.String())
	}

	var tasks []domain.Task
	if err := json.NewDecoder(listRes.Body).Decode(&tasks); err != nil {
		t.Fatal(err)
	}
	if len(tasks) != 1 || tasks[0].Status != domain.StatusCompleted || tasks[0].CompletedAt == nil {
		t.Fatalf("unexpected task list: %+v", tasks)
	}
}

func TestCreateTaskValidatesContent(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{"content":"   "}`))
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", res.Code, res.Body.String())
	}
}

func TestMissingTaskReturns404(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/999", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", res.Code, res.Body.String())
	}
}

func TestAlreadyCompletedTaskReturns409(t *testing.T) {
	handler := newTestServer(t)

	// Create and complete a task
	createReq := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{"content":"done"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	handler.ServeHTTP(createRes, createReq)
	if createRes.Code != http.StatusCreated {
		t.Fatal("failed to create task")
	}

	completeReq := httptest.NewRequest(http.MethodPost, "/tasks/1/complete", nil)
	completeRes := httptest.NewRecorder()
	handler.ServeHTTP(completeRes, completeReq)
	if completeRes.Code != http.StatusOK {
		t.Fatal("failed to complete task")
	}

	// Try to complete again — should return 409 Conflict
	againReq := httptest.NewRequest(http.MethodPost, "/tasks/1/complete", nil)
	againRes := httptest.NewRecorder()
	handler.ServeHTTP(againRes, againReq)
	if againRes.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409 Conflict, body = %s", againRes.Code, againRes.Body.String())
	}
}

func TestInvalidJSONBodyReturns400(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`not json`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 Bad Request, body = %s", res.Code, res.Body.String())
	}
}

func TestGetTaskByID(t *testing.T) {
	handler := newTestServer(t)

	// Create a task
	createReq := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{"content":"fetch me"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	handler.ServeHTTP(createRes, createReq)
	if createRes.Code != http.StatusCreated {
		t.Fatal("failed to create task")
	}

	// Fetch by ID
	getReq := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	getRes := httptest.NewRecorder()
	handler.ServeHTTP(getRes, getReq)
	if getRes.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 OK, body = %s", getRes.Code, getRes.Body.String())
	}

	var task domain.Task
	if err := json.NewDecoder(getRes.Body).Decode(&task); err != nil {
		t.Fatal(err)
	}
	if task.Content != "fetch me" || task.Status != domain.StatusPending {
		t.Fatalf("unexpected task: %+v", task)
	}
}

func TestHealthzEndpoint(t *testing.T) {
	handler := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 OK, body = %s", res.Code, res.Body.String())
	}
}

func TestCORSHeaders(t *testing.T) {
	handler := newTestServer(t)

	// Preflight request
	req := httptest.NewRequest(http.MethodOptions, "/tasks", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("OPTIONS status = %d, want 204 No Content", res.Code)
	}

	allowOrigin := res.Header().Get("Access-Control-Allow-Origin")
	allowMethods := res.Header().Get("Access-Control-Allow-Methods")
	if allowOrigin == "" || allowMethods == "" {
		t.Fatalf("missing CORS headers: Origin=%s, Methods=%s", allowOrigin, allowMethods)
	}
}

func TestMethodNotAllowedForTasksEndpoint(t *testing.T) {
	handler := newTestServer(t)

	// PUT should return 405
	req := httptest.NewRequest(http.MethodPut, "/tasks", nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405 Method Not Allowed", res.Code)
	}
}

func TestDeleteTask(t *testing.T) {
	handler := newTestServer(t)

	// Create a task
	createReq := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(`{"content":"to delete"}`))
	createReq.Header.Set("Content-Type", "application/json")
	createRes := httptest.NewRecorder()
	handler.ServeHTTP(createRes, createReq)
	if createRes.Code != http.StatusCreated {
		t.Fatal("failed to create task")
	}

	// Delete it
	deleteReq := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	deleteRes := httptest.NewRecorder()
	handler.ServeHTTP(deleteRes, deleteReq)
	if deleteRes.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204 No Content, body = %s", deleteRes.Code, deleteRes.Body.String())
	}

	// Verify it's gone
	getReq := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	getRes := httptest.NewRecorder()
	handler.ServeHTTP(getRes, getReq)
	if getRes.Code != http.StatusNotFound {
		t.Fatalf("get after delete status = %d, want 404 Not Found", getRes.Code)
	}
}

func newTestServer(t *testing.T) http.Handler {
	t.Helper()

	repo, err := repository.NewJSONTaskRepository(filepath.Join(t.TempDir(), "tasks.json"))
	if err != nil {
		t.Fatal(err)
	}
	return NewServer(service.NewTaskAppService(repo)).Routes()
}
