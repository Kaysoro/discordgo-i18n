package i18n

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {

	// Must not panic
	mock := NewMock()
	mock.SetDefault(discordgo.ChineseCN)
	assert.NoError(t, mock.LoadBundle(discordgo.SpanishES, ""))
	assert.Empty(t, mock.Get(discordgo.Croatian, "", nil))

	var called bool
	mock.SetDefaultFunc = func(locale discordgo.Locale) {
		called = true
	}
	mock.LoadBundleFunc = func(locale discordgo.Locale, file string) error {
		called = true
		return nil
	}
	mock.GetFunc = func(locale discordgo.Locale, key string, values map[string]interface{}) string {
		called = true
		return ""
	}

	called = false
	mock.SetDefault(discordgo.ChineseCN)
	assert.True(t, called)

	called = false
	assert.NoError(t, mock.LoadBundle(discordgo.SpanishES, ""))
	assert.True(t, called)

	called = false
	assert.Empty(t, mock.Get(discordgo.Croatian, "", nil))
	assert.True(t, called)
}
