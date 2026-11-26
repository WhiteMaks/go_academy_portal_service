package model

import (
	"errors"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterCoachBindingValidatorsFromGin attaches coach length rules to gin's validator engine.
func RegisterCoachBindingValidatorsFromGin() error {
	engine, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("binding validator engine is not *validator.Validate")
	}
	return RegisterCoachBindingValidators(engine)
}

// RegisterCoachBindingValidators registers length rules that read Coach*MaxLen constants.
func RegisterCoachBindingValidators(v *validator.Validate) error {
	if err := v.RegisterValidation("coach_name_max", validateCoachNameMax); err != nil {
		return err
	}
	if err := v.RegisterValidation("coach_specialization_max", validateCoachSpecializationMax); err != nil {
		return err
	}
	if err := v.RegisterValidation("coach_note_max", validateCoachNoteMax); err != nil {
		return err
	}
	if err := v.RegisterValidation("coaches_search_q_max", validateCoachesSearchQMax); err != nil {
		return err
	}
	return v.RegisterValidation("coaches_page_limit", validateCoachesPageLimit)
}

func validateCoachesSearchQMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= CoachesSearchQMaxLen
}

func validateCoachesPageLimit(fl validator.FieldLevel) bool {
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
	return v <= int64(CoachesListMaxLimit)
}

func validateCoachNameMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= CoachNameMaxLen
}

func validateCoachSpecializationMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= CoachSpecializationMaxLen
}

func validateCoachNoteMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= CoachNoteMaxLen
}

// ValidateCoachCreateFieldLengths checks trimmed string lengths (second line after gin binding).
func ValidateCoachCreateFieldLengths(req PostCoachV1Request) error {
	if runeLen(strings.TrimSpace(req.LastName)) > CoachNameMaxLen ||
		runeLen(strings.TrimSpace(req.FirstName)) > CoachNameMaxLen {
		return errors.New("coach field too long")
	}
	if req.Patronymic != nil && runeLen(strings.TrimSpace(*req.Patronymic)) > CoachNameMaxLen {
		return errors.New("coach field too long")
	}
	if runeLen(strings.TrimSpace(req.Specialization)) > CoachSpecializationMaxLen {
		return errors.New("coach field too long")
	}
	if req.Note != nil && runeLen(strings.TrimSpace(*req.Note)) > CoachNoteMaxLen {
		return errors.New("coach field too long")
	}
	return nil
}

// ValidateGetCoachesV1Query repeats list constraints using constants (after gin query binding).
func ValidateGetCoachesV1Query(query GetCoachesV1Query) error {
	if runeLen(strings.TrimSpace(query.Q)) > CoachesSearchQMaxLen {
		return errors.New("search query too long")
	}
	if query.Limit < 0 || query.Limit > CoachesListMaxLimit {
		return errors.New("invalid limit")
	}
	if query.Offset < 0 {
		return errors.New("invalid offset")
	}
	return nil
}
