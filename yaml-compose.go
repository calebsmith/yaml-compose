package main

import (
	"log"

	"gopkg.in/yaml.v2"
)

func main() {
	// TODO: Take filename as CLI argument
	proc_res, err := pre_process_file("examples/main.yaml")
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Remove debug logging of vars
	log.Printf("Vars:\n")
	for k, v := range *proc_res.vars {
		log.Printf("%s=%+v\n", k, v)
	}
	// TODO: Move after YAML verification
	log.Printf(proc_res.out)
	// Parse as yaml
	yaml_in := yaml.MapSlice{}
	if err := yaml.Unmarshal([]byte(proc_res.out), &yaml_in); err != nil {
		log.Fatal(err)
	}
}
