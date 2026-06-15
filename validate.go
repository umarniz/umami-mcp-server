package main

import "fmt"

const reportTypeFunnel = "funnel"

var validReportTypes = map[string]bool{
	"attribution":    true,
	"breakdown":      true,
	reportTypeFunnel: true,
	"goal":           true,
	"goals":          true,
	"insights":       true,
	"journey":        true,
	"performance":    true,
	"retention":      true,
	"revenue":        true,
	"utm":            true,
}

func validateWebsiteID(id string) error {
	if id == "" || len(id) > 36 {
		return fmt.Errorf("invalid website ID")
	}
	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') || c == '-') {
			return fmt.Errorf("invalid website ID")
		}
	}
	return nil
}

func validateReportID(id string) error {
	if err := validateWebsiteID(id); err != nil {
		return fmt.Errorf("invalid report ID")
	}
	return nil
}

func validateReportType(reportType string) error {
	if !validReportTypes[reportType] {
		return fmt.Errorf("invalid report type")
	}
	return nil
}
