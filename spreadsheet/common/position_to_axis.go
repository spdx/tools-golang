package common

import "fmt"

// PositionToAxis takes a column string and a row integer to combines them into an "axis"
// to be used with the Excelize module.
// An "axis" is the word Excelize uses to describe a coordinate position within a spreadsheet,
// e.g. "A1", "B14", etc.
func PositionToAxis(column string, row int) string {
	return fmt.Sprintf("%s%d", column, row)
}
