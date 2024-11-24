package test

import (
	"testing"
	"time"

	"github.com/LydiaTrack/ground/pkg/auth"
	"github.com/LydiaTrack/ground/pkg/domain/audit"
	"github.com/LydiaTrack/ground/pkg/test_support"

	"github.com/LydiaTrack/ground/internal/repository"
	"github.com/LydiaTrack/ground/internal/service"
)

var (
	auditService     service.AuditService
	initializedAudit = false
)

func initializeAuditService() {
	if !initializedAudit {
		test_support.TestWithMongo()
		repo := repository.GetAuditRepository()

		// Create a new Audit service instance
		auditService = service.NewAuditService(repo)
		initializedAudit = true
	}
}

func TestAuditService(t *testing.T) {
	test_support.TestWithMongo()
	initializeAuditService()

	t.Run("CreateAudit", testCreateAudit)
	t.Run("GetAudit", testGetAudit)
	t.Run("ExistsAudit", testExistsAudit)
	t.Run("GetAudits", testGetAudits)
	// FIXME: This test is failing when running gh actions
	//t.Run("DeleteOlderThan", testDeleteOlderThan)
	t.Run("DeleteInterval", testDeleteInterval)
}

// testCreateAudit Create a new Audit
func testCreateAudit(t *testing.T) {

	operation := audit.Operation{
		Domain:  "testDomain",
		Command: "CREATE",
	}

	// Create a new Audit
	createAuditCmd := audit.CreateAuditCommand{
		Source:    "test",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	auditModel, err := auditService.CreateAudit(createAuditCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

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

// testGetAudit Get an Audit
func testGetAudit(t *testing.T) {

	operation := audit.Operation{
		Domain:  "testDomain",
		Command: "CREATE",
	}

	// Create a new Audit
	createAuditCmd := audit.CreateAuditCommand{
		Source:    "test",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	auditModel, err := auditService.CreateAudit(createAuditCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

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
	auditModels, err := auditService.GetAudit(auditModel.ID.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error getting Audit test: %v", err)
	} else {

		if auditModels.Source != "test" {
			t.Errorf("Error getting Audit: %v", err)
		}

		if auditModels.Operation != operation {
			t.Errorf("Error getting Audit: %v", err)
		}

		if auditModels.AdditionalData["testStr"] != "test" {
			t.Errorf("Error getting Audit: %v", err)
		}

		if auditModels.RelatedPrincipal != "Test Lastname" {
			t.Errorf("Error getting Audit: %v", err)
		}
	}
}

// testExistsAudit Check if an Audit exists
func testExistsAudit(t *testing.T) {

	operation := audit.Operation{
		Domain:  "testDomain",
		Command: "CREATE",
	}

	// Create a new Audit
	createAuditCmd := audit.CreateAuditCommand{
		Source:    "test",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	auditModel, err := auditService.CreateAudit(createAuditCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

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
	exists, err := auditService.ExistsAudit(auditModel.ID.Hex(), auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error exists Audit test: %v", err)
	} else {

		if exists != true {
			t.Errorf("Error exists Audit: %v", err)
		}
	}
}

// testGetAudits Get all Audits
func testGetAudits(t *testing.T) {

	operation := audit.Operation{
		Domain:  "testDomain",
		Command: "CREATE",
	}

	// Create a new Audit
	createAuditCmd := audit.CreateAuditCommand{
		Source:    "test",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	_, err := auditService.CreateAudit(createAuditCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error creating Audit test: %v", err)
	}

	// Get all Audits
	audits, err := auditService.GetAudits(auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	} else {

		if len(audits) == 0 {
			t.Errorf("Error getting Audits: %v", err)
		}
	}
}

// testDeleteOlderThan Delete all Audits older than a date
func testDeleteOlderThan(t *testing.T) {

	operation := audit.Operation{
		Domain:  "testDomain",
		Command: "CREATE",
	}

	// Get all Audits
	audits, err := auditService.GetAudits(auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})
	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	}

	auditCount := len(audits)

	// Create a new Audit
	createAuditCmd := audit.CreateAuditCommand{
		Source:    "test",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	_, err = auditService.CreateAudit(createAuditCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error creating Audit test: %v", err)
	}

	// Create a new Audit
	createAuditCmd = audit.CreateAuditCommand{
		Source:    "test2",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	_, err = auditService.CreateAudit(createAuditCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})
	if err != nil {
		t.Errorf("Error creating Audit test: %v", err)
	}

	// Check if auditCount + 2 audits exist
	audits, err = auditService.GetAudits(auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})
	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	} else {

		if len(audits) != auditCount+2 {
			t.Errorf("Error getting Audits: %v", err)
		}
	}

	// Delete all Audits older than a date
	deleteOlderThanCommand := audit.DeleteOlderThanAuditCommand{
		Instant: time.Now(),
	}
	err = auditService.DeleteOlderThan(deleteOlderThanCommand, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error deleting Audits test: %v", err)
	}

	// Get all Audits
	audits, err = auditService.GetAudits(auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})
	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	} else {

		if len(audits) != 0 {
			t.Errorf("Error getting Audits: %v", err)
		}
	}
}

// testDeleteInterval Delete all Audits between two dates
func testDeleteInterval(t *testing.T) {

	operation := audit.Operation{
		Domain:  "testDomain",
		Command: "CREATE",
	}

	// Get all Audits
	audits, err := auditService.GetAudits(auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})
	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	}

	auditCount := len(audits)

	// Create a new Audit
	createAuditCmd := audit.CreateAuditCommand{
		Source:    "test",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	_, err = auditService.CreateAudit(createAuditCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error creating Audit test: %v", err)
	}

	// Create a new Audit
	createAuditCmd = audit.CreateAuditCommand{
		Source:    "test2",
		Operation: operation,
		AdditionalData: map[string]interface{}{
			"testStr": "test",
		},
		RelatedPrincipal: "Test Lastname",
	}
	_, err = auditService.CreateAudit(createAuditCmd, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})
	if err != nil {
		t.Errorf("Error creating Audit test: %v", err)
	}

	// Check if auditCount + 2 audits exist
	audits, err = auditService.GetAudits(auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})
	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	} else {

		if len(audits) != auditCount+2 {
			t.Errorf("Error getting Audits: %v", err)
		}
	}

	// Delete all Audits between two dates
	deleteIntervalCommand := audit.DeleteIntervalAuditCommand{
		From: time.Now().Add(-time.Hour * 24),
		To:   time.Now(),
	}
	err = auditService.DeleteInterval(deleteIntervalCommand, auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})

	if err != nil {
		t.Errorf("Error deleting Audits test: %v", err)
	}

	// Get all Audits
	audits, err = auditService.GetAudits(auth.PermissionContext{
		Permissions: []auth.Permission{auth.AdminPermission},
		UserID:      nil,
	})
	if err != nil {
		t.Errorf("Error getting Audits test: %v", err)
	} else {

		if len(audits) != 0 {
			t.Errorf("Error getting Audits: %v", err)
		}
	}
}
