package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Location struct {
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	NextBusStopId string `json:"next_bus_stop_id"`
	IsStop        string `json:"is_stop"`
}

type BusPostion struct {
	BusId         string `json:"bus_id"`
	Latitude      string `json:"latitude"`
	Longitude     string `json:"longitude"`
	NextBusStopId string `json:"next_bus_stop_id"`
	IsBusStop     bool   `json:"is_bus_stop"`
}

func waitForHub(url string, retries int, delay time.Duration) {
    for i := 0; i < retries; i++ {
        resp, err := http.Get(url)
        if err == nil && resp.StatusCode == 200 {
            log.Println("Hub is ready")
            return
        }
        log.Println("Hub not ready, retrying...")
        time.Sleep(delay)
    }
    log.Fatal("Hub service not available")
}

func main() {
	data, err := ioutil.ReadFile("dataset.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var locations []Location

	if err := json.Unmarshal(data, &locations); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	waitForHub("http://hub:9090/hub/health", 20, 2*time.Second)

	// Dev environment: url := "http://localhost:9090/bus/position"
	url := "http://hub:9090/hub/bus/position"

	for _, loc := range locations {
		fmt.Printf("Latitude: %s, Longitude: %s\n", loc.Latitude, loc.Longitude)
		isBusStop, err := strconv.ParseBool(loc.IsStop)
		if err != nil {
			panic(err)
		}
		payload := BusPostion{
			BusId:         "492",
			Latitude:      loc.Latitude,
			Longitude:     loc.Longitude,
			NextBusStopId: loc.NextBusStopId,
			IsBusStop:     isBusStop,
		}
		jsonData, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println("Status:", resp.Status)
		fmt.Println("Response:", string(body))
		time.Sleep(1 * time.Second)
	}
}
