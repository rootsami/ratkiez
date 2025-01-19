// internal/aws/organizations.go
package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/sts"
)

type AccountInfo struct {
	ID   string
	Name string
}

func (c *Client) getOrganizationAccounts() ([]AccountInfo, error) {
	if c.session == nil {
		return nil, fmt.Errorf("invalid session")
	}

	orgsSvc := organizations.New(c.session)

	// First, get the management account ID to skip the role assumption
	describeOrg, err := orgsSvc.DescribeOrganization(&organizations.DescribeOrganizationInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to describe organization: %v", err)
	}
	managementAccountID := *describeOrg.Organization.MasterAccountId

	var accounts []AccountInfo
	input := &organizations.ListAccountsInput{}

	err = orgsSvc.ListAccountsPages(input, func(page *organizations.ListAccountsOutput, lastPage bool) bool {
		for _, account := range page.Accounts {
			// Skip the management account
			if *account.Id == managementAccountID {
				continue
			}

			if account.Id != nil && account.Name != nil {
				accounts = append(accounts, AccountInfo{
					ID:   *account.Id,
					Name: *account.Name,
				})
			}
		}
		return !lastPage
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == organizations.ErrCodeAWSOrganizationsNotInUseException {
				return nil, fmt.Errorf("AWS Organizations is not enabled for this account")
			}
		}
		return nil, fmt.Errorf("failed to list organization accounts: %v", err)
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("no member accounts found in the organization")
	}

	return accounts, nil
}

// assumeRole assumes the specified role in a target account and returns a new session
func (c *Client) assumeRole(accountID, roleName string) (*session.Session, error) {
	if c.session == nil {
		return nil, fmt.Errorf("invalid session")
	}

	stsSvc := sts.New(c.session)

	roleARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, roleName)
	sessionName := fmt.Sprintf("RatkiezScan-%s", c.profile)

	result, err := stsSvc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(roleARN),
		RoleSessionName: aws.String(sessionName),
		DurationSeconds: aws.Int64(900), // 15 minutes
	})
	if err != nil {
		return nil, fmt.Errorf("failed to assume role in account %s: %v", accountID, err)
	}

	// Create a new session with the assumed role credentials
	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			*result.Credentials.AccessKeyId,
			*result.Credentials.SecretAccessKey,
			*result.Credentials.SessionToken,
		),
		Region: c.session.Config.Region,
	})
}

// createMemberAccountClients creates IAM clients for all member accounts
func (c *Client) createMemberAccountClients(roleName string) ([]*Client, error) {
	if c.session == nil {
		return nil, fmt.Errorf("invalid session")
	}

	accounts, err := c.getOrganizationAccounts()
	if err != nil {
		return nil, err
	}

	var clients []*Client
	for _, account := range accounts {
		sess, err := c.assumeRole(account.ID, roleName)
		if err != nil {
			fmt.Printf("Warning: Failed to assume role in account %s (%s): %v\n", account.Name, account.ID, err)
			continue
		}

		clients = append(clients, &Client{
			session:     sess,
			iam:         iam.New(sess),
			profile:     fmt.Sprintf("member-of-%s", c.profile),
			accountID:   account.ID,
			accountName: account.Name,
		})
	}

	if len(clients) == 0 {
		return nil, fmt.Errorf("no member account clients could be created")
	}

	return clients, nil
}
