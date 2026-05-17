package ui

type UI struct {
	events []Event
}

func New() *UI {
	return &UI{}
}

func (u *UI) Record(event Event) {
	u.events = append(u.events, event)
}

func (u *UI) Events() []Event {
	return append([]Event(nil), u.events...)
}
