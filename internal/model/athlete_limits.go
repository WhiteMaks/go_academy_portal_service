package model

// Лимиты полей спортсмена — см. .docs/field-limits.md, migration 4, react api/constants/athleteFieldLimits.ts.
// Используются в athlete_validators.go и athlete_service (не дублировать magic numbers в binding/sql).
const (
	AthleteNameMaxLen  = 100
	AthleteNoteMaxLen  = 2000
	AthletePhoneMaxLen = 16 // E.164: '+' и до 15 цифр

	AthletesListDefaultLimit = ListQueryDefaultLimit
	AthletesListMaxLimit     = ListQueryMaxLimit
	AthletesSearchQMaxLen    = 200
)
