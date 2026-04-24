package modules

import (
	"testing"
	"time"
)

// helpers

func sshtarpit() SshTarpit {
	return SshTarpit{debug: false}
}

func newTempMap() map[int64]SshTarpitItem {
	return make(map[int64]SshTarpitItem)
}

// connectedLine builds a valid "connected" log line.
// format: <date> <time> INFO Client from ('<ip>', <port>) connected
func connectedLine(date, ip string, port int) string {
	return date + " INFO Client from ('" + ip + "', " + itoa(port) + ") connected"
}

// disconnectedLine builds a valid "disconnected" log line.
func disconnectedLine(date, ip string, port int) string {
	return date + " INFO Client from ('" + ip + "', " + itoa(port) + ") disconnected"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	return digits
}

// --- processLine: irrelevant lines ---

func TestSshTarpit_processLine_IrrelevantLine_ReturnsFalse(t *testing.T) {
	m := newTempMap()
	item, err := sshtarpit().processLine(m, "2023-06-15 10:00:00 INFO Starting server")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Success {
		t.Error("expected Success=false for irrelevant line")
	}
}

func TestSshTarpit_processLine_EmptyLine_ReturnsFalse(t *testing.T) {
	m := newTempMap()
	item, err := sshtarpit().processLine(m, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Success {
		t.Error("expected Success=false for empty line")
	}
}

// --- processLine: connected (first entry) ---

func TestSshTarpit_processLine_Connected_AddsToMap_ReturnsFalse(t *testing.T) {
	m := newTempMap()
	line := connectedLine("2023-06-15 10:00:00", "185.220.101.47", 22)

	item, err := sshtarpit().processLine(m, line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Success {
		t.Error("expected Success=false — item not complete yet")
	}
	if _, exists := m[22]; !exists {
		t.Error("expected port 22 to be added to tempMap")
	}
	if m[22].ip != "185.220.101.47" {
		t.Errorf("expected IP 185.220.101.47 in tempMap, got %s", m[22].ip)
	}
}

// --- processLine: connected then disconnected ---

func TestSshTarpit_processLine_ConnectThenDisconnect_ReturnsCompleteItem(t *testing.T) {
	m := newTempMap()
	r := sshtarpit()

	connectLine := connectedLine("2023-06-15 10:00:00", "185.220.101.47", 22)
	disconnectLine := disconnectedLine("2023-06-15 10:02:00", "185.220.101.47", 22)

	_, _ = r.processLine(m, connectLine)
	item, err := r.processLine(m, disconnectLine)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !item.Success {
		t.Fatal("expected Success=true after connect+disconnect")
	}
	if item.Ip != "185.220.101.47" {
		t.Errorf("expected IP 185.220.101.47, got %s", item.Ip)
	}

	expectedBegin := time.Date(2023, 6, 15, 10, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2023, 6, 15, 10, 2, 0, 0, time.UTC)

	if !item.Begin.Equal(expectedBegin) {
		t.Errorf("expected Begin %v, got %v", expectedBegin, item.Begin)
	}
	if !item.End.Equal(expectedEnd) {
		t.Errorf("expected End %v, got %v", expectedEnd, item.End)
	}

	expectedDuration := float32(120) // 2 minutes in seconds
	if item.Duration != expectedDuration {
		t.Errorf("expected Duration %v, got %v", expectedDuration, item.Duration)
	}
}

func TestSshTarpit_processLine_ConnectThenDisconnect_PortRemovedFromMap(t *testing.T) {
	m := newTempMap()
	r := sshtarpit()

	_, _ = r.processLine(m, connectedLine("2023-06-15 10:00:00", "185.220.101.47", 22))
	_, _ = r.processLine(m, disconnectedLine("2023-06-15 10:02:00", "185.220.101.47", 22))

	if _, exists := m[22]; exists {
		t.Error("expected port 22 to be removed from tempMap after disconnect")
	}
}

// --- processLine: private IP filtering ---

func TestSshTarpit_processLine_PrivateIP_Connected_ReturnsFalse(t *testing.T) {
	m := newTempMap()
	line := connectedLine("2023-06-15 10:00:00", "10.0.0.5", 22)

	item, err := sshtarpit().processLine(m, line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Success {
		t.Error("expected Success=false for private IP")
	}
	if _, exists := m[22]; exists {
		t.Error("private IP should not be added to tempMap")
	}
}

func TestSshTarpit_processLine_PrivateIP_192168_ReturnsFalse(t *testing.T) {
	m := newTempMap()
	line := connectedLine("2023-06-15 10:00:00", "192.168.1.50", 22)
	item, err := sshtarpit().processLine(m, line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Success {
		t.Error("expected Success=false for private IP 192.168.x.x")
	}
}

// --- processLine: second connected on same port overwrites ---

func TestSshTarpit_processLine_ReconnectOnSamePort_OverwritesEntry(t *testing.T) {
	m := newTempMap()
	r := sshtarpit()

	_, _ = r.processLine(m, connectedLine("2023-06-15 10:00:00", "185.220.101.47", 22))
	_, _ = r.processLine(m, connectedLine("2023-06-15 10:01:00", "45.33.32.156", 22))

	if m[22].ip != "45.33.32.156" {
		t.Errorf("expected overwritten IP 45.33.32.156, got %s", m[22].ip)
	}

	expectedStart := time.Date(2023, 6, 15, 10, 1, 0, 0, time.UTC)
	if !m[22].start.Equal(expectedStart) {
		t.Errorf("expected overwritten start %v, got %v", expectedStart, m[22].start)
	}
}

// --- processLine: disconnected without prior connected ---

func TestSshTarpit_processLine_DisconnectWithoutConnect_AddsToMap(t *testing.T) {
	m := newTempMap()
	// No prior connected — falls into else branch, adds entry with disconnect time as start
	line := disconnectedLine("2023-06-15 10:05:00", "185.220.101.47", 22)
	item, err := sshtarpit().processLine(m, line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Success {
		t.Error("expected Success=false when no prior connect exists")
	}
	if _, exists := m[22]; !exists {
		t.Error("expected orphan disconnect entry to be added to tempMap")
	}
}

// --- processLine: multiple ports independently tracked ---

func TestSshTarpit_processLine_MultiplePortsTrackedIndependently(t *testing.T) {
	m := newTempMap()
	r := sshtarpit()

	_, _ = r.processLine(m, connectedLine("2023-06-15 10:00:00", "1.2.3.4", 100))
	_, _ = r.processLine(m, connectedLine("2023-06-15 10:01:00", "5.6.7.8", 200))

	item1, err := r.processLine(m, disconnectedLine("2023-06-15 10:03:00", "1.2.3.4", 100))
	if err != nil || !item1.Success {
		t.Errorf("expected Success=true for port 100: err=%v success=%v", err, item1.Success)
	}

	item2, err := r.processLine(m, disconnectedLine("2023-06-15 10:04:00", "5.6.7.8", 200))
	if err != nil || !item2.Success {
		t.Errorf("expected Success=true for port 200: err=%v success=%v", err, item2.Success)
	}

	if item1.Ip != "1.2.3.4" {
		t.Errorf("wrong IP for port 100: %s", item1.Ip)
	}
	if item2.Ip != "5.6.7.8" {
		t.Errorf("wrong IP for port 200: %s", item2.Ip)
	}
}

// --- processLine: malformed lines ---

func TestSshTarpit_processLine_InvalidDate_ReturnsError(t *testing.T) {
	m := newTempMap()
	line := "BADDATE INFO Client from ('1.2.3.4', 22) connected"
	item, err := sshtarpit().processLine(m, line)
	if err == nil {
		t.Error("expected error for invalid date")
	}
	if item.Success {
		t.Error("expected Success=false")
	}
}

// --- mapToImportItem ---

func TestSshTarpit_mapToImportItem_MapsFieldsCorrectly(t *testing.T) {
	r := sshtarpit()
	start := time.Date(2023, 6, 15, 10, 0, 0, 0, time.UTC)
	end := time.Date(2023, 6, 15, 10, 5, 0, 0, time.UTC)

	value := SshTarpitItem{
		start:    start,
		end:      end,
		duration: 300,
		ip:       "185.220.101.47",
		port:     22,
	}

	item, err := r.mapToImportItem(value)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !item.Success {
		t.Error("expected Success=true")
	}
	if item.Ip != "185.220.101.47" {
		t.Errorf("expected IP 185.220.101.47, got %s", item.Ip)
	}
	if !item.Begin.Equal(start) {
		t.Errorf("expected Begin %v, got %v", start, item.Begin)
	}
	if !item.End.Equal(end) {
		t.Errorf("expected End %v, got %v", end, item.End)
	}
	if item.Duration != 300 {
		t.Errorf("expected Duration 300, got %v", item.Duration)
	}
}
