// internal/aws/clients.go
package aws

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ratkiez/internal/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"gopkg.in/ini.v1"
)

type Client struct {
	iam     *iam.IAM
	profile string
}

type CommandConfig struct {
	UserList []string
	KeyList  []string
	ScanCmd  string
	UserCmd  string
	KeyCmd   string
}

// NewClients creates multiple AWS clients for given profiles
func NewClients(profiles []string, region string) ([]*Client, error) {
	var clients []*Client
	for _, profile := range profiles {
		client, err := newClient(profile, region)
		if err != nil {
			return nil, fmt.Errorf("failed to create client for profile %s: %v", profile, err)
		}
		clients = append(clients, client)
	}
	return clients, nil
}

// newClient creates a single AWS client
func newClient(profile, region string) (*Client, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		Profile: profile,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return &Client{
		iam:     iam.New(sess),
		profile: profile,
	}, nil
}

// GetProfiles returns list of AWS profiles to use
func GetProfiles(useAllProfiles bool, specifiedProfiles []string) ([]string, error) {
	if useAllProfiles {
		return getAllProfiles()
	}
	return specifiedProfiles, nil
}

// TODO: make aws config path configurable
// getAllProfiles reads all profiles from AWS config
func getAllProfiles() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to get home directory: %v", err)
	}

	configPath := filepath.Join(homeDir, ".aws", "config")
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config file: %v", err)
	}

	var profiles []string
	for _, section := range cfg.Sections() {
		name := section.Name()
		if name == "default" {
			profiles = append(profiles, "default")
		} else if strings.HasPrefix(name, "profile ") {
			profiles = append(profiles, strings.TrimPrefix(name, "profile "))
		}
	}

	if len(profiles) == 0 {
		return nil, fmt.Errorf("no AWS profiles found")
	}
	return profiles, nil
}

// ExecuteCommand executes the specified command across all clients
func ExecuteCommand(cmd string, clients []*Client, config *CommandConfig) (types.KeyDetailsSlice, error) {
	var allData types.KeyDetailsSlice

	for _, client := range clients {
		data, err := client.executeCommand(cmd, config)
		if err != nil {
			// Log error but continue with other clients
			fmt.Printf("Warning: Error processing profile %s: %v\n", client.profile, err)
			continue
		}
		allData = append(allData, data...)
	}

	if len(allData) == 0 {
		return nil, fmt.Errorf("no data collected from any profile")
	}

	return allData, nil
}

// executeCommand handles command execution for a single client
func (c *Client) executeCommand(cmd string, config *CommandConfig) (types.KeyDetailsSlice, error) {
	switch cmd {
	case config.ScanCmd:
		return c.scanCommand()
	case config.UserCmd:
		return c.userCommand(config.UserList)
	case config.KeyCmd:
		return c.keyCommand(config.KeyList)
	default:
		return nil, fmt.Errorf("unknown command: %s", cmd)
	}
}
