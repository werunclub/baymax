package util

import (
	"context"
	"os"

	"github.com/nicksnyder/go-i18n/i18n"

	"baymax/rpc/helpers"
)

var (
	// T 翻译方法
	T i18n.TranslateFunc
)

func init() {
	defaultLanguage := os.Getenv("DEFAULTLANGUAGE")
	if defaultLanguage == "" {
		defaultLanguage = "zh-Hans"
	}

	T, _ = i18n.Tfunc(defaultLanguage)
}

// TfuncForRPC 返回翻译方法
func TfuncForRPC(ctx context.Context, defaultLanguage string) (i18n.TranslateFunc, error) {
	meta := helpers.NewMetaDataFormContext(ctx)
	acceptLang := meta.Get("lang")
	return i18n.Tfunc(acceptLang, defaultLanguage)
}
