package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type Config struct {
	AppConf      `yaml:"app"`
	ResourceConf `yaml:"resource"`
}

type AppConf struct {
	Name          string `yaml:"name"`
	Slogan        string `yaml:"slogan"`
	Url           string `yaml:"url"`
	Port          int32  `yaml:"port"`
	Icp           string `yaml:"icp"`
	GithubSecret  string `yaml:"github-secret"`
	GithubAddress string `yaml:"github-address"`
	PageLimit     int    `yaml:"page-limit"`
}

type ResourceConf struct {
	RootDir string `yaml:"root-dir"`
	DocDir  string `yaml:"doc-dir"`
	ImgDir  string `yaml:"img-dir"`
}

var Conf Config

func init() {
	var env string
	if !flag.Parsed() {
		flag.StringVar(&env, "env", "dev", "the server env")
		flag.Parse()
	}
	dir, _ := os.Getwd()
	configFile := fmt.Sprintf("%s/config/conf_%s.yml", dir, env)
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("Config file (%s) not exist", configFile))
	}
	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse yaml config file (%s)", configFile))
	}
	fmt.Println(Conf)
}
