// pkg/types/types.go
package types

type KeyDetails struct {
	User         string   `json:"user"`
	KeyID        string   `json:"key_id"`
	CreationDate string   `json:"creation_date"`
	LastUsedDate string   `json:"last_used_date"`
	Policies     []string `json:"policies"`
	Profile      string   `json:"profile"`
	AccountID    string   `json:"account_id"`
	AccountName  string   `json:"account_name"`
}

type KeyDetailsSlice []KeyDetails
