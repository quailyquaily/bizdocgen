package i18n

import (
	"embed"
	"fmt"
	"log/slog"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var LocaleFS embed.FS

type (
	I18nBundle struct {
		bundle     *i18n.Bundle
		localizers map[string]*i18n.Localizer
	}
)

func New() *I18nBundle {
	localizers := make(map[string]*i18n.Localizer)
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	supported := []struct {
		key      string
		filename string
		tag      string
	}{
		{key: "en", filename: "en.toml", tag: "en"},
		{key: "ja", filename: "ja.toml", tag: "ja"},
		{key: "zh_cn", filename: "zh_cn.toml", tag: "zh-CN"},
		{key: "zh_tw", filename: "zh_tw.toml", tag: "zh-TW"},
	}

	for _, lang := range supported {
		_, err := bundle.LoadMessageFileFS(LocaleFS, fmt.Sprintf("locales/%s", lang.filename))
		if err != nil {
			slog.Warn("[i18n] failed to load locale file", "error", err, "lang", lang.key)
			continue
		}
		localizer := i18n.NewLocalizer(bundle, lang.tag, lang.key)
		localizers[lang.key] = localizer
		localizers[lang.tag] = localizer
	}
	return &I18nBundle{
		bundle:     bundle,
		localizers: localizers,
	}
}

func (i *I18nBundle) Localizer(lang string) *i18n.Localizer {
	if l, ok := i.localizers[normalizeLang(lang)]; ok {
		return l
	}
	return i.localizers["en"] // fallback to English
}

func normalizeLang(lang string) string {
	lang = strings.TrimSpace(lang)
	lang = strings.ReplaceAll(lang, "-", "_")
	lang = strings.ToLower(lang)
	switch lang {
	case "zh_cn", "zh_hans", "zh_hans_cn", "zh_chs", "zh_sg":
		return "zh_cn"
	case "zh_tw", "zh_hant", "zh_hant_tw", "zh_cht", "zh_hk":
		return "zh_tw"
	case "ja_jp":
		return "ja"
	case "en_us", "en_gb":
		return "en"
	default:
		return lang
	}
}

func (i *I18nBundle) T(lang, msgID string, data any) (string, error) {
	localizer := i.Localizer(lang)
	output, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    msgID,
		TemplateData: data,
	})
	if err != nil {
		return "", err
	}
	return output, nil
}

func (i *I18nBundle) MusT(lang, msgID string, data any) string {
	output, _ := i.T(lang, msgID, data)
	return output
}
