package main

import (
	"fmt"
	"os"
	"strings"
)

var (
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

// createFullPath create the full path to the outputs
func createFullPath(f string) error {
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		os.MkdirAll(outputPath, 0644)
	}
	return nil
}

// createFile creates the file in an specific path
func createFile(f string) {
	// create the outputs folder before write the file
	err1 := createFullPath(f)
	if err1 != nil {
		fmt.Println("ERROR creating the file: ", f, " ->", err1)
	}

	// add the file extension
	f = ChangeFileExtension(f)

	file, err := os.Create(outputPath + f)
	if err != nil {
		fmt.Println("ERROR creating the file: ", f, " ->", err)
	}

	defer file.Close()
}

// writeToFile self description
func writeToFile(f string, content []string) {
	// Change the file format to .yaml
	f = ChangeFileExtension(f)

	// if file exsits, remove it first
	if _, err := os.Stat(outputPath + f); err == nil {
		os.Remove(outputPath + f)
	}

	// create or verifyt that the file exists
	createFile(f)

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
