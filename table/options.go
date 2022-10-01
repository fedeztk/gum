package table

type Options struct {
	Text      []string `arg:"" optional:"" help:"CSV text to choose from."`
	Delimiter string   `default:"," help:"Delimiter of CSV text."`
	CellIndex int      `default:"-1" help:"Column cell's index to select from. Negative values select all cells."`
}
