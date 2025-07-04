package cli

import (
	"testing"
	"github.com/urfave/cli/v2"
)

func TestColumnOrdering(t *testing.T) {
	// Test that applyColumnOrder returns correct number of columns
	ctx := &cli.Context{}
	
	// Test with valid column order
	order := []string{"Size", "Name"}
	result := applyColumnOrder(ctx, order)
	
	if len(result) != 2 {
		t.Errorf("Expected 2 columns for order %v, got %d", order, len(result))
	}
	
	// Test with empty order (should return default)
	emptyOrder := []string{}
	defaultResult := applyColumnOrder(ctx, emptyOrder)
	
	if len(defaultResult) == 0 {
		t.Error("Expected default columns for empty order, got none")
	}
	
	// Test with invalid column names (should be filtered out)
	invalidOrder := []string{"InvalidColumn", "Size"}
	invalidResult := applyColumnOrder(ctx, invalidOrder)
	
	if len(invalidResult) != 1 {
		t.Errorf("Expected 1 valid column from %v, got %d", invalidOrder, len(invalidResult))
	}
}

func TestShouldIncludeColumn(t *testing.T) {
	ctx := &cli.Context{}
	
	// Test that owner column is included by default
	if !shouldIncludeColumn(ctx, "Owner") {
		t.Error("Owner column should be included by default")
	}
	
	// Test that group column is included by default  
	if !shouldIncludeColumn(ctx, "Group") {
		t.Error("Group column should be included by default")
	}
	
	// Test other columns are always included
	if !shouldIncludeColumn(ctx, "Size") {
		t.Error("Size column should always be included")
	}
	
	if !shouldIncludeColumn(ctx, "Name") {
		t.Error("Name column should always be included")
	}
}