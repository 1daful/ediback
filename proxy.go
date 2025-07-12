// proxy.go
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type ProxyRequestPayload struct {
	Targets []struct {
		URL    string                 `json:"url"`
		Method string                 `json:"method"`
		Body   map[string]interface{} `json:"body,omitempty"`
	} `json:"targets"`
}

func ProxyHandler(cfg Config) http.HandlerFunc {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
		IdleConnTimeout:     90 * time.Second,
	}
	client := &http.Client{Transport: transport, Timeout: 30 * time.Second}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload ProxyRequestPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		var wg sync.WaitGroup
		responses := make([]map[string]interface{}, len(payload.Targets))

		for i, target := range payload.Targets {
			wg.Add(1)
			go func(i int, targetURL, method string, body map[string]interface{}) {
				defer wg.Done()

				parsedURL, err := url.Parse(targetURL)
				if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
					responses[i] = map[string]interface{}{"url": targetURL, "error": "Invalid URL"}
					return
				}

				allowed := false
				for _, host := range cfg.AllowedHosts {
					if strings.Contains(parsedURL.Host, host) {
						allowed = true
						break
					}
				}
				if !allowed {
					responses[i] = map[string]interface{}{"url": targetURL, "error": "Host not allowed"}
					return
				}

				var reqBody io.Reader
				if body != nil {
					bodyBytes, _ := json.Marshal(body)
					reqBody = bytes.NewBuffer(bodyBytes)
				}

				req, err := http.NewRequest(method, targetURL, reqBody)
				if err != nil {
					responses[i] = map[string]interface{}{"url": targetURL, "error": "Request creation failed"}
					return
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Accept", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					responses[i] = map[string]interface{}{"url": targetURL, "error": err.Error()}
					return
				}
				defer resp.Body.Close()

				respBytes, _ := io.ReadAll(resp.Body)
				var respJSON interface{}
				json.Unmarshal(respBytes, &respJSON)

				responses[i] = map[string]interface{}{
					"url":     targetURL,
					"status":  resp.StatusCode,
					"headers": resp.Header,
					"body":    respJSON,
				}
			}(i, target.URL, target.Method, target.Body)
		}

		wg.Wait()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responses)
	}
}
