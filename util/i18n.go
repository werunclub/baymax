package util

import (
	"context"

	"github.com/nicksnyder/go-i18n/i18n"

	"baymax/rpc/helpers"
)

// TfuncForRPC 返回翻译方法
func TfuncForRPC(ctx context.Context, defaultLanguage string) (i18n.TranslateFunc, error) {
	meta := helpers.NewMetaDataFormContext(ctx)
	acceptLang := meta.Get("lang")
	return i18n.Tfunc(acceptLang, defaultLanguage)
}
