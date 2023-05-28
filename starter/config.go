package starter

import (
	"fmt"
)

func init() {
	Register(&Config{})
}

type Config struct{}

func (c *Config) Init(ctx StarterContext) {
	fmt.Println("Config Init")
}

func (c *Config) Setup(ctx StarterContext) {
	fmt.Println("Config Setup")
}

func (c *Config) Start(ctx StarterContext) {
	fmt.Println("Config Start")
}

func (c *Config) StartBlocking() bool {
	fmt.Println("Config StartBlocking")
	return true
}

func (c *Config) Stop(ctx StarterContext) {
	fmt.Println("Config Stop")
}
