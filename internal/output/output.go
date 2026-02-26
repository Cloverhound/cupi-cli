package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// Print renders data in the specified format (json, table, csv, raw)
func Print(data interface{}, format string) error {
	switch format {
	case "json":
		return printJSON(data)
	case "table":
		return printTable(data)
	case "csv":
		return printCSV(data)
	case "raw":
		return printRaw(data)
	default:
		return fmt.Errorf("unknown output format: %s", format)
	}
}

func printJSON(data interface{}) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(b))
	return nil
}

func printTable(data interface{}) error {
	switch v := data.(type) {
	case []map[string]string:
		return printTableFromMaps(v)
	case []map[string]interface{}:
		return printTableFromInterfaceMaps(v)
	case map[string]interface{}:
		return printTableFromSingleMap(v)
	case string:
		fmt.Println(v)
		return nil
	default:
		b, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(b))
		return nil
	}
}

func printTableFromMaps(data []map[string]string) error {
	if len(data) == 0 {
		fmt.Println("No results")
		return nil
	}

	var columns []string
	for k := range data[0] {
		columns = append(columns, k)
	}
	sort.Strings(columns)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(columns)
	table.SetAutoFormatHeaders(true)
	table.SetAutoWrapText(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, row := range data {
		var values []string
		for _, col := range columns {
			values = append(values, row[col])
		}
		table.Append(values)
	}

	table.Render()
	return nil
}

func printTableFromInterfaceMaps(data []map[string]interface{}) error {
	var stringMaps []map[string]string
	for _, m := range data {
		sm := make(map[string]string)
		for k, v := range m {
			sm[k] = fmt.Sprintf("%v", v)
		}
		stringMaps = append(stringMaps, sm)
	}
	return printTableFromMaps(stringMaps)
}

func printTableFromSingleMap(data map[string]interface{}) error {
	rows := []map[string]string{}
	for k, v := range data {
		rows = append(rows, map[string]string{
			"key":   k,
			"value": fmt.Sprintf("%v", v),
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i]["key"] < rows[j]["key"]
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"key", "value"})
	table.SetAutoFormatHeaders(true)
	table.SetAutoWrapText(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, row := range rows {
		table.Append([]string{row["key"], row["value"]})
	}

	table.Render()
	return nil
}

func printCSV(data interface{}) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	switch v := data.(type) {
	case []map[string]string:
		if len(v) == 0 {
			return nil
		}

		var columns []string
		for k := range v[0] {
			columns = append(columns, k)
		}
		sort.Strings(columns)

		writer.Write(columns)

		for _, row := range v {
			var values []string
			for _, col := range columns {
				values = append(values, row[col])
			}
			writer.Write(values)
		}

	case []map[string]interface{}:
		if len(v) == 0 {
			return nil
		}

		var columns []string
		for k := range v[0] {
			columns = append(columns, k)
		}
		sort.Strings(columns)

		writer.Write(columns)

		for _, row := range v {
			var values []string
			for _, col := range columns {
				values = append(values, fmt.Sprintf("%v", row[col]))
			}
			writer.Write(values)
		}

	case map[string]interface{}:
		var columns []string
		for k := range v {
			columns = append(columns, k)
		}
		sort.Strings(columns)

		writer.Write(columns)
		var values []string
		for _, col := range columns {
			values = append(values, fmt.Sprintf("%v", v[col]))
		}
		writer.Write(values)

	default:
		return fmt.Errorf("unsupported data type for CSV")
	}

	return nil
}

func printRaw(data interface{}) error {
	switch v := data.(type) {
	case string:
		fmt.Println(v)
	case []string:
		for _, s := range v {
			fmt.Println(s)
		}
	default:
		fmt.Println(v)
	}
	return nil
}

// FormatKeyValue creates a formatted key-value display
func FormatKeyValue(data map[string]string) string {
	var lines []string
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		lines = append(lines, fmt.Sprintf("%s: %s", k, data[k]))
	}
	return strings.Join(lines, "\n")
}
