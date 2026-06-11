package rollup

import (
	"testing"
	"time"

	"github.com/kentik/ktranslate"
	"github.com/kentik/ktranslate/pkg/eggs/logger"
	lt "github.com/kentik/ktranslate/pkg/eggs/logger/testing"
	"github.com/kentik/ktranslate/pkg/kt"
)

func TestCacheRollup(t *testing.T) {
	// Create test configuration
	cfg := &ktranslate.RollupConfig{
		JoinKey:          "^",
		TopK:             5,
		KeepUndefined:    false,
		MaxMemoryMB:      1,  // Small limit for testing
		MaxKeys:          10, // Small limit for testing
		EmergencyCleanup: true,
	}

	// Create test logger
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()

	// Create rollup definition for sum aggregation
	rd := RollupDef{
		Method:     Sum,
		Metrics:    []string{"bytes"},
		Dimensions: []string{"src_addr"},
	}

	// Create cache rollup
	rollup, err := newCacheRollup(l, rd, cfg, false)
	if err != nil {
		t.Fatalf("Failed to create cache rollup: %v", err)
	}

	// Test data - simulate network flows
	testData := []map[string]interface{}{
		{
			"bytes":       int64(100),
			"src_addr":    "192.168.1.1",
			"sample_rate": int64(1),
			"provider":    kt.Provider("pp"),
		},
		{
			"bytes":       int64(200),
			"src_addr":    "192.168.1.1", // Same source
			"sample_rate": int64(1),
			"provider":    kt.Provider("pp"),
		},
		{
			"bytes":       int64(150),
			"src_addr":    "192.168.1.2", // Different source
			"sample_rate": int64(1),
			"provider":    kt.Provider("pp"),
		},
	}

	// Add test data
	rollup.Add(testData)

	// Export results
	results := rollup.Export()

	// Verify results
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Check aggregation for first source (should be 300 total)
	found1 := false
	found2 := false
	for _, result := range results {
		if result.Dimension == "192.168.1.1" {
			if result.Metric != 300.0 {
				t.Errorf("Expected 300 for 192.168.1.1, got %f", result.Metric)
			}
			if result.Count != 2 {
				t.Errorf("Expected count 2 for 192.168.1.1, got %d", result.Count)
			}
			found1 = true
		}
		if result.Dimension == "192.168.1.2" {
			if result.Metric != 150.0 {
				t.Errorf("Expected 150 for 192.168.1.2, got %f", result.Metric)
			}
			if result.Count != 1 {
				t.Errorf("Expected count 1 for 192.168.1.2, got %d", result.Count)
			}
			found2 = true
		}
	}

	if !found1 {
		t.Error("Results for 192.168.1.1 not found")
	}
	if !found2 {
		t.Error("Results for 192.168.1.2 not found")
	}
}

func TestCacheRollupUnique(t *testing.T) {
	// Create test configuration
	cfg := &ktranslate.RollupConfig{
		JoinKey:          "^",
		TopK:             5,
		KeepUndefined:    false,
		MaxMemoryMB:      1,
		MaxKeys:          10,
		EmergencyCleanup: true,
	}

	// Create test logger
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()

	// Create rollup definition for unique aggregation
	rd := RollupDef{
		Method:     Unique,
		Metrics:    []string{"dst_addr"},
		Dimensions: []string{"src_addr"},
	}

	// Create unique cache rollup
	rollup, err := newCacheRollup(l, rd, cfg, true)
	if err != nil {
		t.Fatalf("Failed to create unique cache rollup: %v", err)
	}

	// Test data - simulate network flows with different destination addresses
	testData := []map[string]interface{}{
		{
			"dst_addr":    "10.0.0.1",
			"src_addr":    "192.168.1.1",
			"sample_rate": int64(1),
			"provider":    kt.Provider("pp"),
		},
		{
			"dst_addr":    "10.0.0.2",
			"src_addr":    "192.168.1.1", // Same source, different destination
			"sample_rate": int64(1),
			"provider":    kt.Provider("pp"),
		},
		{
			"dst_addr":    "10.0.0.1",
			"src_addr":    "192.168.1.1", // Duplicate destination
			"sample_rate": int64(1),
			"provider":    kt.Provider("pp"),
		},
	}

	// Add test data
	rollup.Add(testData)

	// Export results
	results := rollup.Export()

	// Verify results
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	// Check unique count for source (should be 2 unique destinations)
	if results[0].Dimension != "192.168.1.1" {
		t.Errorf("Expected dimension 192.168.1.1, got %s", results[0].Dimension)
	}
	if results[0].Metric != 2.0 {
		t.Errorf("Expected 2 unique destinations, got %f", results[0].Metric)
	}
}

func TestCacheRollupEmergencyCleanup(t *testing.T) {
	// Create test configuration with very small limits
	cfg := &ktranslate.RollupConfig{
		JoinKey:          "^",
		TopK:             100,
		KeepUndefined:    false,
		MaxMemoryMB:      0, // Disable memory limit for this test
		MaxKeys:          5, // Very small key limit
		EmergencyCleanup: true,
	}

	// Create test logger
	l := lt.NewTestContextL(logger.NilContext, t).GetLogger().GetUnderlyingLogger()

	// Create rollup definition
	rd := RollupDef{
		Method:     Sum,
		Metrics:    []string{"bytes"},
		Dimensions: []string{"src_addr"},
	}

	// Create cache rollup
	rollup, err := newCacheRollup(l, rd, cfg, false)
	if err != nil {
		t.Fatalf("Failed to create cache rollup: %v", err)
	}

	// Add more entries than the limit to trigger emergency cleanup
	for i := 0; i < 10; i++ {
		testData := []map[string]interface{}{
			{
				"bytes":       int64(100),
				"src_addr":    string(rune('A' + i)), // Generate different source addresses
				"sample_rate": int64(1),
				"provider":    kt.Provider("pp"),
			},
		}
		rollup.Add(testData)

		// Add some delay to create age differences
		time.Sleep(1 * time.Millisecond)
	}

	// Check that emergency cleanup occurred
	rollup.mux.RLock()
	cacheSize := len(rollup.cache)
	rollup.mux.RUnlock()

	if cacheSize > cfg.MaxKeys {
		t.Errorf("Cache size %d exceeds limit %d - emergency cleanup should have occurred", cacheSize, cfg.MaxKeys)
	}

	t.Logf("Cache size after emergency cleanup: %d (limit: %d)", cacheSize, cfg.MaxKeys)
}
