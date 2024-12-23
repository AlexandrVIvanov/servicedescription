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

type TypeSecretkey struct {
	Secretkey string `yaml:"secretkey"`
}

// Описываем тип TypeConfigurationServer1c, который будет читаться из файла conf/conf1c.yaml
type TypeConfigurationServer1c struct {
	Сertificateserver1c            string `yaml:"certificateserver1c"`
	Certificatepath1cservicenew    string `yaml:"certificatepath1cservicenew"`
	Сertificatepath1cservicestatus string `yaml:"certificatepath1cservicestatus"`
	Сertificateserver1ctoken       string `yaml:"certificateserver1ctoken"`
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
		log.Printf("Unmarshal config http client: %v", err)
		return nil, err
	}

	return &c, nil
}

func Getconfigsecretkey() (*TypeSecretkey, error) {

	var c TypeSecretkey

	conffile := filepath.Join("conf", "secretkey.yaml")
	yamlFile, err := os.ReadFile(conffile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Printf("Unmarshal config secretkey: %v", err)
		return nil, err
	}

	return &c, nil

}

func GetconfigServer1c() (*TypeConfigurationServer1c, error) {

	var c TypeConfigurationServer1c

	conffile := filepath.Join("conf", "conf1c.yaml")
	yamlFile, err := os.ReadFile(conffile)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Printf("Unmarshal config server 1c: %v", err)
		return nil, err
	}

	return &c, nil

}
