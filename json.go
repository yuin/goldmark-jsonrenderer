package jsonrenderer

import (
	"encoding/json"
	"encoding/json/jsontext"
)

// JSONWriter is a writer for JSON values.
type JSONWriter struct {
	current []map[string]any
}

// Enter enters a new JSON object.
//
// Given map must contains following keys:
// - "kind": string, the kind of the node.
func (w *JSONWriter) Enter(v map[string]any) {
	if len(w.current) == 0 {
		w.current = []map[string]any{v}
		return
	}
	current := w.current[len(w.current)-1]
	if _, has := current["children"]; !has {
		panic("Current node does not support children.")
	}
	children := current["children"].([]any)
	children = append(children, v)
	current["children"] = children
	w.current = append(w.current, v)
}

// Leave leaves the current JSON object.
func (w *JSONWriter) Leave() {
	if len(w.current) > 0 {
		w.current = w.current[:len(w.current)-1]
	}
}

// Map returns the value as a map.
func (w *JSONWriter) Map() map[string]any {
	if len(w.current) == 0 {
		return nil
	}
	return w.current[0]
}

// ToValue converts the value to a jsontext.Value.
func (w *JSONWriter) ToValue() jsontext.Value {
	b, err := json.Marshal(w.Map())
	if err != nil {
		panic(err) // should not happen
	}
	return jsontext.Value(b)
}
