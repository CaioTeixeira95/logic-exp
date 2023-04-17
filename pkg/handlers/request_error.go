package handlers

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "this field is required"
	}
	return ""
}

func ParseRequestError(err error) gin.H {
	if vErrs, ok := err.(validator.ValidationErrors); ok {
		details := gin.H{}
		for _, vErr := range vErrs {
			details[strings.ToLower(vErr.Field())] = msgForTag(vErr.Tag())
		}
		return gin.H{"details": details}
	}

	if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
		return gin.H{
			"details": gin.H{
				jsonErr.Field: "invalid type provided for this field",
			},
		}
	}

	return gin.H{}
}
