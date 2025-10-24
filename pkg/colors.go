package helmutil

import "github.com/fatih/color"

var (
	ChartColor   = color.New(color.FgCyan, color.Bold)
	AliasColor   = color.New(color.FgHiBlue)
	DepColor     = color.New(color.FgHiWhite)
	VersionColor = color.New(color.FgHiBlack)
	SuccessColor = color.New(color.FgGreen)
	WarningColor = color.New(color.FgYellow)
	ErrorColor   = color.New(color.FgRed)
)
