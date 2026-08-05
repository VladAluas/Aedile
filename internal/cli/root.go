// Package cli provides a cli tool for the infrastructure deployment
package cli

import (
	"github.com/VladAluas/Aedile/internal/commands"
)

func Execute() error {
	return commands.Root().Execute()
}
