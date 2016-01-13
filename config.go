package pbc

import (
	"github.com/spf13/viper"
	"os"
)

func getDataDir() string {
	return os.ExpandEnv(viper.GetString("core.datadir"))
}

func init() {
	viper.SetDefault("core.datadir", os.ExpandEnv("${HOME}/.pbc"))
}