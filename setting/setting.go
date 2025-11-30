/*
本文件用于从 conf 包下的conf文件中读取配置项并保存至结构体中，由viper管理
*/
package setting

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

// Conf 操作实例
var Conf = new(AcatConf)

// AcatConf 纳新web网站的配置结构体
type AcatConf struct {
	WebName      string `mapstructure:"name"`
	WebMode      string `mapstructure:"mode"`
	WebPort      string `mapstructure:"port"`
	WebVersion   string `mapstructure:"version"`
	WebStartTime string `mapstructure:"start_time"`
	JwtSecret    string `mapstructure:"jwt_secret"`
	*LogConf     `mapstructure:"log"`
	*MySQLConf   `mapstructure:"mysql"`
	*RedisConf   `mapstructure:"redis"`
	*EmailConf   `mapstructure:"email"`
}

// LogConf 日志文件配置项
type LogConf struct {
	LogLevel      string `mapstructure:"level"`
	LogFileName   string `mapstructure:"filename"`
	LogMaxSize    int    `mapstructure:"max_size"`
	LogMaxAge     int    `mapstructure:"max_age"`
	LogMaxBackups int    `mapstructure:"max_backups"`
}

// MySQLConf MySQL配置项
type MySQLConf struct {
	MySQLHost        string `mapstructure:"host"`
	MySQLPort        string `mapstructure:"port"`
	MySQLUser        string `mapstructure:"user"`
	MySQLPwd         string `mapstructure:"password"`
	MySQLDbName      string `mapstructure:"dbname"`
	MySQLMaxOpenConn int    `mapstructure:"maxOpenConn"` // 匹配 yaml 的 maxOpenConn
	MySQLMaxIdleConn int    `mapstructure:"maxIdleConn"`
}

// RedisConf redis配置项
type RedisConf struct {
	RedisHost     string `mapstructure:"host"`
	RedisPort     string `mapstructure:"port"`
	RedisPwd      string `mapstructure:"password"`
	RedisDb       int    `mapstructure:"db"`
	RedisPoolSize int    `mapstructure:"pool_size"`
}

// EmailConf 邮箱配置项
type EmailConf struct {
	EmailValidEmail string `mapstructure:"ValidEmail"`
	EmailSMTPHost   string `mapstructure:"SmtpHost"`
	EmailSMTPEmail  string `mapstructure:"SmtpEmail"`
	EmailSMTPPass   string `mapstructure:"SmtpPass"`
}

func Init() error {
	viper.SetConfigFile("./conf/config_dev.yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("[PANIC] config_dev.yaml read failed && WHERE setting/setting.go--Init(),error : ", err)
	}
	if err := viper.Unmarshal(Conf); err != nil {
		log.Fatal("[PANIC] config_dev.yaml unmarshal failed && WHERE setting/setting.go--Init(),error : ", err)
	}

	// 从环境变量加载 SMTP 密钥（最小安全加固）
	Conf.EmailSMTPPass = os.Getenv("EMAIL_SMTP_PASS")

	// 开启配置热更新，WatchConfig来监听配置文件
	viper.WatchConfig()
	// 使用fsnotify库监控你设置的配置文件，每当配置文件改变后，viper会重新加载配置文件
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file changer!")
		if err := viper.Unmarshal(Conf); err != nil {
			log.Fatal("[PANIC] config_dev.yaml unmarshal failed && WHERE setting/setting.go--Init(),error : ", err)
		}
		// 热重载后重新注入 SMTP 密钥（保持安全）
		Conf.EmailSMTPPass = os.Getenv("EMAIL_SMTP_PASS")
	})
	return err
}
