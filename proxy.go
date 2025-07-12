// proxy.go (no allowed host restriction)
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Target struct {
	URL    string                 `json:"url"`
	Method string                 `json:"method"`
	Body   map[string]interface{} `json:"body,omitempty"`
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

		var targets []Target
		if err := json.NewDecoder(r.Body).Decode(&targets); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		var wg sync.WaitGroup
		responses := make([]map[string]interface{}, len(targets))

		for i, target := range targets {
			wg.Add(1)
			go func(i int, target Target) {
				defer wg.Done()

				parsedURL, err := url.Parse(target.URL)
				if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
					responses[i] = map[string]interface{}{"url": target.URL, "error": "Invalid URL"}
					return
				}

				var reqBody io.Reader
				if target.Body != nil {
					bodyBytes, _ := json.Marshal(target.Body)
					reqBody = bytes.NewBuffer(bodyBytes)
				}

				method := target.Method
				if method == "" {
					method = "GET"
				}

				req, err := http.NewRequest(method, target.URL, reqBody)
				if err != nil {
					responses[i] = map[string]interface{}{"url": target.URL, "error": "Request creation failed"}
					return
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Accept", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					responses[i] = map[string]interface{}{"url": target.URL, "error": err.Error()}
					return
				}
				defer resp.Body.Close()

				respBytes, _ := io.ReadAll(resp.Body)
				var respJSON interface{}
				json.Unmarshal(respBytes, &respJSON)

				responses[i] = map[string]interface{}{
					"url":     target.URL,
					"status":  resp.StatusCode,
					"headers": resp.Header,
					"body":    respJSON,
				}
			}(i, target)
		}

		wg.Wait()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responses)
	}
}
