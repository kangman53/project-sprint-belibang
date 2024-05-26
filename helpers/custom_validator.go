package helpers

import (
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator"
)

func validatePhoneNumber(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	pattern := `^\+62\d+$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

func validateUrl(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	pattern := `^(?:https?:\/\/)?(?:www\.)?(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(?:\/[^\s]*)?$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

func validateCategory(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	categVal := fl.Param()

	if categVal == "merchant" {
		return checkCategory(MerchantCategory, value)
	}
	return checkCategory(ItemCategory, value)
}

func checkCategory(categories []string, val string) bool {
	for _, categ := range categories {
		if categ == val {
			return true
		}
	}
	return false
}

func validateGeoCoord(fl validator.FieldLevel) bool {
	value := fl.Field().Float()
	checkVal := fl.Param()
	maxCoord := 0.0

	if checkVal == "lat" {
		maxCoord = 90.0
	} else {
		maxCoord = 180.0
	}

	return value >= -maxCoord && value <= maxCoord
}

func validateISO8601DateTime(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		_, err := time.Parse(time.RFC3339, value)
		return err == nil
	}
	return false
}

func validateInt16Length(fl validator.FieldLevel) bool {
	str := strconv.Itoa(fl.Field().Interface().(int))
	return len(str) == 16
}

func RegisterCustomValidator(validator *validator.Validate) {
	// validator.RegisterValidation() -> if you want to create new tags rule to be used on struct entity
	// validator.RegisterStructValidation() -> if you want to create validator then access all fields to the struct entity

	validator.RegisterValidation("validateCategory", validateCategory)
	validator.RegisterValidation("validateUrl", validateUrl)
	validator.RegisterValidation("validateGeoCoord", validateGeoCoord)
	validator.RegisterValidation("ISO8601DateTime", validateISO8601DateTime)
	validator.RegisterValidation("int16length", validateInt16Length)
}
