package cmd

import (
	"flutterterm/pkg/ui"
	"flutterterm/pkg/utils"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// The file to look for in a flutter project
const pubspec = "pubspec.yaml"

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A guided flutter run command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		shouldForce, err := cmd.Flags().GetBool("force")
		if !assertRootPath() && !shouldForce {
			utils.PrintError("pubspec.yaml not found. Make sure you are in a flutter directory")
			return
		}

		if shouldForce {
			utils.PrintWarning("\n** Force is true, bypassing **\n\n")
		}

		utils.PrintInfo(fmt.Sprintf("Flutter directory detected. Getting devices...\n\n"))

		configs, err := utils.GetConfigs()

		// Add a default run config if none exist
		if len(configs) == 0 {
			utils.PrintInfo("No configs found, using default\n\n")
			help := fmt.Sprintf("Try creating a \"%s\" file or adding a config to an already created one\n\n", utils.ConfigPath)
			utils.PrintHelp(help)
			defaultConfig, err := utils.DefaultConfig()
			if err != nil {
				utils.PrintError(err.Error())
				return
			}
			configs = append(configs, defaultConfig)
		}

		p := tea.NewProgram(ui.InitialRunModel(configs), tea.WithAltScreen())

		model, err := p.Run()

		if err != nil {
			utils.PrintError(fmt.Sprintf("Error %s", err.Error()))
			return
		}

		runModel, _ := model.(ui.RunModel)

		if !runModel.IsComplete() {
			return
		}

		setupAndRun(runModel)
	},
}

// Runs command based on the model received
func setupAndRun(m ui.RunModel) {
	fmt.Printf("Running %s on %s\n\n", m.Selected_config.Name, m.Selected_device.Name)

	// Device
	deviceID := m.Selected_device.ID
	config := m.Selected_config

	err := config.AssertConfig()

	if err != nil {
		e := fmt.Sprintf("Invalid configuration: %s", err)
		utils.PrintError(e)
		return
	}

    args := config.GetArgs(deviceID)

	cmd := utils.FlutterRun(args)

	// For color and input handling
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Start()

	if err != nil {
		utils.PrintError(err.Error())
		return
	}

	if err := cmd.Wait(); err != nil {
		s := fmt.Sprintf("Flutterterm finished with error: %s", err)
		utils.PrintError(s)
	} else {
		utils.PrintSuccess("Flutterterm finished successfully")
	}
}

// Check if in a flutter project
func assertRootPath() bool {
	_, err := os.Stat(pubspec)

	if err != nil {
		return false
	}

	return true
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("force", "f", false, "For bypassing flutter directory assertion")
	runCmd.Flags().MarkHidden("force")
}
