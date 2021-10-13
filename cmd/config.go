package cmd

import (
	//"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ingshtrom/dh-debug/pkg/types"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

//go:embed config/dh-debug.json
var defaultConfig []byte


func initConfig() {
	fmt.Println(cfgFile)
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println("error getting home directory: ", err)
			os.Exit(1)
		}

		// Search config in home directory with name ".dh-debug" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".dh-debug")
		
		if err = viper.ReadInConfig(); err != nil {
			// if we get an error reading the config, then just load the default config and 
			// shove that in the default place for the config
			var config *types.Config
			if err := json.Unmarshal(defaultConfig, &config); err != nil {
				fmt.Println("could not parse default config:", err)
				os.Exit(1)
			}

			viper.SetDefault("shellTests", config.ShellTests)

			if err = viper.SafeWriteConfig(); err != nil {
				fmt.Println("failed to write default config:", err)
				os.Exit(1)
			}
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("cannot read config: ", err)
		os.Exit(1)
	}
}
