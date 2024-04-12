package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/showwin/speedtest-go/speedtest"
)

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
	var lesgo []SpeedResult
	// Iterate over servers
	for _, s := range targets {
		// Ping test
		err := s.PingTest(nil)
		if err != nil {
			log.Println("Ping test failed for server:", s.Name)
			continue
		}

		// Download test
		err = s.DownloadTest()
		if err != nil {
			log.Println("Download test failed for server:", s.Name)
			continue
		}

		// Upload test
		err = s.UploadTest()
		if err != nil {
			log.Println("Upload test failed for server:", s.Name)
			continue
		}

		result := SpeedResult{
			DownloadSpeed: formatSpeed(s.DLSpeed),
			UploadSpeed:   formatSpeed(s.ULSpeed),
		}
		lesgo = append(lesgo, result)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)

		// Reset counter
		s.Context.Reset()
	}

	http.Error(w, "No servers available", http.StatusNotFound)
	return // Return after testing one server, you can remove this if you want to test multiple servers

}

func formatSpeed(speed float64) string {
	return fmt.Sprintf("%.2f Mbps", speed)
}

func main() {
	http.HandleFunc("/speedtest", SpeedtestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
