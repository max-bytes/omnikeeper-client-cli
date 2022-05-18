# omnikeeper CLI GraphQL client

Small CLI application for interacting with an omnikeeper instance via GraphQL queries and mutations. Can be used to automate the configuration of an omnikeeper instance.

## Requirements

- (optional) credentials helper: https://github.com/docker/docker-credential-helpers
  - one of the binary credential helpers in https://github.com/docker/docker-credential-helpers/releases, present in PATH

## Installation

~~~bash
wget https://github.com/max-bytes/omnikeeper-client-cli/releases/latest/download/okql.rpm
# either:
sudo rpm -i okql.rpm
# or:
sudo yum localinstall okql.rpm
# or:
sudo alien -i okql.rpm
~~~

## Notes
In terms of authentication, only the oauth flow "Resource Owner Password Credentials Grant" is supported. Other modes, such as the "Device Authorization Flow" are not (yet) supported. That means it is necessary to supply username+password of the user to authenticate.


## Examples

Note: to run from source instead, replace `okql` with `go run cmd/main.go`

Login (password via prompt):
~~~bash
okql -o https://10.0.0.43:45455 -u username
~~~

Login (password via parameter):
~~~bash
okql -o https://10.0.0.43:45455 -u username -p password
~~~

Run query from parameter:
~~~bash
okql -o https://10.0.0.43:45455 -q "query { activeTraits { id } }"
~~~

Run query from stdin:
~~~bash
cat <<'EOF' | okql -o https://10.0.0.43:45455 --stdin
query {
  activeTraits {
    id
  }
}
EOF
~~~

Run query from file:
~~~bash
okql -o https://10.0.0.43:45455 --stdin < file.graphql
~~~