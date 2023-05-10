package calver

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestTokenCal(t *testing.T) {
	tests := []struct {
		token tokenCal
		want  string
	}{
		{tYYYY, "2002"},
		{tYY, "2"},
		{t0Y, "02"},
		{tMM, "2"},
		{t0M, "02"},
		{tWW, "6"},
		{t0W, "06"},
		{tDD, "4"},
		{t0D, "04"},
	}
	for _, tt := range tests {
		t.Run(tt.token.token(), func(t *testing.T) {
			got := tt.token.timeToString(testtime)
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestTokentrimPrefix(t *testing.T) {
	tests := []struct {
		token      token
		value      string
		wantPrefix string
		wantTrimed string
		wantErr    bool
	}{
		{tYYYY, "2033.12.05", "2033", ".12.05", false},
		{tYYYY, "2033", "2033", "", false},
		{tYYYY, "203", "", "", true},
		{tYY, "3.12.05", "3", ".12.05", false},
		{tYY, "12.12.05", "12", ".12.05", false},
		{t0Y, "03.12.05", "03", ".12.05", false},
		{tMAJOR, "3.7.5", "3", ".7.5", false},
		{tMINOR, "7.5", "7", ".5", false},
		{tMICRO, "5-dev", "5", "-dev", false},
		{tMODIFIER, "-dev", "-dev", "", false},
		{tMICRO, "dev", "", "", true},
		{newTokenSep("."), ".5", ".", "5", false},
		{newTokenSep("."), "3.5", "", "", true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s/%s", tt.token.token(), tt.value), func(t *testing.T) {
			gotPrefix, gotTrimed, gotErr := tt.token.trimPrefix(tt.value)
			if gotPrefix != tt.wantPrefix {
				t.Errorf("got %v\nwant %v", gotPrefix, tt.wantPrefix)
			}
			if gotTrimed != tt.wantTrimed {
				t.Errorf("got %v\nwant %v", gotTrimed, tt.wantTrimed)
			}
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("got %v\nwant %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestTokenizeLayout(t *testing.T) {
	tests := []struct {
		layout  string
		want    []token
		wantErr bool
	}{
		{"YY", []token{tYY}, false},
		{"YYYY", []token{tYYYY}, false},
		{"vYY", []token{newTokenSep("v"), tYY}, false},
		{"YYv", []token{tYY, newTokenSep("v")}, false},
		{"YY.0D", []token{tYY, newTokenSep("."), t0D}, false},
		{"YY0D", []token{tYY, t0D}, false},
		{"YY.0D.MICRO", []token{tYY, newTokenSep("."), t0D, newTokenSep("."), tMICRO}, false},
		{"YYYY.0M.MICRO", []token{tYYYY, newTokenSep("."), t0M, newTokenSep("."), tMICRO}, false},
		{"YYY", []token{tYY, newTokenSep("Y")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.layout, func(t *testing.T) {
			got, err := tokenizeLayout(tt.layout)
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
				cmp.AllowUnexported(tokenCal{}, tokenVer{}, tokenSep{}),
				cmpopts.IgnoreFields(tokenCal{}, "timeToString"),
				cmpopts.IgnoreFields(tokenVer{}, "verToString"),
			}
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("%s", diff)
			}
		})
	}
}
