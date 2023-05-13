package calver

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type token interface {
	token() string
	trimPrefix(value string) (prefix string, trimed string, err error)
	String() string
}

var (
	_ token = tokenCal{}
	_ token = tokenVer{}
	_ token = tokenSep{}
)

type tokenCal struct {
	t            string
	timeToString func(time.Time) string
}

func (t tokenCal) token() string {
	return t.t
}

func (t tokenCal) String() string {
	return t.t
}

func (t tokenCal) trimPrefix(value string) (string, string, error) {
	l := len(t.t)
	var expr string
	if l == 2 && !strings.HasPrefix(t.t, "0") {
		expr = "^([1-9][0-9]?)(.*)$"
	} else {
		expr = fmt.Sprintf("^([0-9]{%d})(.*)$", l)
	}
	re := regexp.MustCompile(expr)
	matches := re.FindAllStringSubmatch(value, -1)
	if len(matches) == 0 {
		return "", "", fmt.Errorf("could not get the value of token '%s' from '%s'", t.t, value)
	}
	return matches[0][1], matches[0][2], nil
}

type tokenVer struct {
	t           string
	verToString func(int, int, int, string) string
}

func (t tokenVer) token() string {
	return t.t
}

func (t tokenVer) String() string {
	return t.t
}

func (t tokenVer) trimPrefix(value string) (string, string, error) {
	if t.t == "MODIFIER" {
		return value, "", nil
	}
	expr := "^([0-9]+)(.*)$"
	re := regexp.MustCompile(expr)
	matches := re.FindAllStringSubmatch(value, -1)
	if len(matches) == 0 {
		return "", "", fmt.Errorf("could not get the value of token '%s' from '%s'", t.t, value)
	}
	return matches[0][1], matches[0][2], nil
}

type tokenSep struct {
	t string
}

func newTokenSep(token string) tokenSep {
	return tokenSep{t: token}
}

func (t tokenSep) token() string {
	return t.t
}

func (t tokenSep) String() string {
	return t.t
}

func (t tokenSep) trimPrefix(value string) (string, string, error) {
	if !strings.HasPrefix(value, t.t) {
		return "", "", fmt.Errorf("could not get the value of token '%s' from '%s'", t.t, value)
	}
	return t.t, strings.TrimPrefix(value, t.t), nil
}

var (
	tYYYY = tokenCal{t: "YYYY", timeToString: func(t time.Time) string { return t.Format("2006") }}
	tYY   = tokenCal{t: "YY", timeToString: func(t time.Time) string { return strings.TrimPrefix(t.Format("06"), "0") }}
	t0Y   = tokenCal{t: "0Y", timeToString: func(t time.Time) string { return t.Format("06") }}
	tMM   = tokenCal{t: "MM", timeToString: func(t time.Time) string { return t.Format("1") }}
	t0M   = tokenCal{t: "0M", timeToString: func(t time.Time) string { return t.Format("01") }}
	tWW   = tokenCal{t: "WW", timeToString: func(t time.Time) string {
		_, w := t.ISOWeek()
		return fmt.Sprintf("%d", w)
	}}
	t0W = tokenCal{t: "0W", timeToString: func(t time.Time) string {
		_, w := t.ISOWeek()
		return fmt.Sprintf("%02d", w)
	}}
	tDD = tokenCal{t: "DD", timeToString: func(t time.Time) string { return t.Format("2") }}
	t0D = tokenCal{t: "0D", timeToString: func(t time.Time) string { return t.Format("02") }}

	tMAJOR    = tokenVer{t: "MAJOR", verToString: func(major, minor, micro int, modifier string) string { return fmt.Sprintf("%d", major) }}
	tMINOR    = tokenVer{t: "MINOR", verToString: func(major, minor, micro int, modifier string) string { return fmt.Sprintf("%d", minor) }}
	tMICRO    = tokenVer{t: "MICRO", verToString: func(major, minor, micro int, modifier string) string { return fmt.Sprintf("%d", micro) }}
	tMODIFIER = tokenVer{t: "MODIFIER", verToString: func(major, minor, micro int, modifier string) string { return modifier }}
)

var builtinTokens = []token{
	tYYYY,
	tYY,
	t0Y,
	tMM,
	t0M,
	tWW,
	t0W,
	tDD,
	t0D,
	tMAJOR,
	tMINOR,
	tMICRO,
	tMODIFIER,
}

func tokenizeLayout(layout string) ([]token, error) {
	tokens := []token{}
	splitted := strings.Split(layout, "")
	size := len(splitted)
	pos := 0
	for idx := 0; idx < size; idx++ {
		v := strings.Join(splitted[pos:idx+1], "")
		prev := strings.Join(splitted[pos:idx], "")
		var match token
		var prevMatch token
		prefixMatches := []token{}
		for _, t := range builtinTokens {
			if strings.HasPrefix(t.token(), v) {
				prefixMatches = append(prefixMatches, t)
			}
			if t.token() == v {
				match = t
			}
			if t.token() == prev {
				prevMatch = t
			}
		}
		switch {
		case len(prefixMatches) == 1 && match != nil:
			tokens = append(tokens, match)
			pos = idx + 1
		case len(prefixMatches) == 0 && prevMatch != nil:
			tokens = append(tokens, prevMatch)
			pos = idx
			idx--
		case len(prefixMatches) == 0 && prevMatch == nil:
			tokens = append(tokens, newTokenSep(v))
			pos = idx + 1
		case size == idx+1 && match != nil:
			tokens = append(tokens, match)
		case size == idx+1 && match == nil && prevMatch != nil:
			tokens = append(tokens, prevMatch)
			pos = idx
			idx--
		case size == idx+1 && match == nil && prevMatch == nil:
			tokens = append(tokens, newTokenSep(v))
		}
	}
	if !lessThanOneContains(tokens, []token{tYYYY, tYY, t0Y}) {
		return nil, fmt.Errorf("Only one of %v, %v, %v can be included in the layout", tYYYY, tYY, t0Y)
	}
	if !lessThanOneContains(tokens, []token{tMM, t0M}) {
		return nil, fmt.Errorf("Only one of %v, %v can be included in the layout", tMM, t0M)
	}
	if !lessThanOneContains(tokens, []token{tWW, t0W}) {
		return nil, fmt.Errorf("Only one of %v, %v can be included in the layout", tWW, t0W)
	}
	if !lessThanOneContains(tokens, []token{tDD, t0D}) {
		return nil, fmt.Errorf("Only one of %v, %v can be included in the layout", tDD, t0D)
	}
	if !lessThanOneContains(tokens, []token{tMAJOR}) {
		return nil, fmt.Errorf("Only one %v can be included in the layout", tMAJOR)
	}
	if !lessThanOneContains(tokens, []token{tMINOR}) {
		return nil, fmt.Errorf("Only one %v can be included in the layout", tMINOR)
	}
	if !lessThanOneContains(tokens, []token{tMICRO}) {
		return nil, fmt.Errorf("Only one %v can be included in the layout", tMICRO)
	}
	if !lessThanOneContains(tokens, []token{tMODIFIER}) {
		return nil, fmt.Errorf("Only one %v can be included in the layout", tMODIFIER)
	}

	return tokens, nil
}

func lessThanOneContains(layout, target []token) bool {
	contained := []token{}
	for _, t := range layout {
		if contains(target, t) {
			contained = append(contained, t)
		}
	}
	return len(contained) <= 1
}
