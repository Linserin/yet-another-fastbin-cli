package fetch

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
)

func FetchFile(server string, key string, secret string) ([]byte, error) {
	slog.Debug("Start Download", "url", server+"/"+key)
	fileContent, err := http.Get(server + "/raw/" + key)
	if err != nil {
		slog.Error("Download failed: " + err.Error())
		return nil, err
	}
	if fileContent.StatusCode != http.StatusOK {
		slog.Error("Download failed: " + fileContent.Status)
		err = errors.New("Failed to Download File")
		return nil, err
	}
	slog.Debug("Download Success")
	content, err := io.ReadAll(fileContent.Body)
	if err != nil {
		slog.Error("Read file content failed: " + err.Error())
		return nil, err
	}
	slog.Debug("Read file content Success")
	return content, nil
}
