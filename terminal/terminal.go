package terminal

import (
	"lea/help"
	"strconv"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Fileln struct {
	FP    string
	Done  int
	Total int
	Bar   []rune
}

func (f *Fileln) Update(newVal int) {
	f.Done = newVal
	per := (100 * newVal) / f.Total
	blocks := ((len(f.Bar) - 2) * per) / 100

	for i := int(1); i < blocks; i++ {
		f.Bar[i] = '#'
	}
}

type Rendering struct {
	File *Fileln
}

func (r *Rendering) Run() {
	logo := []string{
		" ___           _______       ________     ",
		"|\\  \\         |\\  ___ \\     |\\   __  \\    ",
		"\\ \\  \\        \\ \\   __/|    \\ \\  \\|\\  \\   ",
		" \\ \\  \\        \\ \\  \\_|/__   \\ \\   __  \\  ",
		"  \\ \\  \\____    \\ \\  \\_|\\ \\   \\ \\  \\ \\  \\ ",
		"   \\ \\_______\\   \\ \\_______\\   \\ \\__\\ \\__\\",
		"    \\|_______|    \\|_______|    \\|__|\\|__|",
		"                                          ",
		"by @KopyTKG " + help.VERSION,
	}
	g := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	g.SetRect(0, termHeight/4-1, termWidth-1, (termHeight/4)*3-1)
	g.Border = true

	g1 := ui.NewGrid()
	g1.SetRect(0, 0, termWidth-1, termHeight/4)
	g1.Border = true
	logoText := widgets.NewParagraph()

	logoText.Text = ""
	for _, line := range logo {
		logoText.Text += line + "\n"
	}
	logoText.Border = false

	fp := widgets.NewParagraph()
	fp.Text = r.File.FP
	fp.SetRect(0, 0, (termWidth-1)/4, 1)
	fp.Border = false

	bar := widgets.NewParagraph()
	bar.Text = string(r.File.Bar)
	bar.SetRect(0, 0, (termWidth-1)/2, 1)
	bar.Border = false

	per := widgets.NewParagraph()
	percentage := (100 * r.File.Done) / r.File.Total
	per.Text = (strconv.Itoa(percentage) + "%")
	per.SetRect(0, 0, (termWidth-1)/4, 1)
	per.Border = false

	g1.Set(
		ui.NewRow(1.0, ui.NewCol(1.0, logoText)),
	)

	row := ui.NewRow(1.0,
		ui.NewCol(0.25, fp),
		ui.NewCol(0.5, bar),
		ui.NewCol(0.25, per),
	)

	g.Set(row)
	ui.Render(g1, g)
}

func push(arr *[]rune, item rune) {
	*arr = append(*arr, item)
}

func BarSetup(w int) []rune {
	progress := []rune("")
	push(&progress, '[')
	for i := 1; i < w; i++ {
		push(&progress, ' ')
	}

	push(&progress, ']')
	push(&progress, ' ')

	return progress
}
