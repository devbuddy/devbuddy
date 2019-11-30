package integration

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/devbuddy/devbuddy/pkg/termui"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

type CompletionScriptProvider interface {
	GenBashCompletion(w io.Writer) error
	GenZshCompletion(w io.Writer) error
}

// Print prints the integration code for the user's shell
func Print(withCompletion bool, completionScriptProvider CompletionScriptProvider) {
	shell, err := DetectShell()
	if err != nil {
		termui.HookShellDetectionError(err)
		return
	}

	script := buildCompletionScript(shell, withCompletion, completionScriptProvider)
	fmt.Println(script)
}

func buildCompletionScript(shell ShellIdentity, withCompletion bool, completionScriptProvider CompletionScriptProvider) string {
	buffer := bytes.NewBufferString(shellSource)

	switch shell {
	case BASH:
		buffer.WriteString(bashSource)
		if withCompletion {
			_ = completionScriptProvider.GenBashCompletion(buffer)
		}
	case ZSH:
		buffer.WriteString(zshSource)
		if withCompletion {
			_ = completionScriptProvider.GenZshCompletion(buffer)
			buffer.WriteString("compdef _bud bud") // interactively define the completion function
		}
	}

	return buffer.String()
}

// AddFinalizerCd declares a "cd" finalizer (change directory)
func AddFinalizerCd(path string) error {
	return addFinalizer("cd", path)
}

func formatError(message string) error {
	return fmt.Errorf(`there is something wrong with the shell integration:

    %s

This usually means that DevBuddy is not setup properly.
Please follow the setup steps: https://github.com/devbuddy/devbuddy/tree/master#setup

If DevBuddy is already setup, then please open an issue on https://github.com/devbuddy/devbuddy/issues/new?labels=bug
You can use "bud --report-issue" to do that.
`, message)
}

func addFinalizer(action, arg string) (err error) {
	content := fmt.Sprintf("%s:%s\n", action, arg)

	finalizerPath := os.Getenv("BUD_FINALIZER_FILE")

	if finalizerPath == "" {
		return formatError("the BUD_FINALIZER_FILE environment variable is missing or empty")
	}

	return utils.AppendOnlyFile(finalizerPath, []byte(content))
}
