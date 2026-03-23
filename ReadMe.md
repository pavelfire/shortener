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

curl -X POST http://localhost:8082/url/url \
  -u "myuser:mypass" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'


// this three is enought

  curl -X POST http://localhost:8082/url/url \
  -u "myuser:mypass" \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com", "alias": "test"}'

  curl -i http://localhost:8082/test

  curl -X DELETE -v http://localhost:8082/url/test \
  -u "myuser:mypass"

## Bruno

Коллекция в папке `bruno/` — открыть в Bruno через **File → Open Collection** (указать папку `bruno`). В коллекции отключён прокси (`proxy.use: false`), чтобы запросы к localhost не шли через прокси и не давали ECONNREFUSED 127.0.0.1:443. Выбери окружение **local** для переменных baseUrl, authUser, authPass.