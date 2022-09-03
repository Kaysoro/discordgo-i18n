package i18n

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func NewMock() *TranslatorMock {
	return &TranslatorMock{}
}

func (mock *TranslatorMock) SetDefault(locale discordgo.Locale) {
	if mock.SetDefaultFunc != nil {
		mock.SetDefaultFunc(locale)
		return
	}

	log.Warn().Msgf("SetDefault not mocked")
}

func (mock *TranslatorMock) LoadBundle(locale discordgo.Locale, file string) error {
	if mock.LoadBundleFunc != nil {
		return mock.LoadBundleFunc(locale, file)
	}

	log.Warn().Msgf("LoadBundle not mocked")
	return nil
}

func (mock *TranslatorMock) Get(locale discordgo.Locale, key string, values map[string]interface{}) string {
	if mock.GetFunc != nil {
		return mock.GetFunc(locale, key, values)
	}

	log.Warn().Msgf("Get not mocked")
	return ""
}
