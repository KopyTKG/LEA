package terminal

import (
	"fmt"
	"lea/help"
	"lea/state"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/kopytkg/golog"
	"golang.org/x/term"
)

type Fileln struct {
	Filename string // Filepath
	Current  int    // Done size
	Total    int    // Total size
	Done     bool
	Bar      []rune // Bar elements to render
}

func (f *Fileln) Update(newVal int) {
	f.Current = newVal
	per := (100 * newVal) / f.Total
	blocks := ((len(f.Bar) - 2) * per) / 100
	for i := 1; i < blocks; i++ {
		f.Bar[i] = '#'
	}
}

type Rendering struct {
	Files []*Fileln
	Total int
	Done  int
	Mutex sync.Mutex // For safe concurrent access to Files
	Run   func()
}

func (r *Rendering) AddFile(file *Fileln) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	(*r).Files = append((*r).Files, file)
}

func (r *Rendering) linuxUI() {
	ui := "\033[H\033[2J"

	logo := []string{
		" ___           _______       ________     ",
		"|\\  \\         |\\  ___ \\     |\\   __  \\    ",
		"\\ \\  \\        \\ \\   __/|    \\ \\  \\|\\  \\   ",
		" \\ \\  \\        \\ \\  \\_|/__   \\ \\   __  \\  ",
		"  \\ \\  \\____    \\ \\  \\_|\\ \\   \\ \\  \\ \\  \\ ",
		"   \\ \\_______\\   \\ \\_______\\   \\ \\__\\ \\__\\",
		"    \\|_______|    \\|_______|    \\|__|\\|__|",
		"                                          ",
	}

	for _, line := range logo {
		ui += line + "\n"
	}

	width := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		width = w
	}

	ui += fmt.Sprintf("by @KopyTKG %20s \n\n", help.VERSION)

	ui += fmt.Sprintf("\033[1m%s\033[0m | \033[1m%s\033[0m | KEY \033[1m%d\033[0m \n\n", state.Mode, strings.ToUpper(state.CYPHERMODE), state.Key.Metadata.KeyLength)

	ui += fmt.Sprintf("%-10s%2d/%-4d] \n\n", "Status  [", r.Done, r.Total)

	for _, f := range r.Files {
		if f.Done {
			continue
		}

		bar := ""
		for _, c := range f.Bar {
			bar += string(c)
		}

		padding := max(width-len(bar)-len(f.Filename), 0)

		line := fmt.Sprintf("%s%*s%s\n", f.Filename, padding, "", bar)

		ui += line
	}
	fmt.Println(ui)
}

func (r *Rendering) windowsUI() {
	exec.Command("powershell", "-Command", "Clear-Host").Run()
	exec.Command("cmd", "/c", "cls").Run()

	ui := ""

	logo := []string{
		" ___           _______       ________     ",
		"|\\  \\         |\\  ___ \\     |\\   __  \\    ",
		"\\ \\  \\        \\ \\   __/|    \\ \\  \\|\\  \\   ",
		" \\ \\  \\        \\ \\  \\_|/__   \\ \\   __  \\  ",
		"  \\ \\  \\____    \\ \\  \\_|\\ \\   \\ \\  \\ \\  \\ ",
		"   \\ \\_______\\   \\ \\_______\\   \\ \\__\\ \\__\\",
		"    \\|_______|    \\|_______|    \\|__|\\|__|",
		"                                          ",
	}

	for _, line := range logo {
		ui += line + "\n"
	}

	width := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		width = w
	}

	ui += fmt.Sprintf("by @KopyTKG %20s \n\n", help.VERSION)

	ui += fmt.Sprintf("%s | %s | KEY %d \n\n", state.Mode, strings.ToUpper(state.CYPHERMODE), state.Key.Metadata.KeyLength)

	ui += fmt.Sprintf("%-10s%2d/%-4d] \n\n", "Status  [", r.Done, r.Total)

	for _, f := range r.Files {
		if f.Done {
			continue
		}

		bar := ""
		for _, c := range f.Bar {
			bar += string(c)
		}

		padding := max(width-len(bar)-len(f.Filename), 0)

		line := fmt.Sprintf("%s%*s%s\n", f.Filename, padding, "", bar)

		ui += line
	}
	fmt.Println(ui)
}

func (r *Rendering) SetupOS() {
	switch runtime.GOOS {
	case "windows":
		r.Run = r.windowsUI
	case "linux":
		r.Run = r.linuxUI
	default:
		golog.Errorf("Unsupported OS (%s) found", runtime.GOOS)
		os.Exit(1)
	}

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
