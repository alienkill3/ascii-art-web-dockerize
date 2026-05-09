package src

import (
	_ "embed"
	"io"
	"strings"
)

type FontType string

const (
	Standard   FontType = "standard"
	Shadow     FontType = "shadow"
	Thinkertoy FontType = "thinkertoy"
)

//go:embed fonts/standard.txt
var standardStr string

//go:embed fonts/shadow.txt
var shadowStr string

//go:embed fonts/thinkertoy.txt
var thinkertoyStr string

func IsValidFont(inputFlag string) bool {
	switch FontType(inputFlag) {
	case Standard, Shadow, Thinkertoy:
		return true
	default:
		return false
	}
}

func GetReader(f FontType) io.Reader {
	switch f {
	case Shadow:
		return strings.NewReader(shadowStr)
	case Thinkertoy:
		return strings.NewReader(thinkertoyStr)
	default:
		return strings.NewReader(standardStr)
	}
}
