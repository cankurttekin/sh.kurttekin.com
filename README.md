# SSH TUI Portfolio

An interactive terminal UI portfolio accessible via SSH.

## Features

- Clean terminal UI interface built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Multiple sections (about, experience, projects, links)
- Interactive navigation between sections
- Clickable links that open in your browser
- Accessible remotely via SSH

## Installation

### Prerequisites

- Go 1.18 or later

### Building from source

1. Clone the repository:

```bash
git clone https://github.com/cankurttekin/sh.kurttekin.com.git
cd sh.kurttekin.com
```

2. Build the project:

```bash
go build -o tuiserver ./cmd/tuiserver
```

## Usage

### Running the server

```bash
./tuiserver
```

By default, the server runs on port 2222. You can customize the port with the `-addr` flag:

```bash
./tuiserver -addr :3333
```

### Connecting to the server

```bash
ssh localhost -p 2222
```

**Note**: Be sure to use the `-t` flag to allocate a pseudo-terminal:

```bash
ssh -t localhost -p 2222
```

### Navigation

- Navigate between sections with `j`/`k` or arrow keys
- When in a section with links, press `TAB` to enter link selection mode
- Navigate between links with `j`/`k` or arrow keys
- Press `ENTER` to open a selected link in your browser
- Press `q` to quit

## Project Structure

```
.
├── cmd/
│   └── tuiserver/         # Main application entry point
├── internal/
│   ├── models/            # Data models
│   ├── server/            # SSH server implementation
│   └── tui/               # Terminal UI components
├── pkg/
│   └── browser/           # Browser utilities
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
└── README.md              # Project documentation
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgements

- [Charm](https://charm.sh/) for their amazing terminal UI libraries
- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lip Gloss](https://github.com/charmbracelet/lipgloss), and [SSH](https://github.com/charmbracelet/ssh) 