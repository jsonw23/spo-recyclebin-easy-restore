/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/jsonw23/spo-recyclebin-easy-restore/recyclebin"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the files in the recycle bin that match your query",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		query := recyclebin.NewQuery(args)

		for _, item := range query.Results().Data() {
			d := item.Data()
			fmt.Println(
				d.Title,
				d.DeletedByName,
				d.DeletedDate.Local(),
			)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
