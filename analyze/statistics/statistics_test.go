package statistics

import (
	"strconv"
	"testing"
)

// --- weekdays slice ---

func TestWeekdays_HasSevenEntries(t *testing.T) {
	if len(weekdays) != 7 {
		t.Errorf("expected 7 weekdays, got %d", len(weekdays))
	}
}

func TestWeekdays_StartsWithSunday(t *testing.T) {
	// SQLite strftime('%w') returns 0 for Sunday
	if weekdays[0] != "Sunday" {
		t.Errorf("expected weekdays[0]='Sunday', got '%s'", weekdays[0])
	}
}

func TestWeekdays_EndsWithSaturday(t *testing.T) {
	if weekdays[6] != "Saturday" {
		t.Errorf("expected weekdays[6]='Saturday', got '%s'", weekdays[6])
	}
}

func TestWeekdays_AllDaysPresent(t *testing.T) {
	expected := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	for i, want := range expected {
		if weekdays[i] != want {
			t.Errorf("weekdays[%d]: expected '%s', got '%s'", i, want, weekdays[i])
		}
	}
}

// --- months slice ---

func TestMonths_HasTwelveEntries(t *testing.T) {
	if len(months) != 12 {
		t.Errorf("expected 12 months, got %d", len(months))
	}
}

func TestMonths_StartsWithJanuary(t *testing.T) {
	// SQLite strftime('%m') returns "01" for January → index months[0]
	if months[0] != "January" {
		t.Errorf("expected months[0]='January', got '%s'", months[0])
	}
}

func TestMonths_EndsWithDecember(t *testing.T) {
	if months[11] != "December" {
		t.Errorf("expected months[11]='December', got '%s'", months[11])
	}
}

func TestMonths_AllMonthsPresent(t *testing.T) {
	expected := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	for i, want := range expected {
		if months[i] != want {
			t.Errorf("months[%d]: expected '%s', got '%s'", i, want, months[i])
		}
	}
}

// --- Month index mapping (strftime returns 1-based "01".."12") ---

func TestMonths_IndexMapping_StrftimeValue01_IsJanuary(t *testing.T) {
	// Simulates the code: val, _ := strconv.Atoi("01"); months[val-1]
	val, _ := strconv.Atoi("01")
	if months[val-1] != "January" {
		t.Errorf("strftime '01' → expected 'January', got '%s'", months[val-1])
	}
}

func TestMonths_IndexMapping_StrftimeValue12_IsDecember(t *testing.T) {
	val, _ := strconv.Atoi("12")
	if months[val-1] != "December" {
		t.Errorf("strftime '12' → expected 'December', got '%s'", months[val-1])
	}
}

func TestMonths_IndexMapping_StrftimeValue06_IsJune(t *testing.T) {
	val, _ := strconv.Atoi("06")
	if months[val-1] != "June" {
		t.Errorf("strftime '06' → expected 'June', got '%s'", months[val-1])
	}
}

// --- TimeStatistic structs (DAY, MONTH, YEAR) ---

func TestTimeStatistic_DAY_HasCorrectLabel(t *testing.T) {
	if DAY.label != "Weekday" {
		t.Errorf("expected DAY.label='Weekday', got '%s'", DAY.label)
	}
}

func TestTimeStatistic_DAY_HasCorrectFormat(t *testing.T) {
	if DAY.format != "%w" {
		t.Errorf("expected DAY.format='%%w', got '%s'", DAY.format)
	}
}

func TestTimeStatistic_MONTH_HasCorrectLabel(t *testing.T) {
	if MONTH.label != "Month" {
		t.Errorf("expected MONTH.label='Month', got '%s'", MONTH.label)
	}
}

func TestTimeStatistic_MONTH_HasCorrectFormat(t *testing.T) {
	if MONTH.format != "%m %Y" {
		t.Errorf("expected MONTH.format='%%m %%Y', got '%s'", MONTH.format)
	}
}

func TestTimeStatistic_YEAR_HasCorrectLabel(t *testing.T) {
	if YEAR.label != "Year" {
		t.Errorf("expected YEAR.label='Year', got '%s'", YEAR.label)
	}
}

func TestTimeStatistic_YEAR_HasCorrectFormat(t *testing.T) {
	if YEAR.format != "%Y" {
		t.Errorf("expected YEAR.format='%%Y', got '%s'", YEAR.format)
	}
}

// --- Weekday index mapping (strftime '%w' returns "0".."6") ---

func TestWeekdays_IndexMapping_StrftimeValue0_IsSunday(t *testing.T) {
	val, _ := strconv.Atoi("0")
	if weekdays[val] != "Sunday" {
		t.Errorf("strftime '0' → expected 'Sunday', got '%s'", weekdays[val])
	}
}

func TestWeekdays_IndexMapping_StrftimeValue6_IsSaturday(t *testing.T) {
	val, _ := strconv.Atoi("6")
	if weekdays[val] != "Saturday" {
		t.Errorf("strftime '6' → expected 'Saturday', got '%s'", weekdays[val])
	}
}

func TestWeekdays_IndexMapping_StrftimeValue3_IsWednesday(t *testing.T) {
	val, _ := strconv.Atoi("3")
	if weekdays[val] != "Wednesday" {
		t.Errorf("strftime '3' → expected 'Wednesday', got '%s'", weekdays[val])
	}
}
