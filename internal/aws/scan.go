// internal/aws/scan.go
package aws

import (
	"fmt"

	"ratkiez/internal/types"

	"github.com/aws/aws-sdk-go/service/iam"
)

func (c *Client) collectByUser(users []*iam.User) (types.KeyDetailsSlice, error) {
	var details types.KeyDetailsSlice

	for _, user := range users {
		userDetails, err := c.getUserDetails(user)
		if err != nil {
			return nil, fmt.Errorf("failed to get details for user %s: %v", *user.UserName, err)
		}
		details = append(details, userDetails...)
	}

	return details, nil
}

func (c *Client) getUserDetails(user *iam.User) (types.KeyDetailsSlice, error) {
	var details types.KeyDetailsSlice

	// Get the attached policies for the user
	policies, err := c.getAttachedPolicies(user)
	if err != nil {
		return nil, err
	}

	// Get the access keys for the user
	keys, err := c.listAccessKeys(user)
	if err != nil {
		return nil, err
	}

	// If the user has no keys, add a placeholder entry
	if len(keys) == 0 {
		details = append(details, types.KeyDetails{
			User:         *user.UserName,
			KeyID:        "N/A",
			CreationDate: "N/A",
			LastUsedDate: "N/A",
			Policies:     policies,
			Profile:      c.profile,
		})
		return details, nil
	}

	for _, key := range keys {
		lastUsed, err := c.getAccessKeyLastUsed(key)
		if err != nil {
			return nil, fmt.Errorf("failed to get last used info for key %s: %v", *key.AccessKeyId, err)
		}

		details = append(details, types.KeyDetails{
			User:         *user.UserName,
			KeyID:        *key.AccessKeyId,
			CreationDate: key.CreateDate.String(),
			LastUsedDate: lastUsed,
			Policies:     policies,
			Profile:      c.profile,
		})
	}

	return details, nil
}
