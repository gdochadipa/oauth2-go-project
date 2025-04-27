package main

import (
	"log"

	"github.com/gdochadipa/oauth2-go-project/internal/repository"
	"github.com/gdochadipa/oauth2-go-project/pkg/configs"
	"github.com/gdochadipa/oauth2-go-project/pkg/database"
	"github.com/gdochadipa/oauth2-go-project/pkg/server"
	"github.com/gdochadipa/oauth2-go-project/pkg/service"
	"github.com/golang/glog"
)

func main() {
	cfg, err := configs.Load()

	if err != nil {
		glog.Fatalf("Config load failed %v", err.Error())
	}
	// dbnya belum aktif sebelum kita ngebuild grpcnya
	// db, err := database.NewPostgressConnection(cfg.Database)
	// db, err := database.GormDB(cfg.Database)
	db, err := database.InitDBClient(cfg.Database)
	if err != nil {
	 	glog.Fatalf("Failed to initialize the databases. Error: %s", err.Error())
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
