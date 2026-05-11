package ssh

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

// Client wraps an SSH connection to a remote host.
type Client struct {
	User    string
	Host    string
	KeyPath string
	conn    *ssh.Client
}

// NewClient creates a new SSH client configuration.
func NewClient(user, host, keyPath string) *Client {
	return &Client{
		User:    user,
		Host:    host,
		KeyPath: keyPath,
	}
}

// Connect establishes the SSH connection.
func (c *Client) Connect() error {
	key, err := loadKey(c.KeyPath)
	if err != nil {
		return fmt.Errorf("load key %s: %w", c.KeyPath, err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("parse private key: %w", err)
	}

	addr := fmt.Sprintf("%s:22", c.Host)
	config := &ssh.ClientConfig{
		User:            c.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("dial %s: %w", addr, err)
	}
	c.conn = conn
	return nil
}

// Run executes a command over SSH and returns its combined output.
func (c *Client) Run(command string) (string, error) {
	if c.conn == nil {
		return "", fmt.Errorf("not connected; call Connect() first")
	}

	session, err := c.conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return string(output), fmt.Errorf("run command: %w", err)
	}
	return string(output), nil
}

// Close terminates the SSH connection.
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func loadKey(path string) ([]byte, error) {
	// Expand ~ to home directory
	if len(path) > 0 && path[0] == '~' {
		usr, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("get current user: %w", err)
		}
		path = filepath.Join(usr.HomeDir, path[1:])
	}
	return os.ReadFile(path)
}
