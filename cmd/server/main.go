package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/devSLAVUS/yagometrix22/internal/server/config"
	"github.com/devSLAVUS/yagometrix22/internal/server/handlers"
	"github.com/devSLAVUS/yagometrix22/internal/server/router"
	"github.com/devSLAVUS/yagometrix22/internal/server/storage"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Printf("Error parsing environment variables: %v\n", err)
		return
	}

	address := flag.String("a", cfg.Address, "server ip:port")
	flag.Parse()
	if *address != cfg.Address {
		cfg.Address = *address
	}

	store := storage.NewMemStorage()
	handler := handlers.NewHandlers(store)

	r := router.NewRouter(handler)

	if err := r.Run(cfg.Address); err != nil {
		panic(err)
	}
	fmt.Println("Server started:", time.Now())
}
