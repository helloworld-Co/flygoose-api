package configs

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

// 全局配置文件实例
var Cfg *Config

type Config struct {
	Http         HttpCfg     `yaml:"http"`
	Database     DatabaseCfg `yaml:"database"`
	ExecuteDir   string
	StaticDir    string
	StaticImgDir string
}

type HttpCfg struct {
	Port int `yaml:"port"`
}

type DatabaseCfg struct {
	Driver   DbDriverType `yaml:"driver"`
	Host     string       `yaml:"host"`
	Name     string       `yaml:"name"`
	Port     int          `yaml:"port"`
	User     string       `yaml:"user"`
	Password string       `yaml:"password"`
}

// 数据库类型
type DbDriverType string

const (
	DbDriverMySQL      DbDriverType = "mysql"
	DbDriverPostgreSQL DbDriverType = "postgresql"
)

// 环境
type EnvModeType string

//const (
//	EnvModeTypeDevelop EnvModeType = "develop" //开发环境
//	EnvModeTypeProd    EnvModeType = "product" //正式环境
//)

func NewConfig(configFilePath string) (*Config, error) {
	// 读取配置文件
	bytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件错误. err: %w", err)
	}

	//解析配置文件到结构体
	var cfg Config
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件错误. err: %w", err)
	}

	//if err = CheckZeroValue(cfg); err != nil {
	//	return nil, fmt.Errorf("检查配置文件错误. env: %s err: %s", configFilePath, err.Error())
	//}

	return &cfg, nil
}

// CheckZeroValue 检验配置文件各个字段不能为空. str必须为一个结构体类型.
func CheckZeroValue(str interface{}) error {
	t := reflect.TypeOf(str)
	if t.Kind() != reflect.Struct {
		return errors.New("non-struct type:" + t.String())
	}
	v := reflect.ValueOf(str)
	for k := 0; k < t.NumField(); k++ {
		fieldType := v.Field(k).Kind()
		if fieldType == reflect.Struct {
			if err := CheckZeroValue(v.Field(k).Interface()); err != nil {
				return err
			}
		}
		if v.Field(k).IsZero() {
			required := t.Field(k).Tag.Get("required")
			if required == "false" {
				fmt.Errorf("%+v %+v is zero, 但根据校验规则，required为false.跳过检查.\n", t, t.Field(k).Name)
				continue
			}
			return fmt.Errorf("%+v %+v can not be zero", t, t.Field(k).Name)
		}
	}
	return nil
}

const Flygoose_Url_Prefix = "/api"

const Flygoose_Admin_Phone = "12345678901"
