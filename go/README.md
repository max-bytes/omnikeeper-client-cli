## Requirements

- credentials helper: https://github.com/docker/docker-credential-helpers
  - one of the binary credential helpers in https://github.com/docker/docker-credential-helpers/releases, present in PATH

## Examples

Login (password via prompt):
~~~bash
go run cmd/main.go -o https://10.0.0.43:45455 -u username
~~~

Login (password via parameter):
~~~bash
go run cmd/main.go -o https://10.0.0.43:45455 -u username -p password
~~~

Run query from parameter:
~~~bash
go run cmd/main.go -o https://10.0.0.43:45455 -q "query { activeTraits { id } }"
~~~

Run query from stdin:
~~~bash
cat <<'EOF' | go run cmd/main.go -o https://10.0.0.43:45455 --stdin
query {
  activeTraits {
    id
  }
}
EOF
~~~