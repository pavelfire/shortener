go get github.com/ilyakaznacheev/cleanenv

cd cmd/shortener
go build
go run

cd ..
cd ..
go run ./cmd/shortener

export CONFIG_PATH=./config/local.yaml
go run ./cmd/shortener