package models

type Category string

const (
	BREAKING_CHANGE Category = "Breaking Changes"
	ADDED           Category = "Added"
	CHANGED         Category = "Changed"
	REFACTOR        Category = "Refactor"
	FIXED           Category = "Fixed"
	DEPENDENCIES    Category = "Dependencies"
	DEPRECATED      Category = "Deprecated"
	REMOVED         Category = "Removed"
	DOCUMENTATION   Category = "Documentation"
	SECURITY        Category = "Security"
)
