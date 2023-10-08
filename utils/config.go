package utils

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationUrl         string        `mapstructure:"MIGRATIONS_URL"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRATION"`
	CorsAllowedOrigin    string        `mapstructure:"CORS_ALLOWED_ORIGIN"`
	Environment          string        `mapstructure:"ENVIRONMENT"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRATION"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	RedisPassword        string        `mapstructure:"REDIS_PASSWORD"`
	ApiUrl               string        `mapstructure:"API_URL"`
	ApiUsername          string        `mapstructure:"AFRICASTALKING_USERNAME"`
	SenderApiKey         string        `mapstructure:"SMS_SENDER_API_KEY"`
	SenderID             string        `mapstructure:"SMS_SENDER_ID"`
	NewRelicAppName      string        `mapstructure:"NEWRELIC_APP_NAME"`
	NewRelicLicenseKey   string        `mapstructure:"NEW_RELIC_LICENSE_KEY"`
	AWSAccessKeyID       string        `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey   string        `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AWSRegion            string        `mapstructure:"AWS_REGION"`
	AWSBucketName        string        `mapstructure:"AWS_BUCKET_NAME"`
	SendLogsToStdOut     bool
}

func LoadConfig(path string) (config Config, err error) {
	config.SendLogsToStdOut = true

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil

}
