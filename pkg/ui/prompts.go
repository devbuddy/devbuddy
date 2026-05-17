package ui

import "errors"

var ErrPromptCancelled = errors.New("prompt cancelled")

type SelectOption struct {
	Value string
	Label string
}

type SelectRequest struct {
	Label   string
	Options []SelectOption
}

type ConfirmRequest struct {
	Label string
}

type Prompts interface {
	Select(SelectRequest) (string, error)
	Confirm(ConfirmRequest) (bool, error)
}
