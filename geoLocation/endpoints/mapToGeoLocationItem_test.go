package endpoints

import (
	"testing"
)

// --- IpApiCom.mapToGeoLocationItem ---

func TestIpApiCom_mapToGeoLocationItem_EmptySlice_ReturnsEmptySlice(t *testing.T) {
	api := IpApiCom{}
	items := []IpApiComItem{}
	result, err := api.mapToGeoLocationItem(&items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %d items", len(result))
	}
}

func TestIpApiCom_mapToGeoLocationItem_MapsAllFields(t *testing.T) {
	api := IpApiCom{}
	items := []IpApiComItem{
		{
			Query:         "185.220.101.47",
			Status:        "success",
			Continent:     "Europe",
			ContinentCode: "EU",
			Country:       "Germany",
			CountryCode:   "DE",
			Region:        "BE",
			RegionName:    "Berlin",
			City:          "Berlin",
			Zip:           "10115",
			Lat:           52.5200,
			Lon:           13.4050,
		},
	}

	result, err := api.mapToGeoLocationItem(&items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}

	r := result[0]
	if r.Ip != "185.220.101.47" {
		t.Errorf("Ip: expected '185.220.101.47', got '%s'", r.Ip)
	}
	if r.Status != "success" {
		t.Errorf("Status: expected 'success', got '%s'", r.Status)
	}
	if r.Continent != "Europe" {
		t.Errorf("Continent: expected 'Europe', got '%s'", r.Continent)
	}
	if r.ContinentCode != "EU" {
		t.Errorf("ContinentCode: expected 'EU', got '%s'", r.ContinentCode)
	}
	if r.Country != "Germany" {
		t.Errorf("Country: expected 'Germany', got '%s'", r.Country)
	}
	if r.CountryCode != "DE" {
		t.Errorf("CountryCode: expected 'DE', got '%s'", r.CountryCode)
	}
	if r.Region != "BE" {
		t.Errorf("Region: expected 'BE', got '%s'", r.Region)
	}
	if r.RegionName != "Berlin" {
		t.Errorf("RegionName: expected 'Berlin', got '%s'", r.RegionName)
	}
	if r.City != "Berlin" {
		t.Errorf("City: expected 'Berlin', got '%s'", r.City)
	}
	if r.Zip != "10115" {
		t.Errorf("Zip: expected '10115', got '%s'", r.Zip)
	}
	if r.Latitude != 52.5200 {
		t.Errorf("Latitude: expected 52.52, got %v", r.Latitude)
	}
	if r.Longitude != 13.4050 {
		t.Errorf("Longitude: expected 13.405, got %v", r.Longitude)
	}
}

func TestIpApiCom_mapToGeoLocationItem_MultipleItems_MapsAll(t *testing.T) {
	api := IpApiCom{}
	items := []IpApiComItem{
		{Query: "1.1.1.1", Status: "success"},
		{Query: "2.2.2.2", Status: "success"},
	}
	result, err := api.mapToGeoLocationItem(&items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 items, got %d", len(result))
	}
	if result[0].Ip != "1.1.1.1" || result[1].Ip != "2.2.2.2" {
		t.Errorf("unexpected IPs: %s, %s", result[0].Ip, result[1].Ip)
	}
}

// --- ReallyFreeGeoIpOrg.mapToGeoLocationItem ---

func TestReallyFreeGeoIpOrg_mapToGeoLocationItem_MapsAllFields(t *testing.T) {
	api := ReallyFreeGeoIpOrg{}
	item := ReallyFreeGeoIpOrgItem{
		Ip:          "185.220.101.47",
		CountryCode: "DE",
		CountryName: "Germany",
		RegionCode:  "BE",
		RegionName:  "Berlin",
		City:        "Berlin",
		ZipCode:     "10115",
		Latitude:    52.52,
		Longitude:   13.405,
	}

	result, err := api.mapToGeoLocationItem(&item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Ip != "185.220.101.47" {
		t.Errorf("Ip: expected '185.220.101.47', got '%s'", result.Ip)
	}
	if result.Status != "success" {
		t.Errorf("Status: expected 'success', got '%s'", result.Status)
	}
	if result.Continent != "not available" {
		t.Errorf("Continent: expected 'not available', got '%s'", result.Continent)
	}
	if result.Country != "Germany" {
		t.Errorf("Country: expected 'Germany', got '%s'", result.Country)
	}
	if result.CountryCode != "DE" {
		t.Errorf("CountryCode: expected 'DE', got '%s'", result.CountryCode)
	}
	if result.RegionName != "Berlin" {
		t.Errorf("RegionName: expected 'Berlin', got '%s'", result.RegionName)
	}
	if result.City != "Berlin" {
		t.Errorf("City: expected 'Berlin', got '%s'", result.City)
	}
	if result.Zip != "10115" {
		t.Errorf("Zip: expected '10115', got '%s'", result.Zip)
	}
}

// --- IpapiCo.mapToGeoLocationItem ---

func TestIpapiCo_mapToGeoLocationItem_MapsAllFields(t *testing.T) {
	api := IpapiCo{}
	item := IpapiCoItem{
		Ip:            "185.220.101.47",
		CountryName:   "Germany",
		CountryCode:   "DE",
		ContinentCode: "EU",
		RegionCode:    "BE",
		Region:        "Berlin",
		City:          "Berlin",
		Postal:        "10115",
		Latitude:      52.52,
		Longitude:     13.405,
	}

	result, err := api.mapToGeoLocationItem(&item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Ip != "185.220.101.47" {
		t.Errorf("Ip: expected '185.220.101.47', got '%s'", result.Ip)
	}
	if result.Status != "success" {
		t.Errorf("Status: expected 'success', got '%s'", result.Status)
	}
	// IpapiCo hardcodes Continent as "not available"
	if result.Continent != "not available" {
		t.Errorf("Continent: expected 'not available', got '%s'", result.Continent)
	}
	if result.ContinentCode != "EU" {
		t.Errorf("ContinentCode: expected 'EU', got '%s'", result.ContinentCode)
	}
	if result.Country != "Germany" {
		t.Errorf("Country: expected 'Germany', got '%s'", result.Country)
	}
	if result.CountryCode != "DE" {
		t.Errorf("CountryCode: expected 'DE', got '%s'", result.CountryCode)
	}
	if result.Region != "BE" {
		t.Errorf("Region: expected 'BE', got '%s'", result.Region)
	}
	if result.RegionName != "Berlin" {
		t.Errorf("RegionName: expected 'Berlin', got '%s'", result.RegionName)
	}
}

// --- IpWhoIsIo.mapToGeoLocationItem ---

func TestIpWhoIsIo_mapToGeoLocationItem_MapsAllFields(t *testing.T) {
	api := IpWhoIsIo{}
	item := IpWhoIsIoItem{
		Ip:            "185.220.101.47",
		Success:       true,
		Continent:     "Europe",
		ContinentCode: "EU",
		Country:       "Germany",
		CountryCode:   "DE",
		Region:        "Berlin",
		City:          "Berlin",
		Latitude:      52.52,
		Longitude:     13.405,
	}

	result, err := api.mapToGeoLocationItem(&item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Ip != "185.220.101.47" {
		t.Errorf("Ip: expected '185.220.101.47', got '%s'", result.Ip)
	}
	// Status is mapped from bool via strconv.FormatBool
	if result.Status != "true" {
		t.Errorf("Status: expected 'true', got '%s'", result.Status)
	}
	// IpWhoIsIo hardcodes Region as "not available" and uses item.Region as RegionName
	if result.Region != "not available" {
		t.Errorf("Region: expected 'not available' (hardcoded), got '%s'", result.Region)
	}
	if result.RegionName != "Berlin" {
		t.Errorf("RegionName: expected 'Berlin', got '%s'", result.RegionName)
	}
	// IpWhoIsIo hardcodes Zip as "-1"
	if result.Zip != "-1" {
		t.Errorf("Zip: expected '-1' (hardcoded), got '%s'", result.Zip)
	}
	if result.Country != "Germany" {
		t.Errorf("Country: expected 'Germany', got '%s'", result.Country)
	}
}

func TestIpWhoIsIo_mapToGeoLocationItem_SuccessFalse_StatusIsFalse(t *testing.T) {
	api := IpWhoIsIo{}
	item := IpWhoIsIoItem{Success: false}
	result, _ := api.mapToGeoLocationItem(&item)
	if result.Status != "false" {
		t.Errorf("expected Status='false', got '%s'", result.Status)
	}
}

// --- GeoPluginCom.mapToGeoLocationItem ---

func TestGeoPluginCom_mapToGeoLocationItem_MapsAllFields(t *testing.T) {
	api := GeoPluginCom{}
	item := GeoPluginComItem{
		GeopluginRequest:       "185.220.101.47",
		GeopluginLatitude:      "52.52",
		GeopluginLongitude:     "13.405",
		GeopluginContinentCode: "EU",
		GeopluginContinentName: "Europe",
		GeopluginCountryName:   "Germany",
		GeopluginCountryCode:   "DE",
		GeopluginRegionCode:    "BE",
		GeopluginRegionName:    "Berlin",
		GeopluginCity:          "Berlin",
	}

	result, err := api.mapToGeoLocationItem(&item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Ip != "185.220.101.47" {
		t.Errorf("Ip: expected '185.220.101.47', got '%s'", result.Ip)
	}
	if result.Status != "success" {
		t.Errorf("Status: expected 'success', got '%s'", result.Status)
	}
	if result.Latitude != 52.52 {
		t.Errorf("Latitude: expected 52.52, got %v", result.Latitude)
	}
	if result.Longitude != 13.405 {
		t.Errorf("Longitude: expected 13.405, got %v", result.Longitude)
	}
	if result.Country != "Germany" {
		t.Errorf("Country: expected 'Germany', got '%s'", result.Country)
	}
	if result.CountryCode != "DE" {
		t.Errorf("CountryCode: expected 'DE', got '%s'", result.CountryCode)
	}
	// NOTE: GeoPluginCom swaps Continent/ContinentCode fields in the mapping:
	// Continent     ← GeopluginContinentCode ("EU")
	// ContinentCode ← GeopluginContinentName ("Europe")
	if result.Continent != "EU" {
		t.Errorf("Continent: expected 'EU' (from ContinentCode field), got '%s'", result.Continent)
	}
	if result.ContinentCode != "Europe" {
		t.Errorf("ContinentCode: expected 'Europe' (from ContinentName field), got '%s'", result.ContinentCode)
	}
}

func TestGeoPluginCom_mapToGeoLocationItem_InvalidLatLon_DefaultsToZero(t *testing.T) {
	api := GeoPluginCom{}
	item := GeoPluginComItem{
		GeopluginRequest:   "1.2.3.4",
		GeopluginLatitude:  "not-a-number",
		GeopluginLongitude: "also-not-a-number",
	}

	result, err := api.mapToGeoLocationItem(&item)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// strconv.ParseFloat returns 0 on error
	if result.Latitude != 0 {
		t.Errorf("expected Latitude=0 for invalid input, got %v", result.Latitude)
	}
	if result.Longitude != 0 {
		t.Errorf("expected Longitude=0 for invalid input, got %v", result.Longitude)
	}
}
