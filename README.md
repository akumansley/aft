Generate a local cert:

`brew install mkcert`
`./make-local-ca.sh`

Run the server:

`go run ./cmd/aft -db <path_to_db_file>`

Run catalog:

`cd ./client/catalog`
`npm install`
`npm run dev`

Run explorer:

`cd ./client/explorer`
`npm install`
`npm start`