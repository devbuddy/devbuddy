package ui

type FakePrompts struct {
	SelectRequests  []SelectRequest
	ConfirmRequests []ConfirmRequest
	SelectValue     string
	SelectErr       error
	ConfirmValue    bool
	ConfirmErr      error
}

func (p *FakePrompts) Select(req SelectRequest) (string, error) {
	p.SelectRequests = append(p.SelectRequests, req)
	return p.SelectValue, p.SelectErr
}

func (p *FakePrompts) Confirm(req ConfirmRequest) (bool, error) {
	p.ConfirmRequests = append(p.ConfirmRequests, req)
	return p.ConfirmValue, p.ConfirmErr
}
