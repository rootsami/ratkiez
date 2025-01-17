// internal/aws/scan.go
package aws

import (
	"fmt"
	"sync"

	"ratkiez/internal/types"

	"github.com/aws/aws-sdk-go/service/iam"
	"golang.org/x/sync/errgroup"
)

func (c *Client) collectByUser(users []*iam.User) (types.KeyDetailsSlice, error) {
	var (
		details types.KeyDetailsSlice
		mu      sync.Mutex
		g       errgroup.Group
	)

	for _, user := range users {
		g.Go(func() error {
			userDetails, err := c.getUserDetails(user)
			if err != nil {
				return fmt.Errorf("failed to get details for user %s: %v", *user.UserName, err)
			}
			mu.Lock()
			details = append(details, userDetails...)
			mu.Unlock()
			return nil
		})

	}

	if err := g.Wait(); err != nil {
		return nil, err
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
			AccountID:    c.accountID,
			AccountName:  c.accountName,
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
			AccountID:    c.accountID,
			AccountName:  c.accountName,
		})
	}

	return details, nil
}
