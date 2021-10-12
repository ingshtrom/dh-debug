package cmd

import (
	"fmt"
	"os"

	"github.com/ingshtrom/dh-debug/pkg"
	"github.com/spf13/cobra"
)

var (
	configFile    string
	debugFile    string
	isDebug     bool
	isSummarize bool
	isPrint bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dh-debug",
	Short: "Docker Hub Debug cli",
	Long:  `Docker cli containing commands to debug network connections to Docker Hub.`,
	Run: func(cmd *cobra.Command, args []string) {
		if configFile == "" {
			fmt.Println("Please specify the -c, --config parameter")
			os.Exit(1)
		}

		if !isDebug && !isSummarize && !isPrint {
			isDebug = true
		}

		if isDebug {
			pkg.RunDebugTests(configFile, debugFile)
		} else if isPrint {
			pkg.PrintDebugTests(debugFile)
		} else if isSummarize {
			pkg.SummarizeTestResults(debugFile)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "./dh-config.json", "File to read test declarations from.")
	rootCmd.Flags().StringVarP(&debugFile, "file", "f", "./dh-debug.json", "File to write runs to or read runs from.")
	rootCmd.Flags().BoolVarP(&isDebug, "debug", "d", false, "(default) Whether to debug the network to Docker Hub.")
	rootCmd.Flags().BoolVarP(&isSummarize, "summarize", "s", false, "Summarize some of the information in the output.")
	rootCmd.Flags().BoolVarP(&isPrint, "print", "p", false, "Human readable print of the output file from a debug or parse run.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
