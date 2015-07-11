# gohome [![Build Status](https://drone.io/github.com/caarlos0/gohome/status.png)](https://drone.io/github.com/caarlos0/gohome/latest) [![Coverage Status](https://coveralls.io/repos/caarlos0/gohome/badge.svg?branch=master&service=github)](https://coveralls.io/github/caarlos0/gohome?branch=master)

Easily get cache and config folders for your app according to each OS spec

## Usage

```go
import "github.com/caarlos0/gohome"

const appName = "my-app"

func main()  {
  config := gohome.Config(appName) // gives you the right config folder for the current OS
  cache  := gohome.Cache(appName)  // gives you the right cache folder for the current OS
  // ...
}
```
