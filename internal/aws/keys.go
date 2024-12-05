// internal/aws/keys.go
package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/iam"
)

func (c *Client) getUsernamesByAccessKeys(keyIDs []string) ([]string, error) {
	var usernames []string
	for _, keyID := range keyIDs {
		result, err := c.iam.GetAccessKeyLastUsed(&iam.GetAccessKeyLastUsedInput{
			AccessKeyId: &keyID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "AccessDenied") {
				// silently ignore keys that we can't access, as they may be from other accounts
				continue
			}
			return nil, fmt.Errorf("failed to get key %s: %v", keyID, err)
		}
		usernames = append(usernames, *result.UserName)
	}
	return usernames, nil
}

func (c *Client) listAccessKeys(user *iam.User) ([]*iam.AccessKeyMetadata, error) {
	result, err := c.iam.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: user.UserName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list access keys for user %s: %v", *user.UserName, err)
	}
	return result.AccessKeyMetadata, nil
}

func (c *Client) getAccessKeyLastUsed(key *iam.AccessKeyMetadata) (string, error) {
	result, err := c.iam.GetAccessKeyLastUsed(&iam.GetAccessKeyLastUsedInput{
		AccessKeyId: key.AccessKeyId,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get last used info: %v", err)
	}

	if result.AccessKeyLastUsed.LastUsedDate != nil {
		return result.AccessKeyLastUsed.LastUsedDate.String(), nil
	}
	return "Never Used", nil
}
