package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var delCmd = &cobra.Command{
	Use:   "del",
	Short: "2fa del <name>",
	Run:   delRun,
}

func delRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		_ = cmd.Help()
		return
	}
	if err := tfa.Remove(args[0]); err != nil {
		log.Fatalln("could not remove key", err)
	}
}
