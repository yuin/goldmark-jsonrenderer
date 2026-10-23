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
                    "value": "Heading 1"
                }
            ],
            "kind": "heading",
            "level": 1
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
                                    "value": "aaaa"
                                }
                            ],
                            "kind": "paragraph"
                        }
                    ],
                    "kind": "listItem"
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
                                    "value": "bbbb"
                                }
                            ],
                            "kind": "paragraph"
                        }
                    ],
                    "kind": "listItem"
                }
            ],
            "isTight": true,
            "kind": "list",
            "marker": "-",
            "start": 0
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
                            "value": "link"
                        }
                    ],
                    "destination": "http://www.example.com",
                    "kind": "link",
                    "title": ""
                }
            ],
            "kind": "paragraph"
        }
    ],
    "kind": "document",
    "metadata": {}
}
```
