package env

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

type Config struct {
	AtlasURI            string `env:"ATLAS_URI,required"`
	ApplicationRootPath string `env:"APPLICATION_ROOT_PATH"`
	ServerBaseURL       string `env:"SERVER_BASE_URL,default=http://localhost"`
	ServerPort          uint16 `env:"SERVER_PORT,default=8081"`
	WorkSpacePath       string `env:"WORK_SPACE_PATH,default=8081"`
	JWTKEY              string `env:"JWTKEY,required"`
}

func LoadEnvironment() (Config, error) {
	var cfg Config

	err := godotenv.Load()
	if err != nil {
		log.Printf("could not load environment file: %v", err)
	}

	err = envdecode.Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	if cfg.ApplicationRootPath == "" {
		cfg.ApplicationRootPath = filepath.Join(os.Getenv("GOPATH"), cfg.WorkSpacePath)
		cfg.ApplicationRootPath = filepath.Join(cfg.ApplicationRootPath, "default-build")
	}

	return cfg, nil
}
