package ui

import (
	"errors"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

type SurveyPrompts struct{}

func (SurveyPrompts) Select(req SelectRequest) (string, error) {
	labels := make([]string, 0, len(req.Options))
	values := map[string]string{}
	for _, opt := range req.Options {
		labels = append(labels, opt.Label)
		values[opt.Label] = opt.Value
	}

	var selected string
	err := survey.AskOne(&survey.Select{
		Message: req.Label,
		Options: labels,
	}, &selected)
	if errors.Is(err, terminal.InterruptErr) {
		return "", ErrPromptCancelled
	}
	if err != nil {
		return "", err
	}
	return values[selected], nil
}

func (SurveyPrompts) Confirm(req ConfirmRequest) (bool, error) {
	var confirmed bool
	err := survey.AskOne(&survey.Confirm{
		Message: req.Label,
	}, &confirmed)
	if errors.Is(err, terminal.InterruptErr) {
		return false, ErrPromptCancelled
	}
	return confirmed, err
}
