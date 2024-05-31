package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "2fa add <name> <key>",
	Run:   addRun,
}

func addRun(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		_ = cmd.Help()
		return
	}
	if err := tfa.Add(args[0], args[1]); err != nil {
		log.Fatalln("could not add key", err)
	}
}
