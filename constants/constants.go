package constants

import (
	"github.com/fatih/color"
)

const (
	Red                        color.Attribute = color.FgRed
	Green                      color.Attribute = color.FgGreen
	Cyan                       color.Attribute = color.FgCyan
	White                      color.Attribute = color.FgWhite
	MaxIdleConnections         int             = 25
	RequestTimeout             int             = 180
	OKHttpStatusResponse       string          = "200 OK"
	AcceptedHttpStatusResponse string          = "202 Accepted"
)

var ValidServices = []string{"blueprint", "postgresql", "mongodb", "redis"}
