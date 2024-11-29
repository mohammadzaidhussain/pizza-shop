package config

import (
	"fmt"
	"os"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/mohammadzaidhussain/pizza-shop/logger"
)

var env ConfigDto

type ConfigDto struct {
	port                    string
	rabbit_mq_host          string
	rabbit_mq_username      string
	rabbit_mq_password      string
	rabbit_mq_port          string
	rabbit_mq_default_queue string
}

func ConfigEnv() {
	LoadEnvVariable()
	env = ConfigDto{
		port:                    os.Getenv("PORT"),
		rabbit_mq_host:          os.Getenv("RABBIT_MQ_HOST"),
		rabbit_mq_username:      os.Getenv("RABBIT_MQ_USERNAME"),
		rabbit_mq_password:      os.Getenv("RABBIT_MQ_PASSWORD"),
		rabbit_mq_port:          os.Getenv("RABBIT_MQ_PORT"),
		rabbit_mq_default_queue: os.Getenv("RABBIT_MQ_DEFAULT_QUEUE"),
	}
}

func init() {
	if env.port == "" {
		ConfigEnv()
	}
}

func accessField(key string) (string, error) {
	v := reflect.ValueOf(env)
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("expected struct")
	}

	_, ok := t.FieldByName(key)
	if !ok {
		return "", fmt.Errorf("field not found")
	}

	f := v.FieldByName(key)
	return f.String(), nil
}

func LoadEnvVariable() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			panic("error loading env")
		}
	} else {
		logger.Log(" no .env file found, its using sytem env conf")
	}
}

func GetEnvProperty(propertyKey string) string {
	if env.port == "" {
		ConfigEnv()
	}
	val, err := accessField(propertyKey)
	if err != nil {
		logger.Log(fmt.Sprintf("error accesing field : %v", propertyKey))
	}
	return val
}
