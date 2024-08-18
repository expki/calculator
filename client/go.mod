module calculator

go 1.23.0

require (
	github.com/coder/websocket v1.8.12
	github.com/expki/calculator/lib v0.0.0-20240818121954-bab18c80aca2
)

require github.com/klauspost/compress v1.17.9 // indirect

replace github.com/expki/calculator/lib => ../lib
