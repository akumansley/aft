# Building Aft 

To build Aft from source, install go and npm, then:

```bash
npm install --prefix ./client/catalog
npm run-script --prefix ./client/catalog build
go get github.com/markbates/pkger/cmd/pkger
go run github.com/markbates/pkger/cmd/pkger -o ./cmd/aft
go build -o ./bin/aft ./cmd/aft
./bin/aft -db <path_to_db_file> -authed=false
```

# Developing Aft

Run the server:

```bash
go run ./cmd/aft -db <path_to_db_file> -authed=false
```

Run catalog:

```bash
cd ./client/catalog
npm install
npm run dev
```