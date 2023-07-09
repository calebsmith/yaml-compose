package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type ProcessResult struct {
	out  string
	vars *map[string]interface{}
}

func pre_process_line_load(res ProcessResult) (ProcessResult, error) {
	// Pattern matches any line with {# filename #}, with filename being assigned to a group.
	// Ignores YAML comments
	pattern := `^[^#/s]*\{#(?P<filename>.*?)#\}`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(res.out)
	if len(match) > 0 {
		groupIndex := re.SubexpIndex("filename")
		if groupIndex != -1 {
			filename := strings.TrimSpace(match[groupIndex])
			var_res, err := pre_process_file_w_prefix(filename, "")
			if err != nil {
				return ProcessResult{
					out:  "",
					vars: nil}, err
			}
			vars := make(map[string]interface{})
			if var_res.vars != nil {
				for k, v := range *var_res.vars {
					vars[k] = v
				}
			}
			if res.vars != nil {
				for k, v := range *res.vars {
					vars[k] = v
				}
			}
			if err := yaml.Unmarshal([]byte(var_res.out), &vars); err != nil {
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
	return res, nil
}

func pre_process_line_inject(line_in string) (ProcessResult, error) {
	// Pattern matches any line with {$ filename $}, with filename being assigned to a group.
	// Ignores YAML comments
	// Replaces pattern with contents of YAML file, maintaining any preceding characters.
	// Also prepends any leading whitespace to injected text to match indentation level of injection point
	pattern := `^(?P<prefix>[^#/s]*?)\{\$(?P<filename>.*?)\$\}`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(line_in)
	if len(match) > 0 {
		filenameGroupIndex := re.SubexpIndex("filename")
		prefixGroupIndex := re.SubexpIndex("prefix")
		if filenameGroupIndex != -1 && prefixGroupIndex != -1 {
			filename := strings.Trim(match[filenameGroupIndex], " ")
			prefix := match[prefixGroupIndex]
			return pre_process_file_w_prefix(filename, prefix)
		}
	}
	// Normal line, return as-is
	return ProcessResult{
		out:  line_in,
		vars: nil}, nil
}

func pre_process_line(line_in string) (ProcessResult, error) {
	res, err := pre_process_line_inject(line_in)
	if err != nil {
		return res, err
	}
	return pre_process_line_load(res)
}

func pre_process_file_w_prefix(filename string, prefix string) (ProcessResult, error) {
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
	content_index := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw_line := scanner.Text()
		var line string
		if content_index == 0 {
			// Maintain original line prefix if it exists for injection point
			line = fmt.Sprintf("%s%s", prefix, raw_line)
		} else {
			// Remaining text is prepended with leading whitespace to match injection point
			line = fmt.Sprintf("%s%s", strings.Repeat(" ", len(prefix)), raw_line)
		}
		proc_line_res, err := pre_process_line(line)
		if err != nil {
			return proc_line_res, err
		}
		if proc_line_res.vars != nil {
			for k, v := range *proc_line_res.vars {
				vars[k] = v
			}
		}
		// Only append and increment if content encountered
		// N.B. - This allows for correct handling of directives within
		// injected content
		if strings.TrimSpace(proc_line_res.out) != "" {
			content_index++
			lines = append(lines, proc_line_res.out)
		}
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

func pre_process_file(filename string) (ProcessResult, error) {
	// Pre-process given file with no prefix
	return pre_process_file_w_prefix(filename, "")
}
