package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env           string              `yaml:"env"`
	GRPC          GRPCConfig          `yaml:"grpc"`
	Websocket     WebsocketConfig     `yaml:"websocket"`
	MessageBroker MessageBrokerConfig `yaml:"message_broker"`
	HTTP          HTTPConfig          `yaml:"HTTP"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type WebsocketConfig struct {
	Port     int    `yaml:"port"`
	KeyFile  string `yaml:"keyfile"`
	CertFile string `yaml:"certfile"`
}

type MessageBrokerConfig struct {
	Port            int    `yaml:"port"`
	TranscribeQueue string `yaml:"transcribe_queue"`
	ProcessQueue    string `yaml:"process_queue"`
}

type HTTPConfig struct {
	Address  string        `yaml:"address"`
	TokenTTL time.Duration `yaml:"tokenTTL"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	err := godotenv.Load(fmt.Sprintf(".env.%s", cfg.Env))
	if err != nil {
		panic("failed to load environment variables:  " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
