go get github.com/ilyakaznacheev/cleanenv

cd cmd/shortener
go build
go run

cd ..
cd ..
go run ./cmd/shortener