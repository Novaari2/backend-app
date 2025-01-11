package cmd

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "authentication-app",
		Short: "Run service",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(apiCmd())

	return command
}
