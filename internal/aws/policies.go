// internal/aws/policies.go
package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
)

func (c *Client) getAttachedPolicies(user *iam.User) ([]string, error) {
	result, err := c.iam.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{
		UserName: user.UserName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list attached policies for user %s: %v", *user.UserName, err)
	}

	policies := make([]string, len(result.AttachedPolicies))
	for i, policy := range result.AttachedPolicies {
		policies[i] = *policy.PolicyName
	}

	return policies, nil
}
