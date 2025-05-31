module calculator

go 1.24.3

require (
	github.com/expki/calculator/lib v0.0.0-20241202065501-3786743a7efc
	github.com/gorilla/websocket v1.5.3
	github.com/klauspost/compress v1.18.0
	go.uber.org/zap v1.27.0
	golang.org/x/net v0.40.0
)

require (
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/text v0.25.0 // indirect
)

replace github.com/expki/calculator/lib => ../lib
