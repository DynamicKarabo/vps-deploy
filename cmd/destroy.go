package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/DynamicKarabo/vps-deploy/internal/config"
	"github.com/DynamicKarabo/vps-deploy/internal/ssh"
	"github.com/spf13/cobra"
)

var (
	destroyVolumes bool
	destroyRMI     bool
	destroyPath    string
)

var destroyCmd = &cobra.Command{
	Use:   "destroy [service]",
	Short: "Destroy a service by running docker compose down",
	Long: `Destroy a service by SSHing into its host and running
docker compose down in the service's directory.

By default runs: cd /root/<service> && docker compose down
Use --path to customize the compose directory.
Use --volumes (-v) to also remove named volumes.
Use --rmi to also remove images (--rmi all).`,
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

		// Build the docker compose down command
		composeDir := destroyPath
		if composeDir == "" {
			composeDir = "/root/" + svcName
		}

		var parts []string
		parts = append(parts, "cd", composeDir, "&&", "docker", "compose", "down")
		if destroyVolumes {
			parts = append(parts, "-v")
		}
		if destroyRMI {
			parts = append(parts, "--rmi", "all")
		}
		composeCmd := strings.Join(parts, " ")

		client := ssh.NewClient(svc.User, svc.Host, svc.KeyPath)
		if err := client.Connect(); err != nil {
			log.Fatalf("ssh connect failed: %v", err)
		}
		defer client.Close()

		fmt.Printf("Destroying %s...\n", svcName)
		fmt.Printf("> %s\n", composeCmd)
		out, err := client.Run(composeCmd)
		if err != nil {
			log.Fatalf("docker compose down failed: %v\noutput: %s", err, out)
		}
		fmt.Println(out)
		fmt.Println("Destroy complete.")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
	destroyCmd.Flags().StringVar(&serviceName, "service", "", "service name to destroy")
	destroyCmd.Flags().BoolVarP(&destroyVolumes, "volumes", "v", false, "also remove named volumes (docker compose down -v)")
	destroyCmd.Flags().BoolVar(&destroyRMI, "rmi", false, "also remove images (docker compose down --rmi all)")
	destroyCmd.Flags().StringVar(&destroyPath, "path", "", "custom path to the docker-compose directory (default: /root/<service>)")
}
