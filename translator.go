package i18n

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"os"
	"strings"
	"text/template"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

const (
	defaultLocale   = discordgo.EnglishUS
	leftDelim       = "{{"
	rightDelim      = "}}"
	executionPolicy = "missingkey=error"
)

func New() *TranslatorImpl {
	return &TranslatorImpl{
		defaultLocale: defaultLocale,
		translations:  make(map[discordgo.Locale]bundle),
		loadedBundles: make(map[string]bundle),
	}
}

func (translator *TranslatorImpl) SetDefault(language discordgo.Locale) {
	translator.defaultLocale = language
}

func (translator *TranslatorImpl) LoadBundle(locale discordgo.Locale, path string) error {
	loadedBundle, found := translator.loadedBundles[path]
	if !found {

		buf, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var newBundle bundle
		err = json.Unmarshal(buf, &newBundle)
		if err != nil {
			return err
		}

		log.Debug().Msgf("Bundle '%s' loaded with '%s' content", locale, path)
		translator.loadedBundles[path] = newBundle
		translator.translations[locale] = newBundle

	} else {
		log.Debug().Msgf("Bundle '%s' already loaded, content now linked to locale %s too", path, locale)
		translator.translations[locale] = loadedBundle
	}

	return nil
}

func (translator *TranslatorImpl) Get(locale discordgo.Locale, key string, variables map[string]interface{}) string {
	bundles, found := translator.translations[locale]
	if !found {
		if locale != translator.defaultLocale {
			log.Warn().Msgf("Bundle '%s' is not loaded, trying to translate key '%s' in '%s'", locale, key, translator.defaultLocale)
			return translator.Get(translator.defaultLocale, key, variables)
		} else {
			log.Warn().Msgf("Bundle '%s' is not loaded, cannot translate '%s', key returned", locale, key)
			return key
		}
	}

	raws, found := bundles[key]
	if !found || len(raws) == 0 {
		if locale != translator.defaultLocale {
			log.Warn().Msgf("No label found for key '%s' in '%s', trying to translate it in %s", key, locale, translator.defaultLocale)
			return translator.Get(translator.defaultLocale, key, variables)
		} else {
			log.Warn().Msgf("No label found for key '%s' in '%s', key returned", locale, key)
			return key
		}
	}

	raw := raws[rand.Intn(len(raws))]

	if variables != nil && strings.Contains(raw, leftDelim) {
		t, err := template.New("").Delims(leftDelim, rightDelim).Option(executionPolicy).Parse(raw)
		if err != nil {
			log.Error().Err(err).Msgf("Cannot parse raw corresponding to key '%s' in '%s', raw returned", locale, key)
			return raw
		}

		var buf bytes.Buffer
		err = t.Execute(&buf, variables)
		if err != nil {
			log.Error().Err(err).Msgf("Cannot inject variables in raw corresponding to key '%s' in '%s', raw returned", locale, key)
			return raw
		}
		return buf.String()
	}

	return raw
}
