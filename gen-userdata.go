package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

func createUserdata(tmplFilename string) *string {
	tmplContent, err := ioutil.ReadFile(tmplFilename)
	if err != nil {
		log.Fatal("Unable to read template:", err)
	}

	templateFuncs := template.FuncMap{
		"mode": func(input uint) uint {
			return input
		},
		"base64": func(filename string) string {
			result, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatal("Cannot read file:", filename, ":", err)
			}

			return fmt.Sprintf("data:text/plain;base64,%s", base64.StdEncoding.EncodeToString(result))
		},
		"escape": func(filename string) string {
			result, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatal(err)
			}

			escapedNewlines := strings.Replace(string(result), "\n", "\\n", -1)
			escapedDoubleQuotes := strings.Replace(escapedNewlines, `"`, `\"`, -1)
			return fmt.Sprintf("%s", escapedDoubleQuotes)
		},
	}

	t, err := template.New("userdata").Funcs(templateFuncs).Parse(string(tmplContent))

	if err != nil {
		log.Fatal("Failed to parse template.", err)
	}

	var templateOutput bytes.Buffer
	err = t.Execute(&templateOutput, nil)

	if err != nil {
		log.Fatal(err)
	}

	userdata := templateOutput.String()
	return &userdata
}

func main() {
	log.Println("Container Linux Ignition Config Generator")
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %v <path/to/userdata.tmpl>\n", os.Args[0])
	}

	userdata := createUserdata(os.Args[1])
	_, err := os.Stdout.WriteString(*userdata)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Ignition config generated.")
}
