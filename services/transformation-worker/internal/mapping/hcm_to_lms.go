package mapping

import (
	"errors"
	"fmt"
	"time"
)

// HCMNewHire represents the raw event from a system like Workday
type HCMNewHire struct {
	EmployeeID   string `json:"id"`
	FullName     string `json:"name"`
	Department   string `json:"dept"`
	HireDate     string `json:"hire_date"`
	Email        string `json:"work_email"`
}

// AristEnrollment represents the schema Arist needs to trigger a learning path
type AristEnrollment struct {
	UserEmail    string    `json:"user_email"`
	CourseID     string    `json:"course_id"`
	TriggerDate  time.Time `json:"trigger_at"`
	Metadata     map[string]string
}

// Transform maps HCM data to LMS requirements with built-in validation
func Transform(input HCMNewHire) (*AristEnrollment, error) {
	// 1. Validation (The "Grit" - catching issues at the edge)
	if input.Email == "" {
		return nil, errors.New("transformation_failed: missing_required_email")
	}

	// 2. Business Logic: Map Department to specific Learning Paths
	courseID := "general_onboarding_101"
	if input.Department == "Engineering" {
		courseID = "eng_security_compliance_v2"
	}

	// 3. Construct the 'Native' Arist Payload
	return &AristEnrollment{
		UserEmail:   input.Email,
		CourseID:    courseID,
		TriggerDate: time.Now().UTC(),
		Metadata: map[string]string{
			"source":      "workday_webhook",
			"external_id": input.EmployeeID,
			"dept":        input.Department,
		},
	}, nil
}
