package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	filePath := flag.String("file", "./test.toml", "Toml file to convert to Kubernetes ConfigMap")
	flag.Parse()

	//filePath := os.Args[1]
	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	for _, eachline := range txtlines {
		//varName = d
		eachline = "    " + eachline
		if strings.Contains(eachline, "=") {
			varName := eachline[:strings.IndexByte(eachline, '=')]
			//fmt.Println(varName)

			varName = strings.ReplaceAll(varName, " ", "")
			eachline = strings.ReplaceAll(eachline, "= ", " = {{ .Values.configMaps.FILENAME."+varName+" | "+"default "+"\"")
			eachline = eachline + "\"" + " | quote }}"

			// broadcast-mode = {{ .Values.configMaps.client_toml.broadcast-mode | default "sync" | quote }}
		}
		fmt.Println(eachline)
	}

}
