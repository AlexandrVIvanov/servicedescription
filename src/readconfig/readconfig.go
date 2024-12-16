package readconfig

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

type TypeSqlConfiguration struct {
	ServerName string `yaml:"servername"`
	UserName   string `yaml:"sqlusername"`
	Password   string `yaml:"sqlpassword"`
	Port       int    `yaml:"sqlport"`
	Database   string `yaml:"sqldatabase"`
}

type TypeHTTPclientConnectionConfig struct {
	URLServerName string `yaml:"urlservername"`
	URLPath       string `yaml:"urlpath"`
}

func Getconfigsqlserver() (*TypeSqlConfiguration, error) {

	var c TypeSqlConfiguration

	conffile := filepath.Join("conf", "sqlconf.yaml")
	yamlFile, err := os.ReadFile(conffile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Printf("Unmarshal config sql server: %v", err)
		return nil, err
	}

	return &c, nil
}

func Getconfighttpclient() (*TypeHTTPclientConnectionConfig, error) {

	var c TypeHTTPclientConnectionConfig

	conffile := filepath.Join("conf", "httpclient.yaml")
	yamlFile, err := os.ReadFile(conffile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Printf("Unmarshal config sql server: %v", err)
		return nil, err
	}

	return &c, nil
}
