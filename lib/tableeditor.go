/*
Copyright © 2025 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package lib

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/tlinden/tablizer/cfg"
)

type Model struct {
	Table table.Model

	// Window dimensions
	totalWidth  int
	totalHeight int

	// Table dimensions
	horizontalMargin int
	verticalMargin   int

	quitting  bool
	unchanged bool
}

const (
	minWidth  = 30
	minHeight = 1

	// Add a fixed margin to account for description & instructions
	fixedVerticalMargin = 0

	HELP = "/:filter esc:clear-filter q:commit c-c:abort space:select a:select-all | "
)

var (
	customBorder = table.Border{
		Top:    "─",
		Left:   "│",
		Right:  "│",
		Bottom: "─",

		TopRight:    "╮",
		TopLeft:     "╭",
		BottomRight: "╯",
		BottomLeft:  "╰",

		TopJunction:    "┬",
		LeftJunction:   "├",
		RightJunction:  "┤",
		BottomJunction: "┴",
		InnerJunction:  "┼",

		InnerDivider: "│",
	}
)

func NewModel(data *Tabdata) Model {
	columns := make([]table.Column, len(data.headers))
	rows := make([]table.Row, len(data.entries))
	lengths := make([]int, len(data.headers))

	// give columns at least the header width
	for idx, header := range data.headers {
		lengths[idx] = len(header)
	}

	// determine max width per column
	for _, entry := range data.entries {
		for i, cell := range entry {
			if len(cell) > lengths[i] {
				lengths[i] = len(cell)
			}
		}
	}

	// setup column data
	for idx, header := range data.headers {
		columns[idx] = table.NewColumn(strings.ToLower(header), header, lengths[idx]+2).
			WithFiltered(true)
	}

	// setup table data
	for idx, entry := range data.entries {
		rowdata := make(table.RowData, len(entry))

		for i, cell := range entry {
			rowdata[strings.ToLower(data.headers[i])] = cell + " "
		}

		rows[idx] = table.NewRow(rowdata)
	}

	keys := table.DefaultKeyMap()
	keys.RowDown.SetKeys("j", "down", "s")
	keys.RowUp.SetKeys("k", "up", "w")

	// our final interactive table filled with our prepared data
	return Model{
		Table: table.New(columns).
			WithRows(rows).
			WithKeyMap(keys).
			Filtered(true).
			Focused(true).
			SelectableRows(true).
			WithSelectedText(" ", "✓").
			WithFooterVisibility(true).
			WithHeaderVisibility(true).
			WithMaxTotalWidth(150).
			WithPageSize(20).
			Border(customBorder),
		horizontalMargin: 10,
	}
}

func (m Model) ToggleSelected() {
	rows := m.Table.GetVisibleRows()
	selected := m.Table.SelectedRows()

	if len(selected) > 0 {
		for i, row := range selected {
			rows[i] = row.Selected(false)
		}
	} else {
		for i, row := range rows {
			rows[i] = row.Selected(true)
		}
	}

	m.Table.WithRows(rows)
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.Table, cmd = m.Table.Update(msg)
	cmds = append(cmds, cmd)

	if !m.Table.GetIsFilterInputFocused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q":
				m.quitting = true
				m.unchanged = false
				cmds = append(cmds, tea.Quit)

			case "ctrl+c":
				m.quitting = true
				m.unchanged = true
				cmds = append(cmds, tea.Quit)

			case "a":
				m.ToggleSelected()

			case "c":
				m.Table.WithFilterInputValue("")
			}
		case tea.WindowSizeMsg:
			m.totalWidth = msg.Width
			m.totalHeight = msg.Height

			m.recalculateTable()
		}
	}
	m.updateFooter()

	return m, tea.Batch(cmds...)
}

func (m *Model) updateFooter() {
	selected := m.Table.SelectedRows()
	footer := fmt.Sprintf("selected: %d", len(selected))

	if m.Table.GetIsFilterInputFocused() {
		footer = fmt.Sprintf("/%s %s", m.Table.GetCurrentFilter(), footer)
	} else if m.Table.GetIsFilterActive() {
		footer = fmt.Sprintf("Filter: %s %s", m.Table.GetCurrentFilter(), footer)
	}

	m.Table = m.Table.WithStaticFooter(HELP + footer)
}

func (m *Model) recalculateTable() {
	m.Table = m.Table.
		WithTargetWidth(m.calculateWidth()).
		WithMinimumHeight(m.calculateHeight()).
		WithPageSize(m.calculateHeight() - 8)
}

func (m Model) calculateWidth() int {
	return m.totalWidth - m.horizontalMargin
}

func (m Model) calculateHeight() int {
	return m.totalHeight - m.verticalMargin - fixedVerticalMargin
}

func (m Model) View() string {
	body := strings.Builder{}

	if !m.quitting {
		body.WriteString(m.Table.View())
	}

	return body.String()
}

func tableEditor(conf *cfg.Config, data *Tabdata) (*Tabdata, error) {
	// we render to STDERR to avoid dead lock when the user redirects STDOUT
	// see https://github.com/charmbracelet/bubbletea/issues/860
	lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(os.Stderr))
	program := tea.NewProgram(
		NewModel(data),
		tea.WithOutput(os.Stderr),
		tea.WithAltScreen())

	m, err := program.Run()

	if err != nil {
		return nil, err
	}

	if m.(Model).unchanged {
		return data, err
	}

	table := m.(Model).Table
	data.entries = make([][]string, len(table.SelectedRows()))

	for pos, row := range m.(Model).Table.SelectedRows() {
		entry := make([]string, len(data.headers))
		for idx, field := range data.headers {
			entry[idx] = row.Data[strings.ToLower(field)].(string)
		}

		data.entries[pos] = entry
	}

	return data, err
}
