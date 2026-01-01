package views

import (
	"bytes"
	"slices"

	"github.com/ProtonMail/go-proton-api"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maisieccino/proton-filters/internal/types"
	"github.com/muesli/reflow/wrap"
)

var listStyle = lipgloss.NewStyle().
	Margin(1, 2)

var docStyle = lipgloss.NewStyle()

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
		docStyle.Render(v.vp.View()),
	)
}

func (v *FilterList) RenderFilter() string {
	var current FilterItem
	if v.list.SelectedItem() != nil {
		current = v.list.SelectedItem().(FilterItem)
	}
	buf := new(bytes.Buffer)
	err := quick.Highlight(buf, current.Sieve, "sieve", "terminal16m", "catppuccin-macchiato")
	if err != nil {
		return current.Sieve
	}
	return wrap.String(buf.String(), v.vp.Width)
}

func (v *FilterList) Init() tea.Cmd {
	return nil
}

func (v *FilterList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case types.FiltersMsg:
		cmd = v.SetFilters(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		x, y := listStyle.GetFrameSize()
		v.list.SetSize((msg.Width/2)-x, msg.Height-y)
		v.list.SetWidth((msg.Width / 2) - x)
		if !v.vpReady {
			v.vp = viewport.New(msg.Width/2, msg.Height)
			v.vpReady = true
		} else {
			v.vp.Width = msg.Width / 2
			v.vp.Height = msg.Height
		}
	}
	v.list, cmd = v.list.Update(msg)
	cmds = append(cmds, cmd)
	v.vp.SetContent(v.RenderFilter())
	v.vp, cmd = v.vp.Update(msg)
	cmds = append(cmds, cmd)

	return v, tea.Batch(cmds...)
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
