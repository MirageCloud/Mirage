package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/showwin/speedtest-go/speedtest" // to be changed later
	"go.uber.org/zap"
)

var logger *zap.Logger

func initLogger() {
	// Create a logger configuration
	config := zap.NewProductionConfig()

	// Initialize the logger
	l, err := config.Build()
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	logger = l
}

type SpeedResult struct {
	DownloadSpeed string `json:"download_speed"`
	UploadSpeed   string `json:"upload_speed"`
}

func SpeedtestHandler(w http.ResponseWriter, r *http.Request) {
	var speedtestClient = speedtest.New()

	// Fetch servers
	serverList, err := speedtestClient.FetchServers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find closest servers
	targets, err := serverList.FindServer([]int{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resultsgo []SpeedResult

	// Iterate over servers
	for _, s := range targets {
		// Ping test
		err := s.PingTest(nil)
		if err != nil {
			logger.Warn("Ping test failed for server", zap.String("server_name", s.Name), zap.Error(err))
			continue
		}

		// Download test
		err = s.DownloadTest()
		if err != nil {
			logger.Warn("Download test failed for server", zap.String("server_name", s.Name), zap.Error(err))
			continue
		}

		// Upload test
		err = s.UploadTest()
		if err != nil {
			logger.Warn("Upload test failed for server", zap.String("server_name", s.Name), zap.Error(err))
			continue
		}

		result := SpeedResult{
			DownloadSpeed: formatSpeed(s.DLSpeed),
			UploadSpeed:   formatSpeed(s.ULSpeed),
		}
		resultsgo = append(resultsgo, result)

		// Reset counter
		s.Context.Reset()
	}

	// If no servers were available, return an error response
	if len(resultsgo) == 0 {
		http.Error(w, "No servers available", http.StatusNotFound)
		return
	}

	// Set the content type and encode the JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resultsgo); err != nil {
		logger.Error("Failed to encode JSON response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func formatSpeed(speed float64) string {
	return fmt.Sprintf("%.2f Mbps", speed)
}

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

func main() {

	// Initialize the logger
	initLogger()
	defer logger.Sync()
	http.HandleFunc("/speedtest", SpeedtestHandler)
	logger.Fatal("HTTP server error", zap.Error(http.ListenAndServe(":8080", nil)))
}
