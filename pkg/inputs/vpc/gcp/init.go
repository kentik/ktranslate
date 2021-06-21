package gcp

import (
	"github.com/json-iterator/go"
)

// fast JSON encoding
var json = jsoniter.ConfigFastest

// jsonSorted is fast, but still sorts keys
var jsonSorted = jsoniter.Config{
	IndentionStep:          4,
	EscapeHTML:             false,
	SortMapKeys:            true,
	ValidateJsonRawMessage: false,
}.Froze()
