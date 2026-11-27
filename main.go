package main

import (
	"fastbin/delete"
	"fastbin/fetch"
	"fastbin/upload"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
)

/*
API Reference:
Raw File:

$ curl http://localhost:3000/raw/slog
[raw file content]

Upload:

$ curl -i -X POST http://localhost:3000/api/documents --data-binary @go.mod
HTTP/1.1 200 OK

{"ok":true,"key":"key","secret":"abcd"}

Delete:

$ curl -i -X DELETE http://localhost:3000/delete/secret
HTTP/1.1 200 OK
*/
func main() {
	for i, cli := range os.Args {
		// Check Server Addr
		server := os.Getenv("FASTBIN_ADDR")
		if server == "" {
			server = "http://localhost:3000"

		}
		// Check if stdin first.
		stat, err := os.Stdin.Stat()
		if err != nil {
			slog.Error("Failed to get stdin stat: " + err.Error())
			os.Exit(1)
		}
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			slog.Debug("Reading from stdin")
			fileContent, err := os.ReadFile("/dev/stdin")
			if err != nil {
				slog.Error("Read from stdin failed: " + err.Error())
				os.Exit(1)
			}
			key, secret, serverAddr, err := upload.UploadFile(fileContent, server)
			if err != nil {
				slog.Error("Upload failed: " + err.Error())
				os.Exit(1)
			}
			fmt.Println("Upload Success!")
			fmt.Println(serverAddr + "/" + key)
			fmt.Println("Secret: " + secret)
			os.Exit(0)
		}

		if cli == "-h" || cli == "--help" || len(os.Args) == 1 {
			fmt.Println("	fastbin")
			fmt.Println("		Yet Another CLI Client for https://github.com/jozsefsallai/fastbin-rewritten")
			fmt.Println("		Usage: ", os.Args[0], " [Args] [file]")
			fmt.Println("		Or You can also use pipe to send data to fastbin")
			fmt.Println()
			fmt.Println("           You should specify the server address with the environment variable FASTBIN_ADDR, default is http://localhost:3000")
			fmt.Println()
			fmt.Println("		e.g. echo \"Hello World\" | ", os.Args[0])
			fmt.Println("		     ", os.Args[0], " -f file.txt")
			fmt.Println("		     ", os.Args[0], " -f file.txt")
			fmt.Println("		     ", os.Args[0], " -rm abcde")
			fmt.Println("		     ", os.Args[0], " -c abcde")
			fmt.Println()
			fmt.Println("	Args:")
			fmt.Println("           -h, --help: Show this help message")
			fmt.Println("           -d, --debug: Set log level as debug")
			fmt.Println("           -f [file], Upload from a file.")
			fmt.Println("           -rm [secret]: Delete the file with the given secret")
			fmt.Println("           -c, Clone/Create a snippet with Text Editor, and then upload it to fastbin (vim/nano/...)")
			fmt.Println("               If left empty, will create a new empty snippet")
			fmt.Println("	Repo Addr: https://github.com/Linserin/yet-another-fastbin-cli")
			os.Exit(0)
		}
		if cli == "-d" || cli == "--debug" {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
			slog.SetDefault(logger)
		} else {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))
			slog.SetDefault(logger)
		}

		// Clone a snippet with vipe
		if os.Args[i] == "-c" {
			slog.Debug("Clone a snippet")
			if i+1 < len(os.Args) && os.Args[i+1] != "-s" && os.Args[i+1] != "--server" && os.Args[i+1] != "-rm" {
				slog.Debug("Clone snippet with key: " + os.Args[i+1])
				key := os.Args[i+1]
				fileContent, err := fetch.FetchFile(server, key, "")
				if err != nil {
					slog.Error("Fetch failed: " + err.Error())
					os.Exit(1)
				}
				// Use vipe to edit the file content
				slog.Debug("Start vipe...")
				vipeCmd := exec.Command("vipe")
				vipeStdin, err := vipeCmd.StdinPipe()
				if err != nil {
					slog.Error("Failed to get vipe stdin: " + err.Error())
					os.Exit(1)
				}
				vipeStdout, err := vipeCmd.StdoutPipe()
				if err != nil {
					slog.Error("Failed to get vipe stdout: " + err.Error())
					os.Exit(1)
				}
				if err := vipeCmd.Start(); err != nil {
					slog.Error("Failed to start vipe: " + err.Error())
					os.Exit(1)
				}
				_, err = vipeStdin.Write(fileContent)
				if err != nil {
					slog.Error("Failed to write to vipe stdin: " + err.Error())
					os.Exit(1)
				}
				vipeStdin.Close()
				editedContent, err := io.ReadAll(vipeStdout)
				if err != nil {
					slog.Error("Failed to read from vipe stdout: " + err.Error())
					os.Exit(1)
				}
				if err := vipeCmd.Wait(); err != nil {
					slog.Error("Failed to wait for vipe: " + err.Error())
					os.Exit(1)
				}
				key, secret, serverAddr, err := upload.UploadFile(editedContent, server)
				if err != nil {
					slog.Error("Upload failed: " + err.Error())
					os.Exit(1)
				}
				fmt.Println("Upload Success!")
				fmt.Println(serverAddr + "/" + key)
				fmt.Println("Secret: " + secret)
				os.Exit(0)
			}
		}

		// No stdin, check if file is provided
		if i+1 < len(os.Args) && os.Args[i] == "-f" {
			slog.Debug("Upload from file: " + os.Args[i+1])
			fileContent, err := os.ReadFile(os.Args[i+1])
			if err != nil {
				slog.Error("Read file failed: " + err.Error())
				os.Exit(1)
			}
			key, secret, serverAddr, err := upload.UploadFile(fileContent, server)
			if err != nil {
				slog.Error("Upload failed: " + err.Error())
				os.Exit(1)
			}
			fmt.Println("Upload Success!")
			fmt.Println(serverAddr + "/" + key)
			fmt.Println("Secret: " + secret)
			os.Exit(0)
		}

		// Delete a file with secret
		if i+1 < len(os.Args) && os.Args[i] == "-rm" {
			slog.Debug("Delete file with secret: " + os.Args[i+1])
			err := delete.DeleteFile(server, os.Args[i+1])
			if err != nil {
				slog.Error("Delete failed: " + err.Error())
				os.Exit(1)
			}
			fmt.Println("Delete Success!")
			os.Exit(0)
		}
	}
}
