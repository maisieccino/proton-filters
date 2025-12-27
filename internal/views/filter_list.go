package views

import (
	"strconv"

	"github.com/ProtonMail/go-proton-api"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maisieccino/proton-filters/internal/types"
)

type FilterList struct {
	list list.Model
}

var _ list.Item = &FilterItem{}

func (v *FilterList) View() string {
	return v.list.View()
}

func (v *FilterList) Init() tea.Cmd {
	v.list = list.New(make([]list.Item, 0), list.NewDefaultDelegate(), 0, 0)
	return nil
}

func (v *FilterList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case types.FiltersMsg:
		_ = v.SetFilters(msg)
		v.list, cmd = v.list.Update(msg)
		return v, cmd
	}
	return v, cmd
}

func (v *FilterList) SetFilters(filters []proton.Filter) tea.Cmd {
	items := make([]list.Item, len(filters))
	for i, f := range filters {
		items[i] = FilterItem{f}
	}
	return v.list.SetItems(items)
}

type FilterItem struct {
	proton.Filter
}

func (i FilterItem) Title() string {
	return i.Name
}

func (i FilterItem) Description() string {
	return strconv.Itoa(i.Priority)
}

func (i FilterItem) FilterValue() string {
	return i.ID
}
