package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	BaseURL    string
	EndPath    string
	WorkersNum int
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := godotenv.Load("./configs/.env"); err != nil {
		return nil, err
	}
	cfg.BaseURL = os.Getenv("BASE_URL")
	cfg.EndPath = os.Getenv("END_PATH")
	var err error
	cfg.WorkersNum, err = strconv.Atoi(os.Getenv("WORKERS_NUM"))
	if err != nil {
		return nil, err
	}
	log.Printf("configs: %#v", cfg)
	return cfg, nil
}
