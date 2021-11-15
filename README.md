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