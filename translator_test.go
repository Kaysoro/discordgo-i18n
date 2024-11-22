package discordgoi18n

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const (
	translatornominalCase1         = "translatorNominalCase1.json"
	translatornominalCase2         = "translatorNominalCase2.json"
	translatorFailedUnmarshallCase = "translatorFailedUnmarshallCase.json"
	translatorFileDoesNotExistCase = "translatorFileDoesNotExistCase.json"

	content1 = `
	{
		"hi": ["this is a {{ .Test }}"],
		"with": ["all"],
		"the": ["elements", "we"],
		"can": ["find"],
		"in": ["a","json"],
		"config": ["file", "! {{ .Author }}"],
		"parse": ["{{if $foo}}{{end}}"]
	}
	`

	content2 = `
	{
		"this": ["is a {{ .Test }}"],
		"with.a.file": ["containing", "less", "variables"],
		"bye": ["see you"],
		"parse2": ["{{if $foo}}{{end}}"]
	}
	`

	badContent = `
	 [
		"content",
		"not",
		"ok",
		"test"
	 ]
	`
)

var (
	//nolint:gochecknoglobals // Acceptable for a test.
	translatorTest *translatorImpl
)

func setUp() {
	translatorTest = newTranslator()
	if err := os.WriteFile(translatornominalCase1, []byte(content1), os.ModePerm); err != nil {
		log.Fatal().Err(err).Msgf("'%s' could not be created, test stopped", translatornominalCase1)
	}
	if err := os.WriteFile(translatornominalCase2, []byte(content2), os.ModePerm); err != nil {
		log.Fatal().Err(err).Msgf("'%s' could not be created, test stopped", translatornominalCase2)
	}
	if err := os.WriteFile(translatorFailedUnmarshallCase, []byte(badContent), os.ModePerm); err != nil {
		log.Fatal().Err(err).Msgf("'%s' could not be created, test stopped", translatorFailedUnmarshallCase)
	}
}

func tearDown() {
	translatorTest = nil
	if err := os.Remove(translatornominalCase1); err != nil {
		log.Warn().Err(err).Msgf("'%s' could not be deleted", translatornominalCase1)
	}
	if err := os.Remove(translatornominalCase2); err != nil {
		log.Warn().Err(err).Msgf("'%s' could not be deleted", translatornominalCase2)
	}
	if err := os.Remove(translatorFailedUnmarshallCase); err != nil {
		log.Warn().Err(err).Msgf("'%s' could not be deleted", translatorFailedUnmarshallCase)
	}
}

func TestNewTranslator(t *testing.T) {
	setUp()
	defer tearDown()

	assert.Empty(t, translatorTest.translations)
	assert.Empty(t, translatorTest.loadedBundles)
}

func TestSetDefault(t *testing.T) {
	setUp()
	defer tearDown()

	assert.Equal(t, defaultLocale, translatorTest.defaultLocale)
	translatorTest.SetDefault(discordgo.Italian)
	assert.Equal(t, discordgo.Italian, translatorTest.defaultLocale)
}

func TestLoadBundle(t *testing.T) {
	setUp()
	defer tearDown()

	// Bad case, file does not exist
	_, err := os.Stat(translatorFileDoesNotExistCase)
	assert.Error(t, os.ErrNotExist, err)
	assert.Error(t, translatorTest.LoadBundle(discordgo.French, translatorFileDoesNotExistCase))
	assert.Empty(t, translatorTest.translations)
	assert.Empty(t, translatorTest.loadedBundles)

	// Bad case, file is not well structured
	assert.Error(t, translatorTest.LoadBundle(discordgo.French, translatorFailedUnmarshallCase))
	assert.Empty(t, translatorTest.translations)
	assert.Empty(t, translatorTest.loadedBundles)

	// Nominal case, load an existing and well structured bundle
	assert.NoError(t, translatorTest.LoadBundle(discordgo.French, translatornominalCase1))
	assert.Equal(t, 1, len(translatorTest.loadedBundles))
	assert.Equal(t, 1, len(translatorTest.translations))
	assert.Equal(t, 7, len(translatorTest.translations[discordgo.French]))

	// Nominal case, reload a bundle
	assert.NoError(t, translatorTest.LoadBundle(discordgo.French, translatornominalCase2))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 1, len(translatorTest.translations))
	assert.Equal(t, 4, len(translatorTest.translations[discordgo.French]))

	// Nominal case, load a bundle already loaded but for another locale
	assert.NoError(t, translatorTest.LoadBundle(discordgo.EnglishGB, translatornominalCase2))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 2, len(translatorTest.translations))
	assert.Equal(t, 4, len(translatorTest.translations[discordgo.EnglishGB]))

	// Nominal case, reload a bundle linked to two locales
	assert.NoError(t, translatorTest.LoadBundle(discordgo.EnglishGB, translatornominalCase1))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 2, len(translatorTest.translations))
	assert.Equal(t, 7, len(translatorTest.translations[discordgo.EnglishGB]))
}

func TestGet(t *testing.T) {
	setUp()
	defer tearDown()

	// Nominal case, get without bundle
	assert.Equal(t, "hi", translatorTest.Get(discordgo.Dutch, "hi", nil))

	// Nominal case, unexisting key with bundle loaded
	assert.NoError(t, translatorTest.LoadBundle(discordgo.Dutch, translatornominalCase1))
	assert.NoError(t, translatorTest.LoadBundle(defaultLocale, translatornominalCase2))
	assert.Equal(t, "does_not_exist", translatorTest.Get(discordgo.Dutch, "does_not_exist", nil))

	// Nominal case, Get existing key from loaded bundle
	assert.Equal(t, "this is a {{ .Test }}", translatorTest.Get(discordgo.Dutch, "hi", nil))
	assert.Equal(t, "this is a test :)", translatorTest.Get(discordgo.Dutch, "hi", Vars{"Test": "test :)"}))

	// Nominal case, Get key not present in loaded bundle but available in default
	assert.Equal(t, "see you", translatorTest.Get(discordgo.Dutch, "bye", nil))

	// Bad case, value is not well structured to be parsed
	assert.Equal(t, "{{if $foo}}{{end}}", translatorTest.Get(discordgo.Dutch, "parse", Vars{}))

	// Bad case, value is well structured but cannot inject value
	assert.Equal(t, "this is a {{ .Test }}", translatorTest.Get(discordgo.Dutch, "hi", Vars{}))
}

func TestGetDefault(t *testing.T) {
	setUp()
	defer tearDown()

	// Nominal case, get without bundle
	assert.Equal(t, "hi", translatorTest.GetDefault("hi", nil))

	// Nominal case, unexisting key with bundle loaded
	assert.NoError(t, translatorTest.LoadBundle(defaultLocale, translatornominalCase2))
	assert.Equal(t, "does_not_exist", translatorTest.GetDefault("does_not_exist", nil))

	// Nominal case, Get existing key from loaded bundle
	assert.Equal(t, "is a {{ .Test }}", translatorTest.GetDefault("this", nil))
	assert.Equal(t, "is a test :)", translatorTest.GetDefault("this", Vars{"Test": "test :)"}))

	// Bad case, value is not well structured to be parsed
	assert.Equal(t, "{{if $foo}}{{end}}", translatorTest.GetDefault("parse2", Vars{}))

	// Bad case, value is well structured but cannot inject value
	assert.Equal(t, "is a {{ .Test }}", translatorTest.GetDefault("this", Vars{}))
}

func TestMapBundleStructure(t *testing.T) {
	setUp()
	defer tearDown()

	tests := []struct {
		Description    string
		Input          map[string]interface{}
		ExpectedBundle bundle
	}{
		{
			Description:    "Nil Input",
			Input:          nil,
			ExpectedBundle: make(bundle),
		},
		{
			Description:    "Empty Input",
			Input:          make(map[string]interface{}),
			ExpectedBundle: make(bundle),
		},
		{
			Description: "Simple string Input",
			Input: map[string]interface{}{
				"simple":       "translation",
				"variabilized": "translation {{ .translation }}",
			},
			ExpectedBundle: bundle{
				"simple":       []string{"translation"},
				"variabilized": []string{"translation {{ .translation }}"},
			},
		},
		{
			Description: "Different types handled",
			Input: map[string]interface{}{
				"pi":                                  3.14,
				"answer_to_ultimate_question_of_life": 42,
				"some_prime_numbers":                  []interface{}{2, "3", 5.0, 7},
			},
			ExpectedBundle: bundle{
				"pi":                                  []string{"3.14"},
				"answer_to_ultimate_question_of_life": []string{"42"},
				"some_prime_numbers":                  []string{"2", "3", "5", "7"},
			},
		},
		{
			Description: "Deep structure",
			Input: map[string]interface{}{
				"command": map[string]interface{}{
					"salutation": map[string]interface{}{
						"hi":  "Hello there!",
						"bye": []interface{}{"Bye {{ .anyone }}!", "See u {{ .anyone }}"},
					},
					"speak": map[string]interface{}{
						"random": []interface{}{"love to talk", "how are u?", "u're so interesting"},
					},
				},
				"panic": "I've panicked!",
			},
			ExpectedBundle: bundle{
				"command.salutation.hi":  []string{"Hello there!"},
				"command.salutation.bye": []string{"Bye {{ .anyone }}!", "See u {{ .anyone }}"},
				"command.speak.random":   []string{"love to talk", "how are u?", "u're so interesting"},
				"panic":                  []string{"I've panicked!"},
			},
		},
	}

	for _, test := range tests {
		result := translatorTest.mapBundleStructure(test.Input)
		assert.True(t, reflect.DeepEqual(test.ExpectedBundle, result),
			fmt.Sprintf("%s:\n\nExpecting: %v\n\nGot      : %v", test.Description, test.ExpectedBundle, result))
	}
}

func TestGetLocalizations(t *testing.T) {
	setUp()
	defer tearDown()

	// Nominal case: empty map when no bundle loaded
	assert.Empty(t, translatorTest.GetLocalizations("hi", Vars{}))

	// Nominal case: two bundles loaded so two translations expected
	assert.NoError(t, translatorTest.LoadBundle(discordgo.Dutch, translatornominalCase1))
	assert.NoError(t, translatorTest.LoadBundle(defaultLocale, translatornominalCase2))
	assert.Equal(t, 2, len(*translatorTest.GetLocalizations("hi", Vars{})))

	// Nominal case: three bundles loaded so three translations expected
	assert.NoError(t, translatorTest.LoadBundle(discordgo.ChineseCN, translatornominalCase1))
	assert.Equal(t, 3, len(*translatorTest.GetLocalizations("hi", Vars{})))
}
