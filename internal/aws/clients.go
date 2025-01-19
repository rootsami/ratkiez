// internal/aws/clients.go
package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type CommandConfig struct {
	UserList []string
	KeyList  []string
	ScanCmd  string
	UserCmd  string
	KeyCmd   string
}

type Client struct {
	session     *session.Session
	iam         *iam.IAM
	profile     string
	accountID   string
	accountName string
}

// NewClients creates multiple AWS clients for given profiles
func NewClients(profiles []string, region string, isOrg bool, orgRole string) ([]*Client, error) {
	var allClients []*Client

	for _, profile := range profiles {
		// Create the main client for this profile
		client, err := newClient(profile, region)
		if err != nil {
			return nil, fmt.Errorf("failed to create client for profile %s: %v", profile, err)
		}

		// Add the main client
		allClients = append(allClients, client)

		// If org flag is set, create clients for all member accounts
		if isOrg {
			memberClients, err := client.createMemberAccountClients(orgRole)
			if err != nil {
				// Log the error but continue with next profile
				fmt.Printf("Warning: Failed to get member accounts for profile %s: %v\n", profile, err)
				continue
			}
			if len(memberClients) > 0 {
				allClients = append(allClients, memberClients...)
			}
		}
	}

	if len(allClients) == 0 {
		return nil, fmt.Errorf("no valid clients could be created")
	}

	return allClients, nil
}

// newClient creates a single AWS client
func newClient(profile, region string) (*Client, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	client := &Client{
		session: sess,
		iam:     iam.New(sess),
		profile: profile,
	}

	return client, nil
}
