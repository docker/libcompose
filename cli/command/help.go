package command

import (
	"os"
	"path"

	"github.com/codegangsta/cli"
)

func init() {
	cli.AppHelpTemplate = `{{.Usage}}

Usage: {{.Name}} {{if .Flags}}[options] {{end}}COMMAND [arg...]

{{if .Flags}}
Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}
Commands:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} COMMAND --help' for more information on a command.
`
	cli.CommandHelpTemplate = `{{.Usage}}{{if .Description}}

{{.Description}}{{end}}

Usage: ` + path.Base(os.Args[0]) + ` {{.Name}}{{if .Flags}} [options]{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{end}}

{{if .Flags}}Options:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

}
