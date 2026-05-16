package ui

type Kind string

const (
	KindDebug             Kind = "debug"
	KindWarning           Kind = "warning"
	KindCommandHeader     Kind = "command_header"
	KindCommandRun        Kind = "command_run"
	KindCommandActed      Kind = "command_acted"
	KindProjectExists     Kind = "project_exists"
	KindJumpProject       Kind = "jump_project"
	KindTaskHeader        Kind = "task_header"
	KindTaskCommand       Kind = "task_command"
	KindTaskShell         Kind = "task_shell"
	KindTaskActed         Kind = "task_acted"
	KindTaskAlreadyOK     Kind = "task_already_ok"
	KindTaskError         Kind = "task_error"
	KindTaskWarning       Kind = "task_warning"
	KindTaskActionHeader  Kind = "task_action_header"
	KindActionHeader      Kind = "action_header"
	KindActionNotice      Kind = "action_notice"
	KindActionDone        Kind = "action_done"
	KindHookActivated     Kind = "hook_activated"
	KindHookFeatureFailed Kind = "hook_feature_failed"
	KindHookDevYMLChanged Kind = "hook_devyml_changed"
	KindShellDetectError  Kind = "shell_detect_error"
)

type Field struct {
	Name  string
	Value string
}

type Event struct {
	Kind   Kind
	Text   string
	Fields []Field
}

func F(name, value string) Field {
	return Field{Name: name, Value: value}
}
