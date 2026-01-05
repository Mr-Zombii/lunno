package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func send(msg any) {
	data, err := json.Marshal(msg)
	if err != nil {
		_, err := fmt.Fprintln(os.Stderr, "Failed to marshal LSP message:", err)
		if err != nil {
			return
		}
		return
	}
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	_, err = os.Stdout.Write([]byte(header))
	if err != nil {
		_, err := fmt.Fprintln(os.Stderr, "Failed to write header:", err)
		if err != nil {
			return
		}
		return
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		_, err := fmt.Fprintln(os.Stderr, "Failed to write data:", err)
		if err != nil {
			return
		}
		return
	}
	err = os.Stdout.Sync()
	if err != nil {
		return
	}
}

func readMessage(r *bufio.Reader) ([]byte, error) {
	headers := map[string]string{}
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	length := 0
	_, err := fmt.Sscanf(headers["Content-Length"], "%d", &length)
	if err != nil {
		return nil, err
	}
	body := make([]byte, length)
	_, err = io.ReadFull(r, body)
	return body, err
}
