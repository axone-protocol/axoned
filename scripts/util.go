package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func writeToFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	w := bufio.NewWriter(file)
	if _, err = w.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	if err = w.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func normalizeMarkdownFiles(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			if err := normalizeMarkdownFiles(path); err != nil {
				return err
			}
		} else if strings.HasSuffix(file.Name(), ".md") {
			if err := normalizeMarkdownFile(path); err != nil {
				return err
			}
		}
	}

	return nil
}

func normalizeMarkdownFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	normalizedContent, err := normalizeMarkdownContent(f)
	if err != nil {
		return err
	}

	return writeToFile(file, normalizedContent)
}

func normalizeMarkdownContent(input io.Reader) (string, error) {
	inCodeBlock := false
	inInlineCode := false
	var output strings.Builder

	reader := bufio.NewReader(input)
	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", err
		}
		switch char {
		case '`':
			output.WriteRune(char)

			if !inCodeBlock {
				inInlineCode = !inInlineCode
			}
		case '\n':
			output.WriteRune(char)

			peekBytes, _ := reader.Peek(3)
			if string(peekBytes) == "```" {
				inCodeBlock = !inCodeBlock

				_, err := io.CopyN(&output, reader, 3)
				if err != nil {
					return "", err
				}
			}
		default:
			if !inCodeBlock && !inInlineCode {
				switch char {
				case '{':
					output.WriteRune('\\')
					output.WriteRune(char)
				case '}':
					output.WriteRune('\\')
					output.WriteRune(char)
				case '<':
					output.WriteString("&lt;")
				case '>':
					output.WriteString("&gt;")
				default:
					output.WriteRune(char)
				}
				continue
			}
			output.WriteRune(char)
		}
	}

	return output.String(), nil
}
