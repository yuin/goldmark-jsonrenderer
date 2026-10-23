goldmark-jsonrenderer
=========================

[![GoDev][godev-image]][godev-url]

[godev-image]: https://pkg.go.dev/badge/github.com/yuin/goldmark-jsonrenderer.svg
[godev-url]: https://pkg.go.dev/github.com/yuin/goldmark-jsonrenderer

goldmark-jsonrenderer is a goldmark extension that renders Markdown documents as JSON.
goldmark-jsonrenderer also provides a `markdown2json` command that converts Markdown documents to JSON and prints them to stdout.

## Compatibility

`github.com/yuin/goldmark-jsonrenderer` is compatible with `goldmark/v2`.

minimum Go version: 1.27

## Usage

```go
func main() {
    p := parser.New(parser.WithExtensions(extension.GFMParser))
    r := jsonrenderer.New(jsonrenderer.WithExtensions(jsonrenderer.Strikethrough, jsonrenderer.Table))
    doc := p.Parse(source)

    var w jsonrenderer.JSONWriter
    err := r.Render(&w, source, doc)
    if err != nil {
        panic(err)
    }
    println(string(w.ToValue()))
}
```

## CLI
### Installation

go get

```bash
go install github.com/yuin/goldmark-jsonrenderer/cmd/markdown2json
```

installer(Linux/MacOS)

```bash
curl -fsSL https://raw.githubusercontent.com/yuin/goldmark-jsonrenderer/main/install.sh | bash -s -- -b "$HOME/.local/bin"
```

### Usage

```text
Usage: markdown2json [options] <input-file>
Description:

  markdown2json is a command-line tool that converts Markdown text to JSON representation of its AST (Abstract Syntax Tree).
  if <input-file> is not specified, the tool will read from standard input.

  JSON representation of the AST is like this:

  {
    "kind": "document",
    "metadata": {},
    "pos": 0,
    "attributes": {},
    "children": [
      {
        "kind": "heading",
        "level": 1,
        "pos": 0,
        "attributes": {},
        "children": [
          {
            "kind": "text",
            "pos": 2,
            "attributes": {},
            "value": "Heading1"
          }
        ]
      }
    ]
  }

Options:
  -p, --pretty-print: Pretty print the output JSON (default: false)
  -e=EXTENSIONS, --extensions=EXTENSIONS: Comma-separated list of extensions to enable. if 'all' is specified, all extensions will be enabled (default: all)

Available extensions:
  gfm(linkify, strikethrough, table, tasklist)
  linkify
  strikethrough
  table
  tasklist
  footnote
  definitionList
  meta
  alert
```

~~~bash
$ markdown2json -p << EOS
# Heading 1
- aaaa
- bbbb

```go
func main() {}
```

[link](http://www.example.com)
EOS
~~~

```json
{
    "attributes": {},
    "children": [
        {
            "attributes": {},
            "children": [
                {
                    "attributes": {},
                    "kind": "text",
                    "pos": 2,
                    "value": "Heading 1"
                }
            ],
            "kind": "heading",
            "level": 1,
            "pos": 0
        },
        {
            "attributes": {},
            "children": [
                {
                    "attributes": {},
                    "children": [
                        {
                            "attributes": {},
                            "children": [
                                {
                                    "attributes": {},
                                    "kind": "text",
                                    "pos": 14,
                                    "value": "aaaa"
                                }
                            ],
                            "kind": "paragraph",
                            "pos": 14
                        }
                    ],
                    "kind": "listItem",
                    "pos": 12
                },
                {
                    "attributes": {},
                    "children": [
                        {
                            "attributes": {},
                            "children": [
                                {
                                    "attributes": {},
                                    "kind": "text",
                                    "pos": 21,
                                    "value": "bbbb"
                                }
                            ],
                            "kind": "paragraph",
                            "pos": 21
                        }
                    ],
                    "kind": "listItem",
                    "pos": 19
                }
            ],
            "isTight": true,
            "kind": "list",
            "marker": "-",
            "pos": 12,
            "start": 0
        },
        {
            "attributes": {},
            "codeBlockKind": "Fenced",
            "info": "go",
            "kind": "codeBlock",
            "language": "go",
            "pos": 27,
            "value": "func main() {}\n"
        },
        {
            "attributes": {},
            "children": [
                {
                    "attributes": {},
                    "children": [
                        {
                            "attributes": {},
                            "kind": "text",
                            "pos": 54,
                            "value": "link"
                        }
                    ],
                    "destination": "http://www.example.com",
                    "kind": "link",
                    "pos": 53,
                    "title": ""
                }
            ],
            "kind": "paragraph",
            "pos": 53
        }
    ],
    "kind": "document",
    "metadata": {},
    "pos": 0
}
```

`markdown2json` follows the philosophy of Unix tools; **Write simple parts connected by clean interfaces** .

Output is just a JSON, so you can use it in your applications or scripts to process Markdown documents programmatically.
Especially, you can use `jq` to extract specific information from the JSON output. LLMs are an expert in using `jq` to extract information from JSON, so you can easily extract information from Markdown documents using `markdown2json` and `jq`.

**extract all links from a Markdown document:**

```bash
$ markdown2json | jq '.. | objects | select(.kind == "link") | .destination'
```

**extract only h1 headings from a Markdown document:**

```bash
$ markdown2json | jq '.. | objects | select(.kind == "heading" and .level == 1)'
```

**extract headings from level 1 to 3 using a range:**

```bash
$ markdown2json | jq '.. | objects | select(.kind == "heading" and .level >= 1 and .level <= 3)'
```

**extract only rust code blocks from a Markdown document:**

```bash
$ markdown2json | jq '.. | objects | select(.kind == "codeBlock" and .language == "rust")'
```
