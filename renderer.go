package jsonrenderer

import (
	"maps"

	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/renderer"
)

// Renderer is a renderer that renders AST nodes as JSON.
type Renderer = renderer.Renderer[*JSONWriter]

// NodeRenderer is a renderer that renders AST nodes as JSON.
type NodeRenderer = renderer.NodeRenderer[*JSONWriter]

// NodeRendererDecorator is a decorator that decorates NodeRenderer.
type NodeRendererDecorator = renderer.NodeRendererDecorator[*JSONWriter]

// NodeRendererFunc adapts f into a NodeRenderer that renders AST nodes as JSON.
func NodeRendererFunc(f func(w *JSONWriter, source []byte,
	n ast.Node, entering bool, rc renderer.Context) (ast.WalkStatus, error)) NodeRenderer {
	return renderer.NodeRendererFunc(f)
}

// Extension is an extension that extends JSON based renderers.
type Extension = renderer.Extension[Config]

// Option is a functional option for JSON based renderers.
type Option = renderer.Option[Config]

// WithNodeRenderers sets a node renderer for the given node kind.
func WithNodeRenderers(nodeRenderers map[ast.NodeKind]NodeRenderer) Option {
	return renderer.WithNodeRenderers[*JSONWriter, Config](nodeRenderers)
}

// WithNodeRenderer sets a node renderer for the given node kind.
func WithNodeRenderer(kind ast.NodeKind, nodeRenderer NodeRenderer) Option {
	return renderer.WithNodeRenderer[*JSONWriter, Config](kind, nodeRenderer)
}

// WithNodeRendererDecorators sets a node renderer decorator for the given node kind.
//
// If a decorator is already set for a node kind, the new decorator will be applied to the existing one.
func WithNodeRendererDecorators(decorators map[ast.NodeKind]NodeRendererDecorator) Option {
	return renderer.WithNodeRendererDecorators[*JSONWriter, Config](decorators)
}

// WithNodeRendererDecorator sets a node renderer decorator for the given node kind.
//
// If a decorator is already set for a node kind, the new decorator will be applied to the existing one.
func WithNodeRendererDecorator(kind ast.NodeKind, decorator NodeRendererDecorator) Option {
	return renderer.WithNodeRendererDecorator[*JSONWriter, Config](kind, decorator)
}

// WithExtensions adds extensions.
func WithExtensions(extensions ...Extension) Option {
	return renderer.WithExtensions[*JSONWriter](extensions...)
}

// A Config struct has configurations for the JSON based renderers.
type Config struct {
	renderer.Config[*JSONWriter, Config]
}

// Default returns a default configuration for the JSON based renderers.
func (c Config) Default() Config {
	return Config{}
}

type jsonRenderer struct {
	*renderer.Helper[*JSONWriter, Config]
}

// New returns a new Renderer with given options.
func New(opts ...Option) Renderer {
	cm := newCommonMark()
	opts = append([]Option{WithExtensions(cm)}, opts...)
	var hb renderer.HelperBuilder[*JSONWriter, Config]
	helper := hb.Options(opts...).Build()
	r := &jsonRenderer{
		Helper: helper,
	}
	return r
}

func addAttributes(data map[string]any, source []byte, node ast.Node) {
	attrs := map[string]string{}
	for _, attr := range node.Attributes() {
		attrs[string(attr.Name)] = attr.Value.Value(source)
	}
	data["attributes"] = attrs
}

// AsMap returns a NodeRenderer that renders AST nodes as JSON using the given map m.
func AsMap(m map[string]any, hasChildren bool) NodeRenderer {
	return NodeRendererFunc(func(
		w *JSONWriter, source []byte, node ast.Node, entering bool, _ renderer.Context) (ast.WalkStatus, error) {
		if entering {
			data := maps.Clone(m)
			if hasChildren {
				data["children"] = []any{}
			}
			addAttributes(data, source, node)
			w.Enter(data)
			return ast.WalkContinue, nil
		}
		w.Leave()
		return ast.WalkContinue, nil
	})
}

// ToMap returns a NodeRenderer that renders AST nodes as JSON using the given function f.
func ToMap[T any](f func(T, []byte) map[string]any, hasChildren bool) NodeRenderer {
	return NodeRendererFunc(func(
		w *JSONWriter, source []byte, node ast.Node, entering bool, _ renderer.Context) (ast.WalkStatus, error) {
		if entering {
			n := node.(T)
			data := f(n, source)
			if hasChildren {
				data["children"] = []any{}
			}
			addAttributes(data, source, node)
			w.Enter(data)
			return ast.WalkContinue, nil
		}
		w.Leave()
		return ast.WalkContinue, nil
	})
}

// ToMapRC returns a NodeRenderer that renders AST nodes as JSON using the given function f.
func ToMapRC[T any](f func(T, []byte, renderer.Context) map[string]any, hasChildren bool) NodeRenderer {
	return NodeRendererFunc(func(
		w *JSONWriter, source []byte, node ast.Node, entering bool, rc renderer.Context) (ast.WalkStatus, error) {
		if entering {
			n := node.(T)
			data := f(n, source, rc)
			if hasChildren {
				data["children"] = []any{}
			}
			addAttributes(data, source, node)
			w.Enter(data)
			return ast.WalkContinue, nil
		}
		w.Leave()
		return ast.WalkContinue, nil
	})
}

func (e *commonMark) renderDocument(w *JSONWriter, source []byte, node ast.Node,
	entering bool, _ renderer.Context) (ast.WalkStatus, error) {
	if entering {
		n := node.(*ast.Document)
		doc := map[string]any{
			"kind":     "document",
			"metadata": n.Metadata(),
			"children": []any{},
		}
		addAttributes(doc, source, node)
		w.Enter(doc)
	}
	return ast.WalkContinue, nil
}
