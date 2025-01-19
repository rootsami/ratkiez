package aws

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

// GetProfiles returns list of AWS profiles to use
func GetProfiles(useAllProfiles bool, specifiedProfiles []string) ([]string, error) {
	if useAllProfiles {
		return getAllProfiles()
	}
	return specifiedProfiles, nil
}

// TODO: make aws config path configurable
// getAllProfiles reads all profiles from AWS config
func getAllProfiles() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to get home directory: %v", err)
	}

	configPath := filepath.Join(homeDir, ".aws", "config")
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config file: %v", err)
	}

	var profiles []string
	for _, section := range cfg.Sections() {
		name := section.Name()
		if name == "default" {
			profiles = append(profiles, "default")
		} else if strings.HasPrefix(name, "profile ") {
			profiles = append(profiles, strings.TrimPrefix(name, "profile "))
		}
	}

	if len(profiles) == 0 {
		return nil, fmt.Errorf("no AWS profiles found")
	}
	return profiles, nil
}
