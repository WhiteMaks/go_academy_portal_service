package model

import (
	"errors"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterBranchBindingValidatorsFromGin() error {
	engine, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("binding validator engine is not *validator.Validate")
	}
	return RegisterBranchBindingValidators(engine)
}

func RegisterBranchBindingValidators(v *validator.Validate) error {
	if err := v.RegisterValidation("branch_name_max", validateBranchNameMax); err != nil {
		return err
	}
	if err := v.RegisterValidation("branch_address_max", validateBranchAddressMax); err != nil {
		return err
	}
	if err := v.RegisterValidation("branch_note_max", validateBranchNoteMax); err != nil {
		return err
	}
	if err := v.RegisterValidation("branches_search_q_max", validateBranchesSearchQMax); err != nil {
		return err
	}
	return v.RegisterValidation("branches_page_limit", validateBranchesPageLimit)
}

func validateBranchesSearchQMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= BranchesSearchQMaxLen
}

func validateBranchesPageLimit(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.Int {
		return false
	}
	v := fl.Field().Int()
	if v == 0 {
		return true
	}
	if v < 0 {
		return false
	}
	return v <= int64(BranchesListMaxLimit)
}

func validateBranchNameMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= BranchNameMaxLen
}

func validateBranchAddressMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= BranchAddressMaxLen
}

func validateBranchNoteMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= BranchNoteMaxLen
}

func ValidateBranchCreateFieldLengths(req PostBranchV1Request) error {
	if runeLen(strings.TrimSpace(req.Name)) > BranchNameMaxLen ||
		runeLen(strings.TrimSpace(req.Address)) > BranchAddressMaxLen {
		return errors.New("branch field too long")
	}
	if req.Note != nil && runeLen(strings.TrimSpace(*req.Note)) > BranchNoteMaxLen {
		return errors.New("branch field too long")
	}
	return nil
}

func ValidateGetBranchesV1Query(query GetBranchesV1Query) error {
	if runeLen(strings.TrimSpace(query.Q)) > BranchesSearchQMaxLen {
		return errors.New("search query too long")
	}
	if query.Limit < 0 || query.Limit > BranchesListMaxLimit {
		return errors.New("invalid limit")
	}
	if query.Offset < 0 {
		return errors.New("invalid offset")
	}
	return nil
}
