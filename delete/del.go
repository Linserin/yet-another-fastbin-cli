package delete

import (
	"fmt"
	"log/slog"
	"net/http"
)

func DeleteFile(server string, secret string) error {
	req, err := http.NewRequest("DELETE", server+"/delete/"+secret, nil)
	if err != nil {
		slog.Error("Delete request failed: " + err.Error())
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Delete request failed: " + err.Error())
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Error("Delete request failed: " + resp.Status)
		return fmt.Errorf("delete request failed: %s", resp.Status)
	}
	return nil
}
