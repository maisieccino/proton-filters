package views

import (
	"github.com/ProtonMail/go-proton-api"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maisieccino/proton-filters/internal/types"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type FilterList struct {
	list   list.Model
	client *proton.Client
}

var _ list.Item = &FilterItem{}

func NewFilterList() (tea.Model, string) {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 1, 1)
	l.Title = "Sieve filters"
	return &FilterList{
		list: l,
	}, "filter-list"
}

func (v *FilterList) View() string {
	var current FilterItem
	if v.list.SelectedItem() != nil {
		current = v.list.SelectedItem().(FilterItem)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top,
		v.list.View(),
		current.Sieve,
	)
}

func (v *FilterList) Init() tea.Cmd {
	return nil
}

func (v *FilterList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case types.FiltersMsg:
		cmd = v.SetFilters(msg)
		return v, cmd
	case tea.WindowSizeMsg:
		x, y := docStyle.GetFrameSize()
		v.list.SetSize(msg.Width-x, msg.Height-y)
	}
	v.list, cmd = v.list.Update(msg)
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
	switch i.Status {
	case 1:
		return "Enabled"
	default:
		return "Disabled"

	}
}

func (i FilterItem) FilterValue() string {
	return i.Name
}
