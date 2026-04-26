package views

import tea "charm.land/bubbletea/v2"

type DefaultScreen struct{}

var _ tea.Model = &DefaultScreen{}

func (v *DefaultScreen) Init() tea.Cmd {
	return nil
}

func (v *DefaultScreen) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (v *DefaultScreen) View() string {
	return "Hi"
}
