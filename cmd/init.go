package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type servicePrompt struct {
	Name            string
	Host            string
	User            string
	KeyPath         string
	DeployCommands  []string
	HealthCheckURL  string
	RollbackCommand string
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a deploy.yaml config interactively",
	Long: `Walks you through creating a deploy.yaml configuration file.
Prompts for each service's details and writes the result to deploy.yaml.`,
	Run: func(cmd *cobra.Command, args []string) {
		outputPath, _ := cmd.Flags().GetString("output")
		if outputPath == "" {
			outputPath = cfgFile
		}

		services := promptServices()

		if len(services) == 0 {
			fmt.Println("No services defined. Exiting.")
			return
		}

		cfg := map[string]interface{}{
			"services": services,
		}

		data, err := yaml.Marshal(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding YAML: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(outputPath, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outputPath, err)
			os.Exit(1)
		}

		fmt.Printf("\n✅ Config written to %s\n", outputPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("output", "o", "", "output path (default: deploy.yaml)")
}

func promptServices() map[string]interface{} {
	services := make(map[string]interface{})
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n━━━ New Service ━━━")
		sp := servicePrompt{}

		sp.Name = prompt(reader, "Service name", "my-service")
		sp.Host = prompt(reader, "Host (IP or domain)", "178.105.76.236")
		sp.User = prompt(reader, "SSH user", "root")
		defaultKey := filepath.Join(os.Getenv("HOME"), ".ssh", "id_ed25519")
		sp.KeyPath = prompt(reader, "SSH key path", defaultKey)
		sp.HealthCheckURL = prompt(reader, "Health check URL (optional)", "")

		fmt.Println("Deploy commands (one per line, empty line to finish):")
		for {
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			sp.DeployCommands = append(sp.DeployCommands, line)
		}

		if len(sp.DeployCommands) == 0 {
			sp.DeployCommands = []string{
				"cd /root/" + sp.Name + " && docker compose pull",
				"cd /root/" + sp.Name + " && docker compose up -d",
			}
		}

		sp.RollbackCommand = prompt(reader, "Rollback command (optional)", "")

		svc := map[string]interface{}{
			"host":     sp.Host,
			"user":     sp.User,
			"key_path": sp.KeyPath,
		}
		if len(sp.DeployCommands) > 0 {
			svc["deploy_commands"] = sp.DeployCommands
		}
		if sp.HealthCheckURL != "" {
			svc["health_check_url"] = sp.HealthCheckURL
		}
		if sp.RollbackCommand != "" {
			svc["rollback_command"] = sp.RollbackCommand
		}

		services[sp.Name] = svc

		ans := prompt(reader, "Add another service?", "n")
		if strings.ToLower(ans) != "y" && strings.ToLower(ans) != "yes" {
			break
		}
	}

	return services
}

func prompt(reader *bufio.Reader, label, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", label, defaultValue)
	} else {
		fmt.Printf("%s: ", label)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
}
