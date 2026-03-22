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