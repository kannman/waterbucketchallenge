//go:build deptools
// +build deptools

package deps

import (
	// keep dependencies here to avoid removal from go.mod
	_ "github.com/swaggo/swag/cmd/swag"
)
