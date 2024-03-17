package models

import (
	"bytes"
	"github.com/go-playground/validator/v10"
	_ "mvdan.cc/xurls/v2"
	"unicode"
)

func setupAlertValidations() {
	if validate == nil {
		validate = validator.New()
	}
	validate.RegisterStructValidation(templateCreateValidation, Alert{})
}

type Alert struct {
	ID            string        `json:"alert_id" validate:"required"`
	UserID        string        `json:"user_id" validate:"required"`
	TransactionID string        `json:"transaction_id" validate:"required"`
	Content       AlterPayloads `json:"payloads" validate:"required"`
}

type AlterPayloads struct {
	SMS SMS `json:"SMS" validate:"required"`
}

type CustomerRequest struct {
	ID         string            `json:"customer_id" validate:"required"`
	Type       string            `json:"type" validate:"required,oneof=customer event"`
	Attributes map[string]string `json:"*"`
}

type AttributesRequest struct {
	Name      string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type SMS struct {
	Recipient  string      `json:"recipient"`
	Attributes interface{} `json:"personalized_attributes"`
}

type AlertSuccessResponse struct {
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func (t *Alert) Validate() error {
	return validate.Struct(t)
}

func (t *Alert) Marshal() (*bytes.Buffer, error) {
	return marshalJSON(t)
}

func templateCreateValidation(sl validator.StructLevel) {
	template, _ := sl.Current().Interface().(Alert)
	validateTemplateId(sl, template)
}

func validateTemplateId(sl validator.StructLevel, template Alert) {
	if !isSnakeCaseOrNum(template.ID) {
		sl.ReportError(template.ID, "id", "Id", "idnotsnakecase", "")
	}
}

func isSnakeCaseOrNum(s string) bool {
	for _, r := range s {
		if !unicode.IsNumber(r) && !unicode.IsLower(r) && r != '_' {
			return false
		}
	}
	return true
}
