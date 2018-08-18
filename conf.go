package labeltile

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

// JWTConf provides configuration for JSON Web Token
type JWTConf struct {
	SecretKey     string `toml:"secret_key"`
	ExpireHour    uint32 `toml:"expire_hour"`
	SigningMethod string `toml:"signing_method"`
}

// IsValid returns whether JWTConf data structure is valid or not
func (c JWTConf) IsValid() bool {
	if c.SecretKey == "" || c.SigningMethod == "" {
		return false
	}
	return c.ExpireHour != 0
}

// ServerConf provides configuration for Web Server
type ServerConf struct {
	Port     uint16 `toml:"port"`
	Template string `toml:"template"`
}

// IsValid returns whether ServerConf data structure is valid or not
func (c ServerConf) IsValid() bool {
	if c.Port == 0 || c.Template == "" {
		return false
	}
	if f, err := os.Stat(c.Template); os.IsNotExist(err) || !f.IsDir() {
		return false
	}
	return true
}

// Conf is configuration binder
type Conf struct {
	Server ServerConf
	JWT    JWTConf
}

// NewConf returns Conf object with validation
func NewConf(path string) (*Conf, error) {
	labeltileConf := &Conf{}
	if _, err := toml.DecodeFile(path, labeltileConf); err != nil {
		return nil, err
	}

	if !labeltileConf.JWT.IsValid() {
		return nil, errors.New("invalid JWT section")
	}
	if !labeltileConf.Server.IsValid() {
		return nil, errors.New("invalid Server section")
	}

	return labeltileConf, nil
}
