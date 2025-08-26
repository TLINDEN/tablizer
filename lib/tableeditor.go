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
	"strings"

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

	id string
}

const (
	minWidth  = 30
	minHeight = 1

	// Add a fixed margin to account for description & instructions
	fixedVerticalMargin = 0

	columnKeyID          = "id"
	columnKeyName        = "name"
	columnKeyDescription = "description"
	columnKeyCount       = "count"
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
	var id string

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
		if id == "" {
			id = strings.ToLower(header)
		}
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
		id:               id,
		horizontalMargin: 10,
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
		case "ctrl+c", "q", "esc":
			// FIXME: feed the reprocessed data to printData(), then call tea.Quit
			// FIXME: we need to reset the screen somehow, otherwise printing doesn't work.
			fmt.Println()
			for _, row := range m.Table.SelectedRows() {
				// selectedIDs = append(selectedIDs, row.Data[m.id].(string))
				//repr.Println(row.Data)
				fmt.Printf("Selected: %s\n", row.Data[m.id].(string))
			}

			cmds = append(cmds, tea.Quit)
		}
	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width
		m.totalHeight = msg.Height

		m.recalculateTable()
	}

	m.updateFooter()

	return m, tea.Batch(cmds...)
}

func (m *Model) updateFooter() {
	selected := m.Table.SelectedRows()
	footer := fmt.Sprintf("selected: %d", len(selected))

	// highlightedRow := m.Table.HighlightedRow()

	// footerText := fmt.Sprintf(
	// 	"Pg. %d/%d - Currently looking at ID: %s",
	// 	m.Table.CurrentPage(),
	// 	m.Table.MaxPages(),
	// 	highlightedRow.Data[m.id],
	// )

	m.Table = m.Table.WithStaticFooter(footer)
}

func (m *Model) recalculateTable() {
	m.Table = m.Table.
		WithTargetWidth(m.calculateWidth()).
		WithMinimumHeight(m.calculateHeight()).
		WithPageSize(m.calculateHeight() - 7)
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
