/*
Copyright © 2024 Jason Williams <jwilliamstx@protonmail.com>
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "spo-recyclebin",
	Short: "Find and restore files from the SharePoint recycle bin",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var SiteURL string

func init() {
	rootCmd.PersistentFlags().String("siteUrl", "", "SharePoint Site URL")
	rootCmd.PersistentFlags().String("before", "", "List files deleted before YYYY-MM-DDTHH:MM:SS")
	rootCmd.PersistentFlags().String("after", "", "List files deleted after YYYY-MM-DDTHH:MM:SS")
	rootCmd.PersistentFlags().String("by", "", "List files deleted by {user's full name}")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatal(err)
	}
}
