package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

/*
Upload:

$ curl -i -X POST http://localhost:3000/api/documents --data-binary @go.mod
HTTP/1.1 200 OK

{"ok":true,"key":"key","secret":"abcd"}
*/

// uploadResponse mirrors the JSON body returned by the server.
type uploadResponse struct {
	OK     bool   `json:"ok"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
	Error  string `json:"error,omitempty"`
}

func UploadFile(fileContent []byte, serverAddr string) (string, string, string, error) {
	slog.Debug("Uploading file to server: " + serverAddr)
	slog.Debug("File Content" + string(fileContent))

	// Allow serverAddr to be passed either as "localhost:3000" or as a full URL.
	base := serverAddr
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}
	url := strings.TrimRight(base, "/") + "/api/documents"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(fileContent))
	if err != nil {
		return "", "", "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("send request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result uploadResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", "", fmt.Errorf("decode response JSON %q: %w", string(body), err)
	}
	if !result.OK {
		return "", "", "", fmt.Errorf("server reported failure: %s", result.Error)
	}
	if result.Key == "" || result.Secret == "" {
		return "", "", "", fmt.Errorf("server response missing key/secret: %s", string(body))
	}

	return result.Key, result.Secret, serverAddr, nil
}
