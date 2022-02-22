/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sftcli",
	Short: "cli interface to scaleft (oktaasa) API",
	Long: `There's a lot of things that the okta asa (scaleft) api can do that the
sft command cant, most notably delete a host!`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sftcli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config",
		".sftcli.yaml",
		"config file (default is sftcli.yaml)")
	//err := viper.WriteConfig()
	//if err != nil {
	//	log.Println("%#v copuldn't write config", err)
	//}
}

func initConfig() {
	viper.SetConfigName(".sftcli")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	//log.Println("We invoked initconfig")
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file: ", viper.ConfigFileUsed())
	} else {
		log.Println("No Config file.")
	}
}
