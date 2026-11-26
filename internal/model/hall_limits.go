package model

// Лимиты полей зала — см. .docs/field-limits.md, migration 7.
const (
	HallNameMaxLen     = 120
	HallSubtitleMaxLen = 200
	HallNoteMaxLen     = 2000
	HallCapacityMin    = 1
	HallCapacityMax    = 999

	HallsListDefaultLimit = ListQueryDefaultLimit
	HallsListMaxLimit     = ListQueryMaxLimit
	HallsSearchQMaxLen    = 200
)
