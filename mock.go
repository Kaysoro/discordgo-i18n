package discordgoi18n

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func newMock() *translatorMock {
	return &translatorMock{}
}

func (mock *translatorMock) SetDefault(locale discordgo.Locale) {
	if mock.SetDefaultFunc != nil {
		mock.SetDefaultFunc(locale)
		return
	}

	log.Warn().Msgf("SetDefault not mocked")
}

func (mock *translatorMock) LoadBundle(locale discordgo.Locale, file string) error {
	if mock.LoadBundleFunc != nil {
		return mock.LoadBundleFunc(locale, file)
	}

	log.Warn().Msgf("LoadBundle not mocked")
	return nil
}

func (mock *translatorMock) Get(locale discordgo.Locale, key string, variables Vars) string {
	if mock.GetFunc != nil {
		return mock.GetFunc(locale, key, variables)
	}

	log.Warn().Msgf("Get not mocked")
	return ""
}

func (mock *translatorMock) GetLocalizations(key string, variables Vars) *map[discordgo.Locale]string {
	if mock.GetFunc != nil {
		return mock.GetLocalizationsFunc(key, variables)
	}

	log.Warn().Msgf("GetLocalizations not mocked")
	return nil
}
