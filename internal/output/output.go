// internal/output/output.go
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"ratkiez/internal/types"

	"github.com/olekukonko/tablewriter"
)

// Formatter handles different output formats
type Formatter struct {
	format string
}

// NewFormatter creates a new output formatter
func NewFormatter(format string) (*Formatter, error) {
	return &Formatter{
		format: format,
	}, nil
}

// Print outputs the data in the specified format
func (f *Formatter) Print(data types.KeyDetailsSlice) error {
	switch f.format {
	case "table":
		return f.printTable(data)
	case "json":
		return f.printJSON(data)
	case "csv":
		return f.printCSV(data)
	default:
		return fmt.Errorf("unsupported output format: %s", f.format)
	}
}

func (f *Formatter) printTable(data types.KeyDetailsSlice) error {
	headers := []string{"USERNAME", "KEY-ID", "CREATION-DATE", "LAST-USED-DATE", "POLICIES", "PROFILE"}
	var rows [][]string

	for _, v := range data {
		rows = append(rows, []string{
			v.User,
			v.KeyID,
			v.CreationDate,
			v.LastUsedDate,
			strings.Join(v.Policies, ", "),
			v.Profile,
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(rows)
	table.Render()

	return nil
}

func (f *Formatter) printJSON(data types.KeyDetailsSlice) error {
	output, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	fmt.Println(string(output))
	return nil
}

func (f *Formatter) printCSV(data types.KeyDetailsSlice) error {

	// Write headers
	headers := []string{"USERNAME", "KEY-ID", "CREATION-DATE", "LAST-USED-DATE", "POLICIES", "PROFILE"}
	fmt.Println(strings.Join(headers, ","))

	// Write data
	for _, v := range data {
		record := []string{
			v.User,
			v.KeyID,
			v.CreationDate,
			v.LastUsedDate,
			strings.Join(v.Policies, "; "),
			v.Profile,
		}
		fmt.Println(strings.Join(record, ","))
	}

	return nil
}
