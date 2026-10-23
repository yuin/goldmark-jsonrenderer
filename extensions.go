package jsonrenderer

import (
	alert "github.com/yuin/goldmark-alert"
	"github.com/yuin/goldmark/v2/ast"
	east "github.com/yuin/goldmark/v2/extension/ast"
)

type footnote struct {
}

// NewFootnote returns a new footnote extension.
func NewFootnote() Extension {
	return &footnote{}
}

// Footnote is a default footnote extension.
var Footnote = NewFootnote()

func (e *footnote) RendererOptions(_ *Config) []Option {
	return []Option{
		WithNodeRenderers(map[ast.NodeKind]NodeRenderer{
			east.KindFootnoteReference: ToMap(func(n *east.FootnoteReference, source []byte) map[string]any {
				return map[string]any{
					"kind":     "footnoteReference",
					"label":    n.Label.Value(source),
					"index":    n.Index,
					"refIndex": n.RefIndex,
				}
			}, true),
			east.KindFootnoteDefinition: ToMap(func(n *east.FootnoteDefinition, source []byte) map[string]any {
				return map[string]any{
					"kind":  "footnoteDefinition",
					"label": n.Label.Value(source),
				}
			}, true),
		}),
	}
}

type strikethrough struct {
}

// NewStrikethrough returns a new strikethrough extension.
func NewStrikethrough() Extension {
	return &strikethrough{}
}

// Strikethrough is a default strikethrough extension.
var Strikethrough = NewStrikethrough()

func (e *strikethrough) RendererOptions(_ *Config) []Option {
	return []Option{
		WithNodeRenderers(map[ast.NodeKind]NodeRenderer{
			east.KindStrikethrough: AsMap(map[string]any{"kind": "strikethrough"}, true),
		}),
	}
}

type table struct {
}

// NewTable returns a new table extension.
func NewTable() Extension {
	return &table{}
}

// Table is a default table extension.
var Table = NewTable()

func (e *table) RendererOptions(_ *Config) []Option {
	return []Option{
		WithNodeRenderers(map[ast.NodeKind]NodeRenderer{
			east.KindTable:       AsMap(map[string]any{"kind": "table"}, true),
			east.KindTableHeader: AsMap(map[string]any{"kind": "tableHeader"}, true),
			east.KindTableBody:   AsMap(map[string]any{"kind": "tableBody"}, true),
			east.KindTableRow:    AsMap(map[string]any{"kind": "tableRow"}, true),
			east.KindTableCell: ToMap(func(n *east.TableCell, _ []byte) map[string]any {
				return map[string]any{
					"kind":  "tableCell",
					"align": n.Alignment.String(),
				}
			}, true),
		}),
	}
}

type definitionList struct {
}

// NewDefinitionList returns a new definition list extension.
func NewDefinitionList() Extension {
	return &definitionList{}
}

// DefinitionList is a default definition list extension.
var DefinitionList = NewDefinitionList()

func (e *definitionList) RendererOptions(_ *Config) []Option {
	return []Option{
		WithNodeRenderers(map[ast.NodeKind]NodeRenderer{
			east.KindDefinitionList:        AsMap(map[string]any{"kind": "definitionList"}, true),
			east.KindDefinitionTerm:        AsMap(map[string]any{"kind": "definitionTerm"}, true),
			east.KindDefinitionDescription: AsMap(map[string]any{"kind": "definitionDescription"}, true),
		}),
	}
}

type alertr struct {
}

// NewAlert returns a new alert extension.
func NewAlert() Extension {
	return &alertr{}
}

// Alert is a default alert extension.
var Alert = NewAlert()

func (e *alertr) RendererOptions(_ *Config) []Option {
	return []Option{
		WithNodeRenderers(map[ast.NodeKind]NodeRenderer{
			alert.KindAlert: ToMap(func(n *alert.Alert, source []byte) map[string]any {
				return map[string]any{
					"kind":      "alert",
					"alertKind": n.AlertKind.Value(source),
					"title":     n.Title.Value(source),
				}
			}, true),
		}),
	}
}
