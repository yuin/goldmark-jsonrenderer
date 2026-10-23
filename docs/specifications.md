# Specifications

## Overview

This document outlines the specifications of the JSON representation of a
Markdown document. The `markdown2json` command enables all supported
extensions by default.

## JSON Structure

All Markdown elements are represented as JSON objects. Every node has a
`kind` and an `attributes` object. Nodes that can contain other nodes also
have a `children` array.

```json
{
  "kind": "document",
  "attributes": {},
  "metadata": {},
  "children": [
    {
      "kind": "heading",
      "attributes": {},
      "level": 1,
      "children": [
        {
          "kind": "text",
          "attributes": {},
          "value": "Heading 1"
        }
      ]
    }
  ]
}
```

**Common Properties:**

| property | type | description |
|----------|------|-------------|
| `kind` | string | A string indicating the type of the element. |
| `attributes` | object | Attributes associated with the element. Attribute names and values are strings. It is `{}` when the element has no attributes. |
| `children` | array | Child elements. This property is present only on elements that can have children. |

## Elements

### Document

| kind | has children |
|------|--------------|
| document | true |

Root element of the JSON representation. It contains all other elements as
children.

```json
{
  "kind": "document",
  "attributes": {},
  "metadata": {
    "title": "Example"
  },
  "children": []
}
```

| property | type | description |
|----------|------|-------------|
| `metadata` | object | Metadata set by the `meta` extension. It is `{}` when no metadata is present. |

### Heading

| kind | has children |
|------|--------------|
| heading | true |

Represents an ATX or setext heading.

```json
{
  "kind": "heading",
  "attributes": {},
  "level": 1,
  "children": []
}
```

| property | type | description |
|----------|------|-------------|
| `level` | number | The heading level, from 1 through 6. |

### Blockquote

| kind | has children |
|------|--------------|
| blockquote | true |

Represents a block quote.

```json
{
  "kind": "blockquote",
  "attributes": {},
  "children": []
}
```

### Code Block

| kind | has children |
|------|--------------|
| codeBlock | false |

Represents an indented or fenced code block.

```json
{
  "kind": "codeBlock",
  "attributes": {},
  "codeBlockKind": "Fenced",
  "info": "go example",
  "language": "go",
  "value": "fmt.Println(\"hello\")\n"
}
```

| property | type | description |
|----------|------|-------------|
| `codeBlockKind` | string | `"Fenced"` for a fenced code block or `"Indented"` for an indented code block. |
| `info` | string | The complete info string. It is empty for an indented code block. |
| `language` | string | The first word of `info`. This property is omitted when the info string has no language. |
| `value` | string | The code block content. |

### HTML Block

| kind | has children |
|------|--------------|
| htmlBlock | false |

Represents a block HTML element.

```json
{
  "kind": "htmlBlock",
  "attributes": {},
  "htmlBlockKind": "Kind6",
  "value": "<div>block HTML</div>\n"
}
```

| property | type | description |
|----------|------|-------------|
| `htmlBlockKind` | string | The CommonMark spec HTML-block kind. Its value depends on the parsed HTML block form. |
| `value` | string | The raw HTML block content. |

### List

| kind | has children |
|------|--------------|
| list | true |

Represents an ordered or unordered list.

```json
{
  "kind": "list",
  "attributes": {},
  "marker": ".",
  "isTight": true,
  "start": 1,
  "children": [
    {
      "kind": "listItem",
      "attributes": {},
      "children": []
    }
  ]
}
```

| property | type | description |
|----------|------|-------------|
| `marker` | string | The list marker. It is `"-"`, `"+"`, or `"*"` for unordered lists and `"."` or `")"` for ordered lists. |
| `isTight` | boolean | Whether the list is tight. |
| `start` | number | The first item number for an ordered list; `0` for an unordered list. |

### List Item

| kind | has children |
|------|--------------|
| listItem | true |

Represents an item in a list.

```json
{
  "kind": "listItem",
  "attributes": {},
  "children": []
}
```

Task-list items also use `listItem`. Their task status is included in
`attributes`:

```json
{
  "kind": "listItem",
  "attributes": {
    "task-status": "completed"
  },
  "children": []
}
```

The value of `task-status` is `"completed"` for checked items and `"active"`
for unchecked items.

### Paragraph

| kind | has children |
|------|--------------|
| paragraph | true |

Represents a paragraph.

```json
{
  "kind": "paragraph",
  "attributes": {},
  "children": []
}
```

### Thematic Break

| kind | has children |
|------|--------------|
| thematicBreak | false |

Represents a thematic break.

```json
{
  "kind": "thematicBreak",
  "attributes": {}
}
```

### Link Reference Definition

| kind | has children |
|------|--------------|
| linkReferenceDefinition | false |

Represents a link reference definition.

```json
{
  "kind": "linkReferenceDefinition",
  "attributes": {},
  "label": "example",
  "destination": "https://www.example.com",
  "title": "Example title"
}
```

| property | type | description |
|----------|------|-------------|
| `label` | string | The reference label. |
| `destination` | string | The link destination. |
| `title` | string | The link title. It is an empty string when no title is provided. |

### Auto Link

| kind | has children |
|------|--------------|
| autoLink | false |

Represents an autolink, including an email autolink.

```json
{
  "kind": "autoLink",
  "attributes": {},
  "destination": "https://www.example.com",
  "label": "https://www.example.com",
  "text": "<https://www.example.com>"
}
```

| property | type | description |
|----------|------|-------------|
| `destination` | string | The resolved link destination. |
| `label` | string | The link label without angle brackets. |
| `text` | string | The original autolink text, including angle brackets. |

### Code Span

| kind | has children |
|------|--------------|
| codeSpan | false |

Represents inline code.

```json
{
  "kind": "codeSpan",
  "attributes": {},
  "value": "code"
}
```

| property | type | description |
|----------|------|-------------|
| `value` | string | The inline code content. |

### Emphasis

| kind | has children |
|------|--------------|
| emphasis | true |

Represents emphasized text.

```json
{
  "kind": "emphasis",
  "attributes": {},
  "children": []
}
```

### Strong

| kind | has children |
|------|--------------|
| strong | true |

Represents strongly emphasized text.

```json
{
  "kind": "strong",
  "attributes": {},
  "children": []
}
```

### Image

| kind | has children |
|------|--------------|
| image | true |

Represents an image.

```json
{
  "kind": "image",
  "attributes": {},
  "destination": "image.png",
  "title": "Image title",
  "reference": {
    "referenceLinkKind": "Full",
    "value": "image-reference"
  },
  "children": [
    {
      "kind": "text",
      "attributes": {},
      "value": "alternative text"
    }
  ]
}
```

| property | type | description |
|----------|------|-------------|
| `destination` | string | The image destination. |
| `title` | string | The image title. It is an empty string when no title is provided. |
| `reference` | object | Reference-link information. It is omitted for inline images. |
| `reference.referenceLinkKind` | string | `"Full"`, `"Collapsed"`, or `"Shortcut"`. |
| `reference.value` | string | The reference label. |

### Link

| kind | has children |
|------|--------------|
| link | true |

Represents a link.

```json
{
  "kind": "link",
  "attributes": {},
  "destination": "https://www.example.com",
  "title": "Example title",
  "reference": {
    "referenceLinkKind": "Full",
    "value": "example"
  },
  "children": [
    {
      "kind": "text",
      "attributes": {},
      "value": "label"
    }
  ]
}
```

| property | type | description |
|----------|------|-------------|
| `destination` | string | The link destination. |
| `title` | string | The link title. It is an empty string when no title is provided. |
| `reference` | object | Reference-link information. It is omitted for inline links. |
| `reference.referenceLinkKind` | string | `"Full"`, `"Collapsed"`, or `"Shortcut"`. |
| `reference.value` | string | The reference label. |

### Raw HTML

| kind | has children |
|------|--------------|
| rawHtml | false |

Represents an inline raw HTML element.

```json
{
  "kind": "rawHtml",
  "attributes": {},
  "value": "<span>"
}
```

| property | type | description |
|----------|------|-------------|
| `value` | string | The raw HTML content. |

### Text

| kind | has children |
|------|--------------|
| text | false |

Represents text.

```json
{
  "kind": "text",
  "attributes": {},
  "value": "line before a soft break",
  "softLineBreak": true
}
```

| property | type | description |
|----------|------|-------------|
| `value` | string | The text content. |
| `softLineBreak` | boolean | Present and `true` when this text node ends at a soft line break. |
| `hardLineBreak` | boolean | Present and `true` when this text node ends at a hard line break. |

Only one of `softLineBreak` and `hardLineBreak` can be present. Both properties
are omitted when the text node does not end at a line break.

### Footnote Reference

| kind | has children |
|------|--------------|
| footnoteReference | true |

Represents a reference to a footnote definition.

```json
{
  "kind": "footnoteReference",
  "attributes": {},
  "label": "note",
  "index": 1,
  "refIndex": 0,
  "children": []
}
```

| property | type | description |
|----------|------|-------------|
| `label` | string | The footnote label. |
| `index` | number | The one-based index of the footnote definition. |
| `refIndex` | number | The zero-based occurrence index for references to the definition. |

### Footnote Definition

| kind | has children |
|------|--------------|
| footnoteDefinition | true |

Represents a footnote definition.

```json
{
  "kind": "footnoteDefinition",
  "attributes": {},
  "label": "note",
  "children": []
}
```

| property | type | description |
|----------|------|-------------|
| `label` | string | The footnote label. |

### Strikethrough

| kind | has children |
|------|--------------|
| strikethrough | true |

Represents strikethrough text.

```json
{
  "kind": "strikethrough",
  "attributes": {},
  "children": []
}
```

### Table

| kind | has children |
|------|--------------|
| table | true |

Represents a table.

```json
{
  "kind": "table",
  "attributes": {},
  "children": []
}
```

### Table Header

| kind | has children |
|------|--------------|
| tableHeader | true |

Represents the header row of a table.

```json
{
  "kind": "tableHeader",
  "attributes": {},
  "children": []
}
```

### Table Body

| kind | has children |
|------|--------------|
| tableBody | true |

Represents the body of a table.

```json
{
  "kind": "tableBody",
  "attributes": {},
  "children": []
}
```

### Table Row

| kind | has children |
|------|--------------|
| tableRow | true |

Represents a table row.

```json
{
  "kind": "tableRow",
  "attributes": {},
  "children": []
}
```

### Table Cell

| kind | has children |
|------|--------------|
| tableCell | true |

Represents a table cell.

```json
{
  "kind": "tableCell",
  "attributes": {},
  "align": "center",
  "children": []
}
```

| property | type | description |
|----------|------|-------------|
| `align` | string | The cell alignment: `"none"`, `"left"`, `"center"`, or `"right"`. |

### Definition List

| kind | has children |
|------|--------------|
| definitionList | true |

Represents a definition list.

```json
{
  "kind": "definitionList",
  "attributes": {},
  "children": []
}
```

### Definition Term

| kind | has children |
|------|--------------|
| definitionTerm | true |

Represents a term in a definition list.

```json
{
  "kind": "definitionTerm",
  "attributes": {},
  "children": []
}
```

### Definition Description

| kind | has children |
|------|--------------|
| definitionDescription | true |

Represents a description in a definition list.

```json
{
  "kind": "definitionDescription",
  "attributes": {},
  "children": []
}
```

### Alert

| kind | has children |
|------|--------------|
| alert | true |

Represents a GitHub-style alert.

```json
{
  "kind": "alert",
  "attributes": {},
  "alertKind": "NOTE",
  "title": "Alert title",
  "children": []
}
```

| property | type | description |
|----------|------|-------------|
| `alertKind` | string | The alert label parsed by `goldmark-alert`. |
| `title` | string | The alert title. It is an empty string when no title is provided. |
