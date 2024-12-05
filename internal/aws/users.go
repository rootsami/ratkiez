// internal/aws/users.go
package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/iam"
)

func (c *Client) listUsers() ([]*iam.User, error) {
	var users []*iam.User
	input := &iam.ListUsersInput{}

	err := c.iam.ListUsersPages(input, func(page *iam.ListUsersOutput, lastPage bool) bool {
		users = append(users, page.Users...)
		return !lastPage
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %v", err)
	}

	return users, nil
}

func (c *Client) getUsersByUsernames(usernames []string) ([]*iam.User, error) {
	var users []*iam.User
	for _, username := range usernames {
		result, err := c.iam.GetUser(&iam.GetUserInput{
			UserName: &username,
		})
		if err != nil {
			if strings.Contains(err.Error(), "NoSuchEntity") {
				// slieently ignore users that don't exist
				continue
			}
			return nil, fmt.Errorf("failed to get user %s: %v", username, err)
		}
		users = append(users, result.User)
	}
	return users, nil
}
