package calver

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/snabb/isoweek"
)

const (
	parsedDefaultYear  = 2000
	parsedDefaultMonth = time.Month(1)
	parsedDefaultDay   = 1
)

type Calver struct {
	major      int
	minor      int
	micro      int
	modifier   string
	ts         time.Time
	loc        *time.Location
	layout     []token
	trimSuffix bool
}

type Calvers []*Calver

var ErrNoVersions = errors.New("no versions")

// New returns *Calver at the current time.
func New(layout string) (*Calver, error) {
	now := time.Now().UTC()
	return NewWithTime(layout, now)
}

// NewWithTime returns *Calver at the given time.
func NewWithTime(layout string, now time.Time) (*Calver, error) {
	tokens, err := tokenizeLayout(layout)
	if err != nil {
		return nil, err
	}
	return &Calver{
		ts:     now, // Do not initialize (zeronize) below hour for In()
		loc:    now.Location(),
		layout: tokens,
	}, nil
}

// In sets *time.Location.
func (cv *Calver) In(loc *time.Location) *Calver {
	ncv := cv.clone()
	ncv.loc = loc
	return ncv
}

// Parse version string using layout.
func (cv *Calver) Parse(value string) (ncv *Calver, err error) {
	org := value
	defer func() {
		if err != nil {
			err = fmt.Errorf("failed to parse '%s' using layout '%s': %w", org, cv.Layout(), err)
		}
	}()
	ncv = cv.clone()
	year := parsedDefaultYear
	month := parsedDefaultMonth
	day := parsedDefaultDay

	var (
		p    string
		week int
	)
	for _, t := range cv.layout {
		switch {
		case contains([]token{tYYYY, tYY, t0Y}, t):
			p, value, err = t.trimPrefix(value)
			if err != nil {
				return nil, err
			}
			year, err = strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			if year < 2000 {
				year += 2000
			}
		case contains([]token{tMM, t0M}, t):
			p, value, err = t.trimPrefix(value)
			if err != nil {
				return nil, err
			}
			m, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			month = time.Month(m)
		case contains([]token{tWW, t0W}, t):
			p, value, err = t.trimPrefix(value)
			if err != nil {
				return nil, err
			}
			week, err = strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
		case contains([]token{tDD, t0D}, t):
			p, value, err = t.trimPrefix(value)
			if err != nil {
				return nil, err
			}
			day, err = strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
		case contains([]token{tMAJOR}, t):
			if value == "" && cv.trimSuffix {
				ncv.major = 0
				continue
			}
			p, value, err = t.trimPrefix(value)
			if err != nil {
				return nil, err
			}
			m, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			ncv.major = m
		case contains([]token{tMINOR}, t):
			if value == "" && cv.trimSuffix {
				ncv.minor = 0
				continue
			}
			p, value, err = t.trimPrefix(value)
			if err != nil {
				return nil, err
			}
			m, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			ncv.minor = m
		case contains([]token{tMICRO}, t):
			if value == "" && cv.trimSuffix {
				ncv.micro = 0
				continue
			}
			p, value, err = t.trimPrefix(value)
			if err != nil {
				return nil, err
			}
			m, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			ncv.micro = m
		case contains([]token{tMODIFIER}, t):
			if value == "" && cv.trimSuffix {
				ncv.modifier = ""
				continue
			}
			p, value, err = t.trimPrefix(value)
			if err != nil {
				return nil, err
			}
			ncv.modifier = p
		default:
			if value == "" && cv.trimSuffix {
				continue
			}
			_, value, err = t.trimPrefix(value)
			if err != nil {
				return nil, err
			}
		}
	}
	if week > 0 {
		year, month, day = isoweek.StartDate(year, week)
	}
	// Initialize (zeronize) hour and below when parsing
	ncv.ts = time.Date(year, month, day, 0, 0, 0, 0, cv.loc)
	if value != "" {
		return nil, errors.New("there are strings that could not be parsed")
	}
	return ncv, nil
}

// Parse version string using layout at the current time.
func Parse(layout, value string) (*Calver, error) {
	cv, err := New(layout)
	if err != nil {
		return nil, err
	}
	return cv.Parse(value)
}

// String returns version string.
func (cv *Calver) String() string {
	var s string
	reversed := reverse(cv.layout)
	trimable := cv.trimSuffix
	for _, t := range reversed {
		switch tt := t.(type) {
		case tokenCal:
			trimable = false
			s = tt.timeToString(cv.ts.In(cv.loc)) + s
		case tokenVer:
			v := tt.verToString(cv.major, cv.minor, cv.micro, cv.modifier)
			if trimable && (v == "0" || v == "") {
				v = ""
			} else {
				trimable = false
			}
			s = v + s
		case tokenSep:
			if !trimable {
				s = tt.String() + s
			}
		}
	}
	return s
}

// Layout returns version layout.
func (cv *Calver) Layout() string {
	var s string
	for _, t := range cv.layout {
		s += t.token()
	}
	return s
}

// Next returns next version *Calver at the current time.
func (cv *Calver) Next() (*Calver, error) {
	now := time.Now()
	return cv.NextWithTime(now)
}

// Next returns next version *Calver at the given time.
func (cv *Calver) NextWithTime(now time.Time) (*Calver, error) {
	if cv.ts.UnixNano() > now.UnixNano() {
		return nil, fmt.Errorf("[%v] is older than the current setting (%v)", now.Truncate(0), cv.ts)
	}
	ncv := cv.clone()
	ncv.ts = now
	if cv.String() != ncv.String() {
		return ncv, nil
	}
	if contains(cv.layout, tMICRO) {
		return cv.Micro()
	}
	if contains(cv.layout, tMINOR) {
		return cv.Minor()
	}
	if contains(cv.layout, tMAJOR) {
		return cv.Major()
	}
	return nil, errors.New("failed to bump up version")
}

// Major returns next major version *Calver.
func (cv *Calver) Major() (*Calver, error) {
	if !contains(cv.layout, tMAJOR) {
		return nil, fmt.Errorf("no 'MAJOR' in the layout '%s'", cv.Layout())
	}
	ncv := cv.clone()
	ncv.major++
	return ncv, nil
}

// Minor returns next minor version *Calver.
func (cv *Calver) Minor() (*Calver, error) {
	if !contains(cv.layout, tMINOR) {
		return nil, fmt.Errorf("no 'MINOR' in the layout '%s'", cv.Layout())
	}
	ncv := cv.clone()
	ncv.minor++
	return ncv, nil
}

// Micro returns next micro version *Calver.
func (cv *Calver) Micro() (*Calver, error) {
	if !contains(cv.layout, tMICRO) {
		return nil, fmt.Errorf("no 'MICRO' in the layout '%s'", cv.Layout())
	}
	ncv := cv.clone()
	ncv.micro++
	return ncv, nil
}

// Modifier returns *Calver with modifier.
func (cv *Calver) Modifier(m string) (*Calver, error) {
	if !contains(cv.layout, tMODIFIER) {
		return nil, fmt.Errorf("no 'MODIFIER' in the layout '%s'", cv.Layout())
	}
	ncv := cv.clone()
	ncv.modifier = m
	return ncv, nil
}

// TrimSuffix returns *Calver enabled/diabled to trim the trailing version of a zero value or an empty string.
func (cv *Calver) TrimSuffix(enable bool) *Calver {
	ncv := cv.clone()
	ncv.trimSuffix = enable
	return ncv
}

func (cv *Calver) clone() *Calver {
	return &Calver{
		major:      cv.major,
		minor:      cv.minor,
		micro:      cv.micro,
		modifier:   cv.modifier,
		ts:         cv.ts,
		loc:        cv.loc,
		layout:     cv.layout,
		trimSuffix: cv.trimSuffix,
	}
}

func (cvs Calvers) Sort() {
	sort.SliceStable(cvs, func(i, j int) bool {
		switch {
		case cvs[i].ts.UnixNano() != cvs[j].ts.UnixNano():
			return cvs[i].ts.UnixNano() > cvs[j].ts.UnixNano()
		case cvs[i].major != cvs[j].major:
			return cvs[i].major > cvs[j].major
		case cvs[i].minor != cvs[j].minor:
			return cvs[i].minor > cvs[j].minor
		case cvs[i].micro != cvs[j].micro:
			return cvs[i].micro > cvs[j].micro
		case cvs[i].modifier == "":
			return true
		case cvs[j].modifier == "":
			return false
		default:
			return cvs[i].modifier > cvs[j].modifier
		}
	})
}

func (cvs Calvers) Latest() (*Calver, error) {
	if len(cvs) == 0 {
		return nil, ErrNoVersions
	}
	cvs.Sort()
	return cvs[0], nil
}

func contains(layout []token, t token) bool {
	for _, tt := range layout {
		if tt.token() == t.token() {
			return true
		}
	}
	return false
}

func reverse(layout []token) []token {
	reversed := []token{}
	for _, t := range layout {
		reversed = append([]token{t}, reversed...)
	}
	return reversed
}
