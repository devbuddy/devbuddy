package ui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

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
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, "%s (y/N): ", req.Label)
		answer, err := reader.ReadString('\n')
		fmt.Fprintln(os.Stderr)
		if err != nil && !errors.Is(err, io.EOF) {
			return false, err
		}

		switch strings.ToLower(strings.TrimSpace(answer)) {
		case "y", "yes":
			return true, nil
		case "", "n", "no":
			return false, nil
		}
	}
}
