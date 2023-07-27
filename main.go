package main

import (
	_ "embed"

	"tynmo/command/root"
	"tynmo/licenses"
)

var (
	//go:embed LICENSE
	license string
)

func main() {
	licenses.SetLicense(license)

	root.NewRootCommand().Execute()
}
