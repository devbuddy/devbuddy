package cmd

import (
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/integration"
	"github.com/devbuddy/devbuddy/pkg/manifest"
	"github.com/devbuddy/devbuddy/pkg/project"
)

var cloneCmd = &cobra.Command{
	Use:          "clone [REMOTE]",
	Short:        "Clone a project from github.com",
	RunE:         cloneRun,
	Args:         onlyOneArg,
	GroupID:      "devbuddy",
	SilenceUsage: true,
}

func cloneRun(_ *cobra.Command, args []string) error {
	ctx, err := context.Load(false)
	if err != nil {
		return err
	}

	proj, err := project.NewFromID(args[0], ctx.Cfg)
	if err != nil {
		return err
	}

	if proj.Exists() {
		ctx.UI.ProjectExists()
	} else {
		if err := proj.Clone(ctx.Executor); err != nil {
			return err
		}
	}

	if !manifest.ExistsIn(proj.Path) {
		prompt := promptui.Prompt{
			Label:     "This project has no dev.yml. Create one",
			IsConfirm: true,
		}
		if _, err := prompt.Run(); err == nil {
			if err := createManifest(ctx.UI, proj.Path, ""); err != nil {
				return err
			}
		}
	}

	ctx.UI.JumpProject(proj.FullName())
	return integration.AddFinalizerCd(proj.Path)
}
