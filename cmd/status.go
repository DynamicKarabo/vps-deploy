package cmd

import (
	"fmt"
	"log"

	"github.com/DynamicKarabo/vps-deploy/internal/config"
	"github.com/DynamicKarabo/vps-deploy/internal/health"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [service]",
	Short: "Check the health status of a deployed service",
	Long: `Check whether a deployed service is healthy by polling its
configured health check URL.`,
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

		if svc.HealthCheckURL == "" {
			fmt.Printf("No health_check_url configured for service %q\n", svcName)
			return
		}

		fmt.Printf("Checking health of %s at %s...\n", svcName, svc.HealthCheckURL)
		checker := health.NewChecker(svc.HealthCheckURL, 10, 2)
		if err := checker.Poll(); err != nil {
			log.Fatalf("service %q is UNHEALTHY: %v", svcName, err)
		}
		fmt.Println("Service is HEALTHY")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().StringVar(&serviceName, "service", "", "service name to check")
}
