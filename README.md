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
