package table

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/style"
	"github.com/charmbracelet/lipgloss"
)

func (o Options) Run() error {
	if len(o.Text) == 0 {
		input, _ := stdin.Read()
		if input == "" {
			return errors.New("no options provided, see `gum table --help`")
		}
		o.Text = strings.Split(strings.TrimSpace(input), "\n")
	}

	header := strings.Split(o.Text[0], o.Delimiter)
	var headerCol []table.Column
	for _, col := range header {
		// TODO: fix width calc
		headerCol = append(headerCol, table.Column{Title: col, Width: len(col)})
	}

	if len(headerCol) <= o.CellIndex {
		return errors.New("cell's index must be between csv columns limits, see `gum table --help`")
	}

	var rows []table.Row
	for _, r := range o.Text[1:] {
		rows = append(rows, strings.Split(r, o.Delimiter))
	}

	tbl := table.New(table.WithColumns(headerCol),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	tbl.SetStyles(s)

	tm, err := tea.NewProgram(model{
		table: tbl,
	}, tea.WithOutput(os.Stderr)).StartReturningModel()

	if err != nil {
		return fmt.Errorf("failed to start tea program: %w", err)
	}

	m := tm.(model)
	if m.aborted {
		return exit.ErrAborted
	}

	if o.CellIndex < 0 {
		fmt.Println(o.Text[m.table.Cursor()+1])
	} else {
		fmt.Println(o.Text[m.table.Cursor()+1][o.CellIndex])
	}

	return nil
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}
