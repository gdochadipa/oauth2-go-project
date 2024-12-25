package main

import (
	"fmt"

	"github.com/gdochadipa/oauth2-go-project/pkg/configs"
	"github.com/gdochadipa/oauth2-go-project/pkg/database"
)

func main() {
	cfg, err := configs.Load()

	if err != nil {
		fmt.Sprintln("Config load failed %v", err)
	}

	// db, err := database.NewPostgressConnection(cfg.Database)
	db, err := database.GormDB(cfg.Database)
	if err != nil {
		fmt.Sprintln("Database load failed %v", err)
	}

	dbSQL, err := db.DB()

	if err != nil {
		fmt.Sprintln("Database dbSQL %v", err)
	}

	defer dbSQL.Close()

}
