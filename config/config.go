package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

const (
	PasswordType   = "password"
	PrivateKeyType = "privateKey"
	Gitee          = "gitee"
	Github         = "github"
)

type Config struct {
	AppConf      `yaml:"app"`
	ResourceConf `yaml:"resource"`
}

type AppConf struct {
	Name                 string `yaml:"name"`
	Slogan               string `yaml:"slogan"`
	Url                  string `yaml:"url"`
	Port                 int32  `yaml:"port"`
	Icp                  string `yaml:"icp"`
	CodeRepoSecret       string `yaml:"code-repo-secret"`
	CodeRepoAddress      string `yaml:"code-repo-address"`
	CodeRepoType         string `yaml:"code-repo-type"`
	CodeRepoAuthType     string `yaml:"code-repo-auth-type"`
	CodeRepoAuthUser     string `yaml:"code-repo-auth-user"`
	CodeRepoAuthPassword string `yaml:"code-repo-auth-password"`
	PageLimit            int    `yaml:"page-limit"`
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
	configFile := fmt.Sprintf("%s/conf_%s.yml", GetConfigDir(), env)
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

func GetConfigDir() string {
	dir, _ := os.Getwd()
	return fmt.Sprintf("%s/config", dir)
}
