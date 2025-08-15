package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	rootDir := "./www"

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			scriptPath := r.Form.Get("scriptPath")
			buttonName := r.Form.Get("buttonName")

			if scriptPath == "" || buttonName == "" {
				http.Error(w, "Missing scriptPath or buttonName parameter", http.StatusBadRequest)
				return
			}

			fullScriptPath := filepath.Join(rootDir, "..", scriptPath)
			output, err := executeScript(fullScriptPath)
			if err != nil {
				output = fmt.Sprintf("[Error executing script: %v]", err)
			}

			htmlContent := fmt.Sprintf(`
                <h1>Welcome to Netgrate</h1>
                <p>Your external IP address is: <<go(scripts/hello.go)>></p>
                <form action="/" method="POST" style="display: inline;">
                    <input type="hidden" name="scriptPath" value="%s">
                    <input type="hidden" name="buttonName" value="%s">
                    <button type="submit" class="btn primary-btn">%s</button>
                </form>
            `, scriptPath, buttonName, output)

			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(htmlContent))
			return
		}
		filePath := filepath.Join(rootDir, r.URL.Path)
		if filePath == rootDir+"/" || filePath == rootDir {
			filePath = filepath.Join(rootDir, "index.html")
		} else if !strings.HasSuffix(filePath, ".html") {
			filePath += ".html"
		}
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			serveErrorPage(w, 404)
			return
		}
		htmlContent, err := readFile(filePath)
		if err != nil {
			serveErrorPage(w, 500)
			return
		}

		processedHTML, err := processGoTags(htmlContent, rootDir)
		if err != nil {
			serveErrorPage(w, 500)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(processedHTML)
	}))

	println("Server started at http://localhost")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		println("Error starting server:", err)
	}
}

func readFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func processGoTags(htmlContent string, rootDir string) ([]byte, error) {
	goTagRe := regexp.MustCompile(`<<go\((.*?)\)>>`)
	htmlContent = goTagRe.ReplaceAllStringFunc(htmlContent, func(match string) string {
		submatches := goTagRe.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return "[Invalid tag format]"
		}
		scriptPath := submatches[1]
		fullScriptPath := filepath.Join(rootDir, "..", scriptPath)

		output, err := executeScript(fullScriptPath)
		if err != nil {
			return fmt.Sprintf("[Error executing script: %v]", err)
		}
		return output
	})

	gobuttonRe := regexp.MustCompile(`<gobutton\s+type="([^"]+)"\s+class="([^"]+)">([^<]+)</gobutton>`)
	htmlContent = gobuttonRe.ReplaceAllStringFunc(htmlContent, func(match string) string {
		submatches := gobuttonRe.FindStringSubmatch(match)
		if len(submatches) < 4 {
			return "[Invalid gobutton format]"
		}

		scriptPath := submatches[1]
		buttonClass := submatches[2]
		buttonText := submatches[3]

		return fmt.Sprintf(`
            <form action="/" method="POST" style="display: inline;">
                <input type="hidden" name="scriptPath" value="%s">
                <input type="hidden" name="buttonName" value="%s">
                <button type="submit" class="%s">%s</button>
            </form>
        `, scriptPath, buttonText, buttonClass, buttonText)
	})

	return []byte(htmlContent), nil
}

func executeScript(scriptPath string) (string, error) {
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return "", fmt.Errorf("script file does not exist: %s", scriptPath)
	}

	cmd := exec.Command("go", "run", scriptPath)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute script: %v, stderr: %s", err, stderr.String())
	}

	return strings.TrimSpace(out.String()), nil
}

func serveErrorPage(w http.ResponseWriter, statusCode int) {
	errorPages := map[int]string{
		404: `
            <!DOCTYPE html>
            <html lang="en">
            <head>
              <meta charset="UTF-8">
              <title>404 Not Found</title>
              <style>
                body {
                  background-color: gray;
                  color: purple;
                  font-family: Arial, sans-serif;
                  display: flex;
                  justify-content: center;
                  align-items: center;
                  height: 100vh;
                  margin: 0;
                }
                img {
                  max-width: 100%;
                  height: auto;
                }
              </style>
            </head>
            <body>
              <div>
                <h1>404 Not Found</h1>
                <p>The page you are looking for does not exist.</p>
                <img src="https://iili.io/3RaXzk7.gif" alt="Animated GIF">
              </div>
            </body>
            </html>
        `,
		500: `
            <!DOCTYPE html>
            <html lang="en">
            <head>
              <meta charset="UTF-8">
              <title>500 Internal Server Error</title>
              <style>
                body {
                  background-color: gray;
                  color: purple;
                  font-family: Arial, sans-serif;
                  display: flex;
                  justify-content: center;
                  align-items: center;
                  height: 100vh;
                  margin: 0;
                }
                img {
                  max-width: 100%;
                  height: auto;
                }
              </style>
            </head>
            <body>
              <div>
                <h1>500 Internal Server Error</h1>
                <p>Something went wrong on the server. Please try again later.</p>
                <img src="https://iili.io/3RaXzk7.gif" alt="Animated GIF">
              </div>
            </body>
            </html>
        `,
		429: `
            <!DOCTYPE html>
            <html lang="en">
            <head>
              <meta charset="UTF-8">
              <title>429 Too Many Requests</title>
              <style>
                body {
                  background-color: gray;
                  color: purple;
                  font-family: Arial, sans-serif;
                  display: flex;
                  justify-content: center;
                  align-items: center;
                  height: 100vh;
                  margin: 0;
                }
                img {
                  max-width: 100%;
                  height: auto;
                }
              </style>
            </head>
            <body>
              <div>
                <h1>429 Too Many Requests</h1>
                <p>You have made too many requests. Please try again later.</p>
                <img src="https://iili.io/3RaXzk7.gif" alt="Animated GIF">
              </div>
            </body>
            </html>
        `,
	}

	htmlContent, exists := errorPages[statusCode]
	if !exists {
		htmlContent = `
            <!DOCTYPE html>
            <html lang="en">
            <head>
              <meta charset="UTF-8">
              <title>Error</title>
              <style>
                body {
                  background-color: gray;
                  color: purple;
                  font-family: Arial, sans-serif;
                  display: flex;
                  justify-content: center;
                  align-items: center;
                  height: 100vh;
                  margin: 0;
                }
                img {
                  max-width: 100%;
                  height: auto;
                }
              </style>
            </head>
            <body>
              <div>
                <h1>An unexpected error occurred.</h1>
                <p>Please contact the administrator.</p>
                <img src="https://iili.io/3RaXzk7.gif" alt="Animated GIF">
              </div>
            </body>
            </html>
        `
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlContent))
}
