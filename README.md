### Gopker is Simple API Wrapper for Docker API, it's simplify your tests with capabilities like NewContainer, RunContainer and StopContainer.

Simple usage example:

```
Container("nginx").
	Port("8080", "80").
	Volume("/var/www").
	Start()
```

Todos:
    - Add Tests
    - List Capabilities
    - Fixed Panics
