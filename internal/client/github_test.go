package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetVariable(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Fatalf(
					"expected GET, got %s",
					r.Method,
				)
			}

			if r.URL.Path != "/repos/test-owner/test-repo/actions/variables/ENVIRONMENT" {
				t.Fatalf("unexpected path: %s", r.Method)
			}

			if r.Header.Get("Authorization") != "Bearer test-token" {
				t.Fatalf("missing authorization header")
			}

			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{
				"name": "ENVIRONMENT",
				"value": "production"
			}`))
		}),
	)

	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	variable, err := client.GetVariable(
		context.Background(),
		"test-owner",
		"test-repo",
		"ENVIRONMENT",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if variable.Name != "ENVIRONMENT" {
		t.Fatalf("expected ENVIRONMENT, got %s", variable.Name)
	}

	if variable.Value != "production" {
		t.Fatalf("expected production, got %s", variable.Value)
	}
}

func TestGetVariableNotFound(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)

			_, _ = w.Write([]byte(`{
			    "message": "Not Found"
			}`))
		}),
	)

	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	_, err := client.GetVariable(
		context.Background(),
		"test-owner",
		"test-repo",
		"DOES_NOT_EXIST",
	)

	if err == nil {
		t.Fatal("expected error got noting")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", apiErr.StatusCode)
	}
}

func TestCreateVariable(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Fatalf("expected POST, got %s", r.Method)
			}

			w.WriteHeader(http.StatusCreated)
		}),
	)

	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	err := client.CreateVariable(
		context.Background(),
		"test-owner",
		"test-repo",
		Variable{
			Name:  "ENVIRONMENT",
			Value: "production",
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateVariable(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPatch {
				t.Fatalf("expected PATCH, got %s", r.Method)
			}

			w.WriteHeader(http.StatusNoContent)
		}),
	)

	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	err := client.UpdateVariable(
		context.Background(),
		"test-owner",
		"test-repo",
		"ENVIRONMENT",
		"staging",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteVariable(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Fatalf("expected DELETE, got %s", r.Method)
			}

			w.WriteHeader(http.StatusNoContent)
		}),
	)

	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	err := client.DeleteVariable(
		context.Background(),
		"test-owner",
		"test-repo",
		"ENVIRONMENT",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
