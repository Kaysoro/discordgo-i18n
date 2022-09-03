# discordgo-i18n ![Build status](https://github.com/kaysoro/discordgo-i18n/workflows/Build/badge.svg) [![Report card](https://goreportcard.com/badge/github.com/kaysoro/discordgo-i18n)](https://goreportcard.com/report/github.com/kaysoro/discordgo-i18n) [![codecov](https://codecov.io/gh/kaysoro/discordgo-i18n/branch/master/graph/badge.svg)](https://codecov.io/gh/kaysoro/discordgo-i18n) [![Sourcegraph](https://sourcegraph.com/github.com/kaysoro/discordgo-i18n/-/badge.svg)](https://sourcegraph.com/github.com/kaysoro/discordgo-i18n?badge)

discordgo-i18n is a simple and lightweight Go package that helps you translate Go programs into [languages supported by Discord](https://discord.com/developers/docs/reference#locales).

- Built to ease usage of [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo)
- Less verbose than [go-i18n](https://raw.githubusercontent.com/nicksnyder/go-i18n)
- Supports multiple strings per key to make your bot "more alive"
- Supports strings with named variables using [text/template](http://golang.org/pkg/text/template/) syntax
- Supports message files of JSON format

# Getting started

## Installing

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

`go get` *will always pull the latest tagged release from the master branch.*

```sh
go get github.com/kaysoro/discordgo-i18n
```

**NOTICE**: this package has been built to ease usage of [bwmarrin/discordgo](https://github.com/bwmarrin/discordgo), it can be used for other projects but will be less practical.

## Usage

TODO

## Contributing

Contributions are very welcomed, however please follow the below guidelines.

- First open an issue describing the bug or enhancement so it can be
discussed.  
- Try to match current naming conventions as closely as possible.  
- Create a Pull Request with your changes against the master branch.

# Licence

discordgo-i18n is available under the MIT license. See the [LICENSE](LICENSE) file for more info.