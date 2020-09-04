package config

import "github.com/spf13/viper"

func Load() {
	viper.SetEnvPrefix("sturdy-octo-meme")

	// logging
	viper.SetDefault("logging.level", "ERROR")
	viper.BindEnv("logging.level", "LOG_LEVEL")

	viper.SetDefault("health.lag.history", 1000)
	viper.SetDefault("health.lag.max", 10000)
	viper.BindEnv("health.lag.history", "HEALTH_LAG_HISTORY")
	viper.BindEnv("health.lag.max", "HEALTH_LAG_MAX")
}
