# Multiplayer Calculator

Multiplayer Calculator is a real-time collaborative calculator where multiple users can interact simultaneously, see each other's cursors, and share a single calculator interface.

## Demo: [https://calculator.vdh.dev/](https://calculator.vdh.dev/)

## Technologies Used

- **WebSockets**  
  Enables real-time communication between users by broadcasting calculator state and cursor positions instantly.

- **Web Workers**
  Enables game to use multiple cores seperating rendering and logic pipelines. This enables game logic and rendering to not affect one another and is displayed in game with GPU and CPU load counter. 

- **Golang Wasm**  
  The calculator client business logic is compiled from Go to WebAssembly for high-performance execution directly in the browser.
  I am never doing this again for a game, its a terrible idea where I ended up in split brain logic where half is in Go and half in Typescript.

- **Canvas**  
  Used to render the calculator UI and display player cursors, allowing for smooth, interactive graphics without a game engine or the complexity of creating OpenGL.

- **Multilingual Stateless Binary Encoding**
  This project implements a stateless binary encoder/decoder that, on average, produces output 35% the size of equivalent JSON encoding. On top of this it supports an optional flag to Zstd compress the data for an even smaller footprint. I based the encoding on ASN.1, but it specifically targets types available in both Go and Typescript enabling cross platform communication.

## Development

This project is a monorepo with three main components:

1. Lib
  The `lib/` directory contains shared type definitions and a stateless binary encoder/decoder used by both the client and server.

2. Server
  The `server/` directory implements the Go web server that clients connect to. It manages shared state, including the calculator and player cursor locations.
```bash
go run .
```

3. Client
  The `client/` directory implements the frontend, including the TypeScript Canvas rendering and Go WebAssembly logic. It renders shared state from the server and submits the client’s cursor position back.
```bash
npm i
npm run dev
```

- Build for production
```bash
./build.sh
```
