package cmd

/* Copyright © 2021 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
  _ "embed"
  "github.com/getlantern/systray"
  "github.com/spf13/cobra"
  "time"
)

//var globe, globeEye image.Image

//go:embed assets/noun_Circle_1776951.ico
var globeICO []byte

//go:embed assets/noun_Eye_4410395.ico
var globeEyeICO []byte

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Watch system proxy, report",
	Long: `Watch the system proxy

report.`,

	Run: func(cmd *cobra.Command, args []string) {
    //globe, _ = png.Decode(bytes.NewReader(globeICO))
    //globeEye, _ = png.Decode(bytes.NewReader(globeEyeICO))

    //systray.Register(onReady, onExit)
    systray.Run(onReady, onExit)
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// proxyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// proxyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func onReady() {
  systray.SetIcon(globeICO)

  quit := systray.AddMenuItem("Quit", "Quit stario-start")

  go func() {
    for {
      select {
      case <- quit.ClickedCh:
        systray.Quit()
        return
      }
    }
  }()

  go func() {
    for {
        isProxy, proxy, problem, err := checkProxy(false)
        if err != nil {
          systray.SetTooltip(err.Error())
        } else if len(problem) > 0 {
          systray.SetTooltip(problem)
        } else if len(proxy) > 0 {
          systray.SetTooltip(proxy)
        } else {
          systray.SetTooltip("Ok - No Proxy")
        }

        if isProxy {
          systray.SetIcon(globeEyeICO)
          systray.SetTooltip("Ok - " + proxy)
        } else {
          systray.SetIcon(globeICO)
          systray.SetTooltip("Ok - No Proxy")
        }

        time.Sleep(2 * time.Second)
      }
  }()
}

func onExit() {

}
