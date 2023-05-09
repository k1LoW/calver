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
		{"0Y.0W.MICROMODIFIER", "02.06.3-dev"},
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
			cv.modifier = "-dev"

			got := cv.String()
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
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
		layout  string
		now     time.Time
		want    string
		wantErr bool
	}{
		{"0Y.0M.MICRO", testtime, "02.02.4", false},
		{"0Y.0W.MICROMODIFIER", testtime, "02.06.4-dev", false},
		{"0Y.0M.MINOR", testtime, "02.02.3", false},
		{"YYYY.0M.0D", testtime, "", true},
		{"YYYY.0M.0D", testtime.AddDate(0, 0, 1), "2002.02.05", false},
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
			cv.modifier = "-dev"

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
	now := time.Now().UTC()
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
				ts:    time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			},
			false,
		},
		{
			"MAJOR.MINOR.MICROMODIFIER", "1.2.3",
			&Calver{
				major: 1,
				minor: 2,
				micro: 3,
				ts:    time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
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
