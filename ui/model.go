package ui

import tea "github.com/charmbracelet/bubbletea"

type Phase int

const (
	PhaseSelectBrands Phase = iota
	PhaseEnterPostcode
	PhaseEnterFuelType
	PhaseDone
)

type Model struct {
	Phase    Phase
	Brands   []string
	Cursor   int
	Selected map[int]struct{}
	FuelType string
	Postcode string
	InputBuf string
}

func InitialModel() Model {
	return Model{
		Phase:    PhaseSelectBrands,
		Brands:   []string{"BP", "Ampol", "Metro Fuel", "EG Ampol", "Enhance", "Shell"},
		Selected: make(map[int]struct{}),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
