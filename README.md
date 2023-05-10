# calver

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
23.5.1
$ calver --layout YY.0M.MICRO --next | calver --layout YY.0M.MICRO --next
23.5.2
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
