package main

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed resources/views/_tests
var testEmbeddedFiles embed.FS

func TestSetupTemplateRegistry(t *testing.T) {
	// Override the global embededFiles variable with our test files
	embededFiles = testEmbeddedFiles

	registry := SetupTemplateRegistry("resources/views/_tests")

	// Test the number of templates
	assert.Len(t, registry.templates, 2, "Expected 2 templates in the registry")

	// Test the contents of the templates map
	expectedTemplates := []string{
		"resources/views/_tests/template1.html",
		"resources/views/_tests/template2.html",
	}
	for _, templateName := range expectedTemplates {
		assert.Contains(t, registry.templates, templateName, "Expected template %s not found in registry", templateName)
	}

	// Test the contents of the baseTemplatePaths map
	expectedBaseTemplatePaths := map[string]string{
		"resources/views/_tests/template1.html": "layout.html",
		"resources/views/_tests/template2.html": "body",
	}
	assert.Equal(t, expectedBaseTemplatePaths, registry.baseTemplatePaths, "BaseTemplatePaths do not match expected values")
}

// TestSetupTemplateRegistryEmptyDirectory tests the behavior when given an empty directory
func TestSetupTemplateRegistryEmptyDirectory(t *testing.T) {
	registry := SetupTemplateRegistry("resources/views/_tests/empty")

	assert.Empty(t, registry.templates, "Expected no templates for an empty directory")
	assert.Empty(t, registry.baseTemplatePaths, "Expected no base template paths for an empty directory")
}

// TestSetupTemplateRegistryInvalidTemplate tests the behavior with an invalid template file
func TestSetupTemplateRegistryInvalidTemplate(t *testing.T) {
	// Override the global embededFiles variable with our test files including an invalid template
	embededFiles = testEmbeddedFiles

	registry := SetupTemplateRegistry("resources/views/_tests/invalid")

	assert.Len(t, registry.templates, 1, "Expected 1 valid template in the registry")
	assert.NotContains(t, registry.templates, "resources/views/_tests/invalid/invalid_template.html", "Invalid template should not be in the registry")
}
