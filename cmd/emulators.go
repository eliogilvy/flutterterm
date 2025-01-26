package cmd

import (
	"flutterterm/pkg/ui"
	"flutterterm/pkg/utils"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// runs flutter emulators
var emulatorsCmd = &cobra.Command{
	Use:   "emulators",
	Short: "Runs 'flutter emulators'",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		isCold, err := cmd.Flags().GetBool("cold")

		p := tea.NewProgram(ui.InitialEmulatorModel(isCold), tea.WithAltScreen())

		_, err = p.Run()

		if err != nil {
			e := fmt.Sprintf("Emulators exited with error: %s", err)
			utils.PrintError(e)
		}
	},
}

func init() {
	rootCmd.AddCommand(emulatorsCmd)
	emulatorsCmd.Flags().BoolP("cold", "c", false, "Cold start the emulator")
}
