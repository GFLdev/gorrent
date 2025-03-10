package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Printf("could not close file: %s\n", err.Error())
	}
}

func flushWriter(writer *bufio.Writer) {
	err := writer.Flush()
	if err != nil {
		fmt.Printf("could not flush writer: %s\n", err.Error())
	}
}

func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %w", path, err)
	}
	defer closeFile(file)

	buffer := make([]byte, 4096)
	content := make([]byte, 0)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("could not read file %s: %w", path, err)
		}
		content = append(content, buffer[:n]...)
	}
	return content, nil
}

func WriteFile(path string, content []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", path, err)
	}
	defer closeFile(file)

	writer := bufio.NewWriter(file)
	defer flushWriter(writer)
	for {
		n, err := writer.Write(content)
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("could not write to file %s: %w", path, err)
		}
		content = content[n:]
	}
	return nil
}
