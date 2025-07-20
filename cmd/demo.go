/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	spinner  spinner.Model
	loading  bool
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.title)
				return m, performAction(m.choice)
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case actionCompleteMsg:
		m.loading = false
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "\n  Goodbye!\n"
	}

	if m.loading {
		return fmt.Sprintf("\n\n   %s Loading %s...\n\n", m.spinner.View(), m.choice)
	}

	if m.choice != "" {
		return fmt.Sprintf("\n  ✓ %s completed successfully!\n\n  Press any key to return to menu or Ctrl+C to exit.\n", m.choice)
	}

	return "\n" + m.list.View()
}

type actionCompleteMsg struct{}

func performAction(action string) tea.Cmd {
	return tea.Tick(time.Millisecond*2000, func(t time.Time) tea.Msg {
		switch action {
		case "Database Demo":
			// Demo database operations
			db, err := gorm.Open(sqlite.Open("demo.db"), &gorm.Config{})
			if err != nil {
				log.Printf("Failed to connect database: %v", err)
			} else {
				type User struct {
					ID   uint `gorm:"primaryKey"`
					Name string
				}
				db.AutoMigrate(&User{})
				db.Create(&User{Name: "Demo User"})
				log.Println("Database demo completed")
			}
		case "SSH Demo":
			log.Println("SSH demo - connection simulation completed")
		case "Form Demo":
			// This would show a huh form in a real implementation
			log.Println("Form demo completed")
		}
		return actionCompleteMsg{}
	})
}

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Interactive TUI demo showcasing Bubble Tea, Huh, Lipgloss, and database features",
	Long: `An interactive demonstration of the TUI capabilities including:
- Bubble Tea terminal UI framework
- Lipgloss styling
- Database operations with GORM
- SSH connection simulation
- Interactive forms with Huh`,
	Run: func(cmd *cobra.Command, args []string) {
		// First show a huh form demo
		showFormDemo()

		// Then show the main TUI demo
		items := []list.Item{
			item{title: "Database Demo", desc: "Test SQLite database operations with GORM"},
			item{title: "SSH Demo", desc: "Simulate SSH connection handling"},
			item{title: "Form Demo", desc: "Show interactive form capabilities"},
			item{title: "Styling Demo", desc: "Demonstrate Lipgloss styling features"},
		}

		const defaultWidth = 20
		const listHeight = 14

		l := list.New(items, list.NewDefaultDelegate(), defaultWidth, listHeight)
		l.Title = "EHRPlus CLI Demo"
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.Styles.Title = titleStyle
		l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
		l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)

		s := spinner.New()
		s.Spinner = spinner.Dot
		s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

		m := model{
			list:    l,
			spinner: s,
		}

		if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
			fmt.Printf("Error running program: %v", err)
			os.Exit(1)
		}
	},
}

func showFormDemo() {
	var name string
	var environment string
	var confirm bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What's your name?").
				Value(&name).
				Placeholder("Enter your name"),

			huh.NewSelect[string]().
				Title("Choose environment").
				Options(
					huh.NewOption("Development", "dev"),
					huh.NewOption("Staging", "staging"),
					huh.NewOption("Production", "prod"),
				).
				Value(&environment),

			huh.NewConfirm().
				Title("Continue with demo?").
				Value(&confirm),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	if !confirm {
		fmt.Println("Demo cancelled.")
		return
	}

	fmt.Printf("Hello %s! Running demo in %s environment.\n\n", name, environment)
}

func init() {
	rootCmd.AddCommand(demoCmd)
}
