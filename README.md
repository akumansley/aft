# Building Aft 

To build Aft from source, install go 1.16+ and npm, then:

```bash
npm install --prefix ./client/catalog
npm run-script --prefix ./client/catalog build
go build -o ./bin/aft ./cmd/aft
./bin/aft -db <path_to_db_file> -authed=false
```

# Developing Aft

Run the server:

```bash
go run ./cmd/aft -db <path_to_db_file> -authed=false -serve_dir client/catalog/public
```

Run catalog:

```bash
cd ./client/catalog
npm install
npm run dev
```
