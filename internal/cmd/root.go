package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ProtonMail/go-proton-api"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maisieccino/proton-filters/internal/client"
	"github.com/maisieccino/proton-filters/internal/views"
	"github.com/spf13/cobra"
)

type appModel struct {
	loading bool
	spinner spinner.Model
	client  *proton.Client
	err     error

	view_FilterList tea.Model
}

func newApp() appModel {
	ctx := context.Background()
	c, err := client.NewClient(ctx)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return appModel{
		loading:         true,
		client:          c,
		spinner:         spinner.New(),
		view_FilterList: &views.FilterList{},
	}
}

type FiltersMsg []proton.Filter

func (m appModel) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg {
		filters, err := m.client.GetAllFilters(context.Background())
		if err != nil {
			return err
		}
		return FiltersMsg(filters)
	},
		m.spinner.Tick,
	)
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "a":
			m.loading = false
			return m, nil
		}
	case error:
		m.err = msg
		return m, nil
	case FiltersMsg:
		m.loading = false
		m.view_FilterList, cmd = m.view_FilterList.Update(msg)
		return m, cmd
	}
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m appModel) View() string {
	if m.loading {
		return m.spinner.View()
	}
	if m.err != nil {
		return "Error fetching filters: " + m.err.Error()
	}
	if m.view_FilterList != nil {
		return m.view_FilterList.View()
	}
	return ""
}

func Root(cmd *cobra.Command, args []string) {
	p := tea.NewProgram(newApp())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
