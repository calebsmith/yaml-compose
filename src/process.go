package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Tracks results of processing a given line or file
type ProcessResult struct {
	out  string
	vars *map[string]interface{}
}

// Handles a "load" directive of a YAML template file during pre-processing.
// Any {# filename #} strings not within a comment should load that file as
// data for template rendering
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
			// Copy any variable data passed down from injection directives
			if res.vars != nil {
				for k, v := range *res.vars {
					vars[k] = v
				}
			}
			// Copy any variable data from pre-processing the given filename
			// Overwrites any key matches since the given file should have higher precedence
			if var_res.vars != nil {
				for k, v := range *var_res.vars {
					vars[k] = v
				}
			}
			// Load any variable data from the given YAML data
			// Overwrites any previous key matches since loading data here
			//  should have the highest precedence.
			if err := yaml.Unmarshal([]byte(var_res.out), &vars); err != nil {
				return ProcessResult{
					out:  "",
					vars: nil}, err
			}
			// YAML variables are loaded, ellide directive from original YAML
			return ProcessResult{
				out:  "",
				vars: &vars}, nil
		}
	}
	// No load directive, return as-is
	return res, nil
}

// Handles an "inject" directive of a YAML template file during pre-processing.
// Any {$ filename $} strings not within a comment should load that files contents
// and inject it into the current file, replacing the directive.
// Any string preceding the directive is maintained and the injected content should match
// the indentation of the injection site.
// E.g.:
// data:
//
//	key1: value1
//	{$ other.yaml $}
//
// and other.yaml containing:
// key2: value2
// ->
// data:
//
//	key1: value1
//	key2: value2
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
	// No injection directive, return line as-is
	return ProcessResult{
		out:  line_in,
		vars: nil}, nil
}

// Pre-process a single line of a YAML template file
func pre_process_line(line_in string) (ProcessResult, error) {
	// Handle any file injection directives first
	res, err := pre_process_line_inject(line_in)
	if err != nil {
		return res, err
	}
	// Handle any load directives
	return pre_process_line_load(res)
}

// Pre-process the given file. If performing an injection, use the string given
// as `prefix` to prepend the initial existing prefix and any indentation as needed
func pre_process_file_w_prefix(filename string, prefix string) (ProcessResult, error) {
	// Load given YAML template or YAML file
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

// pre-process a YAML template file given as `filename` by loading or
// injecting any additional YAML files per the given directives.
func pre_process_file(filename string) (ProcessResult, error) {
	// Pre-process given file with no prefix
	return pre_process_file_w_prefix(filename, "")
}

// Render YAML template with the given variable data
func process(res ProcessResult) (string, error) {
	if res.vars != nil {
		tmp, err := template.New("main").Parse(res.out)
		var buff bytes.Buffer
		tmp.Execute(&buff, res.vars)
		return buff.String(), err
	} else {
		return res.out, nil
	}
}
