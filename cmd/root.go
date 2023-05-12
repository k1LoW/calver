/*
Copyright © 2023 Ken'ichiro Oyama <k1lowxb@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/k1LoW/calver"
	"github.com/k1LoW/calver/version"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var (
	layout string
	next   bool
	major  bool
	minor  bool
	micro  bool
)

var rootCmd = &cobra.Command{
	Use:          "calver",
	Short:        "calver is a tool for manipulating calender versioning",
	Long:         `calver is a tool for manipulating calender versioning.`,
	SilenceUsage: true,
	Version:      version.Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		cv, err := calver.New(layout)
		if err != nil {
			return err
		}
		var versions []string
		switch {
		case len(args) > 0:
			versions = args
		case !isatty.IsTerminal(os.Stdin.Fd()):
			stdin, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			lines := strings.Split(strings.Trim(string(stdin), " \n"), "\n")
			for _, l := range lines {
				versions = append(versions, strings.Split(l, " ")...)
			}
		}

		if len(versions) == 0 && (next || major || minor || micro) {
			return errors.New("--next, --major, --minor, and --micro require a version string to be parsed")
		}

		var errs error
		if len(versions) > 0 {
			cvs := calver.Calvers{}
			for _, v := range versions {
				ccv, err := cv.Parse(v)
				if err != nil {
					errs = errors.Join(errs, err)
					continue
				}
				cvs = append(cvs, ccv)
			}
			cv, err = cvs.Latest()
			if err != nil {
				errs = errors.Join(err, errs)
				return errs
			}
			switch {
			case next:
				cv, err = cv.Next()
				if err != nil {
					return err
				}
			case major:
				cv, err = cv.Major()
				if err != nil {
					return err
				}
			case minor:
				cv, err = cv.Minor()
				if err != nil {
					return err
				}
			case micro:
				cv, err = cv.Micro()
				if err != nil {
					return err
				}
			}
		}

		fmt.Println(cv.String())
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&layout, "layout", "l", "YY.0M.MICRO", "version layout")
	rootCmd.Flags().BoolVarP(&next, "next", "n", false, "show next version of parsed version")
	rootCmd.Flags().BoolVarP(&major, "major", "", false, "show next major version of parsed version")
	rootCmd.Flags().BoolVarP(&minor, "minor", "", false, "show next minor version of parsed version")
	rootCmd.Flags().BoolVarP(&micro, "micro", "", false, "show next micro version of parsed version")
}
