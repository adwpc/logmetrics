package conf

import (
	"flag"
	"fmt"
	"os"

	"github.com/adwpc/logmetrics/zlog"
	"github.com/spf13/viper"
)

var (
	cfg = &Config{}
	log = zlog.Log
)

type Log struct {
	Path string `mapstructure:"path"`
	End  bool   `mapstructure:"end"`
}

type Config struct {
	Logs    map[string]Log `mapstructure:"log"`
	Listen  string         `mapstructure:"listen"`
	CfgFile string
	err     error
}

func ShowHelp() {
	fmt.Sprintf("Usage:%s {params}\n", os.Args[0])
	fmt.Println("      -c {config file}")
	fmt.Println("      -h (show help info)")
}

// func init() {
// conf.Parse()
// }

func GetConfig() *Config {
	return cfg
}

func (c *Config) Load() bool {

	_, c.err = os.Stat(c.CfgFile)
	if c.err != nil {
		fmt.Println(c.CfgFile, " didn't exist!")
		return false
	}

	viper.SetConfigFile(c.CfgFile)
	viper.SetConfigType("toml")

	c.err = viper.ReadInConfig()
	if c.err != nil {
		log.Error().Msgf("config file %s read failed. %v", c.CfgFile, c.err)
		return false
	}
	c.err = viper.GetViper().UnmarshalExact(c)
	if c.err != nil {
		log.Error().Msgf("config file %s loaded failed. %v", c.CfgFile, c.err)
		return false
	}
	if c.Listen == "" {
		log.Error().Msg("config file %s loaded failed. listen=\"\"")
		return false
	}

	log.Info().Msgf("config '%s' load ok!", c.CfgFile)
	return true
}

func (c *Config) Parse() bool {
	flag.StringVar(&c.CfgFile, "c", "conf/conf.toml", "config file")
	help := flag.Bool("h", false, "help info")
	flag.Parse()
	log.Info().Msg(c.CfgFile)

	if !c.Load() || *help {
		ShowHelp()
		return false
	}
	return true
}
