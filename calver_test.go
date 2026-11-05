package calver

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var testtime = time.Date(2002, 2, 4, 0, 0, 0, 0, time.UTC)

func TestString(t *testing.T) {
	tests := []struct {
		layout string
		want   string
	}{
		{"YYYY.0M.0D", "2002.02.04"},
		{"0Y.0M.MICRO", "02.02.3"},
		{"0Y.0W.MICRO-MODIFIER", "02.06.3-dev"},
		{"MAJOR.MINOR.MICRO", "1.2.3"},
	}
	for _, tt := range tests {
		t.Run(tt.layout, func(t *testing.T) {
			cv, err := NewWithTime(tt.layout, testtime)
			if err != nil {
				t.Error(err)
				return
			}
			cv.major = 1
			cv.minor = 2
			cv.micro = 3
			cv.modifier = "dev"

			got := cv.String()
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestLayout(t *testing.T) {
	tests := []struct {
		layout string
	}{
		{"YYYY.0M.0D"},
		{"0Y.0M.MICRO"},
		{"0Y.0W.MICROMODIFIER"},
		{"MAJOR.MINOR.MICRO"},
	}
	for _, tt := range tests {
		t.Run(tt.layout, func(t *testing.T) {
			cv, err := New(tt.layout)
			if err != nil {
				t.Error(err)
				return
			}
			got := cv.Layout()
			if got != tt.layout {
				t.Errorf("got %v\nwant %v", got, tt.layout)
			}
		})
	}
}

func TestIn(t *testing.T) {
	cv, err := NewWithTime("YYYY.0M.0D.MICRO", testtime)
	if err != nil {
		t.Error(err)
	}
	loc := time.FixedZone("UTC-2", -2*60*60)
	if cv.String() == cv.In(loc).String() {
		t.Errorf("got %v\n", cv.In(loc).String())
	}
}

func TestNext(t *testing.T) {
	cv, err := New("YYYY.0M.0D.MICRO")
	if err != nil {
		t.Error(err)
	}
	ncv, err := cv.Next()
	if err != nil {
		t.Error(err)
	}
	if cv.String() == ncv.String() {
		t.Errorf("got %v\n", ncv.String())
	}
}

func TestNextWithTime(t *testing.T) {
	tests := []struct {
		layout   string
		now      time.Time
		modifier string
		want     string
		wantErr  bool
	}{
		{"0Y.0M.MICRO", testtime, "", "02.02.4", false},
		{"0Y.0W.MICRO-MODIFIER", testtime, "", "02.06.4", false},
		{"0Y.0M.MINOR", testtime, "", "02.02.3", false},
		{"YYYY.0M.0D", testtime, "", "", true},
		{"YYYY.0M.0D", testtime.AddDate(0, 0, 1), "", "2002.02.05", false},
		{"0Y.0W.MICRO-MODIFIER", testtime, "dev", "02.06.3", false},
		{"0Y.0M.MICRO", testtime.AddDate(0, 1, 0), "", "02.03", false},
		{"MAJOR.0Y.0M", testtime.AddDate(0, 1, 0), "", "1.02.03", false},
	}
	for _, tt := range tests {
		t.Run(tt.layout, func(t *testing.T) {
			cv, err := New(tt.layout)
			if err != nil {
				t.Error(err)
				return
			}
			cv.ts = testtime
			cv.major = 1
			cv.minor = 2
			cv.micro = 3
			cv.modifier = tt.modifier
			cv.trimSuffix = true

			got, err := cv.NextWithTime(tt.now)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("got error: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Error("want error")
			}
			if got.String() != tt.want {
				t.Errorf("got %v\nwant %v", got.String(), tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		layout  string
		value   string
		want    *Calver
		wantErr bool
	}{
		{
			"YYYY.0M.0D", "2012.12.03",
			&Calver{
				ts: time.Date(2012, time.Month(12), 3, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"YYYY.0W.MICRO", "2002.06.2",
			&Calver{
				micro: 2,
				ts:    time.Date(2002, time.Month(2), 4, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"YYYY.0M", "2012.12.03",
			nil,
			true,
		},
		{
			"MAJOR.MINOR.MICRO", "1.2.3",
			&Calver{
				major: 1,
				minor: 2,
				micro: 3,
				ts:    time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"MAJOR.MINOR.MICROMODIFIER", "1.2.3",
			&Calver{
				major: 1,
				minor: 2,
				micro: 3,
				ts:    time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"YYYY.MMDD.MICRO", "2026.123.0",
			&Calver{
				micro: 0,
				ts:    time.Date(2026, time.Month(12), 3, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"YYYY.MM0D.MICRO", "2026.123.0",
			&Calver{
				micro: 0,
				ts:    time.Date(2026, time.Month(1), 23, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s/%s", tt.layout, tt.value), func(t *testing.T) {
			got, err := Parse(tt.layout, tt.value)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("got error: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Error("want error")
			}
			opts := []cmp.Option{
				cmp.AllowUnexported(Calver{}),
				cmpopts.IgnoreFields(Calver{}, "layout"),
				cmpopts.IgnoreFields(Calver{}, "loc"),
			}
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}

func TestWeekVersions(t *testing.T) {
	tests := []struct {
		date     string
		current  string
		wantNew  string
		wantNext string
	}{
		{"2001-12-31", "2001.52.8", "2002.01.0", "2002.01.0"},
		{"2025-12-29", "2025.52.8", "2026.01.0", "2026.01.0"},
		{"2025-12-30", "2026.01.0", "2026.01.0", "2026.01.1"},
		{"2026-01-02", "2026.01.1", "2026.01.0", "2026.01.2"},
		{"2026-01-05", "2026.01.2", "2026.02.0", "2026.02.0"},
		{"2026-12-28", "2026.52.8", "2026.53.0", "2026.53.0"},
		{"2027-01-01", "2026.53.0", "2026.53.0", "2026.53.1"},
		{"2027-01-02", "2026.53.1", "2026.53.0", "2026.53.2"},
		{"2027-01-04", "2026.53.2", "2027.01.0", "2027.01.0"},
	}
	for _, yearToken := range []string{"YYYY", "YY", "0Y"} {
		for _, weekToken := range []string{"WW", "0W"} {
			layout := yearToken + "." + weekToken + ".MICRO"
			version := func(value string) string {
				var year, week, micro int
				if _, err := fmt.Sscanf(value, "%d.%d.%d", &year, &week, &micro); err != nil {
					t.Fatal(err)
				}
				yearString := fmt.Sprintf("%04d", year)
				if yearToken == "YY" {
					yearString = fmt.Sprintf("%d", year-2000)
				}
				if yearToken == "0Y" {
					year -= 2000
					yearString = fmt.Sprintf("%02d", year)
				}
				if weekToken == "0W" {
					return fmt.Sprintf("%s.%02d.%d", yearString, week, micro)
				}
				return fmt.Sprintf("%s.%d.%d", yearString, week, micro)
			}
			for _, tt := range tests {
				t.Run(layout+"/"+tt.date, func(t *testing.T) {
					now, err := time.Parse("2006-01-02", tt.date)
					if err != nil {
						t.Fatal(err)
					}
					generated, err := NewWithTime(layout, now)
					if err != nil {
						t.Fatal(err)
					}
					if got, want := generated.String(), version(tt.wantNew); got != want {
						t.Fatalf("generated = %s, want %s", got, want)
					}
					parsed, err := Parse(layout, generated.String())
					if err != nil {
						t.Fatal(err)
					}
					if got := parsed.String(); got != generated.String() {
						t.Fatalf("round trip = %s, want %s", got, generated.String())
					}
					current, err := Parse(layout, version(tt.current))
					if err != nil {
						t.Fatal(err)
					}
					next, err := current.NextWithTime(now)
					if err != nil {
						t.Fatal(err)
					}
					if got, want := next.String(), version(tt.wantNext); got != want {
						t.Fatalf("next = %s, want %s", got, want)
					}
					next, err = Parse(layout, next.String())
					if err != nil {
						t.Fatal(err)
					}
					latest, err := (Calvers{current, next}).Latest()
					if err != nil {
						t.Fatal(err)
					}
					if got, want := latest.String(), version(tt.wantNext); got != want {
						t.Fatalf("latest = %s, want %s", got, want)
					}
				})
			}
		}
	}
}

func TestParseInvalidWeek(t *testing.T) {
	for _, value := range []string{"2026.00.0", "2025.53.0", "2026.54.0", "2027.53.0"} {
		t.Run(value, func(t *testing.T) {
			if _, err := Parse("YYYY.0W.MICRO", value); err == nil {
				t.Fatal("want error for invalid ISO week")
			}
		})
	}
}

func TestWeekVersionsWithoutYear(t *testing.T) {
	tests := []struct {
		layout   string
		date     string
		current  string
		wantNew  string
		wantNext string
	}{
		{"WW.MICRO", "2026-12-28", "52.3", "53.0", "53.0"},
		{"WW.MICRO", "2027-01-01", "53.0", "53.0", "53.1"},
		{"WW.MICRO", "2027-01-04", "53.1", "1.0", "1.0"},
		{"0W.MICRO", "2026-12-28", "52.3", "53.0", "53.0"},
		{"0W.MICRO", "2027-01-01", "53.0", "53.0", "53.1"},
		{"0W.MICRO", "2027-01-04", "53.1", "01.0", "01.0"},
	}
	for _, tt := range tests {
		t.Run(tt.layout+"/"+tt.date, func(t *testing.T) {
			now, err := time.Parse("2006-01-02", tt.date)
			if err != nil {
				t.Fatal(err)
			}
			generated, err := NewWithTime(tt.layout, now)
			if err != nil {
				t.Fatal(err)
			}
			if got := generated.String(); got != tt.wantNew {
				t.Fatalf("generated = %s, want %s", got, tt.wantNew)
			}
			parsed, err := Parse(tt.layout, generated.String())
			if err != nil {
				t.Fatal(err)
			}
			if got := parsed.String(); got != tt.wantNew {
				t.Fatalf("round trip = %s, want %s", got, tt.wantNew)
			}
			current, err := Parse(tt.layout, tt.current)
			if err != nil {
				t.Fatal(err)
			}
			next, err := current.NextWithTime(now)
			if err != nil {
				t.Fatal(err)
			}
			if got := next.String(); got != tt.wantNext {
				t.Fatalf("next = %s, want %s", got, tt.wantNext)
			}
		})
	}
}

func TestParseInvalidWeekWithoutYear(t *testing.T) {
	for _, value := range []string{"00.0", "54.0"} {
		t.Run(value, func(t *testing.T) {
			if _, err := Parse("0W.MICRO", value); err == nil {
				t.Fatal("want error for invalid ISO week")
			}
		})
	}
}

func TestWeekVersionIn(t *testing.T) {
	now := time.Date(2025, 12, 29, 0, 0, 0, 0, time.UTC)
	cv, err := NewWithTime("YYYY.0W.MICRO", now)
	if err != nil {
		t.Fatal(err)
	}
	loc := time.FixedZone("UTC-2", -2*60*60)
	cv = cv.In(loc)
	if got := cv.String(); got != "2025.52.0" {
		t.Fatalf("got %s, want 2025.52.0", got)
	}
	parsed, err := cv.Parse(cv.String())
	if err != nil {
		t.Fatal(err)
	}
	next, err := parsed.NextWithTime(now.Add(2 * time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if got := next.String(); got != "2026.01.0" {
		t.Fatalf("got %s, want 2026.01.0", got)
	}
}

func TestCalendarYearVersions(t *testing.T) {
	tests := []struct {
		layout string
		date   string
		want   string
	}{
		{"YYYY.0M.0D", "2025-12-29", "2025.12.29"},
		{"YY.0M.0D", "2025-12-29", "25.12.29"},
		{"0Y.0M.0D", "2027-01-01", "27.01.01"},
		{"YYYY.0M.MICRO", "2027-01-01", "2027.01.0"},
	}
	for _, tt := range tests {
		t.Run(tt.layout+"/"+tt.date, func(t *testing.T) {
			now, err := time.Parse("2006-01-02", tt.date)
			if err != nil {
				t.Fatal(err)
			}
			cv, err := NewWithTime(tt.layout, now)
			if err != nil {
				t.Fatal(err)
			}
			if got := cv.String(); got != tt.want {
				t.Fatalf("got %s, want %s", got, tt.want)
			}
			parsed, err := Parse(tt.layout, cv.String())
			if err != nil {
				t.Fatal(err)
			}
			if got := parsed.String(); got != tt.want {
				t.Fatalf("round trip = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTrimSuffix(t *testing.T) {
	tests := []struct {
		layout     string
		trimSuffix bool
		version    string
		want       string
	}{
		{"YY.0M.MICRO", false, "23.05.0", "23.05.0"},
		{"YY.0M-MODIFIER", false, "23.05-dev", "23.05-dev"},
		{"YY.0M.MICRO", true, "23.05.0", "23.05"},
		{"YY.0M-MODIFIER", true, "23.05-", "23.05"},
		{"YY.0M-MODIFIER", true, "23.05", "23.05"},
		{"YY.0M.MAJOR.MINOR.MICRO", true, "23.05.1", "23.05.1"},
		{"YY.0M.MICRO-MODIFIER", true, "23.05.0-dev", "23.05-dev"},
		{"YY.0M.MICRO-MODIFIER", true, "23.05-dev", "23.05-dev"},
		{"YY.0M.MICROMODIFIER", true, "23.05-dev", "23.05-dev"},
		{"YY.0M.MAJOR.MINOR.MICRO-MODIFIER", true, "23.05-dev", "23.05-dev"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s/%v/%s", tt.layout, tt.trimSuffix, tt.version), func(t *testing.T) {
			cv, err := New(tt.layout)
			if err != nil {
				t.Error(err)
				return
			}
			cv = cv.TrimSuffix(tt.trimSuffix)
			gotcv, err := cv.Parse(tt.version)
			if err != nil {
				t.Error(err)
				return
			}
			got := gotcv.String()
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		})
	}

}
func TestLatest(t *testing.T) {
	tests := []struct {
		layout   string
		versions []string
		want     *Calver
		wantErr  bool
	}{
		{
			"YYYY.0M.0D",
			[]string{"2012.12.03", "2012.12.04"},
			&Calver{
				major: 0,
				minor: 0,
				micro: 0,
				ts:    time.Date(2012, time.Month(12), 4, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"YYYY.0M.MICRO",
			[]string{"2012.12.1", "2012.12.0"},
			&Calver{
				major: 0,
				minor: 0,
				micro: 1,
				ts:    time.Date(2012, time.Month(12), 1, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"YYYY.0M.MICRO",
			[]string{"2012.12.1", "2012.12.20"},
			&Calver{
				major: 0,
				minor: 0,
				micro: 20,
				ts:    time.Date(2012, time.Month(12), 1, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"YYYY.0M.MICROMODIFIER",
			[]string{"2012.12.0", "2012.12.0-dev"},
			&Calver{
				major:    0,
				minor:    0,
				micro:    0,
				modifier: "",
				ts:       time.Date(2012, time.Month(12), 1, 0, 0, 0, 0, time.UTC),
			},
			false,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			cvs := Calvers{}
			cv, err := New(tt.layout)
			if err != nil {
				t.Error(err)
			}
			for _, v := range tt.versions {
				ccv, err := cv.Parse(v)
				if err == nil {
					cvs = append(cvs, ccv)
				}
			}
			got, err := cvs.Latest()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("got error: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Error("want error")
			}
			opts := []cmp.Option{
				cmp.AllowUnexported(Calver{}),
				cmpopts.IgnoreFields(Calver{}, "layout"),
				cmpopts.IgnoreFields(Calver{}, "loc"),
			}
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}

func TestMajor(t *testing.T) {
	tests := []struct {
		layout  string
		wantErr bool
	}{
		{"MAJOR.MINOR.MICRO", false},
		{"YYYY.0M.0D", true},
	}
	for _, tt := range tests {
		t.Run(tt.layout, func(t *testing.T) {
			cv, err := New(tt.layout)
			if err != nil {
				t.Error(err)
				return
			}
			cv.major = 1
			got, err := cv.Major()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("got error: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Error("want error")
			}
			if got.major != 2 {
				t.Errorf("got %v\nwant %v", got.major, 2)
			}
		})
	}
}

func TestModifier(t *testing.T) {
	tests := []struct {
		layout   string
		modifier string
		wantErr  bool
	}{
		{"YYYY.0M.MICRO-MODIFIER", "dev", false},
		{"YYYY.0M.0D", "dev", true},
	}
	for _, tt := range tests {
		t.Run(tt.layout, func(t *testing.T) {
			cv, err := New(tt.layout)
			if err != nil {
				t.Error(err)
				return
			}
			got, err := cv.Modifier(tt.modifier)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("got error: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Error("want error")
			}
			if got.modifier != tt.modifier {
				t.Errorf("got %v\nwant %v", got.modifier, tt.modifier)
			}
		})
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		layout   string
		versions []string
		want     []string
	}{
		{
			"YYYY.0M.MICRO",
			[]string{"2012.12.1", "2012.12.0", "2012.12.2"},
			[]string{"2012.12.2", "2012.12.1", "2012.12.0"},
		},
		{
			"MAJOR.MINOR.MICRO",
			[]string{"1.0.0", "1.1.0", "1.0.1"},
			[]string{"1.1.0", "1.0.1", "1.0.0"},
		},
		{
			"YYYY.0M.MICROMODIFIER",
			[]string{"2012.12.0-dev", "2012.12.0"},
			[]string{"2012.12.0", "2012.12.0-dev"},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			cvs := Calvers{}
			cv, err := New(tt.layout)
			if err != nil {
				t.Error(err)
			}
			for _, v := range tt.versions {
				ccv, err := cv.Parse(v)
				if err == nil {
					cvs = append(cvs, ccv)
				}
			}
			cvs.Sort()
			for i, want := range tt.want {
				if cvs[i].String() != want {
					t.Errorf("got %v\nwant %v", cvs[i].String(), want)
				}
			}
		})
	}
}

func TestLatestEmpty(t *testing.T) {
	cvs := Calvers{}
	_, err := cvs.Latest()
	if err == nil {
		t.Error("want error")
	}
	if !errors.Is(err, ErrNoVersions) {
		t.Errorf("got %v\nwant %v", err, ErrNoVersions)
	}
}

func TestNextWithTimeOlderError(t *testing.T) {
	cv, err := NewWithTime("YYYY.0M.MICRO", testtime)
	if err != nil {
		t.Error(err)
	}
	older := testtime.AddDate(0, 0, -1)
	_, err = cv.NextWithTime(older)
	if err == nil {
		t.Error("want error")
	}
}

func TestMinorError(t *testing.T) {
	cv, err := New("YYYY.0M.0D")
	if err != nil {
		t.Error(err)
	}
	_, err = cv.Minor()
	if err == nil {
		t.Error("want error")
	}
}

func TestMicroError(t *testing.T) {
	cv, err := New("YYYY.0M.0D")
	if err != nil {
		t.Error(err)
	}
	_, err = cv.Micro()
	if err == nil {
		t.Error("want error")
	}
}
