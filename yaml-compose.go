package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

func process(out interface{}) {
	switch i := out.(type) {
	case string:
		// TODO: Modify value as needed
		log.Printf("i is a string: %+v\n", i)
	case nil, bool, int, float32, float64:
		log.Printf("i is non-string scalar: %+v\n", i)
		return
	// Recurse into maps and arrays as needed
	case map[string]interface{}:
		for _, v := range i {
			process(v)
		}
	case []interface{}:
		for _, v := range i {
			process(v)
		}
	default:
		// TODO: Ensure this is unreachable
		log.Fatalf("Unhandled type %T for value %+v", i, i)
	}
}

func main() {

	// Load and parse file
	content, err := ioutil.ReadFile("example.yaml")
	if err != nil {
		log.Fatal(err)
	}
	yaml_in := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(content), &yaml_in); err != nil {
		log.Fatal(err)
	}

	// Process file
	process(yaml_in)

	// Output results
	out, err := yaml.Marshal(yaml_in)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(string(out))
}
