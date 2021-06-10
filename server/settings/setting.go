package settings

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

type SomeSettings struct {
	Location string
}

type Settings struct {
	Some SomeSettings
}

func Load() (*Settings, error) {
	settingsData, err := ioutil.ReadFile("settings.toml")
	if err != nil {
		return nil, err
	}

	var settings Settings
	err = toml.Unmarshal(settingsData, &settings)
	return &settings, err
}

func (s Settings) Valid() error {
	return nil
}
