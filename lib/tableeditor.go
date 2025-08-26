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
	"strings"

	"github.com/alecthomas/repr"
	tea "github.com/charmbracelet/bubbletea"
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
}

const (
	minWidth  = 30
	minHeight = 1

	// Add a fixed margin to account for description & instructions
	fixedVerticalMargin = 0
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
		columns[idx] = table.NewColumn(header, header, lengths[idx]+2).WithFiltered(true)
	}

	// setup table data
	for idx, entry := range data.entries {
		rowdata := make(table.RowData, len(entry))

		for i, cell := range entry {
			rowdata[data.headers[i]] = cell + " "
		}

		rows[idx] = table.Row{Data: rowdata}
	}

	// our final interactive table filled with our prepared data
	return Model{
		Table: table.
			New(columns).
			Filtered(true).
			Focused(true).
			SelectableRows(true).
			WithSelectedText(" ", "✓").
			WithFooterVisibility(true).
			WithHeaderVisibility(true).
			WithMaxTotalWidth(150).
			Border(customBorder).
			WithRows(rows),
	}
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// FIXME: feed the reprocessed data to printData(), then call tea.Quit
			for _, row := range m.Table.SelectedRows() {
				repr.Println(row.Data)
			}

			cmds = append(cmds, tea.Quit)
		case "left":
			if m.calculateWidth() > minWidth {
				m.horizontalMargin++
				m.recalculateTable()
			}

		case "right":
			if m.horizontalMargin > 0 {
				m.horizontalMargin--
				m.recalculateTable()
			}

		case "up":
			if m.calculateHeight() > minHeight {
				m.verticalMargin++
				m.recalculateTable()
			}

		case "down":
			if m.verticalMargin > 0 {
				m.verticalMargin--
				m.recalculateTable()
			}
		}
	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width
		m.totalHeight = msg.Height

		m.recalculateTable()
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) recalculateTable() {
	m.Table = m.Table.
		WithTargetWidth(m.calculateWidth()).
		WithMinimumHeight(m.calculateHeight())
}

func (m Model) calculateWidth() int {
	return m.totalWidth - m.horizontalMargin
}

func (m Model) calculateHeight() int {
	return m.totalHeight - m.verticalMargin - fixedVerticalMargin
}

func (m Model) View() string {
	body := strings.Builder{}

	body.WriteString(m.Table.View())

	return body.String()
}

func tableEditor(conf *cfg.Config, data *Tabdata) error {
	program := tea.NewProgram(NewModel(data))

	_, err := program.Run()

	return err
}
