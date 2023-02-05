package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	// spacesToTab yaml format spaces in a ConfigMap
	spacesToTab = "    "
	// outputPath outputs folder, where all the files will be stored
	outputPath = "./outputs/"
	// configMapKind this is the header of the file, based on a Kubernetes
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
	file := flag.String("file", "", "Toml file to convert to Kubernetes ConfigMap")

	flag.Parse()

	fi, err := os.Stat(*file)
	if err != nil {
		log.Fatalf("failed check the input file: %s", err)
	}

	// filesList -> in case of directory, we'll use this array to store the
	// files on it.
	var filesList []string

	// Check if the input is a directory or a file
	switch mode := fi.Mode(); {
	case mode.IsDir():
		filesList = FindFilesInPath(*file)
	case mode.IsRegular():
		ParseFiles(*file)
	}

	// if there's at least 1, we'll check the files
	if len(filesList) > 0 {
		for _, i := range filesList {
			ParseFiles(i)
		}
	}
}

// FindFilesInPath it returns a list with all the files with .toml extension
func FindFilesInPath(f string) []string {
	// counter - how many files do we have
	c := 0
	var filesToProcess []string

	files, err := ioutil.ReadDir(f)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	// find files in folder with Toml extension
	for _, file := range files {
		// Get the file extension
		ext := filepath.Ext(file.Name())
		// Check if the extension is .toml or .tml
		if ext == ".toml" {
			fmt.Println("File found to process: ", file.Name())
			filesToProcess = append(filesToProcess, file.Name())
			c = c + 1
		}
	}
	fmt.Println("---------------------------")
	fmt.Println("Total files to process: ", c)
	fmt.Println("---------------------------")

	return filesToProcess
}

// ParseFiles process the files and create the ConfigMaps
func ParseFiles(file string) {
	// Open file
	f, err := os.Open(file)
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
	fileNameCleaned := cleanUpFileName(file)

	var fileContent []string

	// Add the Kubernetes ConfigMap structure at top
	configMapKind = strings.ReplaceAll(configMapKind, "CONFIGMAP_NAME", fileNameCleaned)
	fileContent = append(fileContent, configMapKind)
	// Add the data
	dataCM := "  " + file + ": |" + "\n"
	fileContent = append(fileContent, dataCM)

	// Check line by line
	for _, eachline := range txtlines {
		// add the spaces to the beginning of the line - yaml format
		eachline = spacesToTab + eachline

		// if line contains = that means it contains a variable + value
		if strings.Contains(eachline, "=") {
			eachline = FormatLine(eachline, fileNameCleaned)
			eachline = strings.ReplaceAll(eachline, "-", "_")
		}

		eachline = eachline + "\n"
		// Here append the lines to the array in order to bulk the data
		fileContent = append(fileContent, eachline)

		//fmt.Println(eachline)
	}

	// here is where we have to write the content to the new file, we've
	// formated the lines that contains values
	WriteToFile(fileNameCleaned, fileContent)
}

// FormatLine change the current format to Helm Template
func FormatLine(line, fileName string) string {
	// find the = and cut the string from it
	l := line[:strings.IndexByte(line, '=')]

	// Clean the spaces
	l = strings.ReplaceAll(l, " ", "")

	// Replace the value of the variable to a Helm Template format
	// the following lines add the content between quotes, we don't really need
	// it in that way. -> pending to review
	line = strings.ReplaceAll(
		line,
		"= ",
		" = {{ .Values.configMaps."+fileName+"."+l+" | "+"default ")
	line = line + " | quote }}"

	// if contains empty array, replace it
	if strings.Contains(line, "[]") {
		line = strings.ReplaceAll(line, "[]", "\"[]\"")
	}

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

// ChangeFileExtension change the file extension from toml to yaml
func ChangeFileExtension(f string) string {
	// change - to _ in the filename
	f = strings.ReplaceAll(f, "_toml", ".yaml")

	return f
}

// CreateFullPath create the full path to the outputs
func CreateFullPath(f string) error {
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		os.MkdirAll(outputPath, 0744)
	}
	return nil
}

// CreateFile creates the file in an specific path
func CreateFile(f string) {
	// create the outputs folder before write the file
	err1 := CreateFullPath(f)
	if err1 != nil {
		fmt.Println("ERROR creating the file (CreateFullPath): ", f, " ->", err1)
	}

	// add the file extension
	f = ChangeFileExtension(f)

	file, err := os.Create(outputPath + f)
	if err != nil {
		fmt.Println("ERROR creating the file: ", f, " ->", err)
	}

	defer file.Close()
}

// WriteToFile self description
func WriteToFile(f string, content []string) {
	// Change the file format to .yaml
	f = ChangeFileExtension(f)

	// if file exsits, remove it first
	if _, err := os.Stat(outputPath + f); err == nil {
		os.Remove(outputPath + f)
	}

	// create or verifyt that the file exists
	CreateFile(f)

	// Open the file
	file, err := os.OpenFile(outputPath+f, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("ERROR creating the file: ", outputPath+f, " ->", err)
	}

	defer file.Close()

	// write the content to the file
	//size, err2 := file.Write(content)
	//if err2 != nil {
	//	fmt.Println("ERROR creating the file: ", outputPath+f, " ->", err2)
	//}

	//fmt.Printf("Wrote bytes %d to file", size)

	for _, val := range content {
		_, err := file.WriteString(val)
		if err != nil {
			fmt.Println("ERROR creating the file: ", outputPath+f, " ->", err)
		}
	}
}

// main everything starts here!
func main() {
	// Run -> entrypoint
	Run()
}
