package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "vps-deploy",
	Short: "VPS Deploy is a CLI tool for deploying services to VPS hosts via SSH",
	Long: `A CLI tool that reads a deploy.yaml configuration and deploys, rolls back,
or checks the status of services on remote VPS hosts over SSH.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vps-deploy — use a subcommand: deploy, rollback, status")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "deploy.yaml", "path to deploy.yaml config file")
}
