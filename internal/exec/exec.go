package exec

import (
	"fmt"
	"ratkiez/internal/aws"
	"ratkiez/internal/output"
	"ratkiez/internal/types"
	"sync"

	"golang.org/x/sync/errgroup"
)

type CommandConfig struct {
	UserList  []string
	KeyList   []string
	Cmd       string
	OutputFmt string
}

// ExecuteCommand executes the specified command across all clients
func ExecuteCommand(clients []*aws.Client, config CommandConfig) error {
	var (
		allData types.KeyDetailsSlice
		mu      sync.Mutex
		g       errgroup.Group
	)

	for _, c := range clients {
		client := c
		g.Go(func() error {
			data, err := executeCommand(client, config)
			if err != nil {
				return fmt.Errorf("profile %s: %w", c.Profile, err)
			}

			mu.Lock()
			allData = append(allData, data...)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to execute command: %v", err)
	}

	if len(allData) == 0 {
		return fmt.Errorf("No data found")
	}

	formatter, err := output.NewFormatter(config.OutputFmt)
	if err != nil {
		fmt.Errorf("failed to create output formatter: %v", err)
	}
	formatter.Print(allData)

	return nil
}

// executeCommand handles command execution for a single client
func executeCommand(c *aws.Client, config CommandConfig) (types.KeyDetailsSlice, error) {
	switch config.Cmd {
	case "scan":
		return c.ScanCommand()
	case "user":
		return c.UserCommand(config.UserList)
	case "key":
		return c.KeyCommand(config.KeyList)
	default:
		return nil, fmt.Errorf("unknown command: %s", config.Cmd)
	}
}
