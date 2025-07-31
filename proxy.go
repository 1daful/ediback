// proxy/proxy.go
package main

import (
	bytes "bytes"
	json "encoding/json"
	fmt "fmt"
	io "io"
	log "log"
	"net/http"
	urlpkg "net/url"
	"strings"
	"sync"
	"time"
)

// Unified ProxyHandler for batch proxying
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

		var requests []map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		var wg sync.WaitGroup
		responses := make([]map[string]interface{}, len(requests))

		for i, reqData := range requests {
			wg.Add(1)
			go func(i int, reqData map[string]interface{}) {
				defer wg.Done()

				urlStr, _ := reqData["url"].(string)
				method, _ := reqData["method"].(string)
				if method == "" {
					method = http.MethodGet
				} else {
					method = strings.ToUpper(method)
				}

				parsedURL, err := urlpkg.Parse(urlStr)
				if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
					responses[i] = map[string]interface{}{"url": urlStr, "error": "Invalid URL"}
					return
				}

				if !cfg.AllowInsecure && parsedURL.Scheme != "https" {
					responses[i] = map[string]interface{}{"url": urlStr, "error": "Only HTTPS requests are allowed"}
					return
				}

				if params, ok := reqData["params"].(map[string]interface{}); ok && method == http.MethodGet {
					query := parsedURL.Query()
					for k, v := range params {
						query.Set(k, fmt.Sprintf("%v", v))
					}
					parsedURL.RawQuery = query.Encode()
					urlStr = parsedURL.String()
				}

				var reqBody io.Reader
				if method != http.MethodGet {
					if bodyData, ok := reqData["body"].(map[string]interface{}); ok {
						bodyBytes, _ := json.Marshal(bodyData)
						reqBody = bytes.NewBuffer(bodyBytes)
					}
				}

				req, err := http.NewRequest(method, urlStr, reqBody)
				if err != nil {
					responses[i] = map[string]interface{}{"url": urlStr, "error": "Request creation failed"}
					return
				}

				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Accept", "application/json")

				if headers, ok := reqData["headers"].(map[string]interface{}); ok {
					for k, v := range headers {
						req.Header.Set(k, fmt.Sprintf("%v", v))
					}
				}

				log.Printf("[Proxy] Forwarding: METHOD=%s, URL=%s", method, urlStr)
				resp, err := client.Do(req)
				if err != nil {
					responses[i] = map[string]interface{}{"url": urlStr, "error": err.Error()}
					return
				}
				defer resp.Body.Close()

				respBytes, _ := io.ReadAll(resp.Body)
				var respJSON interface{}
				json.Unmarshal(respBytes, &respJSON)

				responses[i] = map[string]interface{}{
					"url":     urlStr,
					"status":  resp.StatusCode,
					"headers": resp.Header,
					"body":    respJSON,
				}
			}(i, reqData)
		}

		wg.Wait()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responses)
	}
}

// ProxySingle enables internal single-call usage without duplicating logic
func ProxySingle(cfg Config, reqData map[string]interface{}) (map[string]interface{}, error) {
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     10,
		IdleConnTimeout:     30 * time.Second,
	}
	client := &http.Client{Transport: transport, Timeout: 30 * time.Second}

	urlStr, _ := reqData["url"].(string)
	method, _ := reqData["method"].(string)
	if method == "" {
		method = http.MethodGet
	} else {
		method = strings.ToUpper(method)
	}

	parsedURL, err := urlpkg.Parse(urlStr)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, fmt.Errorf("invalid URL")
	}

	if !cfg.AllowInsecure && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("only HTTPS requests are allowed")
	}

	if params, ok := reqData["params"].(map[string]interface{}); ok && method == http.MethodGet {
		query := parsedURL.Query()
		for k, v := range params {
			query.Set(k, fmt.Sprintf("%v", v))
		}
		parsedURL.RawQuery = query.Encode()
		urlStr = parsedURL.String()
	}

	var reqBody io.Reader
	if method != http.MethodGet {
		if bodyData, ok := reqData["body"].(map[string]interface{}); ok {
			bodyBytes, _ := json.Marshal(bodyData)
			reqBody = bytes.NewBuffer(bodyBytes)
		}
	}

	req, err := http.NewRequest(method, urlStr, reqBody)
	if err != nil {
		return nil, fmt.Errorf("request creation failed")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if headers, ok := reqData["headers"].(map[string]interface{}); ok {
		for k, v := range headers {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var respJSON interface{}
	json.Unmarshal(respBytes, &respJSON)

	result := map[string]interface{}{
		"url":     urlStr,
		"status":  resp.StatusCode,
		"headers": resp.Header,
		"body":    respJSON,
	}

	return result, nil
}
