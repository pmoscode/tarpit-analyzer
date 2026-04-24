package geoLocation

import (
	"errors"
	"github.com/schollz/progressbar/v3"
	"tarpit-analyzer/geoLocation/structs"
	"testing"
)

// --- mock implementation of QueryGeoLocationAPI ---

type mockApi struct {
	name      string
	canHandle bool
	items     []structs.GeoLocationItem
	err       error
}

func (m *mockApi) QueryGeoLocationAPI(ips *[]string, bar *progressbar.ProgressBar) ([]structs.GeoLocationItem, error) {
	return m.items, m.err
}

func (m *mockApi) Name() string    { return m.name }
func (m *mockApi) CanHandle() bool { return m.canHandle }

// --- nextApi ---

func TestNextApi_AllCanHandle_ReturnsFirst(t *testing.T) {
	a := &mockApi{name: "first", canHandle: true}
	b := &mockApi{name: "second", canHandle: true}
	gl := &GeoLocation{apis: []QueryGeoLocationAPI{a, b}}

	result := gl.nextApi()
	if result == nil {
		t.Fatal("expected non-nil API")
	}
	if result.Name() != "first" {
		t.Errorf("expected 'first', got '%s'", result.Name())
	}
}

func TestNextApi_NoneCanHandle_ReturnsNil(t *testing.T) {
	a := &mockApi{name: "a", canHandle: false}
	b := &mockApi{name: "b", canHandle: false}
	gl := &GeoLocation{apis: []QueryGeoLocationAPI{a, b}}

	if gl.nextApi() != nil {
		t.Error("expected nil when no API can handle")
	}
}

func TestNextApi_EmptyApis_ReturnsNil(t *testing.T) {
	gl := &GeoLocation{apis: []QueryGeoLocationAPI{}}
	if gl.nextApi() != nil {
		t.Error("expected nil for empty API list")
	}
}

func TestNextApi_FirstCannotHandleSecondCan_ReturnsSecond(t *testing.T) {
	a := &mockApi{name: "first", canHandle: false}
	b := &mockApi{name: "second", canHandle: true}
	gl := &GeoLocation{apis: []QueryGeoLocationAPI{a, b}}

	result := gl.nextApi()
	if result == nil {
		t.Fatal("expected non-nil API")
	}
	if result.Name() != "second" {
		t.Errorf("expected 'second', got '%s'", result.Name())
	}
}

func TestNextApi_OnlyOneApi_CanHandle_ReturnsThatApi(t *testing.T) {
	a := &mockApi{name: "only", canHandle: true}
	gl := &GeoLocation{apis: []QueryGeoLocationAPI{a}}

	result := gl.nextApi()
	if result == nil || result.Name() != "only" {
		t.Error("expected the single available API to be returned")
	}
}

// --- ResolveLocations ---

func TestResolveLocations_EmptyIpList_ReturnsEmptySlice(t *testing.T) {
	gl := &GeoLocation{apis: []QueryGeoLocationAPI{
		&mockApi{name: "api", canHandle: true},
	}}

	result, err := gl.ResolveLocations([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*result) != 0 {
		t.Errorf("expected empty result, got %d items", len(*result))
	}
}

func TestResolveLocations_SingleIp_ApiResolvesIt(t *testing.T) {
	resolved := []structs.GeoLocationItem{
		{Ip: "1.2.3.4", Country: "Germany", Status: "success"},
	}
	gl := &GeoLocation{apis: []QueryGeoLocationAPI{
		&mockApi{name: "api", canHandle: true, items: resolved},
	}}

	result, err := gl.ResolveLocations([]string{"1.2.3.4"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(*result))
	}
	if (*result)[0].Ip != "1.2.3.4" {
		t.Errorf("expected IP '1.2.3.4', got '%s'", (*result)[0].Ip)
	}
	if (*result)[0].Country != "Germany" {
		t.Errorf("expected Country 'Germany', got '%s'", (*result)[0].Country)
	}
}

func TestResolveLocations_MultipleIps_AllResolved(t *testing.T) {
	resolved := []structs.GeoLocationItem{
		{Ip: "1.1.1.1", Status: "success"},
		{Ip: "2.2.2.2", Status: "success"},
		{Ip: "3.3.3.3", Status: "success"},
	}
	gl := &GeoLocation{apis: []QueryGeoLocationAPI{
		&mockApi{name: "api", canHandle: true, items: resolved},
	}}

	result, err := gl.ResolveLocations([]string{"1.1.1.1", "2.2.2.2", "3.3.3.3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*result) != 3 {
		t.Errorf("expected 3 items, got %d", len(*result))
	}
}

func TestResolveLocations_NoApiCanHandle_ReturnsEmptySlice(t *testing.T) {
	gl := &GeoLocation{apis: []QueryGeoLocationAPI{
		&mockApi{name: "api", canHandle: false},
	}}

	result, err := gl.ResolveLocations([]string{"1.2.3.4"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// nextApi returns nil immediately → loop breaks → empty result
	if len(*result) != 0 {
		t.Errorf("expected empty result when no API can handle, got %d items", len(*result))
	}
}

func TestResolveLocations_ApiReturnsError_FallsBackToNextApi(t *testing.T) {
	fallbackResolved := []structs.GeoLocationItem{
		{Ip: "1.2.3.4", Status: "success"},
	}
	// First API errors and becomes unavailable (canHandle=false after error in real impl)
	// Here we model the post-error state: first API returns error and CanHandle=false,
	// second API handles it successfully.
	firstApi := &mockApi{name: "failing", canHandle: false, err: errors.New("timeout")}
	secondApi := &mockApi{name: "fallback", canHandle: true, items: fallbackResolved}

	gl := &GeoLocation{apis: []QueryGeoLocationAPI{firstApi, secondApi}}

	result, err := gl.ResolveLocations([]string{"1.2.3.4"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*result) != 1 {
		t.Fatalf("expected 1 item from fallback API, got %d", len(*result))
	}
	if (*result)[0].Ip != "1.2.3.4" {
		t.Errorf("expected IP from fallback API, got '%s'", (*result)[0].Ip)
	}
}

func TestResolveLocations_NoApis_ReturnsEmptySlice(t *testing.T) {
	gl := &GeoLocation{apis: []QueryGeoLocationAPI{}}
	result, err := gl.ResolveLocations([]string{"1.2.3.4"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*result) != 0 {
		t.Errorf("expected empty result, got %d items", len(*result))
	}
}
