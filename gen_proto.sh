protoc --go_out=. --go_opt=paths=source_relative ./internal/data/object.proto
protoc --go_out=. --go_opt=paths=source_relative ./internal/server/services/objects.proto
