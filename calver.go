package calver

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/snabb/isoweek"
)

type Calver struct {
	major    int
	minor    int
	micro    int
	modifier string
	ts       time.Time
	loc      *time.Location
	layout   []token
}

func New(layout string) (*Calver, error) {
	now := time.Now().UTC()
	return NewWithTime(layout, now)
}

func NewWithTime(layout string, now time.Time) (*Calver, error) {
	tokens, err := tokenizeLayout(layout)
	if err != nil {
		return nil, err
	}
	return &Calver{
		ts:     now,
		loc:    now.Location(),
		layout: tokens,
	}, nil
}

func (cv *Calver) In(loc *time.Location) *Calver {
	ncv := cv.clone()
	ncv.loc = loc
	return ncv
}

func (cv *Calver) Parse(value string) (*Calver, error) {
	org := value
	ncv := cv.clone()
	year := ncv.ts.In(ncv.loc).Year()
	month := ncv.ts.In(ncv.loc).Month()
	day := ncv.ts.In(ncv.loc).Day()

	var (
		p    string
		err  error
		week int
	)
	for _, t := range cv.layout {
		switch {
		case contains([]token{tYYYY, tYY, t0Y}, t):
			p, value, err = t.TrimPrefix(value)
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
			p, value, err = t.TrimPrefix(value)
			if err != nil {
				return nil, err
			}
			m, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			month = time.Month(m)
		case contains([]token{tWW, t0W}, t):
			p, value, err = t.TrimPrefix(value)
			if err != nil {
				return nil, err
			}
			week, err = strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
		case contains([]token{tDD, t0D}, t):
			p, value, err = t.TrimPrefix(value)
			if err != nil {
				return nil, err
			}
			day, err = strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
		case contains([]token{tMAJOR}, t):
			p, value, err = t.TrimPrefix(value)
			if err != nil {
				return nil, err
			}
			m, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			ncv.major = m
		case contains([]token{tMINOR}, t):
			p, value, err = t.TrimPrefix(value)
			if err != nil {
				return nil, err
			}
			m, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			ncv.minor = m
		case contains([]token{tMICRO}, t):
			p, value, err = t.TrimPrefix(value)
			if err != nil {
				return nil, err
			}
			m, err := strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
			ncv.micro = m
		case contains([]token{tMODIFIER}, t):
			p, value, err = t.TrimPrefix(value)
			if err != nil {
				return nil, err
			}
			ncv.modifier = p
		default:
			_, value, err = t.TrimPrefix(value)
			if err != nil {
				return nil, err
			}
		}
	}
	if week > 0 {
		year, month, day = isoweek.StartDate(year, week)
	}
	ncv.ts = time.Date(year, month, day, cv.ts.In(cv.loc).Hour(), cv.ts.In(cv.loc).Minute(), cv.ts.In(cv.loc).Second(), cv.ts.In(cv.loc).Nanosecond(), cv.loc)
	if value != "" {
		return nil, fmt.Errorf("failed to parse: %s", org)
	}
	return ncv, nil
}

func Parse(layout, value string) (*Calver, error) {
	cv, err := New(layout)
	if err != nil {
		return nil, err
	}
	return cv.Parse(value)
}

func (cv *Calver) String() string {
	var s string
	for _, t := range cv.layout {
		switch tt := t.(type) {
		case tokenCal:
			s += tt.timeToString(cv.ts.In(cv.loc))
		case tokenVer:
			s += tt.verToString(cv.major, cv.minor, cv.micro, cv.modifier)
		case tokenSep:
			s += tt.String()
		}
	}
	return s
}

func (cv *Calver) Next() (*Calver, error) {
	now := time.Now()
	return cv.NextWithTime(now)
}

func (cv *Calver) NextWithTime(now time.Time) (*Calver, error) {
	if cv.ts.UnixNano() > now.UnixNano() {
		return nil, fmt.Errorf("%v is older than the current setting (%v)", now, cv.ts)
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

func (cv *Calver) Major() (*Calver, error) {
	if !contains(cv.layout, tMAJOR) {
		return nil, errors.New("no 'MAJOR' in the layout")
	}
	ncv := cv.clone()
	ncv.major++
	return ncv, nil
}

func (cv *Calver) Minor() (*Calver, error) {
	if !contains(cv.layout, tMINOR) {
		return nil, errors.New("no 'MINOR' in the layout")
	}
	ncv := cv.clone()
	ncv.minor++
	return ncv, nil
}

func (cv *Calver) Micro() (*Calver, error) {
	if !contains(cv.layout, tMICRO) {
		return nil, errors.New("no 'MICRO' in the layout.")
	}
	ncv := cv.clone()
	ncv.micro++
	return ncv, nil
}

func (cv *Calver) Modifier(m string) *Calver {
	ncv := cv.clone()
	ncv.modifier = m
	return ncv
}

func (cv *Calver) clone() *Calver {
	return &Calver{
		major:    cv.major,
		minor:    cv.minor,
		micro:    cv.micro,
		modifier: cv.modifier,
		ts:       cv.ts,
		loc:      cv.loc,
		layout:   cv.layout,
	}
}

func contains(layout []token, t token) bool {
	for _, tt := range layout {
		if tt.Token() == t.Token() {
			return true
		}
	}
	return false
}
