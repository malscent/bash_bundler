package bundler

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mvdan.cc/sh/v3/syntax"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("bundler")

func CheckError(err error) {
	if err != nil {
		log.Error(err, "error during bundling")
		os.Exit(1)
	}
}

func isShebang(line string) bool {
	return strings.HasPrefix(line, "#!")
}

func deleteEmpty(s []string) []string {
	var r []string

	for _, str := range s {
		if strings.TrimSpace(str) != "" {
			r = append(r, str)
		}
	}

	return r
}

const emptyStringIndicator = 2

func trimQuotes(s string) string {
	if len(s) >= emptyStringIndicator {
		if (s[0] == '"' || s[0] == '\'') &&
			(s[len(s)-1] == '"' || s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}

	return s
}

func header(path string, keepSheBang bool) string {
	var s string
	if keepSheBang {
		s = "#!/bin/bash"
	}

	s += "\n\n"

	breakLine := strings.Repeat("#", 80) + "\n"
	s += breakLine
	s += "#  File:  " + path + "\n"
	s += "#  Bundle Date: " + time.Now().Format("2006-01-02 3:4:5") + "\n"
	s += breakLine + "\n"

	return s
}

func footer(path string) string {
	s := "\n\n"
	breakLine := strings.Repeat("#", 80) + "\n"
	s += breakLine
	s += "#  End File:  " + path + "\n"
	s += breakLine + "\n\n"

	return s
}

func Minify(content string) (string, error) {
	reader := strings.NewReader(content)
	output := "#!/bin/bash\n"
	buffer := bytes.NewBufferString(output)
	printer := syntax.NewPrinter(syntax.Minify(true), syntax.KeepPadding(false))
	parser := syntax.NewParser(syntax.KeepComments(false))
	node, err := parser.Parse(reader, "")

	if err != nil {
		return "", fmt.Errorf("error while minifying: %w", err)
	}

	err = printer.Print(buffer, node)
	if err != nil {
		return "", fmt.Errorf("error while printing: %w", err)
	}

	return buffer.String(), nil
}

const indentSize = 4

func Bundle(path string, keepSheBang bool) (string, error) {
	output := header(path, keepSheBang)

	directory := trimQuotes(filepath.Dir(path))
	if strings.HasPrefix(path, ".") {
		directory = "./" + directory
	}

	buffer := bytes.NewBufferString(output)
	printer := syntax.NewPrinter(syntax.Indent(indentSize), syntax.KeepPadding(true))
	parser := syntax.NewParser(syntax.KeepComments(true))

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}

	in := io.Reader(file)
	err = parser.Stmts(in, func(stmt *syntax.Stmt) bool {
		line := ""
		var internalBuffer = bytes.NewBufferString(line)
		printer.Print(internalBuffer, stmt)
		temp := strings.Split(internalBuffer.String(), "\n")
		for _, s := range temp {
			if strings.Contains(s, "source") && !strings.HasPrefix(strings.TrimSpace(s), "#") {
				set := strings.Split(strings.TrimSpace(s), " ")
				set = deleteEmpty(set)
				embedded := false
				sourcePath := strings.TrimSpace(set[1])
				if strings.HasSuffix(sourcePath, ")") {
					sourcePath = strings.TrimSuffix(sourcePath, ")")
					embedded = true
				}
				sourcePath = trimQuotes(sourcePath)
				if strings.HasPrefix(sourcePath, "./") {
					sourcePath = strings.TrimPrefix(sourcePath, "./")
				}

				subPath := directory + "/" + sourcePath
				log.Info("Bundling Source", "sourceFile", subPath)
				sub, err := Bundle(subPath, false)
				if err != nil {
					log.Error(err, "error during parse")

					return false
				}
				if embedded {
					sub = strings.TrimSuffix(set[0], "source") + sub + ")\n"
				}
				buffer.WriteString(sub)
			} else if !isShebang(s) {
				buffer.WriteString(s + "\n")
			}
		}

		return true
	})

	if err != nil {
		return "", fmt.Errorf("error while parsing statements: %w", err)
	}

	buffer.WriteString(footer(path))

	return buffer.String(), nil
}

func WriteToFile(path string, content string) error {
	err := ioutil.WriteFile(path, []byte(content), 0600)
	if err != nil {
		return fmt.Errorf("error writing to file :%w", err)
	}

	return nil
}
