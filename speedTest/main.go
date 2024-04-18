package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/showwin/speedtest-go/speedtest" // to be changed later
	"go.uber.org/zap"
)

var logger *zap.Logger

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

func formatSpeed(speed float64) string {
	return fmt.Sprintf("%.2f Mbps", speed)
}

func main() {

	// Initialize the logger
	initLogger()
	defer logger.Sync()
	http.HandleFunc("/speedtest", SpeedtestHandler)
	logger.Fatal("HTTP server error", zap.Error(http.ListenAndServe(":8080", nil)))
}
