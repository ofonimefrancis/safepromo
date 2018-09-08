package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/globalsign/mgo"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
	config     Config
)

const (
	DATABASENAME     = "safepromo"
	EVENTSCOLLECTION = "events"
)

//Config Our system configuration
type Config struct {
	MongoDB     string
	MongoServer string
	Port        string
	Session     *mgo.Session
}

//Init Configures the Config struct
func Init() {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{"127.0.0.1:27017", "127.0.0.1:27018"},
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWD"),
	})

	config = Config{}

	if err != nil {
		log.Fatal("Error connecting to Database")
	}

	config.Session = session
	config.MongoDB = DATABASENAME

	if os.Getenv("PORT") == "" {
		config.Port = "3000"
	}
}

//Get Retrieves the Config struct
func Get() *Config {
	return &config
}
