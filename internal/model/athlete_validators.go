package model

import (
	"errors"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// RegisterAthleteBindingValidatorsFromGin attaches athlete length rules to gin's validator engine.
func RegisterAthleteBindingValidatorsFromGin() error {
	engine, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("binding validator engine is not *validator.Validate")
	}
	return RegisterAthleteBindingValidators(engine)
}

// RegisterAthleteBindingValidators registers length rules that read Athlete*MaxLen constants.
func RegisterAthleteBindingValidators(v *validator.Validate) error {
	if err := v.RegisterValidation("athlete_name_max", validateAthleteNameMax); err != nil {
		return err
	}
	if err := v.RegisterValidation("athlete_note_max", validateAthleteNoteMax); err != nil {
		return err
	}
	if err := v.RegisterValidation("athletes_search_q_max", validateAthletesSearchQMax); err != nil {
		return err
	}
	return v.RegisterValidation("athletes_page_limit", validateAthletesPageLimit)
}

func validateAthletesSearchQMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= AthletesSearchQMaxLen
}

func validateAthletesPageLimit(fl validator.FieldLevel) bool {
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
	return v <= int64(AthletesListMaxLimit)
}

func validateAthleteNameMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= AthleteNameMaxLen
}

func validateAthleteNoteMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= AthleteNoteMaxLen
}

func validatorStringField(fl validator.FieldLevel) (string, bool) {
	f := fl.Field()
	if f.Kind() == reflect.Ptr {
		if f.IsNil() {
			return "", true
		}
		f = f.Elem()
	}
	if f.Kind() != reflect.String {
		return "", false
	}
	return f.String(), true
}

// ValidateAthleteCreateFieldLengths checks trimmed string lengths (second line after gin binding).
func ValidateAthleteCreateFieldLengths(req PostAthleteV1Request) error {
	if runeLen(strings.TrimSpace(req.LastName)) > AthleteNameMaxLen ||
		runeLen(strings.TrimSpace(req.FirstName)) > AthleteNameMaxLen {
		return errors.New("athlete field too long")
	}
	if req.Patronymic != nil && runeLen(strings.TrimSpace(*req.Patronymic)) > AthleteNameMaxLen {
		return errors.New("athlete field too long")
	}
	if req.Note != nil && runeLen(strings.TrimSpace(*req.Note)) > AthleteNoteMaxLen {
		return errors.New("athlete field too long")
	}
	if req.ParentLastName != nil && runeLen(strings.TrimSpace(*req.ParentLastName)) > AthleteNameMaxLen {
		return errors.New("athlete field too long")
	}
	if req.ParentFirstName != nil && runeLen(strings.TrimSpace(*req.ParentFirstName)) > AthleteNameMaxLen {
		return errors.New("athlete field too long")
	}
	if req.ParentPatronymic != nil && runeLen(strings.TrimSpace(*req.ParentPatronymic)) > AthleteNameMaxLen {
		return errors.New("athlete field too long")
	}
	return nil
}

// ValidateGetAthletesV1Query repeats list constraints using constants (after gin query binding).
func ValidateGetAthletesV1Query(query GetAthletesV1Query) error {
	if runeLen(strings.TrimSpace(query.Q)) > AthletesSearchQMaxLen {
		return errors.New("search query too long")
	}
	if query.Limit < 0 || query.Limit > AthletesListMaxLimit {
		return errors.New("invalid limit")
	}
	if query.Offset < 0 {
		return errors.New("invalid offset")
	}
	return nil
}

func runeLen(s string) int {
	return utf8.RuneCountInString(s)
}
