package main

import (
	"fmt"
	"log"

	"github.com/gdochadipa/oauth2-go-project/internal/repository"
	"github.com/gdochadipa/oauth2-go-project/pkg/configs"
	"github.com/gdochadipa/oauth2-go-project/pkg/database"
	"github.com/gdochadipa/oauth2-go-project/pkg/server"
	"github.com/gdochadipa/oauth2-go-project/pkg/service"
)

func main() {
	cfg, err := configs.Load()

	if err != nil {
		fmt.Sprintln("Config load failed %v", err)
	}
	// dbnya belum aktif sebelum kita ngebuild grpcnya
	// kayaknya pakein mutex atau gorotine
	db, err := database.NewPostgressConnection(cfg.Database)
	// db, err := database.GormDB(cfg.Database)
	if err != nil {
		fmt.Sprintln("Database load failed %v", err)
	}

	// dbSQL, err := db.DB()

	// if err != nil {
	// 	fmt.Sprintln("Database dbSQL %v", err)
	// }

	r := repository.NewDBRepository(db)
	defer db.Close()

	jwt := service.NewJWTRepository([]byte("testing"))

	log.Println("Listening on port 8080...")
	s := service.NewGrantService(r, jwt)
	log.Fatal(server.ListenGRPC(s, 8080))

}
