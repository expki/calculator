module calculator

go 1.24.3

require (
	github.com/coder/websocket v1.8.13
	github.com/expki/calculator/lib v0.0.0-20241202065501-3786743a7efc
)

require github.com/klauspost/compress v1.18.0 // indirect

replace github.com/expki/calculator/lib => ../lib
