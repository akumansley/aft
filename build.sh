npm run-script --prefix ./client/catalog build
go get github.com/markbates/pkger/cmd/pkger
go run github.com/markbates/pkger/cmd/pkger -o ./cmd/aft
go build -o ./bin/aft ./cmd/aft

