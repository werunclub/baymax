package util

import (
	"context"
	"os"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/sirupsen/logrus"

	"github.com/werunclub/baymax/v2/rpc/helpers"
)

var (
	// T 翻译方法
	T i18n.TranslateFunc
)

func init() {
	Init()
}

// Init 初始化
func Init() {
	languageCode := os.Getenv("DEFAULTLANGUAGE")
	if languageCode == "" {
		languageCode = "zh-Hans"
	}

	logrus.WithField("lang_code", languageCode).Debugf("Init TranslateFunc")
	T, _ = i18n.Tfunc(languageCode)
}

// TfuncForRPC 返回翻译方法
func TfuncForRPC(ctx context.Context, languageCode string) (i18n.TranslateFunc, error) {
	meta := helpers.NewMetaDataFormContext(ctx)
	acceptLang := meta.Get("lang")
	return i18n.Tfunc(acceptLang, languageCode)
}

func GetTfunc(languageCode string, defaultlanguageCode string) (i18n.TranslateFunc, error) {
	return i18n.Tfunc(languageCode, defaultlanguageCode)
}
