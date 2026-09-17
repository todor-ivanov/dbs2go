package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var ddlPath, dialect, emit, outBase string
	flag.StringVar(&ddlPath, "ddl", "", "DDL file to parse")
	flag.StringVar(&dialect, "dialect", "oracle", "SQL dialect: oracle or mysql")
	flag.StringVar(&emit, "emit", "json,markdown", "comma-separated outputs: json,markdown")
	flag.StringVar(&outBase, "out", "", "output basename")
	flag.Parse()

	if ddlPath == "" {
		fatalf("-ddl is required")
	}
	formats, err := parseEmitters(emit)
	if err != nil {
		fatalf("%v", err)
	}
	input, err := os.ReadFile(ddlPath)
	if err != nil {
		fatalf("read %s: %v", ddlPath, err)
	}
	schema, err := parseSchema(filepath.Base(ddlPath), input, dialect)
	if err != nil {
		fatalf("%v", err)
	}
	if outBase == "" {
		outBase = ddlPath + ".schema"
	}
	outputs := make(map[string][]byte, len(formats))
	for _, format := range formats {
		switch format {
		case "json":
			outputs[outBase+".json"], err = renderJSON(schema)
		case "markdown":
			outputs[outBase+".md"] = []byte(renderMarkdown(schema))
		}
		if err != nil {
			fatalf("render %s: %v", format, err)
		}
	}
	if err := writeOutputsAtomically(outputs); err != nil {
		fatalf("write output: %v", err)
	}
	fmt.Printf("tables=%d columns=%d foreign_keys=%d auxiliary=%d\n", len(schema.Tables), columnCount(schema), len(schema.ForeignKeys), len(schema.Auxiliary))
}

func parseEmitters(value string) ([]string, error) {
	seen := map[string]bool{}
	var out []string
	for _, item := range strings.Split(value, ",") {
		item = strings.ToLower(strings.TrimSpace(item))
		if item == "mermaid" {
			item = "markdown"
		}
		if item != "json" && item != "markdown" {
			return nil, fmt.Errorf("unknown emitter %q (supported: json, markdown)", item)
		}
		if !seen[item] {
			seen[item] = true
			out = append(out, item)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("at least one emitter is required")
	}
	return out, nil
}

func columnCount(schema *Schema) int {
	n := 0
	for _, table := range schema.Tables {
		n += len(table.Columns)
	}
	return n
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "dbs-schema: "+format+"\n", args...)
	os.Exit(2)
}
