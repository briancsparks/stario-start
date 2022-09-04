package cmd

/* Copyright © 2021 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
  "fmt"
  "github.com/spf13/cobra"
  "os"
  "strconv"
  "strings"

  homedir "github.com/mitchellh/go-homedir"
  "github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "stario-start",
  Short: "Stuff to run at startup",
  Long: `All the stuff

that would be good to run at Windows startup time.`,

  	Run: func(cmd *cobra.Command, args []string) {
      _ = incFile("C:\\Users\\sparksb\\stario.txt")
    },
}

func incFile(filename string) error {
  data, _ := os.ReadFile(filename)
  //fmt.Printf("read: %s\n", data)

  n, _ := strconv.Atoi(strings.TrimSpace(string(data)))
  n++
  out := strconv.Itoa(n) +"\n"
  //fmt.Printf("cnv: %d %d %s\n", n, n, out)

  _ = os.WriteFile(filename, []byte(out), 0644)

  return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)

  // Here you will define your flags and configuration settings.
  // Cobra supports persistent flags, which, if defined here,
  // will be global for your application.

  rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stario-start.yaml)")

  // Cobra also supports local flags, which will only run
  // when this action is called directly.
  rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
  if cfgFile != "" {
    // Use config file from the flag.
    viper.SetConfigFile(cfgFile)
  } else {
    // Find home directory.
    home, err := homedir.Dir()
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    // Search config in home directory with name ".stario-start" (without extension).
    viper.AddConfigPath(home)
    viper.SetConfigName(".stario-start")
  }

  viper.AutomaticEnv() // read in environment variables that match

  // If a config file is found, read it in.
  if err := viper.ReadInConfig(); err == nil {
    fmt.Println("Using config file:", viper.ConfigFileUsed())
  }
}
