package client_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"open-ve-go-sdk/pkg/client"
	"testing"
)

func TestClient_Check(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/check" {
			t.Errorf("Unexpected URL path: %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got: %s", r.Method)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		defer r.Body.Close()

		var checkRequest client.CheckRequest
		if err := json.Unmarshal(body, &checkRequest); err != nil {
			t.Fatalf("Failed to unmarshal request body: %v", err)
		}

		if len(checkRequest.Validations) == 0 {
			t.Errorf("Expected validations in the request, got none")
		}

		var results []client.CheckResponseResult
		for _, v := range checkRequest.Validations {
			results = append(results, client.CheckResponseResult{
				ID:      v.ID,
				IsValid: true,
				Message: "",
			})
		}

		response := client.CheckResponse{
			Results: results,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	config := client.NewConfig(mockServer.URL, client.VersionOption("v1"))
	testClient := client.NewClient(config)

	checkRequest := &client.CheckRequest{
		Validations: []client.CheckRequestValidation{
			{
				ID: "test-id",
				Variables: map[string]interface{}{
					"name": "test-variable",
				},
			},
		},
	}

	checkResponse, err := testClient.Check(checkRequest)
	if err != nil {
		t.Fatalf("Check method failed: %v", err)
	}

	if len(checkResponse.Results) != 1 {
		t.Errorf("Expected 1 result, got: %d", len(checkResponse.Results))
	}

	result := checkResponse.Results[0]
	if result.ID != "test-id" {
		t.Errorf("Expected result ID 'test-id', got: %s", result.ID)
	}
	if !result.IsValid {
		t.Errorf("Expected result IsValid to be true, got: %v", result.IsValid)
	}
	if result.Message != "" {
		t.Errorf("Expected empty result Message, got: %s", result.Message)
	}
}
