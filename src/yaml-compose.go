package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("Command takes a single filename as only argument")
	}
	filename := args[0]
	// Pre-process injection and loading directives
	pre_res, err := pre_process_file(filename)
	if err != nil {
		log.Printf("Failed to pre-process YAML:\n%s\n", pre_res.out)
		log.Printf("Loaded variables were:\n%+v\n", pre_res.vars)
		log.Fatal(err)
	}
	// Render template variables
	res, err := process(pre_res)
	if err != nil {
		log.Printf("Failed to process YAML template:\n%s\n", res)
		log.Printf("Loaded variables were:\n%+v\n", pre_res.vars)
		log.Fatal(err)
	}
	// Verify resulting content as YAML
	yaml_in := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(res), &yaml_in); err != nil {
		log.Printf("Incorrect YAML generated: %s\n", res)
		log.Fatal(err)
	}
	// Output valid YAML results
	fmt.Println(res)
}
