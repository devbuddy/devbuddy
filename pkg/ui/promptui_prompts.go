package ui

import (
	"errors"

	"github.com/manifoldco/promptui"
)

type PromptUIPrompts struct{}

func (PromptUIPrompts) Select(req SelectRequest) (string, error) {
	items := make([]selectItem, 0, len(req.Options))
	for _, opt := range req.Options {
		items = append(items, selectItem(opt))
	}

	prompt := promptui.Select{
		Label:        req.Label,
		Items:        items,
		HideSelected: true,
		Templates:    promptUISelectTemplates(),
	}

	index, _, err := prompt.Run()
	if errors.Is(err, promptui.ErrInterrupt) || errors.Is(err, promptui.ErrEOF) {
		return "", ErrPromptCancelled
	}
	if err != nil {
		return "", err
	}
	return items[index].Value, nil
}

func (PromptUIPrompts) Confirm(req ConfirmRequest) (bool, error) {
	prompt := promptui.Prompt{
		Label:     req.Label,
		IsConfirm: true,
	}
	_, err := prompt.Run()
	if errors.Is(err, promptui.ErrAbort) || errors.Is(err, promptui.ErrInterrupt) || errors.Is(err, promptui.ErrEOF) {
		return false, nil
	}
	return err == nil, err
}

type selectItem struct {
	Value string
	Label string
}

func promptUISelectTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "🐼 {{ .Label | cyan }}",
		Inactive: "   {{ .Label }}",
		Selected: "🐼 {{ .Label }}",
	}
}
