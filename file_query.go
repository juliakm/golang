package main

import (
    "fmt"
    "log"
    "os"
	"path/filepath"
	"strings"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"encoding/csv"
)

// YamlConfig is the YAML values I am searching for
type YamlConfig struct {
	Manager string  "ms.manager"
	Author string  "ms.author"
}

// Go through folders and look for includes
func visit(files *[]string) filepath.WalkFunc {


    return func(path string, info os.FileInfo, err error) error {

        if err != nil {
            log.Fatal(err)
		}
		 
		if strings.Contains(path, "includes") {
			*files = append(*files, path)
		}
		return nil	
    } 
}

// Parse YAML function to read in YAML files
func parseYAML(fileName string) (string, string) { //define multiple return vals
  //fmt.Println("Parsing YAML file")

  yamlFile, err := ioutil.ReadFile(fileName)

  var yamlConfig YamlConfig
  err = yaml.Unmarshal(yamlFile, &yamlConfig)
  if err != nil {
	  fmt.Printf("Error parsing YAML file: %s\n", err)
  }

  return yamlConfig.Manager, yamlConfig.Author
}

func main() {

	//initalize files
    var files []string

	root := "/Users/juliakm/Code/azure-devops-docs-pr/docs/pipelines"
    err := filepath.Walk(root, visit(&files))
    if err != nil {
        panic(err)
    }
	
	// Create CSV file
	csvfile, err := os.Create("test.csv")
 
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	
	// Establish Writer
	w := csv.NewWriter(csvfile)

	// Writer headers
	if err := w.Write([]string{"File", "Manager","Author"}); err != nil {
        log.Fatalln("error writing record to csv:", err)
    }

    // Iterate through files, Parse YAML 
	for _, file := range files {
	  fmt.Println(file)

	  manager, author := parseYAML(file)

	  // add file row
	  w.Write([]string{file, manager, author})

	}
	w.Flush()
}