package cache

import (
	"tarpit-analyzer/database"
	"tarpit-analyzer/database/schemas"
	geolocationStructs "tarpit-analyzer/geoLocation/structs"
	"testing"
	"time"
)

// setupTestDb initialises the package-level db with an in-memory SQLite instance
// and returns a cleanup function that resets it.
func setupTestDb(t *testing.T) {
	t.Helper()
	memDb, err := database.CreateDbCacheInMemory()
	if err != nil {
		t.Fatalf("could not create in-memory DB: %v", err)
	}
	db = memDb
}

// insertLocation adds a GeoLocationItem to the cache DB and returns the inserted record.
func insertLocation(t *testing.T, item geolocationStructs.GeoLocationItem) schemas.Location {
	t.Helper()
	loc := db.MapToLocation(item)
	_, err := db.AddOrUpdateLocation(loc)
	if err != nil {
		t.Fatalf("could not insert location: %v", err)
	}
	inserted, _ := db.GetLocation(item.Ip)
	return inserted
}

// --- getSavedLocation: NoHit ---

func TestGetSavedLocation_UnknownIp_ReturnsNoHit(t *testing.T) {
	setupTestDb(t)

	_, result := getSavedLocation("99.99.99.99")
	if result != NoHit {
		t.Errorf("expected NoHit, got %v", result)
	}
}

// --- getSavedLocation: Ok (fresh record) ---

func TestGetSavedLocation_FreshRecord_ReturnsOk(t *testing.T) {
	setupTestDb(t)

	insertLocation(t, geolocationStructs.GeoLocationItem{
		Ip:      "1.2.3.4",
		Country: "Germany",
		Status:  "success",
	})

	item, result := getSavedLocation("1.2.3.4")
	if result != Ok {
		t.Errorf("expected Ok, got %v", result)
	}
	if item.Ip != "1.2.3.4" {
		t.Errorf("expected Ip '1.2.3.4', got '%s'", item.Ip)
	}
	if item.Country != "Germany" {
		t.Errorf("expected Country 'Germany', got '%s'", item.Country)
	}
}

// --- getSavedLocation: RecordOutdated (> 96 hours) ---

func TestGetSavedLocation_OutdatedRecord_ReturnsRecordOutdated(t *testing.T) {
	setupTestDb(t)

	insertLocation(t, geolocationStructs.GeoLocationItem{
		Ip:     "5.6.7.8",
		Status: "success",
	})

	// Manually backdate updated_at beyond the 96-hour TTL (4 days + 1 hour = 97 hours ago)
	fiveDaysAgo := time.Now().Add(-97 * time.Hour)
	err := db.ExecRaw(
		"UPDATE locations SET updated_at = ? WHERE ip = ?",
		fiveDaysAgo,
		"5.6.7.8",
	)
	if err != nil {
		t.Fatalf("could not update updated_at: %v", err)
	}

	_, result := getSavedLocation("5.6.7.8")
	if result != RecordOutdated {
		t.Errorf("expected RecordOutdated, got %v", result)
	}
}

func TestGetSavedLocation_RecordJustUnderTTL_ReturnsOk(t *testing.T) {
	setupTestDb(t)

	insertLocation(t, geolocationStructs.GeoLocationItem{
		Ip:     "9.10.11.12",
		Status: "success",
	})

	// 95 hours ago → still within the 96-hour TTL
	recentEnough := time.Now().Add(-95 * time.Hour)
	err := db.ExecRaw(
		"UPDATE locations SET updated_at = ? WHERE ip = ?",
		recentEnough,
		"9.10.11.12",
	)
	if err != nil {
		t.Fatalf("could not update updated_at: %v", err)
	}

	_, result := getSavedLocation("9.10.11.12")
	if result != Ok {
		t.Errorf("expected Ok for record within TTL, got %v", result)
	}
}

// --- TTL boundary: slightly under 96 hours ---

func TestGetSavedLocation_SlightlyUnderTTL_ReturnsOk(t *testing.T) {
	setupTestDb(t)

	insertLocation(t, geolocationStructs.GeoLocationItem{
		Ip:     "20.21.22.23",
		Status: "success",
	})

	// 95h59m ago — slightly under the 96-hour TTL, should still be Ok.
	// We avoid exact 96h to prevent timing races in CI.
	slightlyUnder := time.Now().Add(-(96*time.Hour - 1*time.Minute))
	err := db.ExecRaw(
		"UPDATE locations SET updated_at = ? WHERE ip = ?",
		slightlyUnder,
		"20.21.22.23",
	)
	if err != nil {
		t.Fatalf("could not update updated_at: %v", err)
	}

	_, result := getSavedLocation("20.21.22.23")
	if result != Ok {
		t.Errorf("expected Ok for record at 95h59m (just under TTL), got %v", result)
	}
}

// --- GeoLocationItem fields are correctly mapped on Ok ---

func TestGetSavedLocation_Ok_MapsAllFields(t *testing.T) {
	setupTestDb(t)

	original := geolocationStructs.GeoLocationItem{
		Ip:            "100.200.100.200",
		Status:        "success",
		Country:       "France",
		CountryCode:   "FR",
		Continent:     "Europe",
		ContinentCode: "EU",
		Region:        "IDF",
		RegionName:    "Île-de-France",
		City:          "Paris",
		Zip:           "75001",
		Latitude:      48.8566,
		Longitude:     2.3522,
	}
	insertLocation(t, original)

	item, result := getSavedLocation("100.200.100.200")
	if result != Ok {
		t.Fatalf("expected Ok, got %v", result)
	}

	if item.Country != "France" {
		t.Errorf("Country: expected 'France', got '%s'", item.Country)
	}
	if item.CountryCode != "FR" {
		t.Errorf("CountryCode: expected 'FR', got '%s'", item.CountryCode)
	}
	if item.Latitude != 48.8566 {
		t.Errorf("Latitude: expected 48.8566, got %v", item.Latitude)
	}
	if item.Longitude != 2.3522 {
		t.Errorf("Longitude: expected 2.3522, got %v", item.Longitude)
	}
}
