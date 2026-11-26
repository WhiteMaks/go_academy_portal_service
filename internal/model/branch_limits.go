package model

// Лимиты полей филиала — см. .docs/field-limits.md, migration 7.
const (
	BranchNameMaxLen    = 120
	BranchAddressMaxLen = 300
	BranchNoteMaxLen    = 2000

	BranchesListDefaultLimit = ListQueryDefaultLimit
	BranchesListMaxLimit     = ListQueryMaxLimit
	BranchesSearchQMaxLen    = 200
)
