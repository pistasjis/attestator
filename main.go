// attestator (pistasjis Attestator) is an open-source application which lets you evaluate how apps on your computer respect your privacy and security.
package main

import (
	"embed"

	"github.com/pistasjis/attestator/vars"

	"github.com/pistasjis/attestator/cmd"
)

//go:embed results_template.html
var templates embed.FS

func main() {
	vars.Templates = templates

	cmd.Execute()
}
