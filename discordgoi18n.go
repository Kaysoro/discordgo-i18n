package discordgoi18n

import (
	"github.com/bwmarrin/discordgo"
)

var translator Translator

func init() {
	translator = New()
}

func SetDefault(language discordgo.Locale) {
	translator.SetDefault(language)
}

func LoadBundle(language discordgo.Locale, file string) error {
	return translator.LoadBundle(language, file)
}

func Get(language discordgo.Locale, key string, values ...map[string]interface{}) string {
	args := make(map[string]interface{})

	for _, variables := range values {
		for variable, value := range variables {
			args[variable] = value
		}
	}

	return translator.Get(language, key, args)
}
