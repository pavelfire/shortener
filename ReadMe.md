go get github.com/ilyakaznacheev/cleanenv

cd cmd/shortener
go build
go run

cd ..
cd ..
go run ./cmd/shortener

export CONFIG_PATH=./config/local.yaml
go run ./cmd/shortener

go get "github.com/mattn/go-sqlite3"
go get github.com/go-chi/chi/v5
go get github.com/go-chi/render
go get github.com/go-playground/validator/v10

go mod tidy

curl -X POST http://localhost:8082/url \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'

  curl -X POST http://localhost:8082/url \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com", "alias": "my-link"}'

  curl -X POST http://localhost:8082/url \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com", "alias": "test"}'

  curl -i http://localhost:8082/test