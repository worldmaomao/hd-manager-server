package utils

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func NewValidator() (*validator.Validate, error) {
	uni := ut.New(zh.New())
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		return nil, err
	}
	return validate, nil
}
