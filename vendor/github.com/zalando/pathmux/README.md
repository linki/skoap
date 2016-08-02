# pathmux

Package pathmux is a fork of
[httptreemux](https://github.com/dimfeld/httptreemux). This fork exposes
the interface of the internal tree lookup implementation of the original
httptreemux package, and has the http related wrapper code stripped off.
When the use case is about only method+path matching of http requests,
it is recommended to use the original httptreemux package.

Please, see the documentation at:
https://godoc.org/github.com/zalando/pathmux
