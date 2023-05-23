package ncvoters

import (
	"fmt"
	"strings"
)

func CreateDDL() string {	
	sb := strings.Builder{}
	sb.WriteString("CREATE TABLE IF NOT EXISTS voters (\n")
	for i, colName := range selectedCols {
		comma := ","
		if i == len(selectedCols)-1 {
			comma = ""
		}
		part := fmt.Sprintf("  %-20s TEXT%s\n", colName, comma)
		sb.WriteString(part)
	}
	sb.WriteString(")\n")
	s := sb.String()
	return s
}