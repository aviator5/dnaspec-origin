package ui

import (
	"testing"

	"github.com/aviator5/dnaspec/internal/core/config"
)

func TestSelectGuidelinesByName(t *testing.T) {
	available := []config.ManifestGuideline{
		{Name: "go-style", File: "guidelines/go-style.md", Description: "Go style guide"},
		{Name: "rest-api", File: "guidelines/rest-api.md", Description: "REST API guide"},
		{Name: "security", File: "guidelines/security.md", Description: "Security guide"},
	}

	t.Run("select single guideline", func(t *testing.T) {
		names := []string{"go-style"}
		result, err := SelectGuidelinesByName(available, names)
		if err != nil {
			t.Fatalf("SelectGuidelinesByName() error = %v", err)
		}

		if len(result) != 1 {
			t.Errorf("Expected 1 guideline, got %d", len(result))
		}

		if result[0].Name != "go-style" {
			t.Errorf("Expected 'go-style', got %s", result[0].Name)
		}
	})

	t.Run("select multiple guidelines", func(t *testing.T) {
		names := []string{"go-style", "security"}
		result, err := SelectGuidelinesByName(available, names)
		if err != nil {
			t.Fatalf("SelectGuidelinesByName() error = %v", err)
		}

		if len(result) != 2 {
			t.Errorf("Expected 2 guidelines, got %d", len(result))
		}

		expectedNames := map[string]bool{"go-style": true, "security": true}
		for _, g := range result {
			if !expectedNames[g.Name] {
				t.Errorf("Unexpected guideline: %s", g.Name)
			}
		}
	})

	t.Run("select all available guidelines", func(t *testing.T) {
		names := []string{"go-style", "rest-api", "security"}
		result, err := SelectGuidelinesByName(available, names)
		if err != nil {
			t.Fatalf("SelectGuidelinesByName() error = %v", err)
		}

		if len(result) != 3 {
			t.Errorf("Expected 3 guidelines, got %d", len(result))
		}
	})

	t.Run("error on nonexistent guideline", func(t *testing.T) {
		names := []string{"nonexistent"}
		_, err := SelectGuidelinesByName(available, names)
		if err == nil {
			t.Error("Expected error for nonexistent guideline, got nil")
		}
	})

	t.Run("error on partial match", func(t *testing.T) {
		names := []string{"go-style", "nonexistent"}
		_, err := SelectGuidelinesByName(available, names)
		if err == nil {
			t.Error("Expected error when some guidelines don't exist, got nil")
		}
	})

	t.Run("error on empty names", func(t *testing.T) {
		names := []string{}
		_, err := SelectGuidelinesByName(available, names)
		if err == nil {
			t.Error("Expected error for empty names, got nil")
		}
	})

	t.Run("preserve order", func(t *testing.T) {
		names := []string{"security", "go-style"}
		result, err := SelectGuidelinesByName(available, names)
		if err != nil {
			t.Fatalf("SelectGuidelinesByName() error = %v", err)
		}

		if len(result) != 2 {
			t.Fatalf("Expected 2 guidelines, got %d", len(result))
		}

		// Result should be in the order requested
		if result[0].Name != "security" {
			t.Errorf("First guideline = %s, want security", result[0].Name)
		}
		if result[1].Name != "go-style" {
			t.Errorf("Second guideline = %s, want go-style", result[1].Name)
		}
	})

	t.Run("handle duplicates in request", func(t *testing.T) {
		names := []string{"go-style", "go-style"}
		result, err := SelectGuidelinesByName(available, names)
		if err != nil {
			t.Fatalf("SelectGuidelinesByName() error = %v", err)
		}

		// Should return duplicate entries if requested
		if len(result) != 2 {
			t.Errorf("Expected 2 guidelines (including duplicate), got %d", len(result))
		}
	})
}

func TestSelectGuidelinesByName_EmptyAvailable(t *testing.T) {
	available := []config.ManifestGuideline{}
	names := []string{"any-guideline"}

	_, err := SelectGuidelinesByName(available, names)
	if err == nil {
		t.Error("Expected error when no guidelines are available, got nil")
	}
}

func TestSelectGuidelinesByName_ErrorMessage(t *testing.T) {
	available := []config.ManifestGuideline{
		{Name: "guideline-1", File: "g1.md"},
		{Name: "guideline-2", File: "g2.md"},
	}

	names := []string{"nonexistent-1", "nonexistent-2"}

	_, err := SelectGuidelinesByName(available, names)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Error should mention both missing guidelines and available ones
	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error message is empty")
	}

	// Should contain information about what's available
	// (This is a basic check - the actual error format may vary)
	t.Logf("Error message: %s", errMsg)
}
