package model

import (
	"errors"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterHallBindingValidatorsFromGin() error {
	engine, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("binding validator engine is not *validator.Validate")
	}
	return RegisterHallBindingValidators(engine)
}

func RegisterHallBindingValidators(v *validator.Validate) error {
	if err := v.RegisterValidation("hall_name_max", validateHallNameMax); err != nil {
		return err
	}
	if err := v.RegisterValidation("hall_subtitle_max", validateHallSubtitleMax); err != nil {
		return err
	}
	if err := v.RegisterValidation("hall_note_max", validateHallNoteMax); err != nil {
		return err
	}
	if err := v.RegisterValidation("hall_capacity_range", validateHallCapacityRange); err != nil {
		return err
	}
	if err := v.RegisterValidation("halls_search_q_max", validateHallsSearchQMax); err != nil {
		return err
	}
	return v.RegisterValidation("halls_page_limit", validateHallsPageLimit)
}

func validateHallsSearchQMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= HallsSearchQMaxLen
}

func validateHallsPageLimit(fl validator.FieldLevel) bool {
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
	return v <= int64(HallsListMaxLimit)
}

func validateHallNameMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= HallNameMaxLen
}

func validateHallSubtitleMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= HallSubtitleMaxLen
}

func validateHallNoteMax(fl validator.FieldLevel) bool {
	s, ok := validatorStringField(fl)
	if !ok {
		return false
	}
	s = strings.TrimSpace(s)
	if s == "" {
		return true
	}
	return utf8.RuneCountInString(s) <= HallNoteMaxLen
}

func validateHallCapacityRange(fl validator.FieldLevel) bool {
	if fl.Field().Kind() != reflect.Int {
		return false
	}
	v := fl.Field().Int()
	return v >= HallCapacityMin && v <= HallCapacityMax
}

func ValidateHallCreateFieldLengths(req PostHallV1Request) error {
	if runeLen(strings.TrimSpace(req.Name)) > HallNameMaxLen {
		return errors.New("hall field too long")
	}
	if req.Subtitle != nil && runeLen(strings.TrimSpace(*req.Subtitle)) > HallSubtitleMaxLen {
		return errors.New("hall field too long")
	}
	if req.Note != nil && runeLen(strings.TrimSpace(*req.Note)) > HallNoteMaxLen {
		return errors.New("hall field too long")
	}
	if req.Capacity < HallCapacityMin || req.Capacity > HallCapacityMax {
		return errors.New("invalid capacity")
	}
	return nil
}

func ValidateGetBranchHallsV1Query(query GetBranchHallsV1Query) error {
	if runeLen(strings.TrimSpace(query.Q)) > HallsSearchQMaxLen {
		return errors.New("search query too long")
	}
	if query.Limit < 0 || query.Limit > HallsListMaxLimit {
		return errors.New("invalid limit")
	}
	if query.Offset < 0 {
		return errors.New("invalid offset")
	}
	return nil
}
