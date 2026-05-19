package helpers

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/executor"
	"github.com/devbuddy/devbuddy/pkg/helpers/osidentity"
)

const xcodeCLTInstallMessage = "Xcode Command Line Tools are required. Complete the installer dialog, then re-run bud up."
const xcodeFirstLaunchMessage = "Xcode first-launch setup is required. Run `sudo xcodebuild -runFirstLaunch`, then re-run bud up."

// EnsureXcodeCommandLineTools verifies the macOS native build toolchain needed
// by source-built language installers such as pyenv and rbenv.
func EnsureXcodeCommandLineTools(ctx *context.Context) error {
	return ensureXcodeCommandLineTools(ctx, osidentity.Detect(), term.IsTerminal(int(os.Stdin.Fd())))
}

func ensureXcodeCommandLineTools(ctx *context.Context, osIdent *osidentity.Identity, interactive bool) error {
	if !osIdent.IsMacOS() {
		return nil
	}

	result := ctx.Executor.Capture(executor.New("xcode-select", "-p"))
	if result.Error == nil {
		return ensureXcodeFirstLaunch(ctx, interactive)
	}

	ctx.UI.Warningf(xcodeCLTInstallMessage)
	result = ctx.RunTaskCommand(executor.New("xcode-select", "--install"))
	if result.Error != nil {
		return fmt.Errorf("failed to start Xcode Command Line Tools installer: %w", result.Error)
	}

	return errors.New(xcodeCLTInstallMessage)
}

func ensureXcodeFirstLaunch(ctx *context.Context, interactive bool) error {
	result := ctx.Executor.Capture(executor.New("xcodebuild", "-checkFirstLaunchStatus"))
	if result.Error == nil {
		return nil
	}
	if !interactive {
		ctx.UI.Warningf(xcodeFirstLaunchMessage)
		return errors.New(xcodeFirstLaunchMessage)
	}

	ctx.UI.Warningf("Xcode first-launch setup is required. DevBuddy will run `sudo xcodebuild -runFirstLaunch`; this installs required Xcode components and accepts the Xcode/SDK license.")
	result = ctx.RunTaskCommand(executor.New("sudo", "xcodebuild", "-runFirstLaunch"))
	if result.Error != nil {
		return fmt.Errorf("failed to run Xcode first-launch setup: %w", result.Error)
	}
	return nil
}
