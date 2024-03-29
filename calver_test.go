package calver

import (
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
