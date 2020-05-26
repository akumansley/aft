Generate a local cert:

`brew install mkcert`
`./make-local-ca.sh`

Run the server:

`go run ./cmd/aft -db <path_to_db_file>`

Run the client:

`npm run dev`
