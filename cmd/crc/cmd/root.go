package cmd

import (
	"github.com/code-ready/crc/pkg/crc/errors"
	"github.com/code-ready/crc/pkg/crc/machine"
	"github.com/code-ready/crc/pkg/crc/output"
	"github.com/spf13/cobra"

	cmdConfig "github.com/code-ready/crc/cmd/crc/cmd/config"

	"github.com/code-ready/crc/pkg/crc/config"
	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/pkg/crc/logging"
)

var rootCmd = &cobra.Command{
	Use:   commandName,
	Short: descriptionShort,
	Long:  descriptionLong,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		runPrerun()
	},
	Run: func(cmd *cobra.Command, args []string) {
		runRoot()
		_ = cmd.Help()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		runPostrun()
	},
}

var globalForce bool

func init() {
	if err := constants.EnsureBaseDirExists(); err != nil {
		logging.Fatal(err.Error())
	}
	if err := config.EnsureConfigFileExists(); err != nil {
		logging.Fatal(err.Error())
	}
	if err := config.InitViper(); err != nil {
		logging.Fatal(err.Error())
	}

	setConfigDefaults()

	// subcommands
	rootCmd.AddCommand(cmdConfig.ConfigCmd)

	rootCmd.PersistentFlags().StringVar(&logging.LogLevel, "log-level", constants.DefaultLogLevel, "log level (e.g. \"debug | info | warn | error\")")
	rootCmd.PersistentFlags().BoolVarP(&globalForce, "force", "f", false, "Forcefully perform an action")
}

func runPrerun() {
	// Setting up logrus
	logging.InitLogrus(logging.LogLevel)
	logging.SetupFileHook()
}

func runPostrun() {
	logging.CloseLogging()
}

func runRoot() {
	output.Outln("No command given")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logging.Fatal(err)
	}
}

func setConfigDefaults() {
	config.SetDefaults()
}

func exitIfMachineMissing(name string) {
	exists, err := machine.MachineExists(name)
	if err != nil {
		errors.ExitWithMessage(1, err.Error())
	}
	if !exists {
		errors.ExitWithMessage(1, "Machine \"crc\" does not exist. Use \"crc start\" to add a new one.")
	}
}
