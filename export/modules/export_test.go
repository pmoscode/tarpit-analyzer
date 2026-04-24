package modules

import (
	"encoding/json"
	"strings"
	"tarpit-analyzer/database/schemas"
	"tarpit-analyzer/importData/structs"
	"testing"
	"time"
)

// --- helpers ---

func makeData(ip string, begin, end time.Time, duration float32) schemas.Data {
	return schemas.Data{
		ImportItem: structs.ImportItem{
			Ip:       ip,
			Begin:    begin,
			End:      end,
			Duration: duration,
		},
	}
}

var (
	testBegin = time.Date(2023, 6, 15, 10, 0, 0, 0, time.UTC)
	testEnd   = time.Date(2023, 6, 15, 10, 2, 0, 0, time.UTC)
)

// --- CSV.mapToCSV ---

func TestCSV_mapToCSV_ReturnsCorrectFields(t *testing.T) {
	c := &CSV{Separator: ","}
	data := makeData("185.220.101.47", testBegin, testEnd, 120)

	result := c.mapToCSV(data)

	if len(result) != 4 {
		t.Fatalf("expected 4 fields, got %d", len(result))
	}
	if result[0] != testBegin.Format(time.RFC3339) {
		t.Errorf("Begin: expected %s, got %s", testBegin.Format(time.RFC3339), result[0])
	}
	if result[1] != testEnd.Format(time.RFC3339) {
		t.Errorf("End: expected %s, got %s", testEnd.Format(time.RFC3339), result[1])
	}
	if result[2] != "185.220.101.47" {
		t.Errorf("Ip: expected '185.220.101.47', got '%s'", result[2])
	}
	if result[3] != "120" {
		t.Errorf("Duration: expected '120', got '%s'", result[3])
	}
}

func TestCSV_mapToCSV_FractionalDuration_FormattedCorrectly(t *testing.T) {
	c := &CSV{Separator: ","}
	data := makeData("1.2.3.4", testBegin, testEnd, 60.5)

	result := c.mapToCSV(data)
	if result[3] != "60.5" {
		t.Errorf("Duration: expected '60.5', got '%s'", result[3])
	}
}

func TestCSV_mapToCSV_JoinedWithSeparator(t *testing.T) {
	c := &CSV{Separator: ";"}
	data := makeData("1.2.3.4", testBegin, testEnd, 10)

	joined := strings.Join(c.mapToCSV(data), c.Separator)
	parts := strings.Split(joined, ";")
	if len(parts) != 4 {
		t.Errorf("expected 4 semicolon-separated fields, got %d", len(parts))
	}
}

func TestCSV_mapToCSV_TabSeparator(t *testing.T) {
	c := &CSV{Separator: "\t"}
	data := makeData("9.9.9.9", testBegin, testEnd, 30)

	joined := strings.Join(c.mapToCSV(data), c.Separator)
	if !strings.Contains(joined, "\t") {
		t.Error("expected tab-separated output")
	}
}

// --- JSON.mapToJsonItem ---

func TestJSON_mapToJsonItem_ReturnsCorrectFields(t *testing.T) {
	j := &JSON{}
	data := makeData("185.220.101.47", testBegin, testEnd, 120)

	result := j.mapToJsonItem(data)

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Ip != "185.220.101.47" {
		t.Errorf("Ip: expected '185.220.101.47', got '%s'", result.Ip)
	}
	if result.Begin != testBegin.Format(time.RFC3339) {
		t.Errorf("Begin: expected %s, got %s", testBegin.Format(time.RFC3339), result.Begin)
	}
	if result.End != testEnd.Format(time.RFC3339) {
		t.Errorf("End: expected %s, got %s", testEnd.Format(time.RFC3339), result.End)
	}
	if result.Duration != "120" {
		t.Errorf("Duration: expected '120', got '%s'", result.Duration)
	}
}

func TestJSON_mapToJsonItem_IsJsonSerializable(t *testing.T) {
	j := &JSON{}
	data := makeData("1.2.3.4", testBegin, testEnd, 45.5)

	item := j.mapToJsonItem(data)
	b, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("unexpected JSON marshal error: %v", err)
	}

	var decoded JsonItem
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("unexpected JSON unmarshal error: %v", err)
	}
	if decoded.Ip != "1.2.3.4" {
		t.Errorf("expected Ip '1.2.3.4', got '%s'", decoded.Ip)
	}
	if decoded.Duration != "45.5" {
		t.Errorf("expected Duration '45.5', got '%s'", decoded.Duration)
	}
}

func TestJSON_mapToJsonItem_JsonTagsUsed(t *testing.T) {
	j := &JSON{}
	data := makeData("5.6.7.8", testBegin, testEnd, 10)

	b, _ := json.Marshal(j.mapToJsonItem(data))
	s := string(b)

	for _, key := range []string{`"begin"`, `"end"`, `"ip"`, `"duration"`} {
		if !strings.Contains(s, key) {
			t.Errorf("expected JSON key %s in output: %s", key, s)
		}
	}
}

// --- KML.generateKMLContent ---

func TestKML_generateKMLContent_EmptyData_OnlyWrappers(t *testing.T) {
	k := &KML{CenterGeoLocationLongitude: "13.405", CenterGeoLocationLatitude: "52.52"}
	result := make([]string, 0)
	data := make([]KmlDbItem, 0)

	k.generateKMLContent(&result, &data)

	combined := strings.Join(result, "")
	if !strings.Contains(combined, `<?xml version="1.0"`) {
		t.Error("expected XML declaration in output")
	}
	if !strings.Contains(combined, `<kml xmlns=`) {
		t.Error("expected kml root element")
	}
	if !strings.Contains(combined, `</kml>`) {
		t.Error("expected closing kml tag")
	}
	if strings.Contains(combined, "<Placemark>") {
		t.Error("expected no Placemark elements for empty data")
	}
}

func TestKML_generateKMLContent_SingleItem_ContainsPlacemark(t *testing.T) {
	k := &KML{CenterGeoLocationLongitude: "13.405000", CenterGeoLocationLatitude: "52.520000"}
	result := make([]string, 0)
	data := []KmlDbItem{
		{Country: "Germany", Latitude: 51.165691, Longitude: 10.451526},
	}

	k.generateKMLContent(&result, &data)
	combined := strings.Join(result, "")

	if !strings.Contains(combined, "Germany") {
		t.Error("expected country name 'Germany' in KML output")
	}
	if !strings.Contains(combined, "<Placemark>") {
		t.Error("expected Placemark element")
	}
	if !strings.Contains(combined, "<LineString>") {
		t.Error("expected LineString element")
	}
	if !strings.Contains(combined, "10.451526") {
		t.Error("expected longitude value in KML output")
	}
	if !strings.Contains(combined, k.CenterGeoLocationLongitude) {
		t.Errorf("expected center longitude %s in KML output", k.CenterGeoLocationLongitude)
	}
}

func TestKML_generateKMLContent_MultipleItems_AllPresent(t *testing.T) {
	k := &KML{CenterGeoLocationLongitude: "0.0", CenterGeoLocationLatitude: "0.0"}
	result := make([]string, 0)
	data := []KmlDbItem{
		{Country: "France", Latitude: 46.2276, Longitude: 2.2137},
		{Country: "Japan", Latitude: 36.2048, Longitude: 138.2529},
	}

	k.generateKMLContent(&result, &data)
	combined := strings.Join(result, "")

	if !strings.Contains(combined, "France") {
		t.Error("expected 'France' in output")
	}
	if !strings.Contains(combined, "Japan") {
		t.Error("expected 'Japan' in output")
	}
}

func TestKML_generateKMLContent_StylePresent(t *testing.T) {
	k := &KML{}
	result := make([]string, 0)
	data := make([]KmlDbItem, 0)

	k.generateKMLContent(&result, &data)
	combined := strings.Join(result, "")

	if !strings.Contains(combined, "transBluePoly") {
		t.Error("expected style ID 'transBluePoly' in KML output")
	}
	if !strings.Contains(combined, "<LineStyle>") {
		t.Error("expected LineStyle element in KML output")
	}
}
