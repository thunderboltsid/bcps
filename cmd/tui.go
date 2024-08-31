package cmd

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const maxWidth = 150

var (
	red    = lipgloss.AdaptiveColor{Light: "#FE5F86", Dark: "#FE5F86"}
	indigo = lipgloss.AdaptiveColor{Light: "#5A56E0", Dark: "#7571F9"}
	green  = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
)

type Styles struct {
	Base,
	HeaderText,
	TableText,
	Status,
	StatusHeader,
	Highlight,
	ErrorHeaderText,
	Help lipgloss.Style
}

func NewStyles(lg *lipgloss.Renderer) *Styles {
	s := Styles{}
	s.Base = lg.NewStyle().
		Padding(1, 4, 0, 1)
	s.HeaderText = lg.NewStyle().
		Foreground(indigo).
		Bold(true).
		Padding(0, 1, 0, 2)
	s.TableText = lg.NewStyle().
		Foreground(green).
		Padding(0, 1, 0, 2)
	s.Status = lg.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(indigo).
		PaddingLeft(1).
		MarginTop(1)
	s.StatusHeader = lg.NewStyle().
		Foreground(green).
		Bold(true)
	s.Highlight = lg.NewStyle().
		Foreground(lipgloss.Color("212"))
	s.ErrorHeaderText = s.HeaderText.Copy().
		Foreground(red)
	s.Help = lg.NewStyle().
		Foreground(lipgloss.Color("240"))
	return &s
}

const (
	borrowedSumKey                      = "borrowed-sum"
	borrowedYearKey                     = "borrowed-year"
	sharingPercentageKey                = "sharing-percentage"
	repaymentStartYearKey               = "repayment-start-year"
	startingSalaryKey                   = "starting-salary"
	expectedSalaryIncreasePercentageKey = "expected-salary-increase-percentage"
	expectedCPIIncreasePercentageKey    = "expected-cpi-increase-percentage"
	numberOfRepaymentYearsKey           = "number-of-repayment-years"
)

type Model struct {
	lg     *lipgloss.Renderer
	styles *Styles
	form   *huh.Form
	width  int
}

func NewModel() Model {
	m := Model{width: maxWidth}
	m.lg = lipgloss.DefaultRenderer()
	m.styles = NewStyles(m.lg)

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Key(borrowedYearKey).
				Options(huh.NewOptions(pastYears()...)...).
				Title("Choose the year you agreed to the ISA Contract").
				Description("This will determine the year you borrowed the sum"),

			huh.NewInput().
				Key(borrowedSumKey).
				Title("Enter the sum you borrowed").
				Description("This is the sum you borrowed from 🧠🤑").
				Validate(func(s string) error {
					if _, err := strconv.ParseFloat(s, 64); err != nil {
						return fmt.Errorf("invalid number")
					}

					return nil
				}),

			huh.NewInput().
				Key(startingSalaryKey).
				Title("Enter the gross salary you had the year you started repaying").
				Description("This is the the initial salary based on which the future salaries will be calculated").
				Validate(func(s string) error {
					if _, err := strconv.ParseFloat(s, 64); err != nil {
						return fmt.Errorf("invalid number")
					}

					return nil
				}),

			huh.NewInput().
				Key(sharingPercentageKey).
				Title("Enter the repayment rate").
				Description("This is the percentage of your gross salary you will repay to 🧠🤑 every year").
				Validate(func(s string) error {
					if _, err := strconv.ParseFloat(s, 64); err != nil {
						return fmt.Errorf("invalid number")
					}

					return nil
				}),

			huh.NewInput().
				Key(numberOfRepaymentYearsKey).
				Title("Enter the number of years you have to repay").
				Description("This is the number of years you have to share your income with 🧠🤑").
				Validate(func(s string) error {
					if _, err := strconv.Atoi(s); err != nil {
						return fmt.Errorf("invalid number")
					}

					return nil
				}),

			huh.NewSelect[int]().
				Key(repaymentStartYearKey).
				Options(huh.NewOptions(allYears()...)...).
				Title("Choose the year you started repaying").
				Description("This will be used as a starting point for the payment plan calculations").
				Validate(func(v int) error {
					if v < m.form.GetInt(borrowedYearKey) {
						return fmt.Errorf("repayment start year cannot be before the borrowed year")
					}
					return nil
				}),

			huh.NewInput().
				Key(expectedSalaryIncreasePercentageKey).
				Title("Enter the expected annual salary increase percentage").
				Description("This is the average percentage by which you expect your salary to increase every year").
				Validate(func(s string) error {
					if _, err := strconv.ParseFloat(s, 64); err != nil {
						return fmt.Errorf("invalid number")
					}

					return nil
				}),

			huh.NewInput().
				Key(expectedCPIIncreasePercentageKey).
				Title("Enter the average predicted CPI rate for future years").
				Description("This is the average percentage by which consumer price index (inflation) is expected to increase every year in the future").
				Validate(func(s string) error {
					if _, err := strconv.ParseFloat(s, 64); err != nil {
						return fmt.Errorf("invalid number")
					}

					return nil
				}),

			huh.NewConfirm().
				Key("done").
				Title("All done?").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("welp, finish up then")
					}
					return nil
				}).
				Affirmative("Yep").
				Negative("Wait, no"),
		),
	).
		WithWidth(120).
		WithShowHelp(false).
		WithShowErrors(false)
	return m
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = min(msg.Width, maxWidth) - m.styles.Base.GetHorizontalFrameSize()
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd

	// Process the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		// Quit when the form is done.
		cmds = append(cmds, tea.Quit)
	}

	return m, tea.Batch(cmds...)
}

func setFormValues(m Model) {
	borrowedYear = m.form.GetInt(borrowedYearKey)
	borrowedSum, _ = strconv.ParseFloat(m.form.GetString(borrowedSumKey), 64)
	startingSalary, _ = strconv.ParseFloat(m.form.GetString(startingSalaryKey), 64)
	sharingPercentage, _ = strconv.ParseFloat(m.form.GetString(sharingPercentageKey), 64)
	repaymentStartYear = m.form.GetInt(repaymentStartYearKey)
	expectedSalaryIncreasePercentage, _ = strconv.ParseFloat(m.form.GetString(expectedSalaryIncreasePercentageKey), 64)
	expectedCPIIncreasePercentage, _ = strconv.ParseFloat(m.form.GetString(expectedCPIIncreasePercentageKey), 64)
	repaymentYears, _ = strconv.Atoi(m.form.GetString(numberOfRepaymentYearsKey))
}

func (m Model) View() string {
	s := m.styles

	switch m.form.State {
	case huh.StateCompleted:
		var b strings.Builder
		fmt.Fprintf(&b, "\n🧠🤑 Income Sharing Agreement\n\n")

		paymentSchedule(&b, s)

		return s.Status.Copy().Margin(0, 1).Padding(1, 2).Width(80).Render(b.String()) + "\n\n"
	default:
		setFormValues(m)

		v := strings.TrimSuffix(m.form.View(), "\n\n")
		form := m.lg.NewStyle().Margin(1, 0).Render(v)

		errors := m.form.Errors()
		header := m.appBoundaryView("🧠🤑 Payment Schedule Calculator")
		if len(errors) > 0 {
			header = m.appErrorBoundaryView(m.errorView())
		}
		body := lipgloss.JoinHorizontal(lipgloss.Top, form)

		footer := m.appBoundaryView(m.form.Help().ShortHelpView(m.form.KeyBinds()))
		if len(errors) > 0 {
			footer = m.appErrorBoundaryView("")
		}

		return s.Base.Render(header + "\n" + body + "\n\n" + footer)
	}
}

func (m Model) errorView() string {
	var s string
	for _, err := range m.form.Errors() {
		s += err.Error()
	}
	return s
}

func (m Model) appBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.HeaderText.Render(text),
		lipgloss.WithWhitespaceChars(""),
		lipgloss.WithWhitespaceForeground(indigo),
	)
}

func (m Model) appErrorBoundaryView(text string) string {
	return lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Left,
		m.styles.ErrorHeaderText.Render(text),
		lipgloss.WithWhitespaceChars(""),
		lipgloss.WithWhitespaceForeground(red),
	)
}
