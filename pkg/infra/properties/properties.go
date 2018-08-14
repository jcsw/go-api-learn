package properties

import (
	"io/ioutil"
	"os"

	"github.com/jcsw/go-api-learn/pkg/infra/logger"
	"gopkg.in/yaml.v2"
)

// Properties define the properties values
type Properties struct {
	ServerPort string `yaml:"serverPort"`
}

// AppProperties the loaded properties values
var AppProperties Properties

// LoadProperties load properties in AppProperties
func LoadProperties(env string) {

	pwd, _ := os.Getwd()
	fileProperties, errReadFile := ioutil.ReadFile(pwd + "/properties/" + env + ".yaml")
	if errReadFile != nil {
		logger.Fatal("f=loadProperties errReadFile=%v", errReadFile)
	}

	errUnmarshalStrict := yaml.UnmarshalStrict(fileProperties, &AppProperties)
	if errUnmarshalStrict != nil {
		logger.Fatal("f=loadProperties errUnmarshalStrict=%v", errUnmarshalStrict)
	}

	logger.Info("f=LoadProperties %v", AppProperties)
}
