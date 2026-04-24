package modules

import (
	"testing"
	"time"
)

// helpers

func endlessh() Endlessh {
	return Endlessh{debug: false}
}

// validCloseLine builds a well-formed CLOSE log line.
// format: <RFC3339Nano> CLOSE host=<ip> port=<p> fd=<fd> time=<sec> bytes=<b>
func validCloseLine(ts, ip, duration string) string {
	return ts + " CLOSE host=" + ip + " port=22 fd=5 time=" + duration + " bytes=1234"
}

// --- processLine: non-CLOSE lines ---

func TestEndlessh_processLine_NonCloseLine_ReturnsFalse(t *testing.T) {
	item, err := endlessh().processLine("2023-06-15T10:00:00.000000000Z ACCEPT host=185.220.101.47 port=22 fd=5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Success {
		t.Error("expected Success=false for non-CLOSE line")
	}
}

func TestEndlessh_processLine_EmptyLine_ReturnsFalse(t *testing.T) {
	item, err := endlessh().processLine("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Success {
		t.Error("expected Success=false for empty line")
	}
}

// --- processLine: valid CLOSE line ---

func TestEndlessh_processLine_ValidCloseLine_ReturnsItem(t *testing.T) {
	// End: 2023-06-15T10:02:00Z, duration: 120s → Begin: 2023-06-15T10:00:00Z
	line := validCloseLine("2023-06-15T10:02:00.000000000Z", "185.220.101.47", "120")

	item, err := endlessh().processLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !item.Success {
		t.Fatal("expected Success=true")
	}
	if item.Ip != "185.220.101.47" {
		t.Errorf("expected IP 185.220.101.47, got %s", item.Ip)
	}

	expectedEnd := time.Date(2023, 6, 15, 10, 2, 0, 0, time.UTC)
	if !item.End.Equal(expectedEnd) {
		t.Errorf("expected End %v, got %v", expectedEnd, item.End)
	}

	expectedBegin := time.Date(2023, 6, 15, 10, 0, 0, 0, time.UTC)
	if !item.Begin.Equal(expectedBegin) {
		t.Errorf("expected Begin %v, got %v", expectedBegin, item.Begin)
	}

	if item.Duration != float32(120) {
		t.Errorf("expected Duration 120, got %v", item.Duration)
	}
}

func TestEndlessh_processLine_ValidCloseLine_FractionalDuration(t *testing.T) {
	line := validCloseLine("2023-06-15T10:02:00.500000000Z", "185.220.101.47", "60.5")

	item, err := endlessh().processLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !item.Success {
		t.Fatal("expected Success=true")
	}
	if item.Duration != float32(60.5) {
		t.Errorf("expected Duration 60.5, got %v", item.Duration)
	}
}

// --- processLine: private IP filtering ---

func TestEndlessh_processLine_PrivateIP_10_ReturnsFalse(t *testing.T) {
	line := validCloseLine("2023-06-15T10:02:00.000000000Z", "10.0.0.1", "60")
	item, err := endlessh().processLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Success {
		t.Error("expected Success=false for private IP 10.x.x.x")
	}
}

func TestEndlessh_processLine_PrivateIP_192168_ReturnsFalse(t *testing.T) {
	line := validCloseLine("2023-06-15T10:02:00.000000000Z", "192.168.1.100", "60")
	item, err := endlessh().processLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Success {
		t.Error("expected Success=false for private IP 192.168.x.x")
	}
}

func TestEndlessh_processLine_PrivateIP_Loopback_ReturnsFalse(t *testing.T) {
	line := validCloseLine("2023-06-15T10:02:00.000000000Z", "127.0.0.1", "60")
	item, err := endlessh().processLine(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Success {
		t.Error("expected Success=false for loopback IP")
	}
}

// --- processLine: malformed lines ---

func TestEndlessh_processLine_InvalidTimestamp_ReturnsError(t *testing.T) {
	line := "NOT-A-DATE CLOSE host=185.220.101.47 port=22 fd=5 time=60 bytes=1234"
	item, err := endlessh().processLine(line)
	if err == nil {
		t.Error("expected error for invalid timestamp")
	}
	if item.Success {
		t.Error("expected Success=false")
	}
}

func TestEndlessh_processLine_InvalidDuration_ReturnsError(t *testing.T) {
	line := validCloseLine("2023-06-15T10:02:00.000000000Z", "185.220.101.47", "notanumber")
	item, err := endlessh().processLine(line)
	if err == nil {
		t.Error("expected error for non-numeric duration")
	}
	if item.Success {
		t.Error("expected Success=false")
	}
}

// --- getValue ---

func TestEndlessh_getValue_ValidKeyValue_ReturnsValue(t *testing.T) {
	e := endlessh()
	val, err := e.getValue("host=185.220.101.47")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "185.220.101.47" {
		t.Errorf("expected '185.220.101.47', got '%s'", val)
	}
}

func TestEndlessh_getValue_EmptySource_ReturnsError(t *testing.T) {
	e := endlessh()
	_, err := e.getValue("")
	if err == nil {
		t.Error("expected error for empty source")
	}
}

func TestEndlessh_getValue_NoEqualsSign_ReturnsError(t *testing.T) {
	e := endlessh()
	_, err := e.getValue("hostonly")
	if err == nil {
		t.Error("expected error for source without '='")
	}
}

func TestEndlessh_getValue_MultipleEqualsSign_ReturnsFirstValue(t *testing.T) {
	e := endlessh()
	// split on "=" returns ["key", "val1", "val2"], index [1] = "val1"
	val, err := e.getValue("key=val1=val2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "val1" {
		t.Errorf("expected 'val1', got '%s'", val)
	}
}
