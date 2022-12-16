package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"flag"
	"path"
	"path/filepath"

	"github.com/aymerick/raymond"
)

func usage(dest *os.File) {
	fmt.Fprintln(dest, "Usage:")
	fmt.Fprintln(dest, "  [...] | " + os.Args[0] + " <json-data-file> | [...]")
	fmt.Fprintln(dest, "  " + os.Args[0] + " [-d|--data <json-data-file|->] [-t|--template <tmpl-file|tmpl-dir|->] [-o|--output <output-file|output-dir|->]")
	fmt.Fprintln(dest, `
  -d|--data 		JSON data file containing substitution values (file path or "-" for stdin)
  -t|--template 	Handlebar template file to render, or directory of template files (default: "-", stdin)
  -o|--output 		File to concat rendered templates to, or directory to place rendered templates within (default: "-", stdout)
  --tmpl-token 		Token to use to detect template files in input-directory mode (default: ".tmpl")

  (NOTE: You cannot specify "-"/stdin as a source for both template and data.)
`)
}

func main() {
	// handle flags:
	var jsonPath string
	var inputPath string
	var outputPath string
	var tmplToken string

	flag.StringVar(&jsonPath, "d", "", "")
	flag.StringVar(&jsonPath, "data", "", "")

	flag.StringVar(&inputPath, "t", "", "")
	flag.StringVar(&inputPath, "template", "", "")

	flag.StringVar(&outputPath, "o", "", "")
	flag.StringVar(&outputPath, "output", "", "")

	flag.StringVar(&tmplToken, "tmpl-token", "", "")

	flag.Usage = func() {
		usage(os.Stdout)
		os.Exit(0)
	}
	flag.Parse()

	// pluck defaults out of env variables if they're set but not specified from flags:
	// os.Getenv defaults to blank strings if the value isn't present, conveniently
	if jsonPath == "" {
		jsonPath = os.Getenv("UHANDLES_DATA")
	}

	if inputPath == "" {
		inputPath = os.Getenv("UHANDLES_TEMPLATE")
	}

	if outputPath == "" {
		outputPath = os.Getenv("UHANDLES_OUTPUT")
	}

	if tmplToken == "" {
		tmplToken = os.Getenv("UHANDLES_TMPL_TOKEN")
	}

	// if we're still unset, set to our default...
	if tmplToken == "" {
		tmplToken = ".tmpl"
	}

	// if we still have no paths from env and no other arguments... show usage politely:
	if len(os.Args) == 1 && jsonPath == "" && inputPath == "" && outputPath == "" {
		usage(os.Stdout)
		os.Exit(0)
	}

	// pull JSON data from argv[1] if not specified with a flag above:
	if jsonPath == "" && jsonPath != "-" {
		if len(os.Args) != 2 {
			fmt.Println("No data file specified with -d|--data and no data file found in first argument.")
			usage(os.Stderr)
			os.Exit(1)
		}

		jsonPath = os.Args[1]
	}

	// read JSON data from desired source:
	var jsonData []byte

	if jsonPath != "-" {
		// read from file
		readData, err := ioutil.ReadFile(jsonPath)
		if err != nil {
			fmt.Println("Error reading JSON file:", err)
			usage(os.Stderr)
			os.Exit(1)
		}
		jsonData = readData
	} else {
		// read from stdin
		readData, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Println("Error reading JSON from stdin:", err)
			usage(os.Stderr)
			os.Exit(1)
		}
		jsonData = readData
	}

	// parse read JSON data
	var data interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Error parsing JSON file:", err)
		os.Exit(1)
	}

	// pull out real input and output paths depending on if the path supplied is a directory or not:
	var inputPaths []string
	var outputPaths []string
	var outputPathIsDir bool

	if inputPath == "" || inputPath == "-" {
		// input is either implicitly or explicitly stdin
		inputPaths = []string{"-"}
	} else {
		// sample input path to see if it's a directory
		inputPathInfo, err := os.Stat(inputPath)
		if err != nil {
			fmt.Println("Cannot stat input path:", inputPath, err)
			usage(os.Stderr)
			os.Exit(1)
		}
	
		if inputPathInfo.IsDir() {
			// populate input paths by globbing for template files in the input dir
			inputPaths, err = filepath.Glob(path.Clean(inputPath) + "/*" + tmplToken + "*")
			if err != nil {
				fmt.Println("Cannot glob input path:", inputPath, err)
				usage(os.Stderr)
				os.Exit(1)
			}
		} else {
			// we're just a file
			inputPaths = []string{path.Clean(inputPath)}
		}
	}

	if outputPath == "" || outputPath == "-" {
		// output is either implicitly or explicitly stdin
		// ... just alter the output path to stdin simplify it for us later
		outputPath = "-"
	} else {
		// sample output path to see if it's a directory
		outputPathInfo, err := os.Stat(outputPath)
		if err != nil {
			fmt.Println("Cannot stat output path:", outputPath, err)
			usage(os.Stderr)
			os.Exit(1)
		}

		// retain for later append-vs-write logic
		outputPathIsDir = outputPathInfo.IsDir()

		if outputPathIsDir {
			// if input is a single file or a directory of files, strip the .tmpl from their names and
			// place them within the destination directory
			for i := 0; i < len(inputPaths); i++ {
				baseName := strings.Replace(path.Base(inputPaths[i]), tmplToken, "", -1)
				newName := path.Join(outputPath, baseName)
				outputPaths = append(outputPaths, newName)
			}
		}
	}

	// iterate through all the inputs as directed, render them with the data, and output them as prescribed
	var outputHandle *os.File // needed for appending multiple inputs to a single file
	for idx, inPath := range inputPaths {
		// read template...
		var template []byte

		if inPath != "-" {
			// read template file
			reader, err := os.Open(inPath)
			if err != nil {
				fmt.Println("Error opening template file:", inPath, err)
				usage(os.Stderr)
				os.Exit(1)
			}

			template, err = ioutil.ReadAll(reader)
			if err != nil {
				fmt.Println("Error reading template file:", inPath, err)
				usage(os.Stderr)
				os.Exit(1)
			}
			reader.Close()
		} else {
			// read template from stdin
			if jsonPath == "-" {
				fmt.Println("Can't read template from stdin, already reading JSON from stdin.")
				usage(os.Stderr)
				os.Exit(1)
			}

			readTemplate, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				fmt.Println("Error reading template from stdin:", err)
				usage(os.Stderr)
				os.Exit(1)
			}

			template = readTemplate
		}

		// render template from input with data
		result, err := raymond.Render(string(template), data)
		if err != nil {
			fmt.Println("Error rendering template:", err)
			usage(os.Stderr)
			os.Exit(1)
		}
		
		// output the rendering
		if outputPath == "" || outputPath == "-" {
			// outputing to stdout, just dump it out. it'll automatically "append" if there are more inputs
			fmt.Println(strings.TrimSpace(result))
		} else if outputPathIsDir {
			// snag corresponding input file permissions
			inStat, _ := os.Stat(inPath)

			// outputting to separate file(s) in a dir with the input file permissions
			ioutil.WriteFile(outputPaths[idx], []byte(strings.TrimSpace(result)), inStat.Mode())
		} else {
			// outputting to a single file, appending renders
			if idx == 0 {
				// snag (first) input file permissions
				inStat, _ := os.Stat(inPath)
				
				// if this is the first rendering, create the file with the input permissions
				ioutil.WriteFile(outputPath, []byte(strings.TrimSpace(result)), inStat.Mode())

				// and open it for further append writing...
				outputHandle, _ = os.OpenFile(outputPath, os.O_APPEND|os.O_WRONLY, inStat.Mode())
			} else {
				// if this isn't the first rendering, append to the file
				outputHandle.WriteString("\n" + strings.TrimSpace(result))
			}
		}
	}

	// close single-output appending output handle, if it was used
	if outputHandle != nil {
		outputHandle.Close()
	}
}
