package jsonrenderer

import "github.com/yuin/goldmark/v2/ast"

type commonMark struct {
	opts   []Option
	config *Config
}

func newCommonMark(opts ...Option) Extension {
	return &commonMark{opts: opts}
}

func (e *commonMark) RendererOptions(cfg *Config) []Option {
	if len(e.opts) != 0 {
		thisConfig := *cfg
		for _, opt := range e.opts {
			opt.SetFormatOption(&thisConfig)
		}
		cfg = &thisConfig
	}
	e.config = cfg

	return []Option{
		WithNodeRenderers(map[ast.NodeKind]NodeRenderer{
			ast.KindDocument: NodeRendererFunc(e.renderDocument),
			ast.KindHeading: ToMap(func(n *ast.Heading, _ []byte) map[string]any {
				return map[string]any{
					"kind":  "heading",
					"level": n.Level,
				}
			}, true),
			ast.KindBlockquote: AsMap(map[string]any{"kind": "blockquote"}, true),
			ast.KindCodeBlock: ToMap(func(n *ast.CodeBlock, source []byte) map[string]any {
				data := map[string]any{
					"kind":          "codeBlock",
					"codeBlockKind": n.CodeBlockKind.String(),
					"info":          n.Info.Value(source),
					"value":         n.Value.Str(source),
				}
				if l, ok := n.Language(source); ok {
					data["language"] = l
				}
				return data
			}, false),
			ast.KindHTMLBlock: ToMap(func(n *ast.HTMLBlock, source []byte) map[string]any {
				return map[string]any{
					"kind":          "htmlBlock",
					"htmlBlockKind": n.HTMLBlockKind.String(),
					"value":         n.Value.Str(source),
				}
			}, false),
			ast.KindList: ToMap(func(n *ast.List, _ []byte) map[string]any {
				return map[string]any{
					"kind":    "list",
					"marker":  string(n.Marker),
					"isTight": n.IsTight,
					"start":   n.Start,
				}
			}, true),
			ast.KindListItem:      AsMap(map[string]any{"kind": "listItem"}, true),
			ast.KindParagraph:     AsMap(map[string]any{"kind": "paragraph"}, true),
			ast.KindThematicBreak: AsMap(map[string]any{"kind": "thematicBreak"}, false),
			ast.KindLinkReferenceDefinition: ToMap(func(n *ast.LinkReferenceDefinition, source []byte) map[string]any {
				return map[string]any{
					"kind":        "linkReferenceDefinition",
					"label":       n.Label.Value(source),
					"destination": n.Destination.Value(source),
					"title":       n.Title.Value(source),
				}
			}, false),
			ast.KindAutoLink: ToMap(func(n *ast.AutoLink, source []byte) map[string]any {
				return map[string]any{
					"kind":        "autoLink",
					"destination": n.Destination.Value(source),
					"label":       n.Label.Value(source),
					"text":        n.Text.Value(source),
				}
			}, false),
			ast.KindCodeSpan: ToMap(func(n *ast.CodeSpan, source []byte) map[string]any {
				return map[string]any{
					"kind":  "codeSpan",
					"value": n.Value.Value(source),
				}
			}, false),
			ast.KindEmphasis: AsMap(map[string]any{"kind": "emphasis"}, true),
			ast.KindStrong:   AsMap(map[string]any{"kind": "strong"}, true),
			ast.KindImage: ToMap(func(n *ast.Image, source []byte) map[string]any {
				data := map[string]any{
					"kind":        "image",
					"destination": n.Destination.Value(source),
					"title":       n.Title.Value(source),
				}
				if n.Reference != nil {
					r := map[string]any{
						"referenceLinkKind": n.Reference.ReferenceLinkKind.String(),
						"value":             n.Reference.Value.Value(source),
					}
					data["reference"] = r
				}
				return data
			}, true),
			ast.KindLink: ToMap(func(n *ast.Link, source []byte) map[string]any {
				data := map[string]any{
					"kind":        "link",
					"destination": n.Destination.Value(source),
					"title":       n.Title.Value(source),
				}
				if n.Reference != nil {
					r := map[string]any{
						"referenceLinkKind": n.Reference.ReferenceLinkKind.String(),
						"value":             n.Reference.Value.Value(source),
					}
					data["reference"] = r
				}
				return data
			}, true),
			ast.KindRawHTML: ToMap(func(n *ast.RawHTML, source []byte) map[string]any {
				return map[string]any{
					"kind":  "rawHtml",
					"value": n.Value.Str(source),
				}
			}, false),
			ast.KindText: ToMap(func(n *ast.Text, source []byte) map[string]any {
				text := map[string]any{
					"kind":  "text",
					"value": n.Value.Value(source),
				}
				if n.SoftLineBreak() {
					text["softLineBreak"] = true
				}
				if n.HardLineBreak() {
					text["hardLineBreak"] = true
				}
				return text
			}, false),
		},
		),
	}
}
