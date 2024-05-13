//go:build !doc

// Package doc is used to generate docs
// when `-tags doc` is passed to go build
// Enable = true, and the corresponding logic in the main.go will be executed
// g.md and man will be generated
// by default, Enable = false
package doc

const Enable = false
