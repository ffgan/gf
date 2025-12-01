package cli

import (
	"fmt"
	"strings"

	"github.com/ffgan/gf/internal/utils"
)

// ColorBlocksParams holds all parameters needed for generating color blocks
type ColorBlocksParams struct {
	ColorBlocks bool   // Whether to display color blocks
	BlockWidth  int    // Width of each block in spaces
	BlockHeight int    // Height of blocks in rows
	BlockRange  [2]int // Start and end color range [start, end]
	ColOffset   string // Column offset ("auto" or number)
	TextPadding int    // Text padding value
	Term        string // Terminal type (e.g., "linux")
	BashVersion int    // Bash major version (for compatibility)
	ZWS         string // Zero-width space character
}

// ColorBlocksResult holds the result of color block generation
type ColorBlocksResult struct {
	Output     string // The formatted output string
	InfoHeight int    // Additional height to add to info_height
	Prin       int    // Flag indicating manual printing (1 = printed)
}

// GetCols generates color blocks based on the provided parameters
func GetCols(params ColorBlocksParams, leftMax int) ColorBlocksResult {
	result := ColorBlocksResult{
		Prin: 0,
	}

	if !params.ColorBlocks {
		return result
	}

	var blocks, blocks2 strings.Builder

	// Generate block width string (spaces)
	blockWidth := strings.Repeat(" ", params.BlockWidth)

	// Generate the color blocks - IMPORTANT: colors accumulate in the string
	// Each iteration adds to the existing string, creating a continuous color bar
	for i := params.BlockRange[0]; i <= params.BlockRange[1]; i++ {
		if i >= 0 && i <= 7 {
			// Standard 8 colors (0-7)
			// Format: ESC[3Xm (foreground) ESC[4Xm (background) spaces
			blocks.WriteString(fmt.Sprintf("\033[3%dm\033[4%dm%s", i, i, blockWidth))
		} else {
			// Extended 256 colors (8+)
			// Format: ESC[38;5;Xm (foreground) ESC[48;5;Xm (background) spaces
			blocks2.WriteString(fmt.Sprintf("\033[38;5;%dm\033[48;5;%dm%s", i, i, blockWidth))
		}
	}

	// Workaround for bright background colors in Linux framebuffer console
	if params.BlockRange[1] < 16 && params.Term == "linux" {
		blocks2Str := blocks2.String()
		blocks2.Reset()
		blocks2.WriteString(fmt.Sprintf("\033[5m%s\033[25m", blocks2Str))
	}

	// Get the accumulated color strings
	blocksStr := blocks.String()
	blocks2Str := blocks2.String()

	// Build column output
	var cols strings.Builder

	// Convert height into spaces for repetition
	blockSpaces := strings.Repeat(" ", params.BlockHeight)

	// For each row (determined by block_height), add the color blocks
	// The blocks string is repeated for each "space" in blockSpaces
	if blocksStr != "" {
		// Add reset codes after the blocks: ESC[39;49m (reset fg/bg to default)
		row := strings.ReplaceAll(blockSpaces, " ", blocksStr+"\033[39;49mnl")
		cols.WriteString(row)
	}

	if blocks2Str != "" {
		// Add reset code after the blocks: ESC[0m (reset all attributes)
		// BUG: On some terminals, the blocks may shift left by one column
		if term == "vscode" {
			leftMax += 1
		}
		row := fmt.Sprintf("\033[%dC", leftMax) + strings.ReplaceAll(blockSpaces, " ", blocks2Str+"\033[0mnl")

		cols.WriteString(row)
	}

	// Determine the horizontal offset of the blocks
	var blockOffset int
	if params.ColOffset == utils.AUTO {
		blockOffset = params.TextPadding
	} else {
		fmt.Sscanf(params.ColOffset, "%d", &blockOffset)
	}
	// Process the column string
	colsStr := cols.String()

	// Remove trailing "nl"
	colsStr = strings.TrimSuffix(colsStr, "nl")

	// Replace "nl" with newlines and cursor positioning
	// Format: newline + ESC[XC (move cursor right X columns) + zero-width space
	var offsetStr string
	if blockOffset > 0 {
		offsetStr = fmt.Sprintf("\n\033[%dC", blockOffset)
	} else {
		offsetStr = "\n"
	}
	colsStr = strings.ReplaceAll(colsStr, "nl", offsetStr)

	// Add block height to info height
	heightIncrement := params.BlockHeight + 1
	if params.BlockRange[1] > 7 {
		heightIncrement = params.BlockHeight + 2
	}
	result.InfoHeight = heightIncrement
	// Format the final output with initial positioning
	// result.Output = fmt.Sprintf("\033[%dC%s%s", blockOffset, params.ZWS, colsStr)
	if blockOffset > 0 {
		result.Output = fmt.Sprintf("\033[%dC%s", blockOffset, colsStr)
	} else {
		result.Output = colsStr
	}

	// Tell info() that we printed manually
	result.Prin = 1

	return result
}

func Getcols(leftMax int) string {
	params := ColorBlocksParams{
		ColorBlocks: true,
		BlockWidth:  3,
		BlockHeight: 1,
		BlockRange:  [2]int{0, 15}, // Display colors 0-15
		ColOffset:   utils.AUTO,
		TextPadding: 0,
		Term:        "xterm-256color",
		BashVersion: 5,
		ZWS:         "\u200B", // Zero-width space
	}

	result := GetCols(params, leftMax)
	// if result.Prin == 1 {
	// fmt.Print(result.Output)
	// fmt.Printf("Info height increment: %d\n", result.InfoHeight)
	// }
	return result.Output

}
