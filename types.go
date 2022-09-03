package discordgoi18n

import (
	"github.com/bwmarrin/discordgo"
)

type Translator interface {
	SetDefault(locale discordgo.Locale)
	LoadBundle(locale discordgo.Locale, file string) error
	Get(locale discordgo.Locale, key string, values map[string]interface{}) string
}

type TranslatorImpl struct {
	defaultLocale discordgo.Locale
	translations  map[discordgo.Locale]bundle
	loadedBundles map[string]bundle
}

type TranslatorMock struct {
	SetDefaultFunc func(locale discordgo.Locale)
	LoadBundleFunc func(locale discordgo.Locale, file string) error
	GetFunc        func(locale discordgo.Locale, key string, values map[string]interface{}) string
}

type bundle map[string][]string
