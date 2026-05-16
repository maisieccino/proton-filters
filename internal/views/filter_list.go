package views

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"slices"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/ProtonMail/go-proton-api"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/maisieccino/proton-filters/internal/filter"
	"github.com/maisieccino/proton-filters/internal/types"
	"github.com/muesli/reflow/wrap"
)

var (
	listStyle = lipgloss.NewStyle().
			Margin(1, 2)
	docStyle   = lipgloss.NewStyle()
	styleRed   = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(1))
	styleGreen = lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(10))
)

type FilterList struct {
	list    list.Model
	client  *proton.Client
	vp      viewport.Model
	vpReady bool
	editing bool
}

var _ list.Item = &FilterItem{}

func NewFilterList(client *proton.Client) (tea.Model, string) {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 1, 1)
	l.Title = "Sieve filters"
	l.SetStatusBarItemName("filter", "filters")
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("a"),
				key.WithHelp("a", "toggle filter"),
			),
			key.NewBinding(
				key.WithKeys("e"),
				key.WithHelp("e", "edit filter"),
			),
		}
	}
	return &FilterList{
		list:   l,
		client: client,
	}, "filter-list"
}

func (v *FilterList) View() tea.View {
	view := tea.NewView(lipgloss.JoinHorizontal(lipgloss.Top,
		v.list.View(),
		docStyle.Render(v.vp.View()),
	))
	view.WindowTitle = v.list.Title
	view.AltScreen = true
	return view
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
	return wrap.String(buf.String(), v.vp.Width())
}

func (v *FilterList) Init() tea.Cmd {
	return nil
}

func (v *FilterList) CurrentFilter() (FilterItem, bool) {
	if v.list.SelectedItem() != nil {
		return v.list.SelectedItem().(FilterItem), true
	}
	return FilterItem{}, false
}

func (v *FilterList) ToggleFilter(filter proton.Filter) tea.Cmd {
	return func() tea.Msg {
		// TODO: Cancellable context.
		ctx := context.Background()
		fmt.Fprintf(os.Stderr, "Status: %d\n", filter.Status)
		if filter.Status == 1 {
			if err := v.client.DisableFilter(ctx, filter.ID); err != nil {
				return err
			}
		} else {
			if err := v.client.EnableFilter(ctx, filter.ID); err != nil {
				return err
			}
		}
		return types.ToggleMsg{FilterID: filter.ID}
	}
}

func (v *FilterList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if v.editing {
		return v, nil
	}
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case error:
		cmd = v.list.NewStatusMessage(fmt.Sprintf("Error: %s", msg.Error()))
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			if f, ok := v.CurrentFilter(); ok {
				cmds = append(cmds, v.ToggleFilter(f.Filter))
			}
		case "e":
			if f, ok := v.CurrentFilter(); ok {
				editCmd := editFilter(f)
				v.editing = true
				cmds = append(cmds, editCmd)
			}
		}
	case types.ToggleMsg:
		cmd = v.list.NewStatusMessage(fmt.Sprintf("Filter %s toggled", msg.FilterID))
		cmds = append(cmds, cmd)
		items := []proton.Filter{}
		for _, i := range v.list.Items() {
			items = append(items, i.(FilterItem).Filter)
		}
		idx := slices.IndexFunc(items, func(f proton.Filter) bool {
			return f.ID == msg.FilterID
		})
		if idx > 0 {
			switch items[idx].Status {
			case 1:
				items[idx].Status = 0
			case 0:
				items[idx].Status = 1
			}
		}
		cmd = v.SetFilters(items)
		cmds = append(cmds, cmd)
	case types.FiltersMsg:
		cmd = v.SetFilters(msg)
		cmds = append(cmds, cmd)
	case types.FilterMsg:
		v.editing = false
		// Find idx of filter.
		i := slices.IndexFunc(v.list.Items(), func(i list.Item) bool {
			if f, ok := i.(FilterItem); ok {
				return f.ID == msg.ID
			}
			return false
		})
		cmds = append(cmds, v.list.SetItem(i, FilterItem{proton.Filter(msg)}))
	case tea.WindowSizeMsg:
		x, y := listStyle.GetFrameSize()
		v.list.SetSize((msg.Width/2)-x, msg.Height-y)
		v.list.SetWidth((msg.Width / 2) - x)
		if !v.vpReady {
			v.vp = viewport.New(
				viewport.WithWidth(msg.Width/2),
				viewport.WithWidth(msg.Height),
			)
			v.vpReady = true
		} else {
			v.vp.SetWidth(msg.Width / 2)
			v.vp.SetHeight(msg.Height)
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
	dot := "• "
	switch i.Status {
	case 1:
		return styleGreen.Render(dot) + "Enabled"
	default:
		return styleRed.Render(dot) + "Disabled"
	}
}

func (i FilterItem) FilterValue() string {
	return i.Name
}

func editFilter(f FilterItem) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		newFilter, err := filter.EditFilter(ctx, f.Filter)
		if err != nil {
			return nil
		}
		return types.FilterMsg(newFilter)
	}
}
