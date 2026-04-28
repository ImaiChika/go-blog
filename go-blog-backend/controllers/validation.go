package controllers

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func bindErrorMessage(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) && len(validationErrors) > 0 {
		validationError := validationErrors[0]
		field := validationFieldName(validationError.Field())

		switch validationError.Tag() {
		case "required":
			return fmt.Sprintf("%s不能为空", field)
		case "min":
			if validationError.Kind().String() == "int" {
				return fmt.Sprintf("%s不能小于%s", field, validationError.Param())
			}
			return fmt.Sprintf("%s长度不能少于%s", field, validationError.Param())
		case "max":
			if validationError.Kind().String() == "int" {
				return fmt.Sprintf("%s不能大于%s", field, validationError.Param())
			}
			return fmt.Sprintf("%s长度不能超过%s", field, validationError.Param())
		case "url":
			return fmt.Sprintf("%s必须是合法的 URL", field)
		default:
			return fmt.Sprintf("%s格式不正确", field)
		}
	}

	return "请求参数格式错误"
}

func validationFieldName(field string) string {
	switch field {
	case "Username":
		return "用户名"
	case "Password":
		return "密码"
	case "OldPassword":
		return "旧密码"
	case "NewPassword":
		return "新密码"
	case "Title":
		return "标题"
	case "Content":
		return "内容"
	case "CoverImage":
		return "封面图"
	case "Page":
		return "页码"
	case "PageSize":
		return "每页数量"
	default:
		return field
	}
}
