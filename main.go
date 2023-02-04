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

	var fileContent []byte

	// Check line by line
	for _, eachline := range txtlines {
		// add the spaces to the beginning of the line - yaml format
		eachline = spacesToText + eachline

		// if line contains = that means it contains a variable + value
		if strings.Contains(eachline, "=") {
			eachline = formatLine(eachline, fileNameCleaned)

		}
		// Here append the lines to the array in order to bulk the data
		fileContent = append(fileContent, eachline...)

		//fmt.Println(eachline)
	}
	// here is where we have to write the content to the new file, we've
	// formated the lines that contains values
	//fmt.Println(eachline)
	//fmt.Println(fileContent)

	// TODO
	writeToFile(fileNameCleaned, fileContent)
}

// formatLine change the current format to Helm Template
func formatLine(line, fileName string) string {
	// find the = and cut the string from it
	l := line[:strings.IndexByte(line, '=')]

	// Clean the spaces
	l = strings.ReplaceAll(l, " ", "")

	// Replace the value of the variable to a Helm Template format
	// the following lines add the content between quotes, we don't really need
	// it in that way. -> pending to review
	//line = strings.ReplaceAll(line, "= ", " = {{ .Values.configMaps."+fileName+"."+l+" | "+"default "+"\"")
	//line = line + "\"" + " | quote }}"

	line = strings.ReplaceAll(
		line,
		"= ",
		" = {{ .Values.configMaps."+fileName+"."+l+" | "+"default ")
	line = line + " | quote }}"

	// example line
	// whatever = {{ .Values.configMaps.fileName_toml.whatever | default "sync" | quote }}
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

// changeFileFormat change the file extension from toml to yaml
func changeFileFormat(f string) string {
	// change - to _ in the filename
	f = strings.ReplaceAll(f, "_toml", ".yaml")

	return f
}

// createFile creates the file in an specific path
func createFile(f string) {
	// add the file extension
	f = changeFileFormat(f)

	_, err := os.Create(f)
	if err != nil {
		fmt.Println("ERROR creating the file: ", f, " ->", err)
	}
}

// writeToFile self description
func writeToFile(f string, content []byte) {
	// Change the file format to .yaml
	f = changeFileFormat(f)

	// create or verifyt that the file exists
	createFile(f)

	// Open the file
	file, err := os.OpenFile(f, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("ERROR creating the file: ", f, " ->", err)
	}

	_, err2 := file.Write(content)
	if err2 != nil {
		fmt.Println("ERROR creating the file: ", f, " ->", err2)
	}

	file.Close()


// main everything starts here!
func main() {
	// Run -> entrypoint
	Run()
}
