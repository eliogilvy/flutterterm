package cmd

import (
	"flutterterm/pkg/model"
	"flutterterm/pkg/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var config *model.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flutterterm",
	Short: "A flutter command line tool",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). t only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(setConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setConfig() {
	c, err := model.LoadConfig(model.DefaultConfigPath)

	if err != nil {
		utils.PrintError(fmt.Sprintf("Error loading config, try running flutterterm help init: %s", err.Error()))
	}

	config = c
}
