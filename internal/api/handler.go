package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/numtide/banner-generator/internal/banner"
	"github.com/numtide/banner-generator/internal/config"
	"github.com/numtide/banner-generator/internal/converter"
	"github.com/numtide/banner-generator/internal/github"
)

// Handler handles HTTP requests
type Handler struct {
	svgBuilder   banner.Builder
	githubClient *github.Client
	config       *config.Config
}

// NewHandler creates a new API handler
func NewHandler(svgBuilder banner.Builder, githubClient *github.Client, cfg *config.Config) *Handler {
	return &Handler{
		svgBuilder:   svgBuilder,
		githubClient: githubClient,
		config:       cfg,
	}
}

// HealthCheck returns the health status of the service
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "healthy",
		"version": "1.0.0",
		"time":    time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log error but don't write to response as headers are already sent
		log.Printf("Failed to encode health check response: %v", err)
	}
}

// GenerateBanner generates an SVG banner for a GitHub repository
func (h *Handler) GenerateBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	owner := vars["owner"]
	repo := vars["repo"]

	if owner == "" || repo == "" {
		http.Error(w, "Invalid repository format", http.StatusBadRequest)
		return
	}

	// Check access control
	if !h.config.IsAllowed(owner, repo) {
		http.Error(w, "Access denied: This repository is not allowed", http.StatusForbidden)
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Fetch repository data
	repoData, err := h.githubClient.GetRepositoryData(ctx, owner, repo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch repository data: %v", err), http.StatusNotFound)
		return
	}

	// Generate SVG
	svg, err := h.svgBuilder.BuildBanner(repoData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate banner: %v", err), http.StatusInternalServerError)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(h.config.HTTPCacheDuration.Seconds())))
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Write SVG
	if _, err := w.Write([]byte(svg)); err != nil {
		// Log error but can't send error response as headers are already sent
		log.Printf("Failed to write SVG response: %v", err)
	}
}

// GeneratePNGBanner generates a PNG banner for a GitHub repository
func (h *Handler) GeneratePNGBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	owner := vars["owner"]
	repo := vars["repo"]

	if owner == "" || repo == "" {
		http.Error(w, "Invalid repository format", http.StatusBadRequest)
		return
	}

	// Check access control
	if !h.config.IsAllowed(owner, repo) {
		http.Error(w, "Access denied: This repository is not allowed", http.StatusForbidden)
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Fetch repository data
	repoData, err := h.githubClient.GetRepositoryData(ctx, owner, repo)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch repository data: %v", err), http.StatusNotFound)
		return
	}

	// Generate SVG
	svg, err := h.svgBuilder.BuildBanner(repoData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate banner: %v", err), http.StatusInternalServerError)
		return
	}

	// Convert SVG to PNG
	pngData, err := converter.SVGToPNG([]byte(svg))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to convert to PNG: %v", err), http.StatusInternalServerError)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(h.config.HTTPCacheDuration.Seconds())))
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Write PNG
	if _, err := w.Write(pngData); err != nil {
		// Log error but can't send error response as headers are already sent
		log.Printf("Failed to write PNG response: %v", err)
	}
}

// Index returns a simple landing page
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Numtide Banner Generator</title>
    <style>
        * {
            box-sizing: border-box;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            margin: 0;
            padding: 0;
            background: #0d1117;
            color: #c9d1d9;
            min-height: 100vh;
            display: flex;
            flex-direction: column;
        }
        .container {
            max-width: 900px;
            margin: 0 auto;
            padding: 2rem;
            flex: 1;
        }
        header {
            text-align: center;
            margin-bottom: 3rem;
        }
        h1 { 
            color: #58a6ff;
            font-size: 2.5rem;
            margin-bottom: 0.5rem;
        }
        .tagline {
            font-size: 1.2rem;
            color: #8b949e;
            margin-bottom: 2rem;
        }
        .generator-section {
            background: #161b22;
            border-radius: 8px;
            padding: 2rem;
            margin: 3rem 0;
            border: 1px solid #30363d;
        }
        .input-group {
            display: flex;
            gap: 1rem;
            margin-bottom: 1.5rem;
        }
        input[type="text"] {
            flex: 1;
            padding: 0.8rem;
            font-size: 1rem;
            background: #0d1117;
            border: 1px solid #30363d;
            border-radius: 6px;
            color: #c9d1d9;
        }
        input[type="text"]:focus {
            outline: none;
            border-color: #58a6ff;
        }
        button {
            padding: 0.8rem 1.5rem;
            font-size: 1rem;
            background: #238636;
            color: white;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            transition: background 0.2s;
        }
        button:hover {
            background: #2ea043;
        }
        button.secondary {
            background: #21262d;
            border: 1px solid #30363d;
        }
        button.secondary:hover {
            background: #30363d;
        }
        .format-selector {
            display: flex;
            gap: 2rem;
            justify-content: center;
            margin: 2rem 0;
        }
        .format-selector label {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            cursor: pointer;
        }
        .error {
            color: #f85149;
            margin: 1rem 0;
            padding: 0.8rem;
            background: rgba(248, 81, 73, 0.1);
            border: 1px solid rgba(248, 81, 73, 0.3);
            border-radius: 6px;
            display: none;
        }
        .results-section {
            display: none;
            margin-top: 2rem;
        }
        .results-section.active {
            display: block;
        }
        .preview-container {
            text-align: center;
            margin: 2rem 0;
        }
        .preview-container img {
            max-width: 100%;
            border-radius: 8px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.3);
        }
        code {
            background: #161b22;
            padding: 0.2rem 0.4rem;
            border-radius: 3px;
            font-family: monospace;
        }
        pre {
            background: #0d1117;
            padding: 1rem;
            border-radius: 6px;
            overflow-x: auto;
            border: 1px solid #30363d;
            position: relative;
        }
        .copy-button {
            position: absolute;
            top: 0.5rem;
            right: 0.5rem;
            padding: 0.4rem 0.8rem;
            font-size: 0.875rem;
            background: #30363d;
            color: #c9d1d9;
            border: 1px solid #484f58;
            border-radius: 6px;
            cursor: pointer;
            transition: all 0.2s;
        }
        .copy-button:hover {
            background: #484f58;
            border-color: #6e7681;
        }
        .copy-button.copied {
            background: #238636;
            border-color: #2ea043;
        }
        .action-buttons {
            display: flex;
            gap: 1rem;
            justify-content: center;
            margin-top: 2rem;
        }
        footer {
            background: #010409;
            border-top: 1px solid #30363d;
            padding: 2rem;
            text-align: center;
            margin-top: auto;
        }
        footer p {
            margin: 0;
            color: #8b949e;
        }
        footer a {
            color: #58a6ff;
            text-decoration: none;
        }
        footer a:hover {
            text-decoration: underline;
        }
        .footer-links {
            display: flex;
            gap: 1.5rem;
            justify-content: center;
            align-items: center;
            flex-wrap: wrap;
        }
        .footer-divider {
            color: #30363d;
            user-select: none;
        }
        .section-title {
            font-size: 1.2rem;
            color: #58a6ff;
            margin-bottom: 1rem;
        }
        .modal {
            display: none;
            position: fixed;
            z-index: 1000;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0, 0, 0, 0.8);
            animation: fadeIn 0.2s ease-out;
        }
        .modal.active {
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .modal-content {
            background-color: #161b22;
            padding: 2rem;
            border-radius: 8px;
            border: 1px solid #30363d;
            max-width: 600px;
            width: 90%;
            max-height: 90vh;
            overflow-y: auto;
            position: relative;
            animation: slideIn 0.2s ease-out;
        }
        .modal-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 1.5rem;
        }
        .modal-title {
            font-size: 1.5rem;
            color: #58a6ff;
            margin: 0;
        }
        .close-button {
            background: none;
            border: none;
            color: #8b949e;
            font-size: 1.5rem;
            cursor: pointer;
            padding: 0;
            width: 2rem;
            height: 2rem;
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 4px;
            transition: all 0.2s;
        }
        .close-button:hover {
            background: #30363d;
            color: #c9d1d9;
        }
        .modal-body {
            color: #c9d1d9;
        }
        .modal-body h3 {
            color: #58a6ff;
            margin-top: 1.5rem;
            margin-bottom: 0.5rem;
        }
        .modal-body ol {
            margin: 0.5rem 0;
            padding-left: 1.5rem;
        }
        .modal-body li {
            margin: 0.5rem 0;
        }
        .modal-body pre {
            margin: 1rem 0;
        }
        @keyframes fadeIn {
            from { opacity: 0; }
            to { opacity: 1; }
        }
        @keyframes slideIn {
            from {
                transform: translateY(-20px);
                opacity: 0;
            }
            to {
                transform: translateY(0);
                opacity: 1;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>Numtide Banner Generator</h1>
            <p class="tagline">Generate beautiful, consistent banners for Numtide repositories</p>
        </header>

        <div class="generator-section">
            <div class="input-group">
                <input type="text" id="repoInput" placeholder="Enter repository (e.g., owner/repo)" value="numtide/banner-generator">
                <button onclick="generateBanner()">Generate Banner</button>
            </div>
            
            <div class="error" id="error"></div>
            
            <div class="results-section" id="results">
                <div class="preview-container">
                    <img id="bannerPreview" alt="Generated banner" />
                </div>
                
                <div class="format-selector">
                    <label>
                        <input type="radio" name="format" value="svg" checked>
                        SVG (Vector - Best for README)
                    </label>
                    <label>
                        <input type="radio" name="format" value="png">
                        PNG (Raster - For Social Media)
                    </label>
                </div>
                
                <div>
                    <h3 class="section-title">Usage</h3>
                    <p>Add this to your README.md:</p>
                    <pre><button class="copy-button" onclick="copyToClipboard('usageCode', this)">Copy</button><code id="usageCode"></code></pre>
                    
                    <p>Direct URL:</p>
                    <pre><button class="copy-button" onclick="copyToClipboard('directUrl', this)">Copy</button><code id="directUrl"></code></pre>
                </div>
                
                <div class="action-buttons">
                    <button class="secondary" onclick="downloadBanner('svg')">Download SVG</button>
                    <button class="secondary" onclick="downloadBanner('png')">Download PNG</button>
                </div>
                
                <div style="margin-top: 2rem;">
                    <h3 class="section-title">Setting as GitHub Social Preview</h3>
                    <p>To use this banner as your repository's social media preview:</p>
                    <ol style="margin: 1rem 0; padding-left: 1.5rem; color: #c9d1d9;">
                        <li>Download the PNG version using the "Download PNG" button above</li>
                        <li>Go to your <a id="settingsLink" href="#" target="_blank" style="color: #58a6ff;">repository settings</a></li>
                        <li>Scroll down to the "Social preview" section</li>
                        <li>Click "Edit" and upload the downloaded PNG</li>
                    </ol>
                    <p style="font-size: 0.9rem; color: #8b949e;">The social preview appears when sharing your repository on social media platforms.</p>
                </div>
            </div>
        </div>
    </div>

    <footer>
        <div class="footer-links">
            <span>Powered by <a href="https://github.com/numtide/banner-generator" target="_blank">numtide/banner-generator</a></span>
            <span class="footer-divider">|</span>
            <a href="https://github.com/numtide/banner-generator/issues" target="_blank">Report an Issue</a>
            <span class="footer-divider">|</span>
            <a href="https://numtide.com" target="_blank">Numtide</a>
        </div>
    </footer>

    <script>
        function copyToClipboard(elementId, button) {
            const element = document.getElementById(elementId);
            const text = element.textContent;
            
            navigator.clipboard.writeText(text).then(() => {
                button.textContent = 'Copied!';
                button.classList.add('copied');
                
                setTimeout(() => {
                    button.textContent = 'Copy';
                    button.classList.remove('copied');
                }, 2000);
            }).catch(err => {
                console.error('Failed to copy:', err);
                // Fallback for older browsers
                const textArea = document.createElement('textarea');
                textArea.value = text;
                textArea.style.position = 'fixed';
                textArea.style.left = '-999999px';
                document.body.appendChild(textArea);
                textArea.select();
                try {
                    document.execCommand('copy');
                    button.textContent = 'Copied!';
                    button.classList.add('copied');
                    setTimeout(() => {
                        button.textContent = 'Copy';
                        button.classList.remove('copied');
                    }, 2000);
                } catch (err) {
                    console.error('Fallback copy failed:', err);
                }
                document.body.removeChild(textArea);
            });
        }
        
        function generateBanner() {
            const input = document.getElementById('repoInput');
            const repoPath = input.value.trim();
            const error = document.getElementById('error');
            const results = document.getElementById('results');
            const format = document.querySelector('input[name="format"]:checked').value;
            
            // Reset error
            error.style.display = 'none';
            error.textContent = '';
            
            // Validate input
            if (!repoPath || !repoPath.includes('/')) {
                error.textContent = 'Please enter a valid repository in the format: owner/repo';
                error.style.display = 'block';
                results.classList.remove('active');
                return;
            }
            
            // Get current domain
            const domain = window.location.origin;
            const bannerUrl = domain + '/banner/' + repoPath + '.' + format;
            
            // Update preview
            const img = document.getElementById('bannerPreview');
            img.src = bannerUrl + '?t=' + new Date().getTime();
            
            // Update usage examples
            document.getElementById('usageCode').textContent = 
                '![Banner](' + bannerUrl + ')';
            document.getElementById('directUrl').textContent = bannerUrl;
            
            // Update settings link
            const [owner, repo] = repoPath.split('/');
            document.getElementById('settingsLink').href = 'https://github.com/' + owner + '/' + repo + '/settings';
            
            // Handle load error
            img.onerror = function() {
                error.textContent = 'Failed to generate banner. Please check if the repository exists and is accessible.';
                error.style.display = 'block';
                results.classList.remove('active');
            };
            
            img.onload = function() {
                results.classList.add('active');
            };
        }
        
        function downloadBanner(format) {
            const repoPath = document.getElementById('repoInput').value.trim();
            if (!repoPath || !repoPath.includes('/')) {
                return;
            }
            
            const link = document.createElement('a');
            const domain = window.location.origin;
            link.href = domain + '/banner/' + repoPath + '.' + format;
            link.download = repoPath.replace('/', '-') + '-banner.' + format;
            link.click();
        }
        
        // Update on Enter key
        document.getElementById('repoInput').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                generateBanner();
            }
        });
        
        // Update when format changes
        document.querySelectorAll('input[name="format"]').forEach(radio => {
            radio.addEventListener('change', function() {
                const results = document.getElementById('results');
                if (results.classList.contains('active')) {
                    generateBanner();
                }
            });
        });
        
        // Generate default banner on page load
        document.addEventListener('DOMContentLoaded', function() {
            generateBanner();
        });
    </script>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte(html)); err != nil {
		// Log error but can't send error response as headers are already sent
		log.Printf("Failed to write HTML response: %v", err)
	}
}
