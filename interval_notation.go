package interval_notation

import (
	"fmt"
	"github.com/Masterminds/semver"
	"regexp"
	"strings"
)

var validationRegex = regexp.MustCompile(`(?P<leftBrace>[(\]\[])(?P<first>v?[0-9a-zA-Z.+-]*?)(?P<delimiter>,)?(?P<second>v?[0-9a-zA-Z.+-]*?)(?P<rightBrace>[)\]\[])`)

func Parse(interval string) (*semver.Constraints, error) {
	match := validationRegex.FindStringSubmatch(interval)
	if match == nil || len(match) != 6 || (match[validationRegex.SubexpIndex("first")] == "" && match[validationRegex.SubexpIndex("second")] == "") {
		return nil, fmt.Errorf("invalid interval notation: %s", interval)
	}

	left, right := match[validationRegex.SubexpIndex("first")], match[validationRegex.SubexpIndex("second")]
	leftSet, rightSet := left != "", right != ""
	bothSet := leftSet && rightSet
	hasDelimiter := match[validationRegex.SubexpIndex("delimiter")] != ""
	leftInclusive, rightInclusive := match[validationRegex.SubexpIndex("leftBrace")] == "[", match[validationRegex.SubexpIndex("rightBrace")] == "]"

	var constraint string

	switch {
	case bothSet && hasDelimiter:
		{
			var sb strings.Builder
			sb.WriteString(">")
			if leftInclusive {
				sb.WriteString("=")
			}
			sb.WriteString(" ")
			sb.WriteString(left)
			sb.WriteString(", ")
			sb.WriteString("<")

			if rightInclusive {
				sb.WriteString("=")
			}

			sb.WriteString(" ")
			sb.WriteString(right)
			constraint = sb.String()
		}
	case !hasDelimiter:
		{
			if leftSet {
				constraint = fmt.Sprintf("= %s", left)
			} else {
				constraint = fmt.Sprintf("= %s", right)
			}
		}
	case leftSet && leftInclusive:
		constraint = fmt.Sprintf(">= %s", left)
	case leftSet:
		constraint = fmt.Sprintf("> %s", left)
	case rightSet && rightInclusive:
		constraint = fmt.Sprintf("<= %s", right)
	case rightSet:
		constraint = fmt.Sprintf("< %s", right)
	}
	verConstraint, err := semver.NewConstraint(constraint)
	if err != nil {
		return nil, err
	}

	return verConstraint, nil
}

func InRange(intervalNotation, version string) (bool, []error) {
	constraint, err := Parse(intervalNotation)
	if err != nil {
		return false, []error{err}
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		return false, []error{err}
	}

	return constraint.Validate(v)
}

func IsInRange(intervalNotation, version string) bool {
	inRange, _ := InRange(intervalNotation, version)
	return inRange
}
