package cmd

import (
	dockerApp "github.com/docker/libcompose/cli/docker/app"
	"github.com/spf13/cobra"
	"os"
)

var factory *dockerApp.ProjectFactory

var subCommandTemplate = `Usage:
  {{.CommandPath}} [options] [SERVICE...]
{{if .HasLocalFlags}}
Options:
{{.LocalFlags.FlagUsages | trimRightSpace}}
{{end}}
`
var usageTemplate = `Usage:
{{- if .Runnable -}}
  {{- if .HasFlags}}
 	 {{appendIfNotPresent .UseLine "[flags]"}}
  {{- else}}
  	{{- .UseLine}}
  {{- end}}
{{- end}}
{{- if .HasSubCommands}}
  {{.CommandPath}} [command]
{{- end -}}
{{- if gt .Aliases 0}}
   Aliases:
   {{- .NameAndAliases}}
{{- end}}
{{- if .HasExample}}
   Examples:
 {{- .Example }}
{{- end}}

{{- if .HasAvailableSubCommands}}
{{if .HasLocalFlags}}
Options:
{{.LocalFlags.FlagUsages | trimRightSpace}}
{{end -}}
{{- if .HasInheritedFlags -}}
Global Flags:
   {{.InheritedFlags.FlagUsages | trimRightSpace -}}
{{- end}}
{{- if .HasHelpSubCommands}}
   Additional help topics:{{range .Commands}}
   {{- if .IsHelpCommand}}
     {{- rpad .CommandPath .CommandPathPadding}} {{.Short}}
   {{- end}}
{{- end}}
{{end}}

{{- if .HasAvailableSubCommands}}
Commands:{{range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}
{{- end}}
{{- end}}
{{- end}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "libcompose",
	Short: "Define and run multi-container applications with Docker.",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	factory = &dockerApp.ProjectFactory{}
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	RootCmd.SetUsageTemplate(usageTemplate)
	RootCmd.PersistentFlags().Bool("verbose", false, "Show more output")
	RootCmd.PersistentFlags().StringP("project-name", "p", "", "Specify an alternate project name (default: directory name)")
	RootCmd.PersistentFlags().StringSliceP("file", "f", []string{""}, "Specify an alternate compose file (default: docker-compose.yml)")

	RootCmd.PersistentFlags().String("configdir", "", "Path to docker config dir, default ${HOME}/.docker")
	RootCmd.PersistentFlags().Bool("tls", true, "Use TLS; implied by --tlsverify")
	RootCmd.PersistentFlags().Bool("tlsverify", true, "Use TLS and verify the remote")
	RootCmd.PersistentFlags().String("tlscacert", "", "Trust certs signed only by this CA")
	RootCmd.PersistentFlags().String("tlscert", "", "Path to TLS certificate file")
}
