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

// The context exists outside of the bubble loop, and is being used as
// pointer reciever. That way we can use it as our primary storage
// container.
type Context struct {
	selectedColumn int
	showHelp       bool
	descending     bool
	data           *Tabdata

	// Window dimensions
	totalWidth  int
	totalHeight int

	// Table dimensions
	horizontalMargin int
	verticalMargin   int
}

// Execute tablizer sort function, feed it with fresh config, we do
// NOT use the existing runtime config, because sorting is
// configurable in the UI separately.
func (ctx *Context) Sort(mode string) {
	conf := cfg.Config{
		SortMode:        mode,
		SortDescending:  ctx.descending,
		UseSortByColumn: []int{ctx.selectedColumn + 1},
	}

	ctx.descending = !ctx.descending

	sortTable(conf, ctx.data)
}

// The actual table model, holds the context pointer, a copy of the
// pre-processed data and some flags
type FilterTable struct {
	Table table.Model

	Rows int

	quitting  bool
	unchanged bool

	maxColumns int
	headerIdx  map[string]int

	ctx *Context

	columns []table.Column
}

type HelpLine []string
type HelpColumn []HelpLine

const (
	// header+footer
	ExtraRows = 5

	HelpFooter = "?:help | "
)

var (
	// we use our own custom border style
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

	// Cells in selected columns will be highlighted
	StyleSelected = lipgloss.NewStyle().
			Background(lipgloss.Color("#696969")).
			Foreground(lipgloss.Color("#ffffff")).
			Align(lipgloss.Left)

	StyleHeader = lipgloss.NewStyle().
			Background(lipgloss.Color("#ffffff")).
			Foreground(lipgloss.Color("#696969")).
			Align(lipgloss.Left)

	// help buffer styles
	StyleKey  = lipgloss.NewStyle().Bold(true)
	StyleHelp = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff4500"))

	// the default style
	NoStyle = lipgloss.NewStyle().Align(lipgloss.Left)

	HelpData = []HelpColumn{
		{
			HelpLine{"up", "navigate up"},
			HelpLine{"down", "navigate down"},
			HelpLine{"tab", "navigate columns"},
		},
		{
			HelpLine{"s", "sort alpha-numerically"},
			HelpLine{"n", "sort numerically"},
			HelpLine{"t", "sort by time"},
			HelpLine{"d", "sort by duration"},
		},
		{
			HelpLine{"spc", "[de]select a row"},
			HelpLine{"a", "[de]select all visible rows"},
			HelpLine{"f", "enter fuzzy filter"},
			HelpLine{"esc", "finish filter input"},
		},
		{
			HelpLine{"?", "show help buffer"},
			HelpLine{"q", "commit and quit"},
			HelpLine{"c-c", "discard and quit"},
		},
	}

	// rendered from Help above
	Help = ""

	// number of lines taken by help below, adjust accordingly!
	HelpRows = 0
)

// generate  a  lipgloss  styled  help buffer  consisting  of  various
// columns
func generateHelp() {
	help := strings.Builder{}
	helpcols := []string{}
	maxrows := 0

	for _, col := range HelpData {
		help.Reset()

		// determine max key width to avoid excess spaces between keys and help
		keylen := 0
		for _, line := range col {
			if len(line[0]) > keylen {
				keylen = len(line[0])
			}
		}

		keylenstr := fmt.Sprintf("%d", keylen)

		for _, line := range col {
			// 0: key, 1: help text
			help.WriteString(StyleKey.Render(fmt.Sprintf("%-"+keylenstr+"s", line[0])))
			help.WriteString("  " + StyleHelp.Render(line[1]) + "   \n")
		}

		helpcols = append(helpcols, help.String())

		if len(col) > maxrows {
			maxrows = len(col)
		}
	}

	HelpRows = maxrows + 1
	Help = "\n" + lipgloss.JoinHorizontal(lipgloss.Top, helpcols...)
}

// initializes the table model
func NewModel(data *Tabdata, ctx *Context) FilterTable {
	columns := make([]table.Column, len(data.headers))
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

	// determine flexFactor with base 10, used by flexColumns
	for i, len := range lengths {
		if len <= 10 {
			lengths[i] = 1
		} else {
			lengths[i] = len / 10
		}
	}

	// setup column data with flexColumns
	for idx, header := range data.headers {
		// FIXME: doesn't work
		//columns[idx] = table.NewFlexColumn(strings.ToLower(header), StyleHeader.Render(header),
		columns[idx] = table.NewFlexColumn(strings.ToLower(header), header,
			lengths[idx]).WithFiltered(true).WithStyle(NoStyle)
	}

	// separate variable so we can share the row filling code
	filtertbl := FilterTable{
		maxColumns: len(data.headers),
		Rows:       len(data.entries),
		headerIdx:  hidx,
		ctx:        ctx,
		columns:    columns,
	}

	filtertbl.Table = table.New(columns)
	filtertbl.fillRows()

	// finally construct help buffer
	generateHelp()

	return filtertbl
}

// Applied to every cell on every change (TAB,up,down key, resize
// event etc)
func CellController(input table.StyledCellFuncInput, m FilterTable) lipgloss.Style {
	if m.headerIdx[input.Column.Key()] == m.ctx.selectedColumn {
		return StyleSelected
	}

	return NoStyle
}

// Selects or deselects ALL rows
func (m *FilterTable) ToggleAllSelected() {
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

// ? pressed, display help message
func (m FilterTable) ToggleHelp() {
	m.ctx.showHelp = !m.ctx.showHelp
}

func (m FilterTable) Init() tea.Cmd {
	return nil
}

// Forward call to context sort
func (m *FilterTable) Sort(mode string) {
	m.ctx.Sort(mode)
	m.fillRows()
}

// Fills the table rows with our data. Called once on startup and
// repeatedly if the user changes the sort order in some way
func (m *FilterTable) fillRows() {
	// required to be able to feed the model to the controller
	controllerWrapper := func(input table.StyledCellFuncInput) lipgloss.Style {
		return CellController(input, *m)
	}

	// fill the rows with style
	rows := make([]table.Row, len(m.ctx.data.entries))
	for idx, entry := range m.ctx.data.entries {
		rowdata := make(table.RowData, len(entry))

		for i, cell := range entry {
			rowdata[strings.ToLower(m.ctx.data.headers[i])] =
				table.NewStyledCellWithStyleFunc(cell+" ", controllerWrapper)
		}

		rows[idx] = table.NewRow(rowdata)
	}

	m.Table = m.Table.
		WithRows(rows).
		Filtered(true).
		WithFuzzyFilter().
		Focused(true).
		SelectableRows(true).
		WithSelectedText(" ", "✓").
		WithFooterVisibility(true).
		WithHeaderVisibility(true).
		HighlightStyle(StyleSelected).
		Border(customBorder)
}

// Part of the bubbletea event loop, called every tick
func (m FilterTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.Table, cmd = m.Table.Update(msg)
	cmds = append(cmds, cmd)

	// If the user is about to enter filter text, do NOT respond to
	// key bindings, as they might be part of the filter!
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
				m.ToggleAllSelected()

			case "tab":
				m.SelectNextColumn()

			case "?":
				m.ToggleHelp()
				m.recalculateTable()

			case "s":
				m.Sort("alphanumeric")

			case "n":
				m.Sort("numeric")

			case "d":
				m.Sort("duration")

			case "t":
				m.Sort("time")
			}
		}
	}

	// Happens when the terminal window has been resized
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ctx.totalWidth = msg.Width
		m.ctx.totalHeight = msg.Height

		m.recalculateTable()
	}

	m.updateFooter()

	return m, tea.Batch(cmds...)
}

// Add some info to the footer
func (m *FilterTable) updateFooter() {
	selected := m.Table.SelectedRows()
	footer := fmt.Sprintf("selected: %d ", len(selected))

	if m.Table.GetIsFilterInputFocused() {
		footer = fmt.Sprintf("/%s %s", m.Table.GetCurrentFilter(), footer)
	} else if m.Table.GetIsFilterActive() {
		footer = fmt.Sprintf("Filter: %s %s", m.Table.GetCurrentFilter(), footer)
	}

	m.Table = m.Table.WithStaticFooter(HelpFooter + footer)
}

// Called on resize event (or if help has been toggled)
func (m *FilterTable) recalculateTable() {
	m.Table = m.Table.
		WithTargetWidth(m.calculateWidth()).
		WithMinimumHeight(m.calculateHeight()).
		WithPageSize(m.calculateHeight() - ExtraRows)
}

func (m *FilterTable) calculateWidth() int {
	return m.ctx.totalWidth - m.ctx.horizontalMargin
}

// Take help height into account, if enabled
func (m *FilterTable) calculateHeight() int {
	height := m.Rows + ExtraRows

	if height >= m.ctx.totalHeight {
		height = m.ctx.totalHeight - m.ctx.verticalMargin
	} else {
		height = m.ctx.totalHeight
	}

	if m.ctx.showHelp {
		height = height - HelpRows
	}

	return height
}

// Part of the bubbletable event view, called every tick
func (m FilterTable) View() string {
	body := strings.Builder{}

	if !m.quitting {
		body.WriteString(m.Table.View())

		if m.ctx.showHelp {
			body.WriteString(Help)
		}
	}

	return body.String()
}

// User hit the TAB key
func (m *FilterTable) SelectNextColumn() {
	if m.ctx.selectedColumn == m.maxColumns-1 {
		m.ctx.selectedColumn = 0
	} else {
		m.ctx.selectedColumn++
	}
}

// entry point from outside tablizer into table editor
func tableEditor(conf *cfg.Config, data *Tabdata) (*Tabdata, error) {
	// we render to STDERR to avoid dead lock when the user redirects STDOUT
	// see https://github.com/charmbracelet/bubbletea/issues/860
	//
	// TODO: doesn't work with libgloss v2 anymore!
	lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(os.Stderr))

	ctx := &Context{data: data}

	// Output to  STDERR because  there's a  known bubbletea/lipgloss
	// issue: if  a program with a tui is  expected to write something
	// to STDOUT when the tui is finished, then the styles do not
	// work. So we write to STDERR (which works) and tablizer can
	// still be used inside pipes.
	program := tea.NewProgram(
		NewModel(data, ctx),
		tea.WithOutput(os.Stderr),
		tea.WithAltScreen())

	m, err := program.Run()

	if err != nil {
		return nil, err
	}

	if m.(FilterTable).unchanged {
		return data, err
	}

	// Data has been modified. Extract it, put it back into our own
	// structure and give control back to cmdline tablizer.
	filteredtable := m.(FilterTable)

	data.entries = make([][]string, len(filteredtable.Table.SelectedRows()))
	for pos, row := range m.(FilterTable).Table.SelectedRows() {
		entry := make([]string, len(data.headers))
		for idx, field := range data.headers {
			cell := row.Data[strings.ToLower(field)]
			switch value := cell.(type) {
			case string:
				entry[idx] = value
			case table.StyledCell:
				entry[idx] = value.Data.(string)
			}
		}

		data.entries[pos] = entry
	}

	return data, err
}
