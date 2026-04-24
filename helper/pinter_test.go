package helper

import (
	"testing"
	"time"
)

func TestBool(t *testing.T) {
	v := true
	p := Bool(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestByte(t *testing.T) {
	v := byte(42)
	p := Byte(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestFloat32(t *testing.T) {
	v := float32(3.14)
	p := Float32(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestFloat64(t *testing.T) {
	v := float64(2.718281828)
	p := Float64(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestInt(t *testing.T) {
	v := 99
	p := Int(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestInt8(t *testing.T) {
	v := int8(-10)
	p := Int8(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestInt16(t *testing.T) {
	v := int16(1000)
	p := Int16(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestInt32(t *testing.T) {
	v := int32(100000)
	p := Int32(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestInt64(t *testing.T) {
	v := int64(9999999999)
	p := Int64(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestRune(t *testing.T) {
	v := rune('€')
	p := Rune(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestString(t *testing.T) {
	v := "hello tarpit"
	p := String(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestUint(t *testing.T) {
	v := uint(42)
	p := Uint(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestUint8(t *testing.T) {
	v := uint8(255)
	p := Uint8(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestUint16(t *testing.T) {
	v := uint16(65535)
	p := Uint16(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestUint32(t *testing.T) {
	v := uint32(4294967295)
	p := Uint32(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestUint64(t *testing.T) {
	v := uint64(18446744073709551615)
	p := Uint64(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestTime(t *testing.T) {
	v := time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC)
	p := Time(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if !p.Equal(v) {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

func TestDuration(t *testing.T) {
	v := 5 * time.Minute
	p := Duration(v)
	if p == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *p != v {
		t.Errorf("expected %v, got %v", v, *p)
	}
}

// Verify that pointer functions return independent copies (no aliasing)
func TestPointerIndependence(t *testing.T) {
	a := 10
	p1 := Int(a)
	p2 := Int(a)
	if p1 == p2 {
		t.Error("expected two independent pointers, got same address")
	}
}
