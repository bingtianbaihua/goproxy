package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bingtianbaihua/goproxy/config"
	"github.com/bingtianbaihua/goproxy/server"

	"github.com/spf13/cobra"
)

const (
	version = "0.1.0"
)

var (
	showVersion bool
	cfgFile     string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "", "c", "", "config file of goproxy")
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of goproxy")
}

var rootCmd = &cobra.Command{
	Use:   "goproxy",
	Short: "goproxy is a http/https proxy(https://github.com/bingtianbaihua/goproxy)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version)
			return nil
		}

		content, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		cfg := new(config.Config)
		err = json.Unmarshal(content, cfg)
		if err != nil {
			fmt.Println(err)
			return err
		}

		p, err := server.NewProxyServer(cfg)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return p.ListenAndServe()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
