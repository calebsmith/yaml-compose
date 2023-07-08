package main

import (
	"gopkg.in/yaml.v2"
	"log"
	//"regexp"
	"bufio"
	"os"
	"strings"
)

func pre_process(raw_file string) (string, error) {
	// TBI: No-op for now
	return raw_file, nil
}

/*
func process(out interface{}) {
	switch i := out.(type) {
	case nil, bool, int, float32, float64:
		return
	case string:
		// TODO: Modify value as needed
		log.Printf("i is a string: %+v\n", i)
	// Recurse into maps and arrays as needed
	case yaml.MapItem:
		process(i.Value)
	case yaml.MapSlice:
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
*/

func main() {

	// Load file
	file, err := os.Open("example.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Pre-process linewise
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw_line := scanner.Text()
		processed_line, err := pre_process(raw_line)
		if err != nil {
			log.Fatal(err)
		}
		lines = append(lines, processed_line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	processed_out := strings.Join(lines, "\n")

	log.Printf("Pre-processed:\n")
	log.Printf(processed_out)

	// Parse as yaml
	yaml_in := yaml.MapSlice{}
	if err := yaml.Unmarshal([]byte(processed_out), &yaml_in); err != nil {
		log.Fatal(err)
	}

	// FIXME: Consider removal
	// Process file
	// process(yaml_in)

	// FIXME: Consider just using processed_out directly after verification
	// Output results
	out, err := yaml.Marshal(yaml_in)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(string(out))
}
