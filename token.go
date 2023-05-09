package calver

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type token interface {
	Token() string
	TrimPrefix(value string) (prefix string, trimed string, err error)
}

var (
	_ token = tokenCal{}
	_ token = tokenVer{}
	_ token = tokenSep{}
)

type tokenCal struct {
	token        string
	timeToString func(time.Time) string
}

func (t tokenCal) Token() string {
	return t.token
}

func (t tokenCal) TrimPrefix(value string) (string, string, error) {
	l := len(t.token)
	var expr string
	if l == 2 && !strings.HasPrefix(t.token, "0") {
		expr = "^([1-9][0-9]?)(.*)$"
	} else {
		expr = fmt.Sprintf("^([0-9]{%d})(.*)$", l)
	}
	re := regexp.MustCompile(expr)
	matches := re.FindAllStringSubmatch(value, -1)
	if len(matches) == 0 {
		return "", "", fmt.Errorf("could not get the value of token '%s' from '%s'", t.token, value)
	}
	return matches[0][1], matches[0][2], nil
}

type tokenVer struct {
	token       string
	verToString func(int, int, int, string) string
}

func (t tokenVer) Token() string {
	return t.token
}

func (t tokenVer) TrimPrefix(value string) (string, string, error) {
	if t.token == "MODIFIER" {
		return value, "", nil
	}
	expr := "^([0-9]+)(.*)$"
	re := regexp.MustCompile(expr)
	matches := re.FindAllStringSubmatch(value, -1)
	if len(matches) == 0 {
		return "", "", fmt.Errorf("could not get the value of token '%s' from '%s'", t.token, value)
	}
	return matches[0][1], matches[0][2], nil
}

type tokenSep struct {
	token string
}

func newTokenSep(token string) tokenSep {
	return tokenSep{token: token}
}

func (t tokenSep) Token() string {
	return t.token
}

func (t tokenSep) TrimPrefix(value string) (string, string, error) {
	if !strings.HasPrefix(value, t.token) {
		return "", "", fmt.Errorf("could not get the value of token '%s' from '%s'", t.token, value)
	}
	return t.token, strings.TrimPrefix(value, t.token), nil
}

func (t tokenSep) String() string {
	return t.token
}

var (
	tYYYY = tokenCal{token: "YYYY", timeToString: func(t time.Time) string { return t.Format("2006") }}
	tYY   = tokenCal{token: "YY", timeToString: func(t time.Time) string { return strings.TrimPrefix(t.Format("06"), "0") }}
	t0Y   = tokenCal{token: "0Y", timeToString: func(t time.Time) string { return t.Format("06") }}
	tMM   = tokenCal{token: "MM", timeToString: func(t time.Time) string { return t.Format("1") }}
	t0M   = tokenCal{token: "0M", timeToString: func(t time.Time) string { return t.Format("01") }}
	tWW   = tokenCal{token: "WW", timeToString: func(t time.Time) string {
		_, w := t.ISOWeek()
		return fmt.Sprintf("%d", w)
	}}
	t0W = tokenCal{token: "0W", timeToString: func(t time.Time) string {
		_, w := t.ISOWeek()
		return fmt.Sprintf("%02d", w)
	}}
	tDD = tokenCal{token: "DD", timeToString: func(t time.Time) string { return t.Format("2") }}
	t0D = tokenCal{token: "0D", timeToString: func(t time.Time) string { return t.Format("02") }}

	tMAJOR    = tokenVer{token: "MAJOR", verToString: func(major, minor, micro int, modifier string) string { return fmt.Sprintf("%d", major) }}
	tMINOR    = tokenVer{token: "MINOR", verToString: func(major, minor, micro int, modifier string) string { return fmt.Sprintf("%d", minor) }}
	tMICRO    = tokenVer{token: "MICRO", verToString: func(major, minor, micro int, modifier string) string { return fmt.Sprintf("%d", micro) }}
	tMODIFIER = tokenVer{token: "MODIFIER", verToString: func(major, minor, micro int, modifier string) string { return modifier }}
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
			if strings.HasPrefix(t.Token(), v) {
				prefixMatches = append(prefixMatches, t)
			}
			if t.Token() == v {
				match = t
			}
			if t.Token() == prev {
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
	return tokens, nil
}
