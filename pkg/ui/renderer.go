package ui

import (
	"fmt"
	"strings"

	color "github.com/logrusorgru/aurora"
)

type Renderer interface {
	Render(Event) string
}

type PlainRenderer struct{}
type TerminalRenderer struct{}

func (PlainRenderer) Render(event Event) string {
	return renderPlain(event)
}

func (TerminalRenderer) Render(event Event) string {
	switch event.Kind {
	case KindDebug:
		return fmt.Sprintf("%s: %s\n", color.Yellow("BUD_DEBUG"), color.Gray(12, event.Text))
	case KindWarning:
		return fmt.Sprintf("%s: %s\n", color.Bold(color.Yellow("WARNING")), event.Text)
	case KindCommandHeader:
		return fmt.Sprintf("🐼  %s %s\n", color.Blue("running"), color.Cyan(event.Text))
	case KindCommandRun:
		return fmt.Sprintf("%s %s\n", color.Bold(color.Cyan(event.Text)), color.Cyan(field(event, "args")))
	case KindCommandActed:
		return fmt.Sprintf("  %s\n", color.Green("Done!"))
	case KindProjectExists:
		return fmt.Sprintf("🐼  %s\n", color.Yellow("project already exists locally"))
	case KindJumpProject:
		return fmt.Sprintf("🐼  %s %s\n", color.Yellow("jumping to"), color.Green(event.Text))
	case KindTaskHeader:
		return renderTaskHeader(event, true)
	case KindTaskCommand:
		return fmt.Sprintf("  Running: %s %s\n", color.Bold(color.Cyan(event.Text)), color.Cyan(field(event, "args")))
	case KindTaskShell:
		return fmt.Sprintf("  Running: %s\n", color.Cyan(event.Text))
	case KindTaskActed:
		return fmt.Sprintf("  %s\n", color.Green("Done!"))
	case KindTaskAlreadyOK:
		return fmt.Sprintf("  %s\n", color.Green("Already OK!"))
	case KindTaskError:
		return fmt.Sprintf("  %s\n", color.Red(event.Text))
	case KindTaskWarning:
		return fmt.Sprintf("  Warning: %s\n", color.Yellow(event.Text))
	case KindTaskActionHeader:
		return fmt.Sprintf("  %s%s\n", color.Yellow("▪︎"), color.Magenta(event.Text))
	case KindActionHeader:
		return fmt.Sprintf("🐼  %s\n", color.Cyan(event.Text))
	case KindActionNotice:
		return fmt.Sprintf("⚠️   %s\n", color.Yellow(event.Text))
	case KindActionDone:
		return fmt.Sprintf("✅  %s\n", color.Green("Done!"))
	case KindHookActivated:
		return fmt.Sprintf("🐼  %s %s\n", color.Cyan("activated:"), renderFeatureList(event, true))
	case KindHookFeatureFailed:
		param := field(event, "param")
		if param != "" {
			param = fmt.Sprintf(" (%s)", color.Yellow(param))
		}
		return fmt.Sprintf("🐼  %s%s\n", color.Red(fmt.Sprintf("failed to activate %s. Try running 'bud up' first!", event.Text)), param)
	case KindHookDevYMLChanged:
		return fmt.Sprintf("🐼  %s\n", color.Yellow("dev.yml changed, run `bud up` to apply"))
	case KindShellDetectError:
		return fmt.Sprintf("%s %s\n", color.Yellow("Could not detect your shell:"), event.Text)
	default:
		return renderPlain(event)
	}
}

func renderPlain(event Event) string {
	switch event.Kind {
	case KindDebug:
		return fmt.Sprintf("BUD_DEBUG: %s\n", event.Text)
	case KindWarning:
		return fmt.Sprintf("WARNING: %s\n", event.Text)
	case KindCommandHeader:
		return fmt.Sprintf("🐼  running %s\n", event.Text)
	case KindCommandRun:
		return fmt.Sprintf("%s %s\n", event.Text, field(event, "args"))
	case KindCommandActed:
		return "  Done!\n"
	case KindProjectExists:
		return "🐼  project already exists locally\n"
	case KindJumpProject:
		return fmt.Sprintf("🐼  jumping to %s\n", event.Text)
	case KindTaskHeader:
		return renderTaskHeader(event, false)
	case KindTaskCommand:
		return fmt.Sprintf("  Running: %s %s\n", event.Text, field(event, "args"))
	case KindTaskShell:
		return fmt.Sprintf("  Running: %s\n", event.Text)
	case KindTaskActed:
		return "  Done!\n"
	case KindTaskAlreadyOK:
		return "  Already OK!\n"
	case KindTaskError:
		return fmt.Sprintf("  %s\n", event.Text)
	case KindTaskWarning:
		return fmt.Sprintf("  Warning: %s\n", event.Text)
	case KindTaskActionHeader:
		return fmt.Sprintf("  ▪︎%s\n", event.Text)
	case KindActionHeader:
		return fmt.Sprintf("🐼  %s\n", event.Text)
	case KindActionNotice:
		return fmt.Sprintf("⚠️   %s\n", event.Text)
	case KindActionDone:
		return "✅  Done!\n"
	case KindHookActivated:
		return fmt.Sprintf("🐼  activated: %s\n", renderFeatureList(event, false))
	case KindHookFeatureFailed:
		param := field(event, "param")
		if param != "" {
			param = fmt.Sprintf(" (%s)", param)
		}
		return fmt.Sprintf("🐼  failed to activate %s. Try running 'bud up' first!%s\n", event.Text, param)
	case KindHookDevYMLChanged:
		return "🐼  dev.yml changed, run `bud up` to apply\n"
	case KindShellDetectError:
		return fmt.Sprintf("Could not detect your shell: %s\n", event.Text)
	default:
		return event.Text + "\n"
	}
}

func renderTaskHeader(event Event, styled bool) string {
	name := event.Text
	param := field(event, "param")
	reason := field(event, "reason")
	if styled {
		name = color.Magenta(name).String()
		if param != "" {
			param = fmt.Sprintf(" (%s)", color.Blue(param))
		}
		if reason != "" {
			reason = fmt.Sprintf(" (%s)", color.Yellow(reason))
		}
		return fmt.Sprintf("%s %s%s%s\n", color.Yellow("◼︎"), name, param, reason)
	}
	if param != "" {
		param = fmt.Sprintf(" (%s)", param)
	}
	if reason != "" {
		reason = fmt.Sprintf(" (%s)", reason)
	}
	return fmt.Sprintf("◼︎ %s%s%s\n", name, param, reason)
}

func renderFeatureList(event Event, styled bool) string {
	parts := make([]string, 0, len(event.Fields))
	for _, feature := range event.Fields {
		if feature.Value == "" || strings.HasPrefix(feature.Value, "{") {
			if styled {
				parts = append(parts, color.Blue(feature.Name).String())
			} else {
				parts = append(parts, feature.Name)
			}
			continue
		}
		if styled {
			parts = append(parts, fmt.Sprintf("%s%s%s%s", color.Blue(feature.Name), color.Gray(12, "["), color.Cyan(feature.Value), color.Gray(12, "]")))
		} else {
			parts = append(parts, fmt.Sprintf("%s[%s]", feature.Name, feature.Value))
		}
	}
	separator := ", "
	if styled {
		separator = color.Gray(12, separator).String()
	}
	return strings.Join(parts, separator)
}

func field(event Event, name string) string {
	for _, field := range event.Fields {
		if field.Name == name {
			return field.Value
		}
	}
	return ""
}
