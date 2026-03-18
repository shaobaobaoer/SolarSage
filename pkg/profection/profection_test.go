package profection

import (
	"testing"

	"github.com/shaobaobaoer/solarsage-mcp/pkg/models"
)

func TestCalcAnnualProfection_Age0(t *testing.T) {
	// Age 0: profected sign = natal ASC sign
	p := CalcAnnualProfection(15.0, nil, 0) // ASC at 15° Aries
	if p.ProfectedSign != "Aries" {
		t.Errorf("Age 0: sign = %s, want Aries", p.ProfectedSign)
	}
	if p.ProfectedHouse != 1 {
		t.Errorf("Age 0: house = %d, want 1", p.ProfectedHouse)
	}
}

func TestCalcAnnualProfection_Age1(t *testing.T) {
	// Age 1: advances to 2nd sign
	p := CalcAnnualProfection(15.0, nil, 1)
	if p.ProfectedSign != "Taurus" {
		t.Errorf("Age 1: sign = %s, want Taurus", p.ProfectedSign)
	}
	if p.ProfectedHouse != 2 {
		t.Errorf("Age 1: house = %d, want 2", p.ProfectedHouse)
	}
}

func TestCalcAnnualProfection_Age12(t *testing.T) {
	// Age 12: full cycle, back to 1st house
	p := CalcAnnualProfection(15.0, nil, 12)
	if p.ProfectedSign != "Aries" {
		t.Errorf("Age 12: sign = %s, want Aries", p.ProfectedSign)
	}
	if p.ProfectedHouse != 1 {
		t.Errorf("Age 12: house = %d, want 1", p.ProfectedHouse)
	}
}

func TestCalcAnnualProfection_TimeLord(t *testing.T) {
	// ASC in Leo (120°), age 0: time lord = Sun (ruler of Leo)
	p := CalcAnnualProfection(120.0, nil, 0)
	if p.TimeLord != models.PlanetSun {
		t.Errorf("Leo time lord = %s, want SUN", p.TimeLord)
	}
}

func TestCalcMonthlyProfections(t *testing.T) {
	months := CalcMonthlyProfections(15.0, 0)
	if len(months) != 12 {
		t.Fatalf("Expected 12 monthly profections, got %d", len(months))
	}

	// First month should be same sign as annual
	if months[0].ProfectedSign != "Aries" {
		t.Errorf("Month 1: sign = %s, want Aries", months[0].ProfectedSign)
	}

	// Month 2 should advance one sign
	if months[1].ProfectedSign != "Taurus" {
		t.Errorf("Month 2: sign = %s, want Taurus", months[1].ProfectedSign)
	}
}

func TestProfectionTimeline(t *testing.T) {
	timeline := ProfectionTimeline(15.0, nil, 0, 5)
	if len(timeline) != 6 {
		t.Fatalf("Expected 6 entries (age 0-5), got %d", len(timeline))
	}
	for i, p := range timeline {
		if p.Age != i {
			t.Errorf("Entry %d: age = %d, want %d", i, p.Age, i)
		}
	}
}
