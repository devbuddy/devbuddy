package ui

type UI struct {
	events  []Event
	prompts Prompts
}

func New() *UI {
	return &UI{prompts: SurveyPrompts{}}
}

func NewTesting() (*FakePrompts, *UI) {
	prompts := &FakePrompts{}
	return prompts, &UI{prompts: prompts}
}

func (u *UI) Record(event Event) {
	u.events = append(u.events, event)
}

func (u *UI) Events() []Event {
	return append([]Event(nil), u.events...)
}

func (u *UI) Prompts() Prompts {
	return u.prompts
}

func (u *UI) SetPrompts(prompts Prompts) {
	u.prompts = prompts
}
