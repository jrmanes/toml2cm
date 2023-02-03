package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	// spacesToText yaml format spaces in a ConfigMap
	spacesToText = "    "
)

// Run start the service here
func Run() {
	// filePath path to the file that we'll scan
	file := flag.String("file", "./test.toml", "Toml file to convert to Kubernetes ConfigMap")

	flag.Parse()

	// Open file
	f, err := os.Open(*file)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	// scan the file line by line
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	// close the file when the function will finish
	defer f.Close()

	// Clean up the file name
	fileNameCleaned := cleanUpFileName(*file)

	// Check line by line
	for _, eachline := range txtlines {
		// add the spaces to the beginning of the line - yaml format
		eachline = spacesToText + eachline

		if strings.Contains(eachline, "=") {
			eachline = formatLine(eachline, fileNameCleaned)
		}

		fmt.Println(eachline)
	}
}

func formatLine(line, fileName string) string {
	l := line[:strings.IndexByte(line, '=')]

	l = strings.ReplaceAll(l, " ", "")
	line = strings.ReplaceAll(line, "= ", " = {{ .Values.configMaps."+fileName+"."+l+" | "+"default "+"\"")
	line = line + "\"" + " | quote }}"

	// example line
	// broadcast-mode = {{ .Values.configMaps.client_toml.broadcast-mode | default "sync" | quote }}
	return line
}

// cleanUpFileName clean format in the file name
func cleanUpFileName(f string) string {
	// change . to _ in the filename
	f = strings.ReplaceAll(f, ".", "_")
	// change - to _ in the filename
	f = strings.ReplaceAll(f, "-", "_")

	return f
}

// main everything starts here!
func main() {
	// Run -> entrypoint
	Run()
}
