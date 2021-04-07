package testutil

import (
	"github.com/spf13/viper"
)

func SetupConfig() {
	viper.SetDefault("paseto.audience", "div-dash")
	viper.SetDefault("paseto.issuer", "div-dash")
	viper.SetDefault("paseto.key", "YELLOW SUBMARINE, BLACK WIZARDRY")
	viper.SetDefault("paseto.tokenValid", 24)

	viper.SetDefault("smtp.server", "localhost")
	viper.SetDefault("smtp.port", 1025)

	viper.ReadInConfig()
}
