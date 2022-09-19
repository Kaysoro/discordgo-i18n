package discordgoi18n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

const (
	defaultLocale   = discordgo.EnglishUS
	leftDelim       = "{{"
	rightDelim      = "}}"
	keyDelim        = "."
	executionPolicy = "missingkey=error"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func new() *translatorImpl {
	return &translatorImpl{
		defaultLocale: defaultLocale,
		translations:  make(map[discordgo.Locale]bundle),
		loadedBundles: make(map[string]bundle),
	}
}

func (translator *translatorImpl) SetDefault(language discordgo.Locale) {
	translator.defaultLocale = language
}

func (translator *translatorImpl) LoadBundle(locale discordgo.Locale, path string) error {
	loadedBundle, found := translator.loadedBundles[path]
	if !found {

		buf, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var jsonContent map[string]interface{}
		err = json.Unmarshal(buf, &jsonContent)
		if err != nil {
			return err
		}

		newBundle := translator.mapBundleStructure(jsonContent)

		log.Debug().Msgf("Bundle '%s' loaded with '%s' content", locale, path)
		translator.loadedBundles[path] = newBundle
		translator.translations[locale] = newBundle

	} else {
		log.Debug().Msgf("Bundle '%s' loaded with '%s' content (already loaded for other locales)", locale, path)
		translator.translations[locale] = loadedBundle
	}

	return nil
}

func (translator *translatorImpl) Get(locale discordgo.Locale, key string, variables Vars) string {
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

func (translator *translatorImpl) GetLocalizations(key string, variables Vars) *map[discordgo.Locale]string {
	localizations := make(map[discordgo.Locale]string)

	for locale := range translator.translations {
		localizations[locale] = translator.Get(locale, key, variables)
	}

	return &localizations
}

func (translator *translatorImpl) mapBundleStructure(jsonContent map[string]interface{}) bundle {
	bundle := make(map[string][]string)
	for key, content := range jsonContent {
		switch v := content.(type) {
		case string:
			bundle[key] = []string{v}
		case []interface{}:
			values := make([]string, 0)
			for _, value := range v {
				values = append(values, fmt.Sprintf("%v", value))
			}
			bundle[key] = values
		case map[string]interface{}:
			subValues := translator.mapBundleStructure(v)
			for subKey, subValue := range subValues {
				bundle[fmt.Sprintf("%s%s%s", key, keyDelim, subKey)] = subValue
			}
		default:
			bundle[key] = []string{fmt.Sprintf("%v", v)}
		}
	}

	return bundle
}
