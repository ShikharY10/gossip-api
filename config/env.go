package config

import (
	"os"

	"github.com/ShikharY10/gbAPI/utils"
	"github.com/joho/godotenv"
)

type ENV struct {
	MongoDBConnectionMethod     string // manual
	MongoDBPort                 string // 27017
	MongoDBHost                 string // 127.0.0.1
	MongoDBUsername             string // rootuser
	MongoDBPassword             string // rootpass
	MongoDBConnectionString     string // mongodb connection string will be used when MongoDBConnectionMethod is set to auto
	RedisHost                   string // 127.0.0.1
	RedisPort                   string // 6379
	RabbitMQHost                string // 127.0.0.1
	RabbitMQPort                string // 5672
	RabbitMQUsername            string // guest
	RabbitMQPassword            string // guest
	APIPort                     string // 6001
	APIName                     string // GT____
	APIMode                     string // debug
	JWT_ACCESS_TOKEN_SECRET_KEY string // abcdefghijklmnopqrstuvwxyz
	LogServerHost               string // 127.0.0.1
	LogServerPort               string // 6002
}

func LoadENV() *ENV {
	godotenv.Load()
	var env ENV

	var value string
	var found bool

	value, found = os.LookupEnv("MONGODB_CONNECTION_METHOD")
	if found {
		env.MongoDBConnectionMethod = value
	} else {
		env.MongoDBConnectionMethod = "manual"
	}

	value, found = os.LookupEnv("MONGODB_PORT")
	if found {
		env.MongoDBPort = value
	} else {
		env.MongoDBPort = "27017"
	}

	value, found = os.LookupEnv("MONGODB_HOST")
	if found {
		env.MongoDBHost = value
	} else {
		env.MongoDBHost = "127.0.0.1"
	}

	value, found = os.LookupEnv("MONGODB_USERNAME")
	if found {
		env.MongoDBUsername = value
	} else {
		env.MongoDBUsername = "rootuser"
	}

	value, found = os.LookupEnv("MONGODB_PASSWORD")
	if found {
		env.MongoDBPassword = value
	} else {
		env.MongoDBPassword = "rootpass"
	}

	value, found = os.LookupEnv("MONGODB_CONNECTION_STRING")
	if found {
		env.MongoDBConnectionString = value
	} else {
		env.MongoDBConnectionString = ""
	}

	value, found = os.LookupEnv("REDIS_HOST")
	if found {
		env.RedisHost = value
	} else {
		env.RedisHost = "127.0.0.1"
	}

	value, found = os.LookupEnv("REDIS_PORT")
	if found {
		env.RedisPort = value
	} else {
		env.RedisPort = "6379"
	}

	value, found = os.LookupEnv("RabbitMQHost")
	if found {
		env.RabbitMQHost = value
	} else {
		env.RabbitMQHost = "127.0.0.1"
	}

	value, found = os.LookupEnv("RABBITMQ_PORT")
	if found {
		env.RabbitMQPort = value
	} else {
		env.RabbitMQPort = "5672"
	}

	value, found = os.LookupEnv("RABBITMQ_USERNAME")
	if found {
		env.RabbitMQUsername = value
	} else {
		env.RabbitMQUsername = "guest"
	}

	value, found = os.LookupEnv("RABBITMQ_PASSWORD")
	if found {
		env.RabbitMQPassword = value
	} else {
		env.RabbitMQPassword = "guest"
	}

	value, found = os.LookupEnv("API_PORT")
	if found {
		env.APIPort = value
	} else {
		env.APIPort = "10221"
	}

	value, found = os.LookupEnv("API_NAME")
	if found {
		env.APIName = value
	} else {
		env.APIName = "API_" + utils.Encode(utils.GenerateAesKey(10))
	}

	value, found = os.LookupEnv("API_MODE")
	if found {
		env.APIMode = value
	} else {
		env.APIMode = "debug"
	}

	value, found = os.LookupEnv("JWT_ACCESS_TOKEN_SECRET_KEY")
	if found {
		env.JWT_ACCESS_TOKEN_SECRET_KEY = value
	} else {
		env.JWT_ACCESS_TOKEN_SECRET_KEY = "982u3923jhdwhe3fjdw30fj02j3ijwef023jfijwjf802j300"
	}

	value, found = os.LookupEnv("LOG_SERVER_HOST")
	if found {
		env.LogServerHost = value
	} else {
		env.LogServerHost = "127.0.0.1"
	}

	value, found = os.LookupEnv("LOG_SERVER_PORT")
	if found {
		env.LogServerPort = value
	} else {
		env.LogServerPort = "10223"
	}

	return &env
}
