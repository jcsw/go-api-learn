package properties

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/jcsw/go-api-learn/pkg/infra/logger"
	"gopkg.in/yaml.v2"
)

// Properties define the properties values
type Properties struct {
	ServerPort string            `yaml:"serverPort"`
	MongoDB    MongoDBProperties `yaml:"mongodb"`
}

// MongoDBProperties define the mongoDB properties values
type MongoDBProperties struct {
	Hosts     []string      `yaml:"hosts"`
	Username  string        `yaml:"username"`
	Password  string        `yaml:"password"`
	Database  string        `yaml:"database"`
	Timeout   time.Duration `yaml:"timeout"`
	PoolLimit int           `yaml:"poolLimit"`
}

// AppProperties the loaded properties values
var AppProperties Properties

// LoadProperties load properties in AppProperties
func LoadProperties(env string) {

	pwd, _ := os.Getwd()
	fileProperties, errReadFile := ioutil.ReadFile(pwd + "/properties/" + env + ".yaml")
	if errReadFile != nil {
		logger.Fatal("f=LoadProperties errReadFile=%v", errReadFile)
	}

	errUnmarshalStrict := yaml.UnmarshalStrict(fileProperties, &AppProperties)
	if errUnmarshalStrict != nil {
		logger.Fatal("f=LoadProperties errUnmarshalStrict=%v", errUnmarshalStrict)
	}

	logger.Info("f=LoadProperties %v", AppProperties)
}
