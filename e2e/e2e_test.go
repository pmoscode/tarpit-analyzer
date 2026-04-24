// Package e2e_test contains end-to-end tests that exercise the full pipeline:
// Import → Analyze → Export.
//
// Each test uses a dedicated temporary SQLite database via SetTestDatabasePath
// so no files are written to the project directory.
package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"tarpit-analyzer/analyze"
	"tarpit-analyzer/cli"
	"tarpit-analyzer/database"
	"tarpit-analyzer/database/schemas"
	"tarpit-analyzer/export"
	geolocationStructs "tarpit-analyzer/geoLocation/structs"
	"tarpit-analyzer/importData"
)

// setupTestDB redirects all database factory functions to a temporary file and
// returns the temp directory. The override is reset automatically via t.Cleanup.
func setupTestDB(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	database.SetTestDatabasePath(filepath.Join(dir, "data"))
	t.Cleanup(func() {
		database.SetTestDatabasePath("")
	})
	return dir
}

func defaultCtx() *cli.Context {
	return &cli.Context{
		Debug:     false,
		StartDate: "unset",
		EndDate:   "unset",
	}
}

func ctxWithTarget(target string) *cli.Context {
	return &cli.Context{
		Debug:     false,
		StartDate: "unset",
		EndDate:   "unset",
		Target:    target,
	}
}

// seedCache inserts dummy geo-location entries for the given IPs into the test
// database so that DoAnalyze can look up countries without making external
// network requests.
func seedCache(t *testing.T, ips []string) {
	t.Helper()
	dbCache, err := database.CreateDbCache(false)
	if err != nil {
		t.Fatalf("seedCache: CreateDbCache failed: %v", err)
	}
	for _, ip := range ips {
		loc := schemas.Location{
			GeoLocationItem: geolocationStructs.GeoLocationItem{
				Ip:            ip,
				Status:        "success",
				Country:       "TestLand",
				CountryCode:   "TL",
				Continent:     "TestContinent",
				ContinentCode: "TC",
			},
		}
		if _, err := dbCache.AddOrUpdateLocation(loc); err != nil {
			t.Fatalf("seedCache: AddOrUpdateLocation(%s) failed: %v", ip, err)
		}
	}
}

// TestE2E_Import_Analyze_ExportCSV runs the full pipeline end-to-end:
//  1. Import an endlessh log fixture (IP resolving skipped).
//  2. Seed the cache with dummy geo-location data for the test IPs.
//  3. Run Analyze and verify the statistics report is written.
//  4. Export to CSV and verify row count and field structure.
func TestE2E_Import_Analyze_ExportCSV(t *testing.T) {
	dir := setupTestDB(t)

	// --- 1. Import ---
	err := importData.DoImport(importData.Endlessh, "testdata/endlessh.log", true, defaultCtx())
	if err != nil {
		t.Fatalf("DoImport failed: %v", err)
	}

	// --- 2. Seed cache so DoAnalyze does not hit external geo-APIs ---
	testIPs := []string{"185.220.101.47", "45.33.32.156", "198.51.100.22"}
	seedCache(t, testIPs)

	// --- 3. Analyze ---
	analyzeOut := filepath.Join(dir, "analysis.txt")
	err = analyze.DoAnalyze(ctxWithTarget(analyzeOut))
	if err != nil {
		t.Fatalf("DoAnalyze failed: %v", err)
	}

	content, err := os.ReadFile(analyzeOut)
	if err != nil {
		t.Fatalf("could not read analysis output: %v", err)
	}
	analyzeText := string(content)

	if !strings.Contains(analyzeText, "Tarpit Analyzer Statistics") {
		t.Errorf("analysis output missing expected header; got:\n%s", analyzeText)
	}
	if !strings.Contains(analyzeText, "Selected date range: All data") {
		t.Errorf("analysis output missing date range line; got:\n%s", analyzeText)
	}

	// --- 4. Export CSV ---
	csvOut := filepath.Join(dir, "export.csv")
	err = export.DoExport(export.CSV, export.Parameters{Separator: ","}, ctxWithTarget(csvOut))
	if err != nil {
		t.Fatalf("DoExport(CSV) failed: %v", err)
	}

	csvBytes, err := os.ReadFile(csvOut)
	if err != nil {
		t.Fatalf("could not read CSV output: %v", err)
	}
	csvText := strings.TrimSpace(string(csvBytes))
	csvLines := strings.Split(csvText, "\n")

	// The fixture contains 3 valid CLOSE lines with public IPs; 1 ACCEPT line
	// and 1 private-IP CLOSE line are both filtered out during import.
	const wantRows = 3
	if len(csvLines) != wantRows {
		t.Errorf("expected %d CSV rows, got %d:\n%s", wantRows, len(csvLines), string(csvBytes))
	}

	// Each CSV line must have exactly 4 fields: begin, end, ip, duration.
	for i, line := range csvLines {
		fields := strings.Split(line, ",")
		if len(fields) != 4 {
			t.Errorf("CSV row %d: expected 4 fields, got %d: %q", i+1, len(fields), line)
		}
	}

	// The imported IPs must appear in the CSV output.
	for _, ip := range testIPs {
		if !strings.Contains(csvText, ip) {
			t.Errorf("CSV output missing expected IP %q", ip)
		}
	}
}

// TestE2E_Import_ExportJSON imports the same fixture and exports to JSON,
// verifying that the output is a JSON array containing the expected IPs.
func TestE2E_Import_ExportJSON(t *testing.T) {
	dir := setupTestDB(t)

	err := importData.DoImport(importData.Endlessh, "testdata/endlessh.log", true, defaultCtx())
	if err != nil {
		t.Fatalf("DoImport failed: %v", err)
	}

	jsonOut := filepath.Join(dir, "export.json")
	err = export.DoExport(export.JSON, export.Parameters{}, ctxWithTarget(jsonOut))
	if err != nil {
		t.Fatalf("DoExport(JSON) failed: %v", err)
	}

	jsonBytes, err := os.ReadFile(jsonOut)
	if err != nil {
		t.Fatalf("could not read JSON output: %v", err)
	}
	jsonText := string(jsonBytes)

	// Output must start and end with a JSON array.
	trimmed := strings.TrimSpace(jsonText)
	if !strings.HasPrefix(trimmed, "[") || !strings.HasSuffix(trimmed, "]") {
		t.Errorf("JSON output is not an array:\n%s", jsonText)
	}

	// All three public IPs from the fixture must appear in the JSON.
	for _, ip := range []string{"185.220.101.47", "45.33.32.156", "198.51.100.22"} {
		if !strings.Contains(jsonText, ip) {
			t.Errorf("JSON output missing expected IP %q", ip)
		}
	}
}
