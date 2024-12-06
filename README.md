# Ratkiez

A CLI tool to rat on all aws keys based on creation date, last used date, and attached policies.

Output is supported in multiple formats: json, table, and csv.

## Prerequisites

- [Go](https://golang.org/doc/install)
- Configured AWS credentials
- Set `export AWS_SDK_LOAD_CONFIG=1` in your shell profile

## Usage

```bash
usage: ratkiez [<flags>] <command> [<args> ...]

A CLI tool to rat on all AWS keys based on creation date and last used date

Flags:
  --help                 Show context-sensitive help (also try --help-long and --help-man).
  --region="us-west-2"   AWS region
  --profile=default ...  AWS profiles, reusable to add more profiles
  --all-profiles         Use all profiles in ~/.aws/config
  --format=table         Output format, json, table or csv

Commands:
  help [<command>...]
    Show help.

  scan
    Scan all AWS keys. ex: ratkiez scan --profile profile1 --profile profile2

  user [<user>...]
    Scan by username(s), ex: ratkiez user john.doe jane.doe --profile profile1

  key [<key>...]
    Scan by key-id(s), ex: ratkiez key AKIA1234 AKIA5678 --all-profiles
```

## Examples

### Scan Single Account Profile
```bash
# Scan all users in the specified aws account
ratkiez scan --profile aws-profile-eu-central-1 --format table

# Scan all users in multiple aws account
ratkiez scan --profile aws-profile-eu-central-1 --profile aws-profile-us-west-2 --format table

```

### Scan All Profiles
```bash
# Scan all users in all aws accounts configured in ~/.aws/config
ratkiez scan --all-profiles --format table
```

Sample output:
```
USERNAME                                KEY-ID                  CREATION-DATE                    LAST-USED-DATE                 POLICIES                               PROFILE
xxxxx-sns-user                          AKIASWXXXXXXXXXXXX      2021-02-15 10:53:57 +0000 UTC    Never Used                     Access_Extension_Lambda                aws-profile-eu-central-1
s3-controller                           AKIASWXXXXXXXXXXXX      2020-05-15 08:07:18 +0000 UTC    2020-10-15 08:30:00 +0000 UTC  AmazonS3FullAccess                     aws-profile-us-west-2
example-lambda-user                     AKIASWXXXXXXXXXXXX      2021-02-15 10:53:57 +0000 UTC    Never Used                     Access_Extension_Lambda                aws-profile-us-west-2
```

### Look Up Specific User
```bash
# Look up specific users in one account
ratkiez user example-lambda-user s3-controller --profile aws-profile-us-west-2 --format json

# Look up user across all profiles
ratkiez user example-lambda-user --all-profiles --format json
```

### Look Up Specific Key
```bash
# Look up a specific access key
ratkiez key AKIASWXXXXXXXXXXXXXX --profile aws-playground-eu --format json

# Look up specific keys across all profiles
ratkiez key AKIASWXXXXXXXXXXXXXX AKIASWXXXXXXXXXXXXXX --all-profiles --format json
```

## Installation

### Binary

Download the binary from the [releases](https://github.com/rootsami/ratkiez/releases)

### Build from source

```bash

$ git clone https://github.com/rootsami/ratkiez.git
$ go build -o ratkiez cmd/ratkiez/main.go

```

## License

## Contributing
