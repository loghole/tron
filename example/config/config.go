// Code generated by tron. Place config helpers here.
package config

import (
	"github.com/spf13/viper"
)

func GetExampleValue() string {
	return viper.GetString(ExampleValue)
}
