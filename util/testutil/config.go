package testutil

import (
	"github.com/spf13/viper"
)

func SetupConfig() {
	viper.SetDefault("token.key", "YELLOW SUBMARINE, BLACK WIZARDRY")
	viper.SetDefault("token.tokenValid", 24)

	viper.SetDefault("smtp.server", "localhost")
	viper.SetDefault("smtp.port", 1025)

	viper.ReadInConfig()
}
