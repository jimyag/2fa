package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "2fa get <name>",
	Run:   getRun,
}

var copyToClipboard bool

func init() {
	getCmd.Flags().BoolVarP(&copyToClipboard, "copy", "c", false, "copy to clipboard")
}

func getRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		_ = cmd.Help()
		return
	}
	key := tfa.Get(args[0])
	if key == "" {
		log.Fatalln("could not find key")
	}

	code, err := GenTOTP(key, time.Now(), 6, 30)
	if err != nil {
		log.Fatalln("could not generate code", err)
	}
	if copyToClipboard {
		if err = clipboard.WriteAll(code); err != nil {
			log.Fatalln("could not copy to clipboard", err)
		}
	} else {
		fmt.Println(code)

	}

}
