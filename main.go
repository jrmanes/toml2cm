package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
)

var (
	// spacesToTab yaml format spaces in a ConfigMap
	spacesToTab = "    "
	// outputPath outputs folder, where all the files will be stored
	outputPath = "./outputs/"
	// configMapKind this is the header of the file, based on a Kubernetes ConfigMap
	configMapKind = `apiVersion: v1
kind: ConfigMap
metadata:
  name: CONFIGMAP_NAME
{{- if .Values.namespace.enabled }}
  namespace: {{ .Values.namespace.name | default "default" }}
{{- end }}
data:
	`
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

	var fileContent []string

	// Add the Kubernetes ConfigMap structure at top
	configMapKind = strings.ReplaceAll(configMapKind, "CONFIGMAP_NAME", fileNameCleaned)
	fileContent = append(fileContent, configMapKind)
	// Add the data
	dataCM := "  " + *file + ": |" + "\n"
	fileContent = append(fileContent, dataCM)

	// Check line by line
	for _, eachline := range txtlines {
		// add the spaces to the beginning of the line - yaml format
		eachline = spacesToTab + eachline

		// if line contains = that means it contains a variable + value
		if strings.Contains(eachline, "=") {
			eachline = formatLine(eachline, fileNameCleaned)
			eachline = strings.ReplaceAll(eachline, "-", "_")
		}

		eachline = eachline + "\n"
		// Here append the lines to the array in order to bulk the data
		fileContent = append(fileContent, eachline)

		//fmt.Println(eachline)
	}
	// here is where we have to write the content to the new file, we've
	// formated the lines that contains values
	//fmt.Println(eachline)
	//fmt.Println(fileContent)

	writeToFile(fileNameCleaned, fileContent)
}

// main everything starts here!
func main() {
	// Run -> entrypoint
	Run()
}
