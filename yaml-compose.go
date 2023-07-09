package main

import (
	"log"

	"gopkg.in/yaml.v2"
)

func main() {
	// TODO: Take filename as CLI argument
	pre_res, err := pre_process_file("examples/main.yaml")
	if err != nil {
		log.Printf("Failed to pre-process YAML:\n%s\n", pre_res.out)
		log.Printf("Loaded variables were:\n%+v\n", pre_res.vars)
		log.Fatal(err)
	}
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
	// Output results for use
	log.Print(res)
}
