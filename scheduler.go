// proxy-server/scheduler.go
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// StartSchedule runs a task once after delaySeconds
func StartSchedule(every int, unit string, delaySeconds int, exec func(string, string, io.Reader) ([]byte, interface{}, error), method, url string, body io.Reader) {
	go func() {
		duration := time.Duration(delaySeconds) * time.Second
		if strings.ToLower(unit) == "minute" {
			duration = time.Duration(delaySeconds) * time.Minute
		}
		time.Sleep(duration)

		log.Println("Scheduled task starting:", method, url)
		resp, _, err := exec(method, url, body)
		if err != nil {
			log.Println("Scheduled task error:", err)
			return
		}
		log.Println("Scheduled task completed:", string(resp))
	}()
}

// Example periodic task
func StartRecurringTask(interval time.Duration, method, url string, body *bytes.Buffer) {
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			resp, _, err := toRun(method, url, body)
			if err != nil {
				log.Println("Recurring task error:", err)
				continue
			}
			log.Println("Recurring task result:", string(resp))
		}
	}()
}

// toRun executes an HTTP request and returns raw and parsed JSON responses.
func toRun(method, rawURL string, body io.Reader) ([]byte, interface{}, error) {
    client := &http.Client{Timeout: 30 * time.Second}

    req, err := http.NewRequest(method, rawURL, body)
    if err != nil {
        log.Println("[toRun] Error creating request:", err)
        return nil, nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")

    resp, err := client.Do(req)
    if err != nil {
        log.Println("[toRun] Error making request:", err)
        return nil, nil, err
    }
    defer resp.Body.Close()

    contentType := resp.Header.Get("Content-Type")
    respBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("[toRun] Error reading response body:", err)
        return nil, nil, err
    }

    var result interface{}
    if strings.Contains(contentType, "application/json") {
        if err := json.Unmarshal(respBytes, &result); err != nil {
            log.Println("[toRun] Error decoding JSON response:", err)
        }
    } else {
        result = string(respBytes)
    }

    return respBytes, result, nil
}

