package models

// Data is a datastruct containing common stuff for the templates
type Data struct {
	Protos []*Proto
	// AppNameUnderscoredUpper is an UPPERCASE underscored project proto
	AppNameUnderscoredUpper string
	// AppName is application name
	AppName string
	// CmdName is a project name
	CmdName string
	// Go Module
	Module string
}
