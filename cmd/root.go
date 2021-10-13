package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ingshtrom/dh-debug/pkg"
	"github.com/ingshtrom/dh-debug/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile          string
	debugResultsFile string
	isDebug          bool
	isSummarize      bool
	isPrint          bool
	printErrorsOnly  bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dh-debug",
	Short: "Docker Hub Debug cli",
	Long:  `Docker cli containing commands to debug network connections to Docker Hub.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !isDebug && !isSummarize && !isPrint {
			isDebug = true
		}

		var config *types.Config
		if err := viper.Unmarshal(&config); err != nil {
			fmt.Println("error unmarshalling config file: ", err)
			os.Exit(1)
		}

		if isDebug {
			pkg.RunDebugTests(config, debugResultsFile)
		} else if isPrint {
			pkg.PrintDebugTests(debugResultsFile, printErrorsOnly)
		} else if isSummarize {
			pkg.SummarizeTestResults(debugResultsFile)
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.dh-debug.json)")

	rootCmd.Flags().StringVarP(&debugResultsFile, "file", "f", "./dh-debug-results.json", "File to write runs to or read runs from.")
	rootCmd.Flags().BoolVarP(&isDebug, "debug", "d", false, "(default) Whether to debug the network to Docker Hub.")
	rootCmd.Flags().BoolVarP(&isSummarize, "summarize", "s", false, "Summarize some of the information in the output.")
	rootCmd.Flags().BoolVarP(&isPrint, "print", "p", false, "Human readable print of the output file from a debug or parse run.")
	rootCmd.Flags().BoolVarP(&printErrorsOnly, "errors-only", "e", false, "For print only, only print out errored tests.")

	var config *types.Config
	if err := json.Unmarshal(defaultConfig, &config); err != nil {
		fmt.Println("could not parse default config:", err)
		os.Exit(1)
	}

	viper.SetDefault("shellTests", config.ShellTests)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
