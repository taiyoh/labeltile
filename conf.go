package labeltile

import (
	"errors"
	"net/url"
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

// OAuth2GoogleConf provides configuration for google oauth2
type OAuth2GoogleConf struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	RedirectURL  string `toml:"redirect_url"`
}

// IsValid returns whether OAuth2GoogleConf data structure is valid or not
func (c OAuth2GoogleConf) IsValid() bool {
	if c.ClientID == "" || c.ClientSecret == "" {
		return false
	}
	if u, err := url.Parse(c.RedirectURL); err != nil {
		return false
	} else if u.Scheme == "http" || u.Scheme == "https" {
		return true
	}
	return false
}

// OAuth2Conf provides oauth2 conf section
type OAuth2Conf struct {
	Google OAuth2GoogleConf `toml:"google"`
}

type Database struct {
	Driver string `toml:"driver"`
	Dsn    string `toml:"dsn"`
}

// Conf is configuration binder
type Conf struct {
	Server   ServerConf `toml:"server"`
	JWT      JWTConf    `toml:"jwt"`
	OAuth2   OAuth2Conf `toml:"oauth2"`
	Database Database   `toml:"database"`
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
	if !labeltileConf.OAuth2.Google.IsValid() {
		return nil, errors.New("invalid OAuth2.Google section")
	}

	return labeltileConf, nil
}
