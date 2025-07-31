// handlers.go
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)




/*func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<title>Ediback Admin Dashboard</title>
	<style>
		body { font-family: Arial, sans-serif; max-width: 800px; margin: 40px auto; padding: 20px; }
		.form-group { margin-bottom: 15px; }
		label { display: block; margin-bottom: 5px; font-weight: bold; }
		input, select, textarea { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 4px; }
		button { background: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 4px; cursor: pointer; }
		button:hover { background: #0056b3; }
		.results { margin-top: 20px; padding: 15px; background: #f8f9fa; border-radius: 4px; }
		.error { color: #dc3545; }
		.success { color: #28a745; }
		.header-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
	</style>
</head>
<body>
	<h2>Ediback Admin Dashboard</h2>
	<form id="taskForm" method="POST" action="/admin/run-task">
		<div class="form-group">
			<label>Endpoint URL:</label>
			<input type="text" name="url" required placeholder="https://api.example.com/endpoint">
		</div>
		
		<div class="form-group">
			<label>Method:</label>
			<select name="method" onchange="toggleBodyField()">
				<option value="GET">GET</option>
				<option value="POST">POST</option>
				<option value="PUT">PUT</option>
				<option value="DELETE">DELETE</option>
				<option value="PATCH">PATCH</option>
			</select>
		</div>
		
		<div class="form-group">
			<label>Headers (JSON format):</label>
			<textarea name="headers" rows="3" placeholder='{"Content-Type": "application/json", "Authorization": "Bearer token"}'></textarea>
		</div>
		
		<div class="form-group" id="bodyField" style="display: none;">
			<label>Request Body:</label>
			<textarea name="body" rows="5" placeholder='{"key": "value"}'></textarea>
		</div>
		
		<div class="form-group">
			<label>Timeout (seconds):</label>
			<input type="number" name="timeout" value="30" min="1" max="300">
		</div>
		
		<button type="submit">Execute Request</button>
	</form>
	
	<div class="results" id="results" style="display: none;">
		<h3>Response</h3>
		<div id="responseContent"></div>
	</div>
	
	<script>
		function toggleBodyField() {
			const method = document.querySelector('select[name="method"]').value;
			const bodyField = document.getElementById('bodyField');
			bodyField.style.display = ['POST', 'PUT', 'PATCH'].includes(method) ? 'block' : 'none';
		}
		
		document.getElementById('taskForm').addEventListener('submit', async function(e) {
			e.preventDefault();
			
			const formData = new FormData(this);
			const results = document.getElementById('results');
			const responseContent = document.getElementById('responseContent');
			
			try {
				const response = await fetch('/admin/run-task', {
					method: 'POST',
					body: formData
				});
				
				const result = await response.text();
				responseContent.innerHTML = '<pre>' + result + '</pre>';
				responseContent.className = response.ok ? 'success' : 'error';
				results.style.display = 'block';
			} catch (error) {
				responseContent.innerHTML = '<div class="error">Error: ' + error.message + '</div>';
				results.style.display = 'block';
			}
		});
	</script>
</body>
</html>`
	
	t, err := template.New("admin").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, "Execution error", http.StatusInternalServerError)
	}
}*/

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<title>Ediback Admin Dashboard</title>
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body {
			font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			min-height: 100vh;
			padding: 20px;
		}
		.container {
			max-width: 1200px;
			margin: 0 auto;
			background: rgba(255, 255, 255, 0.95);
			border-radius: 20px;
			padding: 30px;
			box-shadow: 0 25px 50px rgba(0, 0, 0, 0.15);
			backdrop-filter: blur(10px);
		}
		.header {
			text-align: center;
			margin-bottom: 40px;
			padding-bottom: 20px;
			border-bottom: 2px solid #e0e7ff;
		}
		.header h1 {
			color: #1e293b;
			font-size: 2.5rem;
			font-weight: 700;
			margin-bottom: 10px;
			background: linear-gradient(45deg, #667eea, #764ba2);
			-webkit-background-clip: text;
			-webkit-text-fill-color: transparent;
		}
		.header p {
			color: #64748b;
			font-size: 1.1rem;
		}
		.dashboard-grid {
			display: grid;
			grid-template-columns: 1fr 1fr;
			gap: 30px;
			margin-bottom: 30px;
		}
		.card {
			background: white;
			border-radius: 15px;
			padding: 25px;
			box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
			border: 1px solid #e2e8f0;
			transition: transform 0.3s ease, box-shadow 0.3s ease;
		}
		.card:hover {
			transform: translateY(-5px);
			box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
		}
		.card-title {
			color: #1e293b;
			font-size: 1.3rem;
			font-weight: 600;
			margin-bottom: 20px;
			display: flex;
			align-items: center;
			gap: 10px;
		}
		.icon {
			width: 20px;
			height: 20px;
			fill: #667eea;
		}
		.form-group {
			margin-bottom: 20px;
		}
		.form-group label {
			display: block;
			margin-bottom: 8px;
			color: #374151;
			font-weight: 500;
			font-size: 0.95rem;
		}
		.form-control {
			width: 100%;
			padding: 12px 16px;
			border: 2px solid #e2e8f0;
			border-radius: 10px;
			font-size: 1rem;
			transition: border-color 0.3s ease, box-shadow 0.3s ease;
			background: #f8fafc;
		}
		.form-control:focus {
			outline: none;
			border-color: #667eea;
			box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
			background: white;
		}
		.method-selector {
			display: grid;
			grid-template-columns: repeat(5, 1fr);
			gap: 10px;
			margin-top: 8px;
		}
		.method-btn {
			padding: 8px 12px;
			border: 2px solid #e2e8f0;
			border-radius: 8px;
			background: white;
			cursor: pointer;
			text-align: center;
			font-weight: 500;
			transition: all 0.3s ease;
			font-size: 0.9rem;
		}
		.method-btn:hover {
			border-color: #667eea;
			background: #f0f4ff;
		}
		.method-btn.active {
			background: #667eea;
			color: white;
			border-color: #667eea;
		}
		.execute-btn {
			background: linear-gradient(45deg, #667eea, #764ba2);
			color: white;
			border: none;
			padding: 15px 30px;
			border-radius: 10px;
			font-size: 1.1rem;
			font-weight: 600;
			cursor: pointer;
			transition: all 0.3s ease;
			width: 100%;
			margin-top: 20px;
			position: relative;
			overflow: hidden;
		}
		.execute-btn:hover {
			transform: translateY(-2px);
			box-shadow: 0 10px 25px rgba(102, 126, 234, 0.3);
		}
		.execute-btn:active {
			transform: translateY(0);
		}
		.execute-btn:disabled {
			opacity: 0.6;
			cursor: not-allowed;
		}
		.loading-spinner {
			display: none;
			width: 20px;
			height: 20px;
			border: 2px solid transparent;
			border-top: 2px solid white;
			border-radius: 50%;
			animation: spin 1s linear infinite;
			margin-right: 10px;
		}
		@keyframes spin {
			0% { transform: rotate(0deg); }
			100% { transform: rotate(360deg); }
		}
		.results-card {
			grid-column: 1 / -1;
			margin-top: 20px;
			display: none;
		}
		.response-tabs {
			display: flex;
			gap: 10px;
			margin-bottom: 20px;
		}
		.tab-btn {
			padding: 10px 20px;
			border: none;
			border-radius: 8px;
			background: #f1f5f9;
			cursor: pointer;
			font-weight: 500;
			transition: all 0.3s ease;
		}
		.tab-btn.active {
			background: #667eea;
			color: white;
		}
		.response-content {
			background: #1e293b;
			color: #e2e8f0;
			padding: 20px;
			border-radius: 10px;
			font-family: 'Courier New', monospace;
			font-size: 0.9rem;
			overflow-x: auto;
			white-space: pre-wrap;
			max-height: 400px;
			overflow-y: auto;
		}
		.status-indicator {
			display: inline-block;
			padding: 4px 12px;
			border-radius: 20px;
			font-size: 0.8rem;
			font-weight: 600;
			margin-bottom: 10px;
		}
		.status-success {
			background: #dcfce7;
			color: #166534;
		}
		.status-error {
			background: #fef2f2;
			color: #dc2626;
		}
		.stats-grid {
			display: grid;
			grid-template-columns: repeat(3, 1fr);
			gap: 15px;
			margin-bottom: 20px;
		}
		.stat-item {
			text-align: center;
			padding: 15px;
			background: #f8fafc;
			border-radius: 10px;
		}
		.stat-value {
			font-size: 1.5rem;
			font-weight: 700;
			color: #667eea;
		}
		.stat-label {
			font-size: 0.9rem;
			color: #64748b;
			margin-top: 5px;
		}
		.toggle-switch {
			position: relative;
			display: inline-block;
			width: 50px;
			height: 24px;
		}
		.toggle-switch input {
			opacity: 0;
			width: 0;
			height: 0;
		}
		.slider {
			position: absolute;
			cursor: pointer;
			top: 0;
			left: 0;
			right: 0;
			bottom: 0;
			background-color: #ccc;
			transition: 0.4s;
			border-radius: 24px;
		}
		.slider:before {
			position: absolute;
			content: "";
			height: 16px;
			width: 16px;
			left: 4px;
			bottom: 4px;
			background-color: white;
			transition: 0.4s;
			border-radius: 50%;
		}
		input:checked + .slider {
			background-color: #667eea;
		}
		input:checked + .slider:before {
			transform: translateX(26px);
		}
		@media (max-width: 768px) {
			.dashboard-grid {
				grid-template-columns: 1fr;
				gap: 20px;
			}
			.method-selector {
				grid-template-columns: repeat(3, 1fr);
			}
			.header h1 {
				font-size: 2rem;
			}
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Ediback Admin Dashboard</h1>
			<p>Advanced API Testing & Management Console</p>
		</div>
		
		<div class="dashboard-grid">
			<div class="card">
				<h3 class="card-title">
					<svg class="icon" viewBox="0 0 24 24">
						<path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
					</svg>
					Request Configuration
				</h3>
				
				<form id="taskForm">
					<div class="form-group">
						<label>Endpoint URL</label>
						<input type="text" name="url" class="form-control" required placeholder="https://api.example.com/endpoint">
					</div>
					
					<div class="form-group">
						<label>HTTP Method</label>
						<div class="method-selector">
							<div class="method-btn active" data-method="GET">GET</div>
							<div class="method-btn" data-method="POST">POST</div>
							<div class="method-btn" data-method="PUT">PUT</div>
							<div class="method-btn" data-method="DELETE">DELETE</div>
							<div class="method-btn" data-method="PATCH">PATCH</div>
						</div>
						<input type="hidden" name="method" value="GET">
					</div>
					
					<div class="form-group">
						<label>Request Headers</label>
						<textarea name="headers" class="form-control" rows="3" placeholder='{"Content-Type": "application/json", "Authorization": "Bearer token"}'></textarea>
					</div>
					
					<div class="form-group" id="bodyField" style="display: none;">
						<label>Request Body</label>
						<textarea name="body" class="form-control" rows="4" placeholder='{"key": "value"}'></textarea>
					</div>
					
					<div class="form-group">
						<label>Timeout (seconds)</label>
						<input type="number" name="timeout" class="form-control" value="30" min="1" max="300">
					</div>
					
					<button type="submit" class="execute-btn">
						<span class="loading-spinner"></span>
						<span class="btn-text">Execute Request</span>
					</button>
				</form>
			</div>
			
			<div class="card">
				<h3 class="card-title">
					<svg class="icon" viewBox="0 0 24 24">
						<path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zM9 17H7v-7h2v7zm4 0h-2V7h2v10zm4 0h-2v-4h2v4z"/>
					</svg>
					Request Statistics
				</h3>
				
				<div class="stats-grid">
					<div class="stat-item">
						<div class="stat-value" id="totalRequests">0</div>
						<div class="stat-label">Total Requests</div>
					</div>
					<div class="stat-item">
						<div class="stat-value" id="successCount">0</div>
						<div class="stat-label">Success</div>
					</div>
					<div class="stat-item">
						<div class="stat-value" id="errorCount">0</div>
						<div class="stat-label">Errors</div>
					</div>
				</div>
				
				<div class="form-group">
					<label style="display: flex; align-items: center; gap: 10px;">
						<span>Auto-refresh</span>
						<label class="toggle-switch">
							<input type="checkbox" id="autoRefresh">
							<span class="slider"></span>
						</label>
					</label>
				</div>
			</div>
			
			<div class="card results-card" id="resultsCard">
				<h3 class="card-title">
					<svg class="icon" viewBox="0 0 24 24">
						<path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
					</svg>
					Response Details
				</h3>
				
				<div id="responseStatus"></div>
				<div id="timingDisplay" style="margin-top:5px;color:#64748b;font-size:0.9rem;"></div>
				
				<div class="response-tabs">
					<button class="tab-btn active" data-tab="response">Response</button>
					<button class="tab-btn" data-tab="headers">Headers</button>
					<button class="tab-btn" data-tab="timing">Timing</button>
				</div>

				<pre class="response-content"><code id="responseContent" class="language-json"></code></pre>

			</div>
		</div>
	</div>
	<link href="https://cdn.jsdelivr.net/npm/prismjs/themes/prism-tomorrow.css" rel="stylesheet" />
	<script src="https://cdn.jsdelivr.net/npm/prismjs/prism.js"></script>
	<script src="https://cdn.jsdelivr.net/npm/prismjs/components/prism-json.min.js"></script>

	<script>
		let stats = { total: 0, success: 0, error: 0 };
		let currentResponse = {};
		
		// Method selector
		document.querySelectorAll('.method-btn').forEach(btn => {
			btn.addEventListener('click', function() {
				document.querySelectorAll('.method-btn').forEach(b => b.classList.remove('active'));
				this.classList.add('active');
				document.querySelector('input[name="method"]').value = this.dataset.method;
				toggleBodyField();
			});
		});
		
		function toggleBodyField() {
			const method = document.querySelector('input[name="method"]').value;
			const bodyField = document.getElementById('bodyField');
			bodyField.style.display = ['POST', 'PUT', 'PATCH'].includes(method) ? 'block' : 'none';
		}
		
		// Tab switching
		document.querySelectorAll('.tab-btn').forEach(btn => {
			btn.addEventListener('click', function() {
				document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
				this.classList.add('active');
				showTabContent(this.dataset.tab);
			});
		});
		
		function showTabContent(tab) {
			const content = document.getElementById('responseContent');
			switch(tab) {
				case 'response':
					content.textContent = JSON.stringify(currentResponse.bodyParsed, null, 2);
					Prism.highlightElement(content);
					break;
				case 'headers':
					content.textContent = currentResponse.headers ? JSON.stringify(currentResponse.headers, null, 2) : 'No headers data';
					break;
				case 'timing':
					content.textContent = 'Request Duration: ' + (currentResponse.timing || 'N/A') + 'ms\nTimestamp: ' + (currentResponse.timestamp || 'N/A');
					break;
			}
		}
		
		function updateStats() {
			document.getElementById('totalRequests').textContent = stats.total;
			document.getElementById('successCount').textContent = stats.success;
			document.getElementById('errorCount').textContent = stats.error;
		}
		
		// Form submission
		document.getElementById('taskForm').addEventListener('submit', async function(e) {
			e.preventDefault();
			
			const formData = new FormData(this);
			const btn = document.querySelector('.execute-btn');
			const spinner = document.querySelector('.loading-spinner');
			const btnText = document.querySelector('.btn-text');
			const resultsCard = document.getElementById('resultsCard');
			const responseStatus = document.getElementById('responseStatus');
			const responseContent = document.getElementById('responseContent');
			
			// Show loading state
			btn.disabled = true;
			spinner.style.display = 'inline-block';
			btnText.textContent = 'Executing...';
			
			const startTime = performance.now();
			
			try {
				const response = await fetch('/admin/run-task', {
					method: 'POST',
					body: formData
				});
				
				const endTime = performance.now();
				const duration = Math.round(endTime - startTime);
				
				const result = await response.text();
				
				// Update stats
				stats.total++;
				if (response.ok) {
					stats.success++;
				} else {
					stats.error++;
				}
				updateStats();
				
				// Store response data
				currentResponse = {
					body: result,
					headers: Object.fromEntries(response.headers.entries()),
					timing: duration,
					timestamp: new Date().toISOString(),
					status: response.status,
					statusText: response.statusText
				};
				
				// Show status
				const statusClass = response.ok ? 'status-success' : 'status-error';
				responseStatus.innerHTML = '<span class="' + statusClass + '">' + response.status + ' ' + response.statusText + '</span> - ' + duration + 'ms';
				
				// Show response
				responseContent.textContent = result;
				resultsCard.style.display = 'block';
				
			} catch (error) {
				stats.total++;
				stats.error++;
				updateStats();
				
				responseStatus.innerHTML = '<span class="status-error">Network Error</span>';
				responseContent.textContent = 'Error: ' + error.message;
				resultsCard.style.display = 'block';
			} finally {
				// Reset button state
				btn.disabled = false;
				spinner.style.display = 'none';
				btnText.textContent = 'Execute Request';
			}
		});
		
		// Auto-refresh functionality
		let autoRefreshInterval;
		document.getElementById('autoRefresh').addEventListener('change', function() {
			if (this.checked) {
				autoRefreshInterval = setInterval(() => {
					const form = document.getElementById('taskForm');
					if (form.querySelector('input[name="url"]').value) {
						form.dispatchEvent(new Event('submit'));
					}
				}, 30000); // 30 seconds
			} else {
				clearInterval(autoRefreshInterval);
			}
		});
		
		// Initialize
		toggleBodyField();
		updateStats();
	</script>
</body>
</html>`
	
	t, err := template.New("admin").Parse(tmpl)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "text/html")
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, "Execution error", http.StatusInternalServerError)
	}
}

// AdminRunTask executes a proxy call to the provided URL and displays the response.
/*func AdminRunTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}
	method := r.FormValue("method")
	url := r.FormValue("url")

	respBytes, _, err := toRun(method, url, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error running task: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "<pre>%s</pre>", string(respBytes))
}*/

func AdminRunTask(w http.ResponseWriter, r *http.Request) {
    err := r.ParseMultipartForm(32 << 20) // Parse FormData correctly
    if err != nil {
        http.Error(w, "Invalid multipart form", http.StatusBadRequest)
        log.Println("[AdminRunTask] ParseMultipartForm error:", err)
        return
    }

    method := r.FormValue("method")
    url := r.FormValue("url")

    log.Printf("[AdminRunTask] method=%s url=%s", method, url)

    if url == "" {
        http.Error(w, "URL is required", http.StatusBadRequest)
        return
    }

    respBytes, _, err := toRun(method, url, nil)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error running task: %v", err), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "<pre>%s</pre>", string(respBytes))
}

// HandleSchedule optionally allows scheduling tasks via POST /api/v1/schedule
func HandleSchedule(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	url, _ := data["url"].(string)
	StartSchedule(1, "minute", 5, toRun, http.MethodPost, url, nil)
	w.Write([]byte(`{"message": "Scheduled task to run in 5 minutes"}`))
}