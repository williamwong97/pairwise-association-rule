package models

import (
	"fmt"
	"os"
	"strings"
)

type Arguments struct {
	DataFile   string
	InputItems []string
}

const usage = `Arguments:
  --dataFile file_path     Input dataset in CSV format.
  --inputItems string 	   Items input by respondent`

func ParseArgsFromCommand() Arguments {
	var (
		result = Arguments{}
		args   = os.Args[1:]
	)

	if len(args) == 0 {
		fmt.Println(usage)
		os.Exit(-1)
	}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--dataFile":
			{
				if i+1 > len(args) {
					fmt.Println("Expected --dataFile to be followed by input CSV path.")
					os.Exit(-1)
				}
				result.DataFile = args[i+1]
				i++
			}
		case "--inputItems":
			{
				if i+1 > len(args) {
					fmt.Println("Expected --inputItems to be followed by inputItems.")
					os.Exit(-1)
				}
				result.InputItems = strings.Split(args[i+1], "")
				i++
			}
		}
	}
	if len(result.DataFile) == 0 {
		fmt.Println("Missing required parameter '--dataFile $csv_path")
		os.Exit(-1)
	}
	if len(result.InputItems) == 0 {
		fmt.Println("Missing required parameter '--inputItems $input_items")
		os.Exit(-1)
	}
	return result
}
