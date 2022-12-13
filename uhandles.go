package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
            
	"github.com/aymerick/raymond"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: uhandlebars <json file>")
		os.Exit(1)
	}

	jsonFile := os.Args[1]

	jsonData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		os.Exit(1)
	}

	var data interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Error parsing JSON file:", err)
		os.Exit(1)
	}

	template, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println("Error reading template from stdin:", err)
		os.Exit(1)
	}

	result, err := raymond.Render(string(template), data)
	if err != nil {
		fmt.Println("Error rendering template:", err)
		os.Exit(1)
	}

	fmt.Println(strings.TrimSpace(result))
}
