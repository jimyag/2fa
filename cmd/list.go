package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all 2fa keys",
	Run:   listRun,
}

func listRun(cmd *cobra.Command, args []string) {
	keys := tfa.List()

	tb := table.NewWriter()
	tb.Style().Options.DrawBorder = true
	tb.Style().Options.SeparateRows = true
	tTemp := table.Table{}
	tTemp.Render()
	tb.SetColumnConfigs([]table.ColumnConfig{
		{Name: "name", Align: text.AlignCenter},
		{Name: "totp", Align: text.AlignCenter},
		{Name: "lifetime/s", Align: text.AlignCenter},
		{Name: "next totp", Align: text.AlignCenter},
	})
	tb.AppendHeader(table.Row{"name", "totp", "lifetime/s", "next totp"})
	for name, key := range keys {
		if name == "" {
			continue
		}
		t := time.Now()
		code, err := GenTOTP(key.Key, t, 6, 30)
		if err != nil {
			log.Fatalln("could not generate code", err)
		}
		lifeTime := 30 - t.Unix()%30

		nextCode, err := GenTOTP(key.Key, t.Add(30*time.Second), 6, 30)
		if err != nil {
			log.Fatalln("could not generate code", err)
		}
		tb.AppendRow(table.Row{name, code, lifeTime, nextCode})
	}
	fmt.Println(tb.Render())
}
