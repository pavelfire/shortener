package main

import (
	"fmt"
	"shortener/internal/config"
)

func main(){
	cfg:=config.MustLoad()
	fmt.Println(cfg)
	//TODO: init logger: log/slog
	//TODO: init storage sqlite
	//TODO: init router: chi, "chi render"
	//TODO: start server
}