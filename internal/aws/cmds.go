// internal/aws/cmds.go
package aws

import (
	"ratkiez/internal/types"
)

func (c *Client) ScanCommand() (types.KeyDetailsSlice, error) {
	users, err := c.listUsers()
	if err != nil {
		return nil, err
	}
	return c.collectByUser(users)
}

func (c *Client) UserCommand(usernames []string) (types.KeyDetailsSlice, error) {
	users, err := c.getUsersByUsernames(usernames)
	if err != nil {
		return nil, err
	}
	return c.collectByUser(users)
}

func (c *Client) KeyCommand(keyIDs []string) (types.KeyDetailsSlice, error) {
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

func (c *Client) Profile() string {
	return c.profile
}
