package cmd

import (
	"fmt"
	"log"

	"github.com/DynamicKarabo/vps-deploy/internal/config"
	"github.com/DynamicKarabo/vps-deploy/internal/health"
	"github.com/DynamicKarabo/vps-deploy/internal/ssh"
	"github.com/spf13/cobra"
)

var serviceName string

var deployCmd = &cobra.Command{
	Use:   "deploy [service]",
	Short: "Deploy a service to a VPS host",
	Long: `Deploy a service defined in deploy.yaml to its target VPS host.
Executes the configured deploy commands over SSH, then optionally
polls the health check URL until the service is healthy.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Parse(cfgFile)
		if err != nil {
			log.Fatalf("failed to parse config: %v", err)
		}

		svcName := serviceName
		if svcName == "" && len(args) > 0 {
			svcName = args[0]
		}
		if svcName == "" {
			log.Fatal("service name is required; provide as argument or via --service")
		}

		svc, ok := cfg.Services[svcName]
		if !ok {
			log.Fatalf("service %q not found in config", svcName)
		}

		client := ssh.NewClient(svc.User, svc.Host, svc.KeyPath)
		if err := client.Connect(); err != nil {
			log.Fatalf("ssh connect failed: %v", err)
		}
		defer client.Close()

		for _, cmdStr := range svc.DeployCommands {
			fmt.Printf("> %s\n", cmdStr)
			out, err := client.Run(cmdStr)
			if err != nil {
				log.Fatalf("command failed: %s: %v\noutput: %s", cmdStr, err, out)
			}
			fmt.Println(out)
		}

		if svc.HealthCheckURL != "" {
			fmt.Println("Checking health...")
			checker := health.NewChecker(svc.HealthCheckURL, 30, 5)
			if err := checker.Poll(); err != nil {
				log.Fatalf("health check failed: %v", err)
			}
			fmt.Println("Service is healthy!")
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringVar(&serviceName, "service", "", "service name to deploy")
}
