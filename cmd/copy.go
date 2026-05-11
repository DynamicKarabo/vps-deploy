package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/DynamicKarabo/vps-deploy/internal/config"
)

var copyCmd = &cobra.Command{
	Use:   "copy [service] [local_path] [remote_path]",
	Short: "Copy a file to a service's host via SCP",
	Long: `Uses SCP to copy a local file to the remote host configured for a service.
Uses the SSH key and user from the service's deploy.yaml config.

Example:
  vps-deploy copy filebrowser ./docker-compose.yml /root/filebrowser/docker-compose.yml`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		svcName := args[0]
		localPath := args[1]
		remotePath := args[2]

		cfg, err := config.Parse(cfgFile)
		if err != nil {
			log.Fatalf("failed to parse config: %v", err)
		}

		svc, ok := cfg.Services[svcName]
		if !ok {
			log.Fatalf("service %q not found in config", svcName)
		}

		// Check local file exists
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			log.Fatalf("local file %q does not exist", localPath)
		}

		// Expand ~ in key path
		keyPath := svc.KeyPath
		if len(keyPath) > 0 && keyPath[0] == '~' {
			home, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("could not determine home dir: %v", err)
			}
			keyPath = home + keyPath[1:]
		}

		dest := fmt.Sprintf("%s@%s:%s", svc.User, svc.Host, remotePath)

		fmt.Printf("Copying %s → %s\n", localPath, dest)

		scpCmd := exec.Command("scp", "-i", keyPath, localPath, dest)
		scpCmd.Stdout = os.Stdout
		scpCmd.Stderr = os.Stderr

		if err := scpCmd.Run(); err != nil {
			log.Fatalf("SCP failed: %v", err)
		}

		fmt.Println("✅ Copy complete")
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
