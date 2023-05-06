module github.com/tlinden/tablizer

go 1.18

require (
	github.com/alecthomas/repr v0.1.1
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/gookit/color v1.5.2
	github.com/olekukonko/tablewriter v0.0.5
	github.com/spf13/cobra v1.6.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/lithammer/fuzzysearch v1.1.7 // indirect
	golang.org/x/text v0.8.0 // indirect
)

require (
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect

	// force release. > 0.4. doesnt build everywhere, see:
	// https://github.com/TLINDEN/tablizer/actions/runs/3396457307/jobs/5647544615
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.org/x/sys v0.5.0 // indirect
)
