package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ProtonMail/go-proton-api"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maisieccino/proton-filters/internal/client"
	"github.com/maisieccino/proton-filters/internal/types"
	"github.com/maisieccino/proton-filters/internal/views"
	"github.com/spf13/cobra"
)

type appModel struct {
	client  *proton.Client
	screens map[string]tea.Model
	current string
}

func newApp() appModel {
	ctx := context.Background()
	c, err := client.NewClient(ctx)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	m := appModel{
		client:  c,
		screens: make(map[string]tea.Model),
	}
	v, name := views.NewFilterList()
	m.screens[name] = v
	m.current = name
	return m
}

type FiltersMsg []proton.Filter

func (m appModel) currentScreen() tea.Model {
	if cur, ok := m.screens[m.current]; ok {
		return cur
	}
	return &views.DefaultScreen{}
}

func (m appModel) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg {
		filters, err := m.client.GetAllFilters(context.Background())
		if err != nil {
			return err
		}
		return types.FiltersMsg(filters)
	})
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	}
	return m.currentScreen().Update(msg)
}

func (m appModel) View() string {
	return m.currentScreen().View()
}

func Root(cmd *cobra.Command, args []string) {
	p := tea.NewProgram(newApp())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
