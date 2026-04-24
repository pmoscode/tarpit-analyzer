package helper

import (
	"testing"
	"time"
)

// --- GetDate ---

func TestGetDate_Unset_ReturnsNil(t *testing.T) {
	result := GetDate("unset")
	if result != nil {
		t.Errorf("expected nil for 'unset', got %v", result)
	}
}

func TestGetDate_ValidDate_ReturnsCorrectTime(t *testing.T) {
	result := GetDate("2023-06-15")
	if result == nil {
		t.Fatal("expected non-nil result for valid date")
	}

	expected := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, *result)
	}
}

func TestGetDate_InvalidDate_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid date, but did not panic")
		}
	}()

	GetDate("not-a-date")
}

func TestGetDate_WrongFormat_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for wrong date format, but did not panic")
		}
	}()

	GetDate("15.06.2023")
}

// --- IsBefore ---

func TestIsBefore_NilThen_ReturnsTrue(t *testing.T) {
	now := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if !IsBefore(now, nil) {
		t.Error("expected true when then is nil")
	}
}

func TestIsBefore_NowClearlyBeforeThen_ReturnsTrue(t *testing.T) {
	now := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	then := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	if !IsBefore(now, &then) {
		t.Errorf("expected true: now=%v is before then=%v", now, then)
	}
}

func TestIsBefore_NowEqualThen_ReturnsFalse(t *testing.T) {
	now := time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC)
	then := time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC)
	// now+1 day = 2023-01-06, which is NOT before 2023-01-05 → false
	if IsBefore(now, &then) {
		t.Errorf("expected false: now=%v equals then=%v", now, then)
	}
}

func TestIsBefore_NowAfterThen_ReturnsFalse(t *testing.T) {
	now := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	then := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if IsBefore(now, &then) {
		t.Errorf("expected false: now=%v is after then=%v", now, then)
	}
}

func TestIsBefore_NowOneDayBeforeThen_ReturnsFalse(t *testing.T) {
	now := time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC)
	then := time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC)
	// now+1 day = 2023-01-05, which is NOT before 2023-01-05 (strict) → false
	if IsBefore(now, &then) {
		t.Errorf("expected false: now+1=%v equals then=%v (strict Before)", now.AddDate(0, 0, 1), then)
	}
}

func TestIsBefore_NowTwoDaysBeforeThen_ReturnsTrue(t *testing.T) {
	now := time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)
	then := time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC)
	// now+1 day = 2023-01-04, which IS before 2023-01-05 → true
	if !IsBefore(now, &then) {
		t.Errorf("expected true: now+1=%v is before then=%v", now.AddDate(0, 0, 1), then)
	}
}

// --- IsAfter ---

func TestIsAfter_NilThen_ReturnsTrue(t *testing.T) {
	now := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if !IsAfter(now, nil) {
		t.Error("expected true when then is nil")
	}
}

func TestIsAfter_NowClearlyAfterThen_ReturnsTrue(t *testing.T) {
	now := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	then := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	if !IsAfter(now, &then) {
		t.Errorf("expected true: now=%v is after then=%v", now, then)
	}
}

func TestIsAfter_NowEqualThen_ReturnsFalse(t *testing.T) {
	now := time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC)
	then := time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC)
	// now-1 day = 2023-01-04, which is NOT after 2023-01-05 → false
	if IsAfter(now, &then) {
		t.Errorf("expected false: now=%v equals then=%v", now, then)
	}
}

func TestIsAfter_NowBeforeThen_ReturnsFalse(t *testing.T) {
	now := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	then := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	if IsAfter(now, &then) {
		t.Errorf("expected false: now=%v is before then=%v", now, then)
	}
}

func TestIsAfter_NowOneDayAfterThen_ReturnsFalse(t *testing.T) {
	now := time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC)
	then := time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC)
	// now-1 day = 2023-01-05, which is NOT after 2023-01-05 (strict) → false
	if IsAfter(now, &then) {
		t.Errorf("expected false: now-1=%v equals then=%v (strict After)", now.AddDate(0, 0, -1), then)
	}
}

func TestIsAfter_NowTwoDaysAfterThen_ReturnsTrue(t *testing.T) {
	now := time.Date(2023, 1, 7, 0, 0, 0, 0, time.UTC)
	then := time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC)
	// now-1 day = 2023-01-06, which IS after 2023-01-05 → true
	if !IsAfter(now, &then) {
		t.Errorf("expected true: now-1=%v is after then=%v", now.AddDate(0, 0, -1), then)
	}
}
