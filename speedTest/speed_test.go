package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSpeedtestHandler(t *testing.T) {
	// Create a mock HTTP request
	req, err := http.NewRequest("GET", "/speedtest", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function with the mock request and recorder
	SpeedtestHandler(rr, req)

	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the content type of the response
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}

	// Decode the response body into a slice of SpeedResult
	var results []SpeedResult
	if err := json.NewDecoder(rr.Body).Decode(&results); err != nil {
		t.Errorf("failed to decode JSON response: %v", err)
	}

	// Check if the response contains valid speed results
	if len(results) == 0 {
		t.Error("handler returned empty result set")
	}
}

func TestFormatSpeed(t *testing.T) {
	tests := []struct {
		speed    float64
		expected string
	}{
		{10.5, "10.50 Mbps"},
		{100.123456, "100.12 Mbps"},
		{0, "0.00 Mbps"},
		{999.999, "1000.00 Mbps"},
	}

	for _, test := range tests {
		result := formatSpeed(test.speed)
		if result != test.expected {
			t.Errorf("formatSpeed(%f) returned %s, expected %s", test.speed, result, test.expected)
		}
	}
}

func TestMain(t *testing.T) {
	// simple test case to check if the main function runs without error
	go main()
}
