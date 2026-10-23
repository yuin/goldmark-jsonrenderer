// Package main implements a command-line tool that converts Markdown text to JSON representation of
// its AST (Abstract Syntax Tree).
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	alert "github.com/yuin/goldmark-alert"
	jsonrenderer "github.com/yuin/goldmark-jsonrenderer"
	meta "github.com/yuin/goldmark-meta/v2"
	"github.com/yuin/goldmark/v2/extension"
	"github.com/yuin/goldmark/v2/parser"
)

var usageLines = []string{}

func addOption[T any](opt *T, longName string, shortName string, defaultValue T, varName string, usage string) {
	switch v := any(opt).(type) {
	case *bool:
		flag.BoolVar(v, longName, any(defaultValue).(bool), usage)
		if shortName != "" {
			flag.BoolVar(v, shortName, any(defaultValue).(bool), usage)
			usageLines = append(usageLines, fmt.Sprintf(
				"  -%s, --%s: %s (default: %v)", shortName, longName, usage, defaultValue))
		} else {
			usageLines = append(usageLines, fmt.Sprintf("  --%s: %s (default: %v)", longName, usage, defaultValue))
		}
	case *int:
		flag.IntVar(v, longName, any(defaultValue).(int), usage)
		if shortName != "" {
			flag.IntVar(v, shortName, any(defaultValue).(int), usage)
			usageLines = append(usageLines, fmt.Sprintf("  -%s=%s, --%s=%s: %s (default: %v)",
				shortName, varName, longName, varName, usage, defaultValue))
		} else {
			usageLines = append(usageLines, fmt.Sprintf("  --%s=%s: %s (default: %v)",
				longName, varName, usage, defaultValue))
		}
	case *string:
		flag.StringVar(v, longName, any(defaultValue).(string), usage)
		if shortName != "" {
			flag.StringVar(v, shortName, any(defaultValue).(string), usage)
			usageLines = append(usageLines, fmt.Sprintf("  -%s=%s, --%s=%s: %s (default: %v)",
				shortName, varName, longName, varName, usage, defaultValue))
		} else {
			usageLines = append(usageLines, fmt.Sprintf("  --%s=%s: %s (default: %v)",
				longName, varName, usage, defaultValue))
		}
	default:
		panic("unsupported option type")
	}
}

func main() {
	usageLines = append(usageLines, "Usage: markdown2json [options] <input-file>")
	usageLines = append(usageLines, "Description:")
	usageLines = append(usageLines, `
  markdown2json is a command-line tool that converts Markdown text to JSON representation of 
  its AST (Abstract Syntax Tree). if <input-file> is not specified, the tool will read from 
  standard input.
`)
	usageLines = append(usageLines, "Options:")
	var optPrettyPrint bool
	var optExtensions string

	addOption(&optPrettyPrint, "pretty-print", "p", false, "", "Pretty print the output JSON")
	addOption(&optExtensions, "extensions", "e", "all", "EXTENSIONS",
		"Comma-separated list of extensions to enable. if 'all' is specified, all extensions will be enabled")

	usageLines = append(usageLines, "")
	usageLines = append(usageLines, "Available extensions:")
	availableExtensions := []string{
		"gfm(linkify, strikethrough, table, tasklist)",
		"linkify",
		"strikethrough",
		"table",
		"tasklist",
		"footnote",
		"definitionList",
		"meta",
		"alert",
	}
	for _, ext := range availableExtensions {
		usageLines = append(usageLines, fmt.Sprintf("  %s", ext))
	}

	flag.Usage = func() {
		for _, line := range usageLines {
			fmt.Println(line)
		}
	}
	flag.Parse()
	args := flag.Args()
	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Too many arguments.")
		flag.Usage()
		os.Exit(1)
	}

	var source []byte
	if len(args) == 1 {
		filename := args[0]
		b, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read file: %v\n", err)
			os.Exit(1)
		}
		source = b
	} else {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read from stdin: %v\n", err)
			os.Exit(1)
		}
		source = b
	}

	var pexts []parser.Extension
	var rexts []jsonrenderer.Extension
	if optExtensions == "all" {
		pexts = []parser.Extension{
			extension.GFMParser,
			extension.FootnoteParser,
			extension.DefinitionListParser,
			meta.Parser,
			alert.Parser,
		}
		rexts = []jsonrenderer.Extension{
			jsonrenderer.Footnote,
			jsonrenderer.Strikethrough,
			jsonrenderer.Table,
			jsonrenderer.DefinitionList,
			jsonrenderer.Alert,
		}
	} else {
		for ext := range strings.SplitSeq(optExtensions, ",") {
			switch ext {
			case "gfm":
				pexts = append(pexts, extension.GFMParser)
				rexts = append(rexts, jsonrenderer.Strikethrough, jsonrenderer.Table)
			case "linkify":
				pexts = append(pexts, extension.LinkifyParser)
			case "strikethrough":
				pexts = append(pexts, extension.StrikethroughParser)
				rexts = append(rexts, jsonrenderer.Strikethrough)
			case "table":
				pexts = append(pexts, extension.TableParser)
				rexts = append(rexts, jsonrenderer.Table)
			case "tasklist":
				pexts = append(pexts, extension.TaskListItemParser)
			case "footnote":
				pexts = append(pexts, extension.FootnoteParser)
				rexts = append(rexts, jsonrenderer.Footnote)
			case "definitionList":
				pexts = append(pexts, extension.DefinitionListParser)
				rexts = append(rexts, jsonrenderer.DefinitionList)
			case "meta":
				pexts = append(pexts, meta.Parser)
			case "alert":
				pexts = append(pexts, alert.Parser)
				rexts = append(rexts, jsonrenderer.Alert)
			default:
				fmt.Printf("Unknown extension: %s\n", ext)
				return
			}
		}
	}

	p := parser.New(parser.WithExtensions(pexts...))
	r := jsonrenderer.New(jsonrenderer.WithExtensions(rexts...))
	doc := p.Parse(source)
	var w jsonrenderer.JSONWriter
	err := r.Render(&w, source, doc)
	if err != nil {
		panic(err)
	}
	if optPrettyPrint {
		jv := w.ToValue()
		err = jv.Indent()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jv))
	} else {
		fmt.Println(string(w.ToValue()))
	}

}
