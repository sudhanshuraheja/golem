package config

import (
	"log"

	"github.com/joeshaw/envdecode"
)

type Conf struct {
	LogLevel  int  `env:"GOLEM_DEBUG,default=1"` // trace=-1, debug=0, info=1, warn=2, error=3, fatal=4, panic=5
	LogPretty bool `env:"GOLEM_DEBUG_PRETTY,default=true"`
}

// NewConfig ...
func NewConfig() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}
	return &c
}
