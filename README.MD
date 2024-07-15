<div align="center">
  <h3 align="center">🪨 📄 ✂️</h3>
  <h3 align="center">Rock Paper Scissors Online</h3>
  <p align="center">
    A simple game where you can compete in duels with your friends.
  </p>
</div>

## About The Project

The game has a server and a client, both written in Go, and uses pure TCP sockets for client-server communication. This game was developed for practicing network programming, particularly TCP sockets, handling multiple clients, game state synchronization, and managing user input over a network.

## Getting Started

### Prerequisites

To build manually, you need to have Go 1.22 and [Wails](https://github.com/wailsapp/wails) installed on your computer.

### Build&Run server

```sh
  make server
```

### Build&Run client

```sh
  make client
```
