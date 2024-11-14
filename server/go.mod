module calculator

go 1.23.3

require (
	github.com/expki/calculator/lib v0.0.0-20240825055908-38aacee2dcb1
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.3
	github.com/klauspost/compress v1.17.9
	go.uber.org/zap v1.27.0
	golang.org/x/net v0.25.0
)

require (
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/text v0.15.0 // indirect
)

replace github.com/expki/calculator/lib => ../lib
