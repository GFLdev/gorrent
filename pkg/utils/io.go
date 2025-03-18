package utils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ReadFile reads the contents of a file at the given path and returns the data as a byte slice or an error.
func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", path, err)
	}

	buf := make([]byte, 4096)
	content := make([]byte, 0)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("could not read file %s: %w", path, err)
		}
		content = append(content, buf[:n]...)
	}
	return content, file.Close()
}

// WriteFile writes the provided content into a file, create if it does not exist.
func WriteFile(path string, content []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", path, err)
	}

	w := bufio.NewWriter(file)
	for len(content) > 0 {
		n, err := w.Write(content)
		if err != nil {
			return fmt.Errorf("could not write to file %s: %w", path, err)
		}
		content = content[n:]
	}
	_ = w.Flush()
	return file.Close()
}

// ExtractResponseData reads the HTTP response body and returns it as a byte slice.
func ExtractResponseData(res *http.Response) ([]byte, error) {
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	return data, res.Body.Close()
}
