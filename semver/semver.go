package semver

import (
	"regexp"
	"strconv"
	"strings"
)

func Compare(s1, s2 string) (result int, err error) {
	if s1 == s2 {
		return 0, nil
	}

	if len(s1) == 0 && len(s2) > 0 {
		return -1, nil
	}

	if len(s1) > 0 && len(s2) == 0 {
		return 1, nil
	}

	s1Segments, s2Segments, err := parseArguments(s1, s2)
	if err != nil {
		return 0, err
	}

	semverLength := min(len(s1Segments), len(s2Segments))
	for i := 0; i < semverLength; i++ {
		v1, err := strconv.Atoi(s1Segments[i])
		if err != nil {
			return 0, err
		}

		v2, err := strconv.Atoi(s2Segments[i])
		if err != nil {
			return 0, err
		}

		if v1 > v2 {
			return 1, nil
		}

		if v1 < v2 {
			return -1, nil
		}
	}

	return compareBySegmentLengths(s1Segments, s2Segments), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func parseArguments(s1, s2 string) (s1Segments []string, s2Segments []string, err error) {
	l1 := findAndReplace("\\pL", s1, "") // matches all unicode letters
	r1 := findAndReplace("\\pL", s2, "") // matches all unicode letters
	l2 := findAndReplace("\\.+", l1, ".")
	r2 := findAndReplace("\\.+", r1, ".")
	return strings.Split(l2, "."), strings.Split(r2, "."), nil
}

func findAndReplace(pattern string, src string, replace string) string {
	re, _ := regexp.Compile(pattern)
	return re.ReplaceAllString(src, replace)
}

func compareBySegmentLengths(s1Segments []string, s2Segments []string) int {
	if len(s1Segments) > len(s2Segments) {
		return 1
	}

	if len(s1Segments) < len(s2Segments) {
		return -1
	}

	return 0
}