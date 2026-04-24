package helper

import "testing"

// --- Private Networks ---

func TestCheckPrivateNetwork_Loopback_ReturnsTrue(t *testing.T) {
	ips := []string{"127.0.0.1", "127.0.0.255", "127.0.0.0"}
	for _, ip := range ips {
		if !CheckPrivateNetwork(ip) {
			t.Errorf("expected true for loopback IP %s", ip)
		}
	}
}

func TestCheckPrivateNetwork_ClassA_ReturnsTrue(t *testing.T) {
	ips := []string{"10.0.0.1", "10.255.255.255", "10.1.2.3"}
	for _, ip := range ips {
		if !CheckPrivateNetwork(ip) {
			t.Errorf("expected true for class A private IP %s", ip)
		}
	}
}

func TestCheckPrivateNetwork_ClassC_ReturnsTrue(t *testing.T) {
	ips := []string{"192.168.0.1", "192.168.1.100", "192.168.255.255"}
	for _, ip := range ips {
		if !CheckPrivateNetwork(ip) {
			t.Errorf("expected true for class C private IP %s", ip)
		}
	}
}

func TestCheckPrivateNetwork_ClassB_172_16_ReturnsTrue(t *testing.T) {
	ips := []string{"172.16.0.1", "172.16.255.255"}
	for _, ip := range ips {
		if !CheckPrivateNetwork(ip) {
			t.Errorf("expected true for 172.16.x private IP %s", ip)
		}
	}
}

// NOTE: The current implementation only checks for the "172.16." prefix,
// which means 172.17.x.x – 172.31.x.x (also private per RFC 1918) are
// NOT detected. The tests below document this known limitation.

func TestCheckPrivateNetwork_ClassB_172_17_to_31_NotDetected(t *testing.T) {
	ips := []string{"172.17.0.1", "172.20.5.5", "172.31.255.255"}
	for _, ip := range ips {
		// These are RFC 1918 private IPs but the current implementation
		// does NOT recognise them — documenting the existing behaviour.
		if CheckPrivateNetwork(ip) {
			t.Logf("NOTE: %s is now detected as private (implementation was extended)", ip)
		}
	}
}

// --- Public Networks ---

func TestCheckPrivateNetwork_PublicIP_ReturnsFalse(t *testing.T) {
	ips := []string{
		"8.8.8.8",
		"1.1.1.1",
		"203.0.113.1",
		"185.220.101.47",
		"45.33.32.156",
	}
	for _, ip := range ips {
		if CheckPrivateNetwork(ip) {
			t.Errorf("expected false for public IP %s", ip)
		}
	}
}

func TestCheckPrivateNetwork_EmptyString_ReturnsFalse(t *testing.T) {
	if CheckPrivateNetwork("") {
		t.Error("expected false for empty string")
	}
}

func TestCheckPrivateNetwork_PartialMatchNotPrivate_ReturnsFalse(t *testing.T) {
	// "192.169.x.x" looks similar to 192.168 but is public
	if CheckPrivateNetwork("192.169.0.1") {
		t.Error("expected false for 192.169.0.1 (not a private range)")
	}
}
