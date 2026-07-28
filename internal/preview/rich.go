package preview

import (
	"encoding/csv"
	"fmt"
	"strings"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorCyan   = "\033[36m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBgGray = "\033[47;30m"
)

// FormatMarkdown styles simple markdown syntax.
func FormatMarkdown(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	
	inCodeBlock := false
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			if inCodeBlock {
				result = append(result, ColorYellow+"--- Code Block Start ---"+ColorReset)
			} else {
				result = append(result, ColorYellow+"--- Code Block End ---"+ColorReset)
			}
			continue
		}
		
		if inCodeBlock {
			result = append(result, ColorBgGray+line+ColorReset)
			continue
		}
		
		if strings.HasPrefix(trimmed, "#") {
			result = append(result, ColorBold+ColorCyan+line+ColorReset)
		} else if strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "- ") {
			result = append(result, ColorGreen+line[:2]+ColorReset+line[2:])
		} else if strings.Contains(line, "`") {
			parts := strings.Split(line, "`")
			for i := 1; i < len(parts); i += 2 {
				parts[i] = ColorBgGray + parts[i] + ColorReset
			}
			result = append(result, strings.Join(parts, ""))
		} else {
			result = append(result, line)
		}
	}
	
	return strings.Join(result, "\n")
}

// FormatCSV parses CSV and formats it into a zebra-striped table grid.
func FormatCSV(content string, maxWidth int) string {
	r := csv.NewReader(strings.NewReader(content))
	records, err := r.ReadAll()
	if err != nil || len(records) == 0 {
		return content
	}
	
	colWidths := make([]int, len(records[0]))
	for _, row := range records {
		for i, col := range row {
			if i < len(colWidths) {
				if len(col) > colWidths[i] {
					colWidths[i] = len(col)
				}
			}
		}
	}
	
	dividersWidth := (len(records[0]) - 1) * 3
	totalWidth := dividersWidth
	for _, w := range colWidths {
		totalWidth += w
	}
	
	if maxWidth > 0 && totalWidth > maxWidth && maxWidth > dividersWidth {
		available := maxWidth - dividersWidth
		for i := range colWidths {
			colWidths[i] = (colWidths[i] * available) / (totalWidth - dividersWidth)
			if colWidths[i] < 3 {
				colWidths[i] = 3
			}
		}
	}
	
	var sb strings.Builder
	for rIdx, row := range records {
		bgColor := ""
		if rIdx%2 == 1 {
			bgColor = "\033[48;5;236m"
		}
		
		sb.WriteString(bgColor)
		for i, col := range row {
			if i >= len(colWidths) {
				break
			}
			if i > 0 {
				sb.WriteString(" │ ")
			}
			w := colWidths[i]
			cell := col
			if len(cell) > w {
				cell = cell[:w-1] + "…"
			}
			format := fmt.Sprintf("%%-%ds", w)
			sb.WriteString(fmt.Sprintf(format, cell))
		}
		sb.WriteString(ColorReset)
		sb.WriteString("\n")
	}
	
	return sb.String()
}
