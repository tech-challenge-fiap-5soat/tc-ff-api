package config

import (
	"fmt"
	"github.com/integralist/go-findroot/find"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	runOnce sync.Once
	config  *Config
)

type Config struct {
	MongoCfg MongoConfig `yaml:"mongodb"`
	ApiCfg   Api         `yaml:"api"`
}

type MongoConfig struct {
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
}

type Api struct {
	Port                 string     `yaml:"port"`
	AuthConfig           AuthConfig `yaml:"authConfig"`
	AuthorizationBaseUrl string     `yaml:"authorizationUrl"`
}

type AuthConfig struct {
	UserPoolId string `yaml:"userPoolId"`
	ClientId   string `yaml:"clientId"`
	TokenUse   string `yaml:"tokenUse"`
}

func init() {
	config = setupConfig()
}

func GetMongoCfg() MongoConfig {
	return config.MongoCfg
}

func GetApiCfg() Api {
	return config.ApiCfg
}

func setupConfig() *Config {
	runOnce.Do(func() {
		var appConfig Config

		root, _ := find.Repo()
		configFilePath := path.Join(root.Path, "/src/external/api/infra/config")

		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
		viper.SetConfigName("configs")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(configFilePath)
		viper.AddConfigPath("/app/data/configs")
		err := viper.ReadInConfig()

		//if err = viper.ReadInConfig(); err != nil {
		//	panic(fmt.Errorf("Falha ao carregar as configurações: %w \n", err))
		//}

		if err != nil && !allConfigsAreSet() {
			panic(fmt.Errorf("Falha ao carregar as configurações: %w \n", err))
		}

		for _, key := range viper.AllKeys() {
			value := viper.GetString(key)
			envOrRaw := replaceEnvInConfig([]byte(value))
			viper.Set(key, string(envOrRaw))
		}

		if err = viper.Unmarshal(&config); err != nil {
			panic(err)
		}

		if allConfigsAreSet() { // load envs from infra
			appConfig.ApiCfg.Port = viper.Get("api.port").(string)
			appConfig.MongoCfg.Host = viper.Get("mongodb.host").(string)
			appConfig.MongoCfg.Database = viper.Get("mongodb.database").(string)
			appConfig.MongoCfg.User = viper.Get("mongodb.user").(string)
			appConfig.MongoCfg.Pass = viper.Get("mongodb.pass").(string)
			appConfig.ApiCfg.AuthConfig.UserPoolId = viper.Get("api.authConfig.userPoolId").(string)
			appConfig.ApiCfg.AuthConfig.ClientId = viper.Get("api.authConfig.clientId").(string)
			appConfig.ApiCfg.AuthConfig.TokenUse = viper.Get("api.authConfig.tokenUse").(string)

		}

		config = &appConfig
	})

	return config
}

func allConfigsAreSet() bool {
	return viper.Get("mongodb.host") != nil &&
		viper.Get("mongodb.database") != nil &&
		viper.Get("mongodb.user") != nil &&
		viper.Get("api.port") != nil &&
		viper.Get("mongodb.pass") != nil &&
		viper.Get("api.authConfig.userPoolId") != nil &&
		viper.Get("api.authConfig.clientId") != nil &&
		viper.Get("api.authConfig.tokenUse") != nil
}

func replaceEnvInConfig(body []byte) []byte {
	search := regexp.MustCompile(`\$\{([^{}]+)\}`)
	replacedBody := search.ReplaceAllFunc(body, func(b []byte) []byte {
		group1 := search.ReplaceAllString(string(b), `$1`)

		envValue := os.Getenv(group1)
		if len(envValue) > 0 {
			return []byte(envValue)
		}
		return []byte("")
	})

	log.Println(string(replacedBody))
	return replacedBody
}
