package config

import (
	"context"
	"cro_test/pkg/logger"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/vx416/sqlxx"
	"go.uber.org/zap"
)

func Init(cfgName string) error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		_, f, _, _ := runtime.Caller(0)
		dir := filepath.Dir(f)
		configPath = filepath.Join(dir, "../../configs")
	}
	configName := os.Getenv("CONFIG_NAME")
	if configName != "" {
		cfgName = configName
	}
	if cfgName == "" {
		cfgName = "prod"
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigName(cfgName)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	Dump()
	return nil
}

func GetJwtSecret() []byte {
	return []byte(viper.GetString("jwt.secret"))
}

func Dump() {
	keys := viper.AllKeys()

	sort.Strings(keys)
	fmt.Println("==================CONFIG==================")
	for _, key := range keys {
		if strings.Contains(key, "password") || strings.Contains(key, "secret") || strings.Contains(key, "key") || strings.Contains(key, "pass") || strings.Contains(key, "pem") {
			continue
		}
		fmt.Printf("%s: %+v\n", key, viper.Get(key))
	}
	fmt.Println("==================CONFIG==================")
}

func InitSqlxx() (*sqlxx.Sqlxx, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.db"),
		viper.GetString("mysql.options"),
	)
	sqlxDB, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = sqlxDB.Ping()
	if err != nil {
		return nil, err
	}

	viper.SetDefault("mysql.pool.maxidletime", "5s")
	viper.SetDefault("mysql.pool.maxlifetime", "5m")
	viper.SetDefault("mysql.pool.maxidleconns", 50)
	viper.SetDefault("mysql.pool.maxopenconns", 50)
	sqlxDB.DB.SetConnMaxIdleTime(viper.GetDuration("mysql.pool.maxidletime"))
	sqlxDB.DB.SetConnMaxLifetime(viper.GetDuration("mysql.pool.maxidletime"))
	sqlxDB.DB.SetMaxIdleConns(viper.GetInt("mysql.pool.maxidleconns"))
	sqlxDB.DB.SetMaxOpenConns(viper.GetInt("mysql.pool.maxopenconns"))
	sqlxxDB := sqlxx.NewWith(sqlxDB)
	return sqlxxDB, nil
}

func InitRedis() (*redis.Client, error) {
	conn := redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.host") + ":" + viper.GetString("port"),
		DB:   viper.GetInt("redis.db"),
	})
	err := conn.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func InitLogger() (logger.Logger, error) {
	logLevel := logger.Level(viper.GetString("app.log.level"))
	var zapConfig zap.Config
	if viper.GetString("app.env") == "dev" {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = logger.ColorfulLevelEncoder
		zapConfig.EncoderConfig.EncodeCaller = logger.ColorizeCallerEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
	}
	zapConfig.Level = zap.NewAtomicLevelAt(logLevel.ZapLevel())
	if viper.GetString("app.log.outputs") != "" {
		outputs := strings.Split(viper.GetString("app.log.outputs"), ",")
		zapConfig.OutputPaths = append(zapConfig.OutputPaths, outputs...)
		zapConfig.ErrorOutputPaths = append(zapConfig.OutputPaths, outputs...)
	}

	field := zap.String("app_name", viper.GetString("app.name"))
	zaplog, err := zapConfig.Build(zap.Fields(field), zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	log := logger.NewZapAdapter(zaplog)
	logger.SetGlobal(log)
	return log, nil
}
