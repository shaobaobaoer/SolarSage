package sweph

import "testing"

func TestHeliacalUT_Venus(t *testing.T) {
	jdStart := JulDay(2000, 1, 1, 0.0, true)
	result, err := HeliacalUT(jdStart, -0.1278, 51.5074, 0, "venus", SE_HELIACAL_RISING)
	if err != nil {
		t.Fatalf("HeliacalUT: %v", err)
	}
	if result.JDStart <= jdStart {
		t.Errorf("JDStart %f should be after search start %f", result.JDStart, jdStart)
	}
}

func TestHeliacalUT_Mercury(t *testing.T) {
	jdStart := JulDay(2000, 6, 1, 0.0, true)
	result, err := HeliacalUT(jdStart, -0.1278, 51.5074, 0, "mercury", SE_HELIACAL_RISING)
	if err != nil {
		t.Fatalf("HeliacalUT mercury: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}
