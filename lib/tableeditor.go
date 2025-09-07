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

// This one exists outside of the bubble loop, and is referred to as
// pointer reciever. That way we can use it as our primary storage
// container.
type Container struct {
	selectedColumn int
	showHelp       bool
	descending     bool
	data           *Tabdata
}

func (container *Container) Sort(mode string) {
	conf := cfg.Config{
		SortMode:        mode,
		SortDescending:  container.descending,
		UseSortByColumn: []int{container.selectedColumn + 1},
	}

	container.descending = !container.descending

	sortTable(conf, container.data)
}

type FilterTable struct {
	Table table.Model

	// Window dimensions
	totalWidth  int
	totalHeight int

	// Table dimensions
	horizontalMargin int
	verticalMargin   int

	Rows int

	quitting  bool
	unchanged bool

	maxColumns int
	headerIdx  map[string]int

	container *Container

	columns []table.Column
}

const (
	// Add a fixed margin to account for description & instructions
	fixedVerticalMargin = 0

	ExtraRows = 5

	HELP = "?:help | "

	HelpRows = 7

	LongHelp = `
Key bindings usable in interactive filter table:
q         commit and quit           s     sort alphanumerically			  
ctrl-c    discard and quit          d     sort by duration
up|down   navigate rows             t     sort by time
TAB       navigage columns          f     fuzzy filter
space     [de]select row            ESC   finish filter input
a         [de]select all rows       ?     show help buffer
`
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

	StyleSelected = lipgloss.NewStyle().Background(lipgloss.Color("#696969")).Foreground(lipgloss.Color("#ffffff"))
	NoStyle       = lipgloss.NewStyle()
)

func NewModel(data *Tabdata, container *Container) FilterTable {
	columns := make([]table.Column, len(data.headers))
	rows := make([]table.Row, len(data.entries))
	lengths := make([]int, len(data.headers))
	hidx := make(map[string]int, len(data.headers))

	// give columns at least the header width
	for idx, header := range data.headers {
		lengths[idx] = len(header)
		hidx[strings.ToLower(header)] = idx
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

	filtertbl := FilterTable{
		horizontalMargin: 10,
		maxColumns:       len(data.headers),
		Rows:             len(data.entries),
		headerIdx:        hidx,
		container:        container,
		columns:          columns,
	}

	controllerWrapper := func(input table.StyledCellFuncInput) lipgloss.Style {
		return CellController(input, filtertbl)
	}

	// setup table data
	for idx, entry := range data.entries {
		rowdata := make(table.RowData, len(entry))

		for i, cell := range entry {
			rowdata[strings.ToLower(data.headers[i])] =
				table.NewStyledCellWithStyleFunc(cell+" ", controllerWrapper)
		}

		rows[idx] = table.NewRow(rowdata)
	}

	// keys := table.DefaultKeyMap()
	// keys.RowDown.SetKeys("j", "down", "s")
	// keys.RowUp.SetKeys("k", "up", "w")

	// our final interactive table filled with our prepared data
	filtertbl.Table = table.New(columns).
		WithRows(rows).
		//WithKeyMap(keys).
		Filtered(true).
		WithFuzzyFilter().
		Focused(true).
		SelectableRows(true).
		WithSelectedText(" ", "✓").
		WithFooterVisibility(true).
		WithHeaderVisibility(true).
		Border(customBorder)

	return filtertbl
}

func CellController(input table.StyledCellFuncInput, m FilterTable) lipgloss.Style {
	if m.headerIdx[input.Column.Key()] == m.container.selectedColumn {
		return StyleSelected
	}

	return NoStyle
}

func (m FilterTable) ToggleSelected() {
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

func (m FilterTable) ToggleHelp() {
	m.container.showHelp = !m.container.showHelp
}

func (m FilterTable) Init() tea.Cmd {
	return nil
}

// FIXME: possibly re-use it in NewModel() to avoid duplicate code
func (m *FilterTable) Sort(mode string) {
	m.container.Sort(mode)

	controllerWrapper := func(input table.StyledCellFuncInput) lipgloss.Style {
		return CellController(input, *m)
	}

	rows := make([]table.Row, len(m.container.data.entries))
	for idx, entry := range m.container.data.entries {
		rowdata := make(table.RowData, len(entry))

		for i, cell := range entry {
			rowdata[strings.ToLower(m.container.data.headers[i])] =
				table.NewStyledCellWithStyleFunc(cell+" ", controllerWrapper)
		}

		rows[idx] = table.NewRow(rowdata)
	}

	m.Table = table.New(m.columns).
		WithRows(rows).
		Filtered(true).
		WithFuzzyFilter().
		Focused(true).
		SelectableRows(true).
		WithSelectedText(" ", "✓").
		WithFooterVisibility(true).
		WithHeaderVisibility(true).
		Border(customBorder)

	//m.Table.WithRows(rows)
}

func (m FilterTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

			case "tab":
				m.SelectNextColumn()

			case "?":
				m.ToggleHelp()

			case "s":
				m.Sort("alphanumeric")

			case "n":
				m.Sort("numeric")

			case "d":
				m.Sort("duration")

			case "t":
				m.Sort("time")
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

func (m *FilterTable) updateFooter() {
	selected := m.Table.SelectedRows()
	footer := fmt.Sprintf("selected: %d", len(selected))

	if m.Table.GetIsFilterInputFocused() {
		footer = fmt.Sprintf("/%s %s", m.Table.GetCurrentFilter(), footer)
	} else if m.Table.GetIsFilterActive() {
		footer = fmt.Sprintf("Filter: %s %s", m.Table.GetCurrentFilter(), footer)
	}

	m.Table = m.Table.WithStaticFooter(HELP + footer)
}

func (m *FilterTable) recalculateTable() {
	m.Table = m.Table.
		WithTargetWidth(m.calculateWidth()).
		WithMinimumHeight(m.calculateHeight()).
		WithPageSize(m.calculateHeight() - ExtraRows)
}

func (m FilterTable) calculateWidth() int {
	return m.totalWidth - m.horizontalMargin
}

func (m FilterTable) calculateHeight() int {
	height := m.Rows + ExtraRows
	if m.container.showHelp {
		height = height + HelpRows
	}

	if height >= m.totalHeight {
		// FIXME: avoid full screen somehow
		height = m.totalHeight - m.verticalMargin - fixedVerticalMargin
	}

	return height
}

func (m FilterTable) View() string {
	body := strings.Builder{}

	if !m.quitting {
		body.WriteString(m.Table.View())

		if m.container.showHelp {
			body.WriteString(LongHelp)
		}
	}

	return body.String()
}

// FIXME: has no effect since FilterTable is being copied in Update()
// for the time being we're using a global variable. Maybe we can use
// the new GlobalMetadata field and store this kind of stuff there.
func (m *FilterTable) SelectNextColumn() {
	if m.container.selectedColumn == m.maxColumns-1 {
		m.container.selectedColumn = 0
	} else {
		m.container.selectedColumn++
	}
}

func tableEditor(conf *cfg.Config, data *Tabdata) (*Tabdata, error) {
	// we render to STDERR to avoid dead lock when the user redirects STDOUT
	// see https://github.com/charmbracelet/bubbletea/issues/860
	lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(os.Stderr))

	container := &Container{data: data}

	program := tea.NewProgram(
		NewModel(data, container),
		tea.WithOutput(os.Stderr),
		tea.WithAltScreen())

	m, err := program.Run()

	if err != nil {
		return nil, err
	}

	if m.(FilterTable).unchanged {
		return data, err
	}

	filteredtable := m.(FilterTable)

	data.entries = make([][]string, len(filteredtable.Table.SelectedRows()))
	for pos, row := range m.(FilterTable).Table.SelectedRows() {
		entry := make([]string, len(data.headers))
		for idx, field := range data.headers {
			cell := row.Data[strings.ToLower(field)]
			switch cell.(type) {
			case string:
				entry[idx] = cell.(string)
			case table.StyledCell:
				entry[idx] = cell.(table.StyledCell).Data.(string)
			}
		}

		data.entries[pos] = entry
	}

	return data, err
}
