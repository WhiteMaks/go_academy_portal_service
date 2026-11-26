package model

// Лимиты полей тренера — см. .docs/field-limits.md, migration 5.
const (
	CoachNameMaxLen           = 100
	CoachSpecializationMaxLen = 200
	CoachNoteMaxLen           = 2000
	CoachPhoneMaxLen          = 16 // E.164: '+' и до 15 цифр

	CoachesListDefaultLimit = ListQueryDefaultLimit
	CoachesListMaxLimit     = ListQueryMaxLimit
	CoachesSearchQMaxLen    = 200
)
