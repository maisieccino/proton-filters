package views

import (
	"bytes"
	"os"
	"slices"

	"github.com/ProtonMail/go-proton-api"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maisieccino/proton-filters/internal/types"
	"github.com/muesli/reflow/wordwrap"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type FilterList struct {
	list    list.Model
	client  *proton.Client
	vp      viewport.Model
	vpReady bool
}

var _ list.Item = &FilterItem{}

func NewFilterList() (tea.Model, string) {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 1, 1)
	l.Title = "Sieve filters"
	l.SetStatusBarItemName("filter", "filters")
	return &FilterList{
		list: l,
	}, "filter-list"
}

func (v *FilterList) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		v.list.View(),
		v.vp.View(),
	)
}

func (v *FilterList) RenderFilter() string {
	var current FilterItem
	if v.list.SelectedItem() != nil {
		current = v.list.SelectedItem().(FilterItem)
	}
	buf := new(bytes.Buffer)
	err := quick.Highlight(os.Stdout, current.Sieve, "sieve", "terminal256", "monokai")
	if err != nil {
		return current.Sieve
	}
	return wordwrap.String(buf.String(), v.vp.Width)
}

func (v *FilterList) Init() tea.Cmd {
	return nil
}

func (v *FilterList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd, vpCmd tea.Cmd
	switch msg := msg.(type) {
	case types.FiltersMsg:
		cmd = v.SetFilters(msg)
		return v, cmd
	case tea.WindowSizeMsg:
		x, y := docStyle.GetFrameSize()
		v.list.SetSize((msg.Width-x)/2, msg.Height-y)
		v.vp.Width = msg.Width / 2
		v.vpReady = true
	}
	v.list, cmd = v.list.Update(msg)
	v.vp, vpCmd = v.vp.Update(msg)
	cmd = tea.Batch(cmd, vpCmd)
	v.vp.SetContent(v.RenderFilter())

	return v, cmd
}

func (v *FilterList) SetFilters(filters []proton.Filter) tea.Cmd {
	items := make([]list.Item, len(filters))
	for i, f := range filters {
		items[i] = FilterItem{f}
	}
	slices.SortFunc(items, func(a, b list.Item) int {
		return a.(FilterItem).Priority - b.(FilterItem).Priority
	})
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
