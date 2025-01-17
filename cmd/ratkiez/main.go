// cmd/ratkiez/main.go
package main

import (
	"fmt"
	"log"
	"os"

	"ratkiez/internal/aws"
	"ratkiez/internal/output"

	"github.com/alecthomas/kingpin"
)

var (
	app = kingpin.New("ratkiez", "A CLI tool to rat on all AWS keys based on creation date and last used date")

	// Global flags
	region      = app.Flag("region", "AWS region").Default("us-west-2").String()
	profiles    = app.Flag("profile", "AWS profiles, reusable to add more profiles").Default("default").Strings()
	allProfiles = app.Flag("all-profiles", "Use all profiles in ~/.aws/config").Bool()
	outputFmt   = app.Flag("format", "Output format, json, table or csv").Default("table").Enum("table", "json", "csv")
	isOrg       = app.Flag("org", "Scan all organization member accounts").Bool()
	orgRole     = app.Flag("role-name", "Role name to assume in organization member accounts").Default("OrganizationAccountAccessRole").String()

	// Commands
	scan     = app.Command("scan", "Scan all AWS keys. ex: ratkiez scan --profile profile1 --profile profile2")
	user     = app.Command("user", "Scan by username(s), ex: ratkiez user john.doe jane.doe --profile profile1")
	key      = app.Command("key", "Scan by key-id(s), ex: ratkiez key AKIA1234 AKIA5678 --all-profiles")
	userList = user.Arg("user", "List of users to scan").Strings()
	keyList  = key.Arg("key", "List of keys to scan").Strings()
)

func executeCommand(cmd string, clients []*aws.Client) error {
	data, err := aws.ExecuteCommand(cmd, clients, &aws.CommandConfig{
		UserList: *userList,
		KeyList:  *keyList,
		ScanCmd:  scan.FullCommand(),
		UserCmd:  user.FullCommand(),
		KeyCmd:   key.FullCommand(),
	})
	if err != nil {
		return fmt.Errorf("command execution failed: %v", err)
	}

	formatter, err := output.NewFormatter(*outputFmt)
	if err != nil {
		return fmt.Errorf("failed to create output formatter: %v", err)
	}

	return formatter.Print(data)
}

func main() {
	cmd := kingpin.MustParse(app.Parse(os.Args[1:]))

	// prevent using both --all-profiles and --org flags
	if *allProfiles && *isOrg {
		log.Fatal("Cannot use --all-profiles and --org flags together")
	}

	profiles, err := aws.GetProfiles(*allProfiles, *profiles)
	if err != nil {
		log.Fatalf("Failed to get profiles: %v", err)
	}

	clients, err := aws.NewClients(profiles, *region, *isOrg, *orgRole)
	if err != nil {
		log.Fatalf("Failed to create clients: %v", err)
	}

	if err := executeCommand(cmd, clients); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
