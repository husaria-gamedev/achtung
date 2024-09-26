module github.com/jsierp/achtung

go 1.22.6

require (
	github.com/gorilla/websocket v1.5.3
	github.com/tfriedel6/canvas v0.12.1
)

require (
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/veandco/go-sdl2 v0.4.40 // indirect
	golang.org/x/image v0.19.0 // indirect
)

replace github.com/tfriedel6/canvas => github.com/jsierp/go-canvas v0.12.1-fix
