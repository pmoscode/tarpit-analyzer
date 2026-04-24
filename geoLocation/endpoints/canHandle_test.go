package endpoints

import (
	"testing"
	"time"
)

// CanHandle() depends on lastExecutionFinished (unexported field).
// A zero-value struct means "never executed" → should always be ready.
// A just-executed struct means "executed now" → should not be ready.

// --- IpApiCom (cooldown: > 1 minute) ---

func TestIpApiCom_CanHandle_FreshStruct_ReturnsTrue(t *testing.T) {
	api := IpApiCom{}
	if !api.CanHandle() {
		t.Error("expected CanHandle=true for fresh IpApiCom (never executed)")
	}
}

func TestIpApiCom_CanHandle_JustExecuted_ReturnsFalse(t *testing.T) {
	api := IpApiCom{lastExecutionFinished: time.Now()}
	if api.CanHandle() {
		t.Error("expected CanHandle=false when IpApiCom was just executed")
	}
}

func TestIpApiCom_CanHandle_ExecutedOverOneMinuteAgo_ReturnsTrue(t *testing.T) {
	api := IpApiCom{lastExecutionFinished: time.Now().Add(-2 * time.Minute)}
	if !api.CanHandle() {
		t.Error("expected CanHandle=true when last execution was >1 minute ago")
	}
}

// --- ReallyFreeGeoIpOrg (cooldown: > 1 minute) ---

func TestReallyFreeGeoIpOrg_CanHandle_FreshStruct_ReturnsTrue(t *testing.T) {
	api := ReallyFreeGeoIpOrg{}
	if !api.CanHandle() {
		t.Error("expected CanHandle=true for fresh ReallyFreeGeoIpOrg")
	}
}

func TestReallyFreeGeoIpOrg_CanHandle_JustExecuted_ReturnsFalse(t *testing.T) {
	api := ReallyFreeGeoIpOrg{lastExecutionFinished: time.Now()}
	if api.CanHandle() {
		t.Error("expected CanHandle=false when ReallyFreeGeoIpOrg was just executed")
	}
}

func TestReallyFreeGeoIpOrg_CanHandle_ExecutedOverOneMinuteAgo_ReturnsTrue(t *testing.T) {
	api := ReallyFreeGeoIpOrg{lastExecutionFinished: time.Now().Add(-2 * time.Minute)}
	if !api.CanHandle() {
		t.Error("expected CanHandle=true when last execution was >1 minute ago")
	}
}

// --- IpapiCo (cooldown: > 24 hours) ---

func TestIpapiCo_CanHandle_FreshStruct_ReturnsTrue(t *testing.T) {
	api := IpapiCo{}
	if !api.CanHandle() {
		t.Error("expected CanHandle=true for fresh IpapiCo")
	}
}

func TestIpapiCo_CanHandle_JustExecuted_ReturnsFalse(t *testing.T) {
	api := IpapiCo{lastExecutionFinished: time.Now()}
	if api.CanHandle() {
		t.Error("expected CanHandle=false when IpapiCo was just executed")
	}
}

func TestIpapiCo_CanHandle_ExecutedOver24HoursAgo_ReturnsTrue(t *testing.T) {
	api := IpapiCo{lastExecutionFinished: time.Now().Add(-25 * time.Hour)}
	if !api.CanHandle() {
		t.Error("expected CanHandle=true when last execution was >24 hours ago")
	}
}

func TestIpapiCo_CanHandle_ExecutedUnder24HoursAgo_ReturnsFalse(t *testing.T) {
	api := IpapiCo{lastExecutionFinished: time.Now().Add(-23 * time.Hour)}
	if api.CanHandle() {
		t.Error("expected CanHandle=false when last execution was <24 hours ago")
	}
}

// --- GeoPluginCom (cooldown: > 60 minutes) ---

func TestGeoPluginCom_CanHandle_FreshStruct_ReturnsTrue(t *testing.T) {
	api := GeoPluginCom{}
	if !api.CanHandle() {
		t.Error("expected CanHandle=true for fresh GeoPluginCom")
	}
}

func TestGeoPluginCom_CanHandle_JustExecuted_ReturnsFalse(t *testing.T) {
	api := GeoPluginCom{lastExecutionFinished: time.Now()}
	if api.CanHandle() {
		t.Error("expected CanHandle=false when GeoPluginCom was just executed")
	}
}

func TestGeoPluginCom_CanHandle_ExecutedOver60MinutesAgo_ReturnsTrue(t *testing.T) {
	api := GeoPluginCom{lastExecutionFinished: time.Now().Add(-61 * time.Minute)}
	if !api.CanHandle() {
		t.Error("expected CanHandle=true when last execution was >60 minutes ago")
	}
}

func TestGeoPluginCom_CanHandle_ExecutedUnder60MinutesAgo_ReturnsFalse(t *testing.T) {
	api := GeoPluginCom{lastExecutionFinished: time.Now().Add(-30 * time.Minute)}
	if api.CanHandle() {
		t.Error("expected CanHandle=false when last execution was <60 minutes ago")
	}
}

// --- IpWhoIsIo (cooldown: > 30 days) ---

func TestIpWhoIsIo_CanHandle_FreshStruct_ReturnsTrue(t *testing.T) {
	api := IpWhoIsIo{}
	if !api.CanHandle() {
		t.Error("expected CanHandle=true for fresh IpWhoIsIo")
	}
}

func TestIpWhoIsIo_CanHandle_JustExecuted_ReturnsFalse(t *testing.T) {
	api := IpWhoIsIo{lastExecutionFinished: time.Now()}
	if api.CanHandle() {
		t.Error("expected CanHandle=false when IpWhoIsIo was just executed")
	}
}

func TestIpWhoIsIo_CanHandle_ExecutedOver30DaysAgo_ReturnsTrue(t *testing.T) {
	api := IpWhoIsIo{lastExecutionFinished: time.Now().Add(-31 * 24 * time.Hour)}
	if !api.CanHandle() {
		t.Error("expected CanHandle=true when last execution was >30 days ago")
	}
}

func TestIpWhoIsIo_CanHandle_ExecutedUnder30DaysAgo_ReturnsFalse(t *testing.T) {
	api := IpWhoIsIo{lastExecutionFinished: time.Now().Add(-29 * 24 * time.Hour)}
	if api.CanHandle() {
		t.Error("expected CanHandle=false when last execution was <30 days ago")
	}
}
