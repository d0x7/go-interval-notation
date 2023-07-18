package interval_notation

import (
	"github.com/Masterminds/semver"
	"testing"
)

type notationTests struct {
	intervalNotation string
	version          []string
	expected         []bool
}

func validTxt(valid bool) string {
	if valid {
		return "valid"
	}
	return "invalid"
}

func testNotations(t *testing.T, tests []notationTests) {
	for _, test := range tests {
		if len(test.version) != len(test.expected) {
			t.Errorf("Test data is invalid: version and expected must have the same length")
		}
		constraint, err := Parse(test.intervalNotation)
		if err != nil {
			t.Errorf("Error parsing interval: %v", err)
		}

		for i, v := range test.version {
			version, err := semver.NewVersion(v)
			if err != nil {
				t.Errorf("Error parsing version: %v", err)
			}

			valid, errors := constraint.Validate(version)

			if valid != test.expected[i] {
				t.Errorf("FAIL: Expected %s to be %s for %s but is %v", v, validTxt(test.expected[i]), test.intervalNotation, validTxt(valid))
				for _, m := range errors {
					t.Errorf(m.Error())
				}
			} else {
				t.Logf("PASS: %s is %s like it was expected for notation %s", v, validTxt(valid), test.intervalNotation)
			}

		}
	}
}

func TestLeftSetOnly(t *testing.T) {
	tests := []notationTests{
		{"(1.2.3,)", []string{"1.2.3", "1.2.4", "1.3.0", "2.0.0"}, []bool{false, true, true, true}},
		{"(1.2.3,]", []string{"1.2.3", "1.2.4", "1.3.0", "2.0.0"}, []bool{false, true, true, true}},
		{"(1.2.3)", []string{"1.3.0", "2.0.0", "1.2.3", "1.2.4"}, []bool{false, false, true, false}},
		{"[1.2.3,)", []string{"1.2.3", "1.2.4", "1.3.0", "2.0.0"}, []bool{true, true, true, true}},
		{"[1.2.3,]", []string{"1.2.3", "1.2.4", "1.3.0", "2.0.0"}, []bool{true, true, true, true}},
		{"[1.2.3)", []string{"1.2.3", "1.2.4", "1.3.0", "2.0.0"}, []bool{true, false, false, false}},
	}
	testNotations(t, tests)
}

func TestRightSetOnly(t *testing.T) {
	tests := []notationTests{
		{"(,1.2.3)", []string{"1.2.3", "1.2.4", "1.0.0", "1.1.3"}, []bool{false, false, true, true}},
		{"(,1.2.3]", []string{"1.2.3", "1.2.4", "1.0.0", "1.1.3"}, []bool{true, false, true, true}},
		{"(1.2.3)", []string{"1.3.0", "2.0.0", "1.2.3", "1.2.4"}, []bool{false, false, true, false}},
		{"[,1.2.3)", []string{"1.2.3", "1.2.4", "1.0.0", "1.1.3"}, []bool{false, false, true, true}},
		{"[,1.2.3]", []string{"1.2.3", "1.2.4", "1.0.0", "1.1.3"}, []bool{true, false, true, true}},
		{"[1.2.3)", []string{"1.2.3", "1.2.4", "1.3.0", "2.0.0"}, []bool{true, false, false, false}},
	}
	testNotations(t, tests)
}

func TestLeftAndRightSet(t *testing.T) {
	tests := []notationTests{
		{"(1.2.3,1.3.5)", []string{"1.2.3", "1.2.4", "1.2.5", "1.3.5"}, []bool{false, true, true, false}},
		{"(1.2.3,1.3.5]", []string{"1.2.3", "1.2.4", "1.2.5", "1.3.5"}, []bool{false, true, true, true}},
		{"[1.2.3,1.3.5)", []string{"1.2.3", "1.2.4", "1.2.5", "1.3.5"}, []bool{true, true, true, false}},
		{"[1.2.3,1.3.5]", []string{"1.2.3", "1.2.4", "1.2.5", "1.3.5"}, []bool{true, true, true, true}},
	}
	testNotations(t, tests)
}

func TestWithVersionPrefix(t *testing.T) {
	tests := []notationTests{
		{"(v1.2.3,v1.3.5)", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{false, true, true, false}},
		{"(v1.2.3,v1.3.5]", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{false, true, true, true}},
		{"[v1.2.3,v1.3.5)", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{true, true, true, false}},
		{"[v1.2.3,v1.3.5]", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{true, true, true, true}},
		{"(v1.2.3,1.3.5)", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{false, true, true, false}},
		{"(v1.2.3,1.3.5]", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{false, true, true, true}},
		{"[v1.2.3,1.3.5)", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{true, true, true, false}},
		{"[v1.2.3,1.3.5]", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{true, true, true, true}},
		{"(1.2.3,v1.3.5)", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{false, true, true, false}},
		{"(1.2.3,v1.3.5]", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{false, true, true, true}},
		{"[1.2.3,v1.3.5)", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{true, true, true, false}},
		{"[1.2.3,v1.3.5]", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{true, true, true, true}},
		{"(,v1.3.5)", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{true, true, true, false}},
		{"(,v1.3.5]", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{true, true, true, true}},
		{"[,v1.3.5)", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{true, true, true, false}},
		{"[,v1.3.5]", []string{"1.2.3", "v1.2.4", "1.2.5", "1.3.5"}, []bool{true, true, true, true}},
		{"(v1.3.5,)", []string{"1.2.3", "v1.5.4", "1.2.5", "1.3.5"}, []bool{false, true, false, false}},
		{"(v1.3.5,]", []string{"1.2.3", "v1.5.4", "1.2.5", "1.3.5"}, []bool{false, true, false, false}},
		{"[v1.3.5,)", []string{"1.2.3", "v1.5.4", "2.0.0", "1.3.5"}, []bool{false, true, true, true}},
		{"[v1.3.5,]", []string{"1.2.3", "v1.5.4", "2.0.0", "1.3.5"}, []bool{false, true, true, true}},
	}
	testNotations(t, tests)
}

func TestWithPrereleaseAndBuildMetadata(t *testing.T) {
	tests := []notationTests{
		{"(1.2.3-alpha.1,1.3.5-0)", []string{"1.2.3-alpha.1", "1.2.4-alpha.2", "1.2.4-beta.1", "1.2.3-rc.1", "1.2.4-rc1", "1.2.4", "1.2.5", "1.2.6", "1.3.5"}, []bool{false, true, true, true, true, true, true, true, false}},
	}
	testNotations(t, tests)
}
