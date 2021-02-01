package main

import (
	"mvdan.cc/sh/v3/syntax"
	"os"
	"io"
	"fmt"
	"bytes"
	"strings"
	"time"
	"github.com/thatisuday/commando"
	"path/filepath"
	"io/ioutil"
)

func main() {
	// configure commando
	commando.SetExecutableName("bash_bundler").
			 SetVersion("1.0.0").
			 SetDescription("This simple tool bundles bash files into a single bash file.")
	
	commando.Register(nil).
			SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
				fmt.Printf("Printing options of the `info` command...\n\n")

				// print arguments
				for k, v := range args {
					fmt.Printf("arg -> %v: %v(%T)\n", k, v.Value, v.Value)
				}

				// print flags
				for k, v := range flags {
					fmt.Printf("flag -> %v: %v(%T)\n", k, v.Value, v.Value)
				}
			})
	
	commando.Register("bundle").
			 SetShortDescription("bundle a bash script").
			 SetDescription("Takes an entry bash script and bundles it and all its sources into a single output file.").
			 AddFlag("entry,e", "The entrypoint to the bash script to bundle.", commando.String, nil).
			 AddFlag("output,o", "The output file to write to", commando.String, nil).
			 AddFlag("minify,m", "Minify the output file", commando.Bool, false).
			 SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
				fmt.Printf("Performing bundling on entrypoint: %v\n", flags["entry"].Value)
				fmt.Printf("Bundling output to: %v\n", flags["output"].Value)
				entry, err := flags["entry"].GetString()
				output, err :=  flags["output"].GetString()
				min, err := flags["minify"].GetBool()
				if (err != nil) {
					fmt.Println("Error reading parameters.")
				}
				content, err := bundle(entry, true)
				if err != nil {
					fmt.Println("ERROR:  " + err.Error())
					return
				}
				if min {
					content, err = minify(content)
					if err != nil {
						fmt.Println("ERROR:  " + err.Error())
						return
					}
				}
				err = writeToFile(output, content)
				if err != nil {
					fmt.Println("ERROR:  " + err.Error())
					return
				}

				
			 })
	
	commando.Parse(nil)
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

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' || s[0] == '\'') && 
		   (s[len(s) -1] == '"' || s[len(s) -1] == '\'' ) {
			return s[1 : len(s) - 1]
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
	var breakLine string = strings.Repeat("#", 80) + "\n"
	s += breakLine
	s += "#  File:  " + path + "\n"
	s += "#  Bundle Date: " + time.Now().Format("2006-01-02 3:4:5") + "\n"
	s += breakLine + "\n"
	return s
}

func footer(path string) string {
	var s string = "\n\n"
	var breakLine string = strings.Repeat("#", 80) + "\n"
	s += breakLine
	s += "#  End File:  " + path + "\n"
	s += breakLine + "\n\n"
	return s
}

func minify(content string) (string, error) {
	var reader = strings.NewReader(content)
	var output string = "#!/bin/bash\n"
	var buffer = bytes.NewBufferString(output)
	var printer = syntax.NewPrinter(syntax.Minify(true), syntax.KeepPadding(false))
	var parser = syntax.NewParser(syntax.KeepComments(false))
	node, err := parser.Parse(reader, "")
	if err != nil {
		return "", err
	}
	err = printer.Print(buffer, node)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func bundle(path string, keepSheBang bool) (string, error) {
	var output string = header(path, keepSheBang)
	var directory string = filepath.Dir(path)
	var buffer = bytes.NewBufferString(output)
	var printer = syntax.NewPrinter(syntax.Indent(4), syntax.KeepPadding(true))
	var parser = syntax.NewParser(syntax.KeepComments(true))
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	in := io.Reader(file)
	err = parser.Stmts(in, func(stmt *syntax.Stmt) bool {
		var line string = ""
		var internalBuffer = bytes.NewBufferString(line)
		printer.Print(internalBuffer, stmt)
		temp := strings.Split(internalBuffer.String(), "\n") 
		for _, s := range temp {
			if strings.Contains(s, "source") && !strings.HasPrefix(strings.TrimSpace(s), "#") {
				set := strings.Split(strings.TrimSpace(s), " ")
				set = deleteEmpty(set)
				subPath := directory + "/" + trimQuotes(strings.TrimSpace(set[1]))
				fmt.Println("Bundling Source: " + subPath)
				sub, err := bundle(subPath, false)
				if (err != nil) {
					fmt.Println(err.Error())
					return false
				}
				buffer.WriteString(sub)
			} else if !isShebang(s) {
				buffer.WriteString(s + "\n")
			}
		}
		return true
	})
	buffer.WriteString(footer(path))
	return buffer.String(), nil
}

func writeToFile(path string, content string) error {
	err := ioutil.WriteFile(path, []byte(content), 0644)
	return err
}