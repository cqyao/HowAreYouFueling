package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ---------- Styles ----------
var (
	style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228")).
		BorderBackground(lipgloss.Color("63"))
)

// ---------- Bubble Tea Update ----------

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch m.Phase {

		case PhaseSelectBrands:
			switch msg.String() {
			case "q":
				return m, tea.Quit

			case "up", "k":
				if m.Cursor > 0 {
					m.Cursor--
				}

			case "down", "j":
				if m.Cursor < len(m.Brands)-1 {
					m.Cursor++
				}

			case " ":
				if _, ok := m.Selected[m.Cursor]; ok {
					delete(m.Selected, m.Cursor)
				} else {
					m.Selected[m.Cursor] = struct{}{}
				}

			case "enter":
				m.Phase = PhaseEnterFuelType
				m.Cursor = 0
			}

		case PhaseEnterFuelType:
			switch msg.String() {
			case "q":
				return m, tea.Quit

			case "up", "k":
				if m.Cursor > 0 {
					m.Cursor--
				}

			case "down", "j":
				if m.Cursor < 1 {
					m.Cursor++
				}

			case " ":
				if _, ok := m.Selected[m.Cursor]; ok {
					delete(m.Selected, m.Cursor)
				} else {
					m.Selected[m.Cursor] = struct{}{}
				}

			case "enter":
				m.Phase = PhaseDone
			}

		case PhaseDone:
			if msg.String() == "q" {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

// ---------- Bubble Tea View ----------

func (m Model) View() string {
	switch m.Phase {
	case PhaseSelectBrands:
		s := "Choose your service stations:\n\n"
		for i, brand := range m.Brands {
			cursor := " "
			if m.Cursor == i {
				cursor = ">"
			}
			checked := " "
			if _, ok := m.Selected[i]; ok {
				checked = "x"
			}
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, brand)
		}
		s += "\nPress SPACE to select, ENTER to continue, or Q to quit."
		return style.Render(s)

	case PhaseEnterFuelType:
		fuelTypes := []string{"U91", "E10"}
		s := "Choose your fuel type:\n\n"
		for i, f := range fuelTypes {
			cursor := " "
			if m.Cursor == i {
				cursor = ">"
			}
			checked := " "
			if _, ok := m.Selected[i]; ok {
				checked = "x"
			}
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, f)
		}
		s += "\nPress SPACE to select, ENTER to continue, or Q to quit."
		return style.Render(s)

	case PhaseDone:
		return style.Render("Selections complete! Press 'q' to quit and view results.")

	default:
		return style.Render("Unknown phase.")
	}
}
