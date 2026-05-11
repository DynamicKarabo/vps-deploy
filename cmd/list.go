package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/DynamicKarabo/vps-deploy/internal/config"
	"github.com/DynamicKarabo/vps-deploy/internal/ssh"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show status of all configured services",
	Long: `Connects to each service's host via SSH and checks if the service
is running. Displays a table with Service, Host, and Status columns.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Parse(cfgFile)
		if err != nil {
			log.Fatalf("failed to parse config: %v", err)
		}

		if len(cfg.Services) == 0 {
			fmt.Println("No services configured in", cfgFile)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tHOST\tSTATUS")
		fmt.Fprintln(w, "------\t----\t------")

		for name, svc := range cfg.Services {
			status := checkServiceStatus(svc)
			fmt.Fprintf(w, "%s\t%s\t%s\n", name, svc.Host, status)
		}

		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func checkServiceStatus(svc *config.Service) string {
	client := ssh.NewClient(svc.User, svc.Host, svc.KeyPath)
	if err := client.Connect(); err != nil {
		return fmt.Sprintf("SSH error: %v", err)
	}
	defer client.Close()

	query := fmt.Sprintf(
		"docker ps --filter name=^/%s$ --format '{{.Status}}' 2>/dev/null || echo 'not running'",
		svc.Name,
	)
	out, err := client.Run(query)
	if err != nil {
		return "SSH failed"
	}

	out = trimOutput(out)
	if out == "" {
		return "not running"
	}
	return out
}

func trimOutput(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r' || s[len(s)-1] == ' ') {
		s = s[:len(s)-1]
	}
	return s
}
