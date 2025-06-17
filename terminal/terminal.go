package terminal

import (
	"strconv"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Fileln struct {
	FP    string // Filepath
	Done  int    // Done size
	Total int    // Total size
	Bar   []rune // Bar elements to render
}

func (f *Fileln) Update(newVal int) {
	f.Done = newVal
	per := (100 * newVal) / f.Total
	blocks := ((len(f.Bar) - 2) * per) / 100
	for i := 1; i < blocks; i++ {
		f.Bar[i] = '#'
	}
}

type Rendering struct {
	Files *[]Fileln
	Mutex sync.Mutex // For safe concurrent access to Files
}

func (r *Rendering) AddFile(file Fileln) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	*r.Files = append(*r.Files, file)
}

func (r *Rendering) UpdateFileProgress(filepath string, done int) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	for i, f := range *r.Files {
		if f.FP == filepath {
			(*r.Files)[i].Update(done)
			break
		}
	}
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
		"by @KopyTKG " + "VERSION", // Replace with your version variable
	}
	g := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	fileGridHeight := termHeight / 2
	r.Mutex.Lock()
	if len(*r.Files) > 0 {
		fileGridHeight = termHeight/4 + len(*r.Files)*2
		if fileGridHeight > termHeight-2 {
			fileGridHeight = termHeight - 2
		}
	}
	r.Mutex.Unlock()
	g.SetRect(0, termHeight/4, termWidth-1, fileGridHeight)
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
	g1.Set(ui.NewRow(1.0, ui.NewCol(1.0, logoText)))

	r.Mutex.Lock()
	var fileRows []interface{}
	for _, f := range *r.Files {
		fp := widgets.NewParagraph()
		fp.Text = f.FP
		fp.Border = false

		bar := widgets.NewParagraph()
		bar.Text = string(f.Bar)
		bar.Border = false

		per := widgets.NewParagraph()
		percentage := 0
		if f.Total > 0 {
			percentage = (100 * f.Done) / f.Total
		}
		per.Text = strconv.Itoa(percentage) + "%"
		per.Border = false

		row := ui.NewRow(1.0/float64(len(*r.Files)+1),
			ui.NewCol(0.25, fp),
			ui.NewCol(0.5, bar),
			ui.NewCol(0.25, per),
		)
		fileRows = append(fileRows, row)
	}
	r.Mutex.Unlock()

	g.Set(fileRows...)
	ui.Render(g1, g)
}

func BarSetup(w int) []rune {
	progress := make([]rune, 0, w+2)
	progress = append(progress, '[')
	for i := 0; i < w; i++ {
		progress = append(progress, ' ')
	}
	progress = append(progress, ']')
	progress = append(progress, ' ')
	return progress
}
