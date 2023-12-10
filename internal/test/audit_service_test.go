package test

import (
	"lydia-track-base/internal/domain/audit"
	"lydia-track-base/internal/domain/audit/command"
	"lydia-track-base/internal/domain/auth"
	"lydia-track-base/internal/repository"
	"lydia-track-base/internal/service"
	"lydia-track-base/internal/test_support"
	"testing"
	"time"
)

// TestNewUserService Create a new Audit service instance with AuditMongoRepository
func TestNewAuditService(t *testing.T) {
	test_support.TestWithMongo()
	repo := repository.GetAuditRepository()

	// Create a new Audit service instance
	service.NewAuditService(repo)
}

// TestCreateAudit Create a new Audit
func TestCreateAudit(t *testing.T) {
	test_support.TestWithMongo()

	// Create a new Audit service instance
	auditService := service.NewAuditService(repository.GetAuditRepository())
	operation := audit.Operation{
		Domain:  "testDomain",
		Command: "CREATE",
	}

	// Create a new Audit
	createAuditCmd := commands.CreateAuditCommand{
		Source:    "test",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	auditModel, err := auditService.CreateAudit(createAuditCmd, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error creating Audit test: %v", err)
	} else {

		if auditModel.Source != "test" {
			t.Errorf("Error creating Audit: %v", err)
		}

		if auditModel.Operation != operation {
			t.Errorf("Error creating Audit: %v", err)
		}

		if auditModel.AdditionalData["testStr"] != "test" {
			t.Errorf("Error creating Audit: %v", err)
		}

		if auditModel.RelatedPrincipal != "Test Lastname" {
			t.Errorf("Error creating Audit: %v", err)
		}
	}
}

// TestGetAudit Get an Audit
func TestGetAudit(t *testing.T) {
	test_support.TestWithMongo()

	// Create a new Audit service instance
	auditService := service.NewAuditService(repository.GetAuditRepository())
	operation := audit.Operation{
		Domain:  "testDomain",
		Command: "CREATE",
	}

	// Create a new Audit
	createAuditCmd := commands.CreateAuditCommand{
		Source:    "test",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	auditModel, err := auditService.CreateAudit(createAuditCmd, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error creating Audit test: %v", err)
	} else {

		if auditModel.Source != "test" {
			t.Errorf("Error creating Audit: %v", err)
		}

		if auditModel.Operation != operation {
			t.Errorf("Error creating Audit: %v", err)
		}

		if auditModel.AdditionalData["testStr"] != "test" {
			t.Errorf("Error creating Audit: %v", err)
		}

		if auditModel.RelatedPrincipal != "Test Lastname" {
			t.Errorf("Error creating Audit: %v", err)
		}
	}

	// Get an Audit
	audit, err := auditService.GetAudit(auditModel.ID.Hex(), []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error getting Audit test: %v", err)
	} else {

		if audit.Source != "test" {
			t.Errorf("Error getting Audit: %v", err)
		}

		if audit.Operation != operation {
			t.Errorf("Error getting Audit: %v", err)
		}

		if audit.AdditionalData["testStr"] != "test" {
			t.Errorf("Error getting Audit: %v", err)
		}

		if audit.RelatedPrincipal != "Test Lastname" {
			t.Errorf("Error getting Audit: %v", err)
		}
	}
}

// TestExistsAudit Check if an Audit exists
func TestExistsAudit(t *testing.T) {
	test_support.TestWithMongo()

	// Create a new Audit service instance
	auditService := service.NewAuditService(repository.GetAuditRepository())
	operation := audit.Operation{
		Domain:  "testDomain",
		Command: "CREATE",
	}

	// Create a new Audit
	createAuditCmd := commands.CreateAuditCommand{
		Source:    "test",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	auditModel, err := auditService.CreateAudit(createAuditCmd, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error creating Audit test: %v", err)
	} else {

		if auditModel.Source != "test" {
			t.Errorf("Error creating Audit: %v", err)
		}

		if auditModel.Operation != operation {
			t.Errorf("Error creating Audit: %v", err)
		}

		if auditModel.AdditionalData["testStr"] != "test" {
			t.Errorf("Error creating Audit: %v", err)
		}

		if auditModel.RelatedPrincipal != "Test Lastname" {
			t.Errorf("Error creating Audit: %v", err)
		}
	}

	// Check if an Audit exists
	exists, err := auditService.ExistsAudit(auditModel.ID.Hex(), []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error exists Audit test: %v", err)
	} else {

		if exists != true {
			t.Errorf("Error exists Audit: %v", err)
		}
	}
}

// TestGetAudits Get all Audits
func TestGetAudits(t *testing.T) {
	test_support.TestWithMongo()

	// Create a new Audit service instance
	auditService := service.NewAuditService(repository.GetAuditRepository())
	operation := audit.Operation{
		Domain:  "testDomain",
		Command: "CREATE",
	}

	// Create a new Audit
	createAuditCmd := commands.CreateAuditCommand{
		Source:    "test",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	_, err := auditService.CreateAudit(createAuditCmd, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error creating Audit test: %v", err)
	}

	// Get all Audits
	audits, err := auditService.GetAudits([]auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	} else {

		if len(audits) == 0 {
			t.Errorf("Error getting Audits: %v", err)
		}
	}
}

// TestDeleteOlderThan Delete all Audits older than a date
func TestDeleteOlderThan(t *testing.T) {
	test_support.TestWithMongo()

	// Create a new Audit service instance
	auditService := service.NewAuditService(repository.GetAuditRepository())
	operation := audit.Operation{
		Domain:  "testDomain",
		Command: "CREATE",
	}

	// Get all Audits
	audits, err := auditService.GetAudits([]auth.Permission{auth.AdminPermission})
	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	}

	auditCount := len(audits)

	// Create a new Audit
	createAuditCmd := commands.CreateAuditCommand{
		Source:    "test",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	_, err = auditService.CreateAudit(createAuditCmd, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error creating Audit test: %v", err)
	}

	// Create a new Audit
	createAuditCmd = commands.CreateAuditCommand{
		Source:    "test2",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	_, err = auditService.CreateAudit(createAuditCmd, []auth.Permission{auth.AdminPermission})
	if err != nil {
		t.Errorf("Error creating Audit test: %v", err)
	}

	// Check if auditCount + 2 audits exist
	audits, err = auditService.GetAudits([]auth.Permission{auth.AdminPermission})
	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	} else {

		if len(audits) != auditCount+2 {
			t.Errorf("Error getting Audits: %v", err)
		}
	}

	// Delete all Audits older than a date
	deleteOlderThanCommand := commands.DeleteOlderThanAuditCommand{
		Instant: time.Now(),
	}
	err = auditService.DeleteOlderThan(deleteOlderThanCommand, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error deleting Audits test: %v", err)
	}

	// Get all Audits
	audits, err = auditService.GetAudits([]auth.Permission{auth.AdminPermission})
	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	} else {

		if len(audits) != 0 {
			t.Errorf("Error getting Audits: %v", err)
		}
	}
}

// TestDeleteInterval Delete all Audits between two dates
func TestDeleteInterval(t *testing.T) {
	test_support.TestWithMongo()

	// Create a new Audit service instance
	auditService := service.NewAuditService(repository.GetAuditRepository())
	operation := audit.Operation{
		Domain:  "testDomain",
		Command: "CREATE",
	}

	// Get all Audits
	audits, err := auditService.GetAudits([]auth.Permission{auth.AdminPermission})
	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	}

	auditCount := len(audits)

	// Create a new Audit
	createAuditCmd := commands.CreateAuditCommand{
		Source:    "test",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	_, err = auditService.CreateAudit(createAuditCmd, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error creating Audit test: %v", err)
	}

	// Create a new Audit
	createAuditCmd = commands.CreateAuditCommand{
		Source:    "test2",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	_, err = auditService.CreateAudit(createAuditCmd, []auth.Permission{auth.AdminPermission})
	if err != nil {
		t.Errorf("Error creating Audit test: %v", err)
	}

	// Check if auditCount + 2 audits exist
	audits, err = auditService.GetAudits([]auth.Permission{auth.AdminPermission})
	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	} else {

		if len(audits) != auditCount+2 {
			t.Errorf("Error getting Audits: %v", err)
		}
	}

	// Delete all Audits between two dates
	deleteIntervalCommand := commands.DeleteIntervalAuditCommand{
		From: time.Now().Add(-time.Hour * 24),
		To:   time.Now(),
	}
	err = auditService.DeleteInterval(deleteIntervalCommand, []auth.Permission{auth.AdminPermission})

	if err != nil {
		t.Errorf("Error deleting Audits test: %v", err)
	}

	// Get all Audits
	audits, err = auditService.GetAudits([]auth.Permission{auth.AdminPermission})
	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	} else {

		if len(audits) != 0 {
			t.Errorf("Error getting Audits: %v", err)
		}
	}
}
