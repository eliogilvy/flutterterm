package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

const (
	ConfigPath  = ".fterm_config.json"
	mainPath    = "main.dart"
	mainLibPath = "lib/main.dart"
)

// --mode
var flutterModes = []string{
	"debug", "profile", "release",
}

// main.dart paths to look for
var mainPaths = []string{
	mainPath, mainLibPath,
}

type FlutterRunConfig struct {
	Name               string `json:"name"`
	Mode               string `json:"mode"`
	Flavor             string `json:"flavor"`
	Target             string `json:"target"`
	DartDefineFromFile string `json:"dart_define_from_file"`
}

func (config FlutterRunConfig) GetArgs(deviceID string) []string {

	args := []string{"run"}

	if deviceID != "" {
		args = append(args, "-d", deviceID)
	}
	if config.Target != "" {
		args = append(args, "-t", config.Target)
	}
	if config.Mode != "" {
		arg := fmt.Sprintf("--%s", config.Mode)
		args = append(args, arg)
	}
	if config.Flavor != "" {
		args = append(args, "--flavor", config.Flavor)
	}
	if config.DartDefineFromFile != "" {
		args = append(args, "--dart-define-from-file", config.DartDefineFromFile)
	}

	return args
}

func (c FlutterRunConfig) FilterValue() string {
	return c.Name
}

func (c FlutterRunConfig) Title() string {
	return c.Name
}

func (c FlutterRunConfig) Description() string {
	return strings.Join(c.GetArgs(""), " ")
}

// Makes sure config is properly configured
func (config FlutterRunConfig) AssertConfig() error {
	if !assertFlutterMode(config.Mode) {
		e := fmt.Sprintf("Invalid mode: %s", config.Mode)
		return errors.New(e)
	}
	return nil
}

// Verify proper mode being used
func assertFlutterMode(m string) bool {
	// Empty mode is ok
	if m == "" {
		return true
	}
	m = strings.ToLower(m)
	for _, mode := range flutterModes {
		if mode == m {
			return true
		}
	}
	return false
}

func (config FlutterRunConfig) ToString() string {
	var s string
	s = fmt.Sprintf("Config: %s\n", config.Name)
	s += fmt.Sprintf("Mode: %s\n", config.Mode)
	s += fmt.Sprintf("Flavor: %s\n", config.Flavor)
	s += fmt.Sprintf("Target: %s\n", config.Target)
	s += fmt.Sprintf("Dart define file: %s\n", config.DartDefineFromFile)
	return s
}

func DefaultConfig() (FlutterRunConfig, error) {
	target, err := findDefaultTarget()
	if err != nil {
		return FlutterRunConfig{}, err
	}
	return FlutterRunConfig{
		Name:   "Default",
		Mode:   "debug",
		Target: target,
	}, nil
}

func GetConfigs() ([]FlutterRunConfig, error) {
	var configs []FlutterRunConfig

	config_file, err := os.Open(ConfigPath)

	if err != nil {
		return configs, err
	}

	defer config_file.Close()

	// Read file

	bytes, err := io.ReadAll(config_file)

	if err != nil {
		return configs, err
	}

	err = json.Unmarshal(bytes, &configs)

	if err != nil {
		fmt.Println(err)
		return configs, err
	}

	for i := 0; i < len(configs); i++ {
		configs[i].ToString()
	}

	return configs, nil
}

// Looks for main.dart files in default config
func findDefaultTarget() (string, error) {
	for _, path := range mainPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	err := errors.New("main.dart file not found")
	return "", err
}

func GetConfigList(configs []FlutterRunConfig) list.Model {
	var items []list.Item
	for _, device := range configs {
		item := device
		items = append(items, item)
	}

	l := GetList(items)
	l.Title = "Configs"

	return l
}
