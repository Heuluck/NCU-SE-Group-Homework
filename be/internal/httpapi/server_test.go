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

func newTestServer(t *testing.T) http.Handler {
	t.Helper()

	repo, err := repository.NewJSONTaskRepository(filepath.Join(t.TempDir(), "tasks.json"))
	if err != nil {
		t.Fatal(err)
	}
	return NewServer(service.NewTaskAppService(repo)).Routes()
}
