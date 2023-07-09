package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type ProcessResult struct {
	out  string
	vars *map[string]interface{}
}

func pre_process_line(line_in string) (ProcessResult, error) {
	// Pattern matches any line with {# text #}, with text being assigned to a group.
	// Ignores YAML comments
	pattern := `^[^#/s]*\{#(?P<filename>.*?)#\}`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(line_in)
	if len(match) > 0 {
		groupIndex := re.SubexpIndex("filename")
		if groupIndex != -1 {
			filename := strings.Trim(match[groupIndex], " ")
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				return ProcessResult{
					out:  "",
					vars: nil}, err
			}
			vars := make(map[string]interface{})
			if err := yaml.Unmarshal([]byte(content), &vars); err != nil {
				return ProcessResult{
					out:  "",
					vars: nil}, err
			}
			// Loaded YAML variables, ellide directive from original YAML
			return ProcessResult{
				out:  "",
				vars: &vars}, nil
		}
	}
	// Normal line, return as-is
	return ProcessResult{
		out:  line_in,
		vars: nil}, nil
}

func pre_process_file(filename string) (ProcessResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return ProcessResult{
			out:  "",
			vars: nil}, err
	}
	defer file.Close()

	// Pre-process linewise
	vars := make(map[string]interface{})
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw_line := scanner.Text()
		proc_line_res, err := pre_process_line(raw_line)
		if err != nil {
			return proc_line_res, err
		}
		if proc_line_res.vars != nil {
			for k, v := range *proc_line_res.vars {
				vars[k] = v
			}
		}
		lines = append(lines, proc_line_res.out)
	}
	if err := scanner.Err(); err != nil {
		return ProcessResult{
			out:  "",
			vars: nil}, err
	}
	return ProcessResult{
		out:  strings.Join(lines, "\n"),
		vars: &vars}, nil
}
