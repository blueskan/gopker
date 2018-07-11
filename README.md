### Gopker is Simple API Wrapper for Docker API, it's simplify your tests with capabilities like NewContainer, StartContainer and StopContainer.

Simple usage example:

All operations are blocking.

*Create Container with Port and Volume bindings (port|volume bindings are fluent, you can chain it without worry)*

Install Package:
```
go get github.com/blueskan/gopker
```

Import
```
import(
    . "github.com/blueskan/gopker"
)
```

*Start container*

```
containerSetup, err := Container("nginx")

if err != nil {
    panic(err)
}

container, err := containerSetup.
	Port("8080", "80").
	Volume("/var/www").
	Start()
```

*Stop container*

```
container.Stop()
```

*Util: List Containers*

```
containers, err := gopker.Containers()

// just do whatever you want
```
