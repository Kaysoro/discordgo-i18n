package discordgoi18n

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {

	// Must not panic
	mock := newMock()
	mock.SetDefault(discordgo.ChineseCN)
	assert.NoError(t, mock.LoadBundle(discordgo.SpanishES, ""))
	assert.Empty(t, mock.Get(discordgo.Croatian, "", nil))
	assert.Nil(t, mock.GetLocalizations("", nil))

	var called bool
	mock.SetDefaultFunc = func(locale discordgo.Locale) {
		called = true
	}
	mock.LoadBundleFunc = func(locale discordgo.Locale, file string) error {
		called = true
		return nil
	}
	mock.GetFunc = func(locale discordgo.Locale, key string, values Vars) string {
		called = true
		return ""
	}
	mock.GetLocalizationsFunc = func(key string, values Vars) *map[discordgo.Locale]string {
		called = true
		return nil
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

	called = false
	assert.Empty(t, mock.GetLocalizations("", nil))
	assert.True(t, called)
}
