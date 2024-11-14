module calculator

go 1.23.3

require (
	github.com/coder/websocket v1.8.12
	github.com/expki/calculator/lib v0.0.0-20241114094423-5dd7244a93a0
)

require github.com/klauspost/compress v1.17.9 // indirect

replace github.com/expki/calculator/lib => ../lib
