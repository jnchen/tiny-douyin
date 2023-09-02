package config

import (
	"fmt"
	"github.com/spf13/viper" // 配置管理
	"log"
	"os"
	"reflect"
	"strings"
)

var Conf = new(TotalConfig)

type TotalConfig struct {
	*MySQL   `mapstructure:"mysql"`
	*Storage `mapstructure:"storage"`
}

type MySQL struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
	Options  string `mapstructure:"options"`
	Port     int    `mapstructure:"port"`
}

type Storage struct {
	Local struct {
		Path    string `mapstructure:"path"`
		BaseURL string `mapstructure:"base_url"`
	} `mapstructure:"local"`
	OSS struct {
		Enable          bool   `mapstructure:"enable"`
		Endpoint        string `mapstructure:"endpoint"`
		BucketName      string `mapstructure:"bucket_name"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		AccessKeySecret string `mapstructure:"access_key_secret"`
	} `mapstructure:"oss"`
}

func init() {
	log.Println("解析配置")

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("配置读取失败：%v\n", err)
	}

	// 模板解析
	for i, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if !strings.HasPrefix(value, "${") || !strings.HasSuffix(value, "}") {
			continue
		}

		envKey, defaultValue, defaultExists := parseEnvTemplate(value)
		envVal, envSet := os.LookupEnv(envKey)

		log.Printf("%d.\t%s 绑定环境变量 %s", i, k, envKey)
		log.Println("\t读取环境变量\t", envKey)
		if !envSet || len(envVal) == 0 {
			if defaultExists {
				log.Println("\t默认值\t\t", defaultValue)
				viper.Set(k, defaultValue)
			} else {
				log.Panicf("环境变量 %s 不存在或为空\n", envKey)
			}
		}

		viper.MustBindEnv(k, envKey)
	}

	if err := viper.Unmarshal(Conf); err != nil {
		log.Panicf("配置解析失败：%v\n", err)
	}

	fmt.Println(formatConfig("config", viper.AllSettings(), 0))
}

func parseEnvTemplate(envTemplate string) (
	envKey string,
	defaultValue string,
	defaultExists bool,
) {
	// 去除 "${" 和 "}"
	envTemplate = strings.TrimSuffix(
		strings.TrimPrefix(envTemplate, "${"),
		"}",
	)
	// 获取环境变量名和默认值
	envKey, defaultValue, defaultExists = strings.Cut(envTemplate, ":")
	// 去除默认值首尾空格和双引号
	defaultValue = strings.Trim(
		strings.TrimSpace(defaultValue),
		"\"",
	)
	return
}

func formatConfig(key string, value interface{}, indentLevel int) string {
	indent := strings.Repeat("  ", indentLevel)
	valType := reflect.TypeOf(value).Kind()

	switch valType {
	case reflect.String:
		return fmt.Sprintf("%s%s: %s", indent, key, value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%s%s: %d", indent, key, value)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%s%s: %f", indent, key, value)
	case reflect.Bool:
		return fmt.Sprintf("%s%s: %v", indent, key, value)
	case reflect.Slice, reflect.Array:
		items := reflect.ValueOf(value)
		var formattedItems []string
		for i := 0; i < items.Len(); i++ {
			formattedItems = append(formattedItems, formatConfig("", items.Index(i).Interface(), indentLevel+1))
		}
		return fmt.Sprintf("%s%s:\n%s", indent, key, strings.Join(formattedItems, "\n"))
	case reflect.Map:
		return formatMapConfig(key, value.(map[string]interface{}), indentLevel)
	default:
		return fmt.Sprintf("%s%s: %v", indent, key, value)
	}
}

func formatMapConfig(key string, value map[string]interface{}, indentLevel int) string {
	indent := strings.Repeat("  ", indentLevel)
	var formattedItems []string
	for k, v := range value {
		formattedItems = append(formattedItems, formatConfig(k, v, indentLevel+1))
	}
	return fmt.Sprintf("%s%s:\n%s", indent, key, strings.Join(formattedItems, "\n"))
}
