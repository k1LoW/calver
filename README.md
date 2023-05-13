# calver [![CI](https://github.com/k1LoW/calver/actions/workflows/ci.yml/badge.svg)](https://github.com/k1LoW/calver/actions/workflows/ci.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/k1LoW/calver.svg)](https://pkg.go.dev/github.com/k1LoW/calver) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/calver/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/calver/ratio.svg)

`calver` is a package/tool provides the ability to work with [Calendar Versioning](https://calver.org/) in Go.

## Usage

### As a package

``` go
package main

import (
	"fmt"

	"github.com/k1LoW/calver"
)

func main() {
	cv, _ := calver.Parse("YY.0M.MICRO", "23.05.1")
	fmt.Println(cv) // Output: 23.05.1
	ncv, _ := cv.Next()
	fmt.Println(ncv) // Output: 23.05.2
}
```

### As a tool

``` console
$ date
Tue May  9 13:04:09 UTC 2023
$ calver --layout YY.0M.MICRO
23.5.0
$ calver --layout YY.0M.MICRO --next
23.5.0
$ calver --layout YY.0M.MICRO | calver --layout YY.0M.MICRO --next
23.5.1
```

#### Example: Get latest version and generate next version

``` console
$ date
Tue May  9 13:04:09 UTC 2023
$ gh release list | cut -f 1
v1.1.2
v1.1.1
2023.05.1
2023.05.0
v1.1.0
v1.0.1
2023.03.1
2023.03.0
v1.0.0
2023.02.0
v0.0.4
v0.0.3
v0.0.2
v0.0.1
$ gh release list | cut -f 1 | calver --layout YY.0M.MICRO
2023.05.1
$ gh release list | cut -f 1 | calver --layout YY.0M.MICRO --next
2023.05.2
```

## Install

### As a package

```console
$ go get github.com/k1LoW/calver
```

### As a tool

**deb:**

``` console
$ export CALVER_VERSION=X.X.X
$ curl -o calver.deb -L https://github.com/k1LoW/calver/releases/download/v$CALVER_VERSION/calver_$CALVER_VERSION-1_amd64.deb
$ dpkg -i calver.deb
```

**RPM:**

``` console
$ export CALVER_VERSION=X.X.X
$ yum install https://github.com/k1LoW/calver/releases/download/v$CALVER_VERSION/calver_$CALVER_VERSION-1_amd64.rpm
```

**apk:**

``` console
$ export CALVER_VERSION=X.X.X
$ curl -o calver.apk -L https://github.com/k1LoW/calver/releases/download/v$CALVER_VERSION/calver_$CALVER_VERSION-1_amd64.apk
$ apk add calver.apk
```

**homebrew tap:**

```console
$ brew install k1LoW/tap/calver
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/calver/releases)

**go install:**

```console
$ go install github.com/k1LoW/calver/cmd/calver@latest
```
