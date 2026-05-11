package cmd

import (
	"fmt"
	"log"

	"github.com/DynamicKarabo/vps-deploy/internal/config"
	"github.com/DynamicKarabo/vps-deploy/internal/ssh"
	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback [service]",
	Short: "Rollback a service to its previous state",
	Long: `Rollback a service by executing its configured rollback command
over SSH on the target VPS host.`,
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

		if svc.RollbackCommand == "" {
			log.Fatalf("no rollback_command configured for service %q", svcName)
		}

		client := ssh.NewClient(svc.User, svc.Host, svc.KeyPath)
		if err := client.Connect(); err != nil {
			log.Fatalf("ssh connect failed: %v", err)
		}
		defer client.Close()

		fmt.Printf("Rolling back %s...\n", svcName)
		fmt.Printf("> %s\n", svc.RollbackCommand)
		out, err := client.Run(svc.RollbackCommand)
		if err != nil {
			log.Fatalf("rollback command failed: %v\noutput: %s", err, out)
		}
		fmt.Println(out)
		fmt.Println("Rollback complete.")
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
	rollbackCmd.Flags().StringVar(&serviceName, "service", "", "service name to rollback")
}
