package config

import (
	"log"
	"os"
	"time"

	"analyze-service/internal/storage/postgres"
	"analyze-service/internal/storage/redis"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Database        postgres.Database `yaml:"database"`
		Redis           redis.Redis       `yaml:"redis"`
		Env             string            `yaml:"env" env-default:"development"`
		MigrationPath   string            `yaml:"migration_path" env-required:"true"`
		MoralisAPIKey   string            `env:"MORALIS_API_KEY"`
		MLServiceAdress string            `yaml:"ml_service_address" env-required:"true"`
		HTTPServer      `yaml:"http_server"`
		CORS            `yaml:"cors"`
	}
	HTTPServer struct {
		Host        string        `yaml:"host" env:"HTTP-IP" env-default:"0.0.0.0"`
		Port        int           `yaml:"port" env:"HTTP-PORT" env-default:"8080"`
		Address     string        `yaml:"address" env-default:"0.0.0.0:8080"`
		Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	}
	CORS struct {
		AllowedMethods     []string `yaml:"allowed-methods" env:"HTTP-CORS-ALLOWED-METHODS"`
		AllowedOrigins     []string `yaml:"allowed-origins"`
		AllowCredentials   bool     `yaml:"allow-credentials"`
		AllowedHeaders     []string `yaml:"allowed-headers"`
		OptionsPassthrough bool     `yaml:"options-passthrough"`
		ExposedHeaders     []string `yaml:"exposed-headers"`
		Debug              bool     `yaml:"debug"`
	}
)

//Также обращаю внимание,
//что путь до конфиг-файла я получаю из переменной окружения CONFIG_PATH,
//дефолтный путь не предусмотрен.
//	Чтобы передать значение такой переменной,
//можно запустить приложение следующей командой:
//CONFIG_PATH=./config/local.yaml ./your-app

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: no .env file found: %v", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("error reading env: %s", err)
	}

	if cfg.MoralisAPIKey == "" {
		log.Fatal("MORALIS_API_KEY is not set")
	}

	return &cfg
}
