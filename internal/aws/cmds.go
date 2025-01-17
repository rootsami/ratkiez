// internal/aws/cmds.go
package aws

import (
	"fmt"
	"ratkiez/internal/types"
	"sync"

	"golang.org/x/sync/errgroup"
)

// ExecuteCommand executes the specified command across all clients
func ExecuteCommand(cmd string, clients []*Client, config *CommandConfig) (types.KeyDetailsSlice, error) {
	var (
		allData types.KeyDetailsSlice
		mu      sync.Mutex
		g       errgroup.Group
	)

	for _, c := range clients {
		g.Go(func() error {
			data, err := c.executeCommand(cmd, config)
			if err != nil {
				return fmt.Errorf("profile %s: %w", c.profile, err)
			}

			mu.Lock()
			allData = append(allData, data...)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
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

func (c *Client) scanCommand() (types.KeyDetailsSlice, error) {
	users, err := c.listUsers()
	if err != nil {
		return nil, err
	}
	return c.collectByUser(users)
}

func (c *Client) userCommand(usernames []string) (types.KeyDetailsSlice, error) {
	users, err := c.getUsersByUsernames(usernames)
	if err != nil {
		return nil, err
	}
	return c.collectByUser(users)
}

func (c *Client) keyCommand(keyIDs []string) (types.KeyDetailsSlice, error) {
	usernames, err := c.getUsernamesByAccessKeys(keyIDs)
	if err != nil {
		return nil, err
	}
	users, err := c.getUsersByUsernames(usernames)
	if err != nil {
		return nil, err
	}
	return c.collectByUser(users)
}
