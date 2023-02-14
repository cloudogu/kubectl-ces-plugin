package config

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "config [edit | delete] [args]...",
	}

	cmd.AddCommand(
		editCmd(),
		deleteCmd(),
	)

	return cmd
}

func editCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "edit <dogu-name> <registry-key> <value>",
		// TODO add completion, dogu name validation
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("not implemented")
		},
	}

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete <dogu-name> <registry-key>",
		// TODO add completion, dogu name validation
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("not implemented")
		},
	}

	return cmd
}
