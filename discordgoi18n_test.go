package discordgoi18n

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestFacade(t *testing.T) {
	var expectedFile, expectedKey = "File", "Key"
	var expectedValues map[string]interface{}
	var called bool

	mock := NewMock()
	mock.SetDefaultFunc = func(locale discordgo.Locale) {
		assert.Equal(t, discordgo.Italian, locale)
		called = true
	}
	mock.LoadBundleFunc = func(locale discordgo.Locale, file string) error {
		assert.Equal(t, discordgo.French, locale)
		assert.Equal(t, expectedFile, file)
		called = true
		return nil
	}
	mock.GetFunc = func(locale discordgo.Locale, key string, values map[string]interface{}) string {
		assert.Equal(t, discordgo.ChineseCN, locale)
		assert.Equal(t, expectedValues, values)
		assert.Equal(t, expectedKey, key)
		called = true
		return ""
	}

	translator = mock

	called = false
	SetDefault(discordgo.Italian)
	assert.True(t, called)

	called = false
	assert.NoError(t, LoadBundle(discordgo.French, expectedFile))
	assert.True(t, called)

	called = false
	expectedValues = make(map[string]interface{})
	Get(discordgo.ChineseCN, expectedKey)
	assert.True(t, called)

	called = false
	expectedValues = map[string]interface{}{
		"Hi": "There",
	}
	Get(discordgo.ChineseCN, expectedKey, expectedValues)
	assert.True(t, called)

	called = false
	expectedValues = map[string]interface{}{
		"Hi":  "There",
		"Bye": "See u",
	}
	Get(discordgo.ChineseCN, expectedKey, map[string]interface{}{"Hi": "There"}, map[string]interface{}{"Bye": "See u"})
	assert.True(t, called)
}
