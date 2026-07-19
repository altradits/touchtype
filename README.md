# Touch Typing Master

Touch Typing Master is a local, offline-ready typing game. It teaches the home row, letters, numbers, punctuation, and Go-oriented code through 30 progressive lessons. The falling prompt always matches the exact next lesson key, so learners can build speed without guessing.

## What you need

- A laptop or desktop with a physical keyboard and a current web browser.
- Go 1.26 or newer to run from source.
- Git to clone the repository (or download its ZIP instead).
- Docker Desktop or Docker Engine is optional; it is the easiest way to run the packaged app.

The app runs entirely on your machine. After the initial install, it does not need an internet connection.

## Choose a setup method

| Method | Best for | Needs |
| --- | --- | --- |
| Go from source | Learners and developers who want to edit or test the app | Go and Git |
| Docker | A clean, repeatable local install | Docker |
| Pre-built binary | A release download supplied by the project maintainer | The matching binary for your operating system |

## 1. Get the source code

Open Terminal, PowerShell, or a Linux shell and run the following. Replace `<repository-url>` with this project’s Git URL.

```bash
git clone <repository-url> touchtype
cd touchtype
```

If you downloaded a ZIP file instead, extract it, open a terminal in the extracted `touchtype` folder, and continue at the relevant operating-system section below.

## 2. Run with Go (recommended for development)

After installation, use the same commands on every operating system:

```bash
go version
go run .
```

Open [http://localhost:8080](http://localhost:8080). Keep the terminal open while playing; press `Ctrl+C` there to stop the app.

To use a different port:

```bash
# macOS, Linux, or ChromeOS Linux terminal
PORT=3000 go run .

# Windows PowerShell
$env:PORT=3000; go run .
```

Then open `http://localhost:3000`.

### Windows 10 or 11

1. Install [Git for Windows](https://git-scm.com/download/win). Accept the option that adds Git to the command line.
2. Install Go from [go.dev/dl](https://go.dev/dl/) using the Windows installer. Close and reopen PowerShell after it finishes.
3. Open **PowerShell** and check the installation:

   ```powershell
   git --version
   go version
   ```

4. Clone or extract the project, then run it:

   ```powershell
   git clone <repository-url> touchtype
   Set-Location touchtype
   go run .
   ```

5. Visit [http://localhost:8080](http://localhost:8080) in Edge, Chrome, or Firefox.

Windows alternative: install Go with `winget install GoLang.Go` if Winget is available, then reopen PowerShell and run `go version`.

### macOS (Apple Silicon and Intel)

1. Download and install the macOS package for your chip from [go.dev/dl](https://go.dev/dl/): use **arm64** for Apple Silicon (M1/M2/M3/M4) and **amd64** for Intel.
2. Install [Git for macOS](https://git-scm.com/download/mac) if `git --version` is not already available. This guide does not require Xcode or Xcode Command Line Tools.
3. Open **Terminal**, verify Go and Git, then clone and run the app:

   ```bash
   git --version
   go version
   git clone <repository-url> touchtype
   cd touchtype
   go run .
   ```

4. Open [http://localhost:8080](http://localhost:8080) in Safari, Chrome, or Firefox.

### Linux (Ubuntu, Debian, Fedora, Arch, and similar)

1. If you have `sudo` access, install Git with your distribution package manager:

   ```bash
   # Ubuntu/Debian
   sudo apt update && sudo apt install -y git

   # Fedora
   sudo dnf install -y git

   # Arch
   sudo pacman -S git
   ```

2. If you do **not** have `sudo`, ask your administrator to provide Git and Go, or install them in your home directory. Download the current Linux archive for your CPU from [go.dev/dl](https://go.dev/dl/), extract it under a directory you own, and add its `bin` directory to your shell path:

   ```bash
   mkdir -p "$HOME/opt"
   tar -C "$HOME/opt" -xzf go<version>.linux-<arch>.tar.gz
   export PATH="$HOME/opt/go/bin:$PATH"
   printf '\nexport PATH="$HOME/opt/go/bin:$PATH"\n' >> "$HOME/.profile"
   ```

   Replace `<version>` and `<arch>` with the downloaded file’s values, then start a new terminal.

3. Install Go 1.26 or newer using your distribution’s current Go package, or use the same official archive method above. Verify it before continuing:

   ```bash
   go version
   ```

4. Clone and run:

   ```bash
   git clone <repository-url> touchtype
   cd touchtype
   go run .
   ```

5. Open [http://localhost:8080](http://localhost:8080). On a remote Linux machine, open the port only on a trusted network and browse to `http://<server-ip>:8080`.

### ChromeOS laptop

ChromeOS runs the app through its Linux development environment.

1. Enable **Linux development environment** in **Settings → About ChromeOS → Developers**. Open the new Terminal app.
2. Install Git and Go:

   ```bash
   sudo apt update
   sudo apt install -y git golang
   go version
   ```

3. If the installed Go version is older than 1.26, install the current Linux archive from [go.dev/dl](https://go.dev/dl/) instead.
4. Clone and run the project:

   ```bash
   git clone <repository-url> touchtype
   cd touchtype
   go run .
   ```

5. Open [http://localhost:8080](http://localhost:8080) in Chrome. ChromeOS normally forwards local Linux ports automatically.

## 3. Run with Docker

Docker avoids installing Go. Install [Docker Desktop](https://docs.docker.com/desktop/) on Windows or macOS. On Linux, choose either Docker Engine managed with `sudo` or [Rootless Docker](https://docs.docker.com/engine/security/rootless/), which runs entirely as your normal user.

### Windows and macOS Docker Desktop

1. Install [Docker Desktop](https://docs.docker.com/desktop/) and start it.
2. In PowerShell (Windows) or Terminal (macOS), confirm Docker is available:

```bash
docker version
```

### Linux with sudo-managed Docker Engine

1. Install Docker Engine using your Linux distribution’s Docker instructions. This normally requires an administrator.
2. Confirm that Docker is running:

   ```bash
   sudo docker version
   ```

3. Prefix the Docker commands below with `sudo`, unless your administrator has added your account to the `docker` group. Membership in that group grants root-equivalent access; only accept it on a machine you trust.

### Linux with Rootless Docker (no sudo required to run containers)

1. Follow the [Rootless Docker setup guide](https://docs.docker.com/engine/security/rootless/) once. The initial setup may require an administrator to configure user namespaces; if it does, ask them to perform only that prerequisite.
2. Start Docker’s user service and set the socket for your current shell:

   ```bash
   systemctl --user start docker
   export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/docker.sock
   printf '\nexport DOCKER_HOST=unix://$XDG_RUNTIME_DIR/docker.sock\n' >> "$HOME/.profile"
   ```

   If `systemctl --user` is unavailable, use the rootless setup tool’s instructions to start `dockerd-rootless.sh` in your user session.
3. Confirm that no elevated access is needed:

   ```bash
   docker version
   ```

### Build and run the app

From the project folder, build and run the image. Use `sudo docker` instead of `docker` only for a sudo-managed Linux installation that requires it:

```bash
docker build -t touchtype-master:latest .
docker run --rm --name touchtype -p 8080:8080 touchtype-master:latest
```

Open [http://localhost:8080](http://localhost:8080). Stop the container with `Ctrl+C` in its terminal. The image includes a health check at [http://localhost:8080/healthz](http://localhost:8080/healthz).

## Publish OS-specific download buttons

The website includes download buttons for Windows (Intel/AMD and ARM), macOS (Intel and Apple Silicon), and Linux (Intel/AMD and ARM). Build their exact assets with:

```bash
make downloads
```

Upload every file from `dist/downloads/` to a public HTTPS folder, CDN, or release-assets location. Configure the deployed app with the URL of that folder, without a trailing file name:

```bash
# Docker example
docker run --rm -p 8080:8080 \
  -e DOWNLOAD_BASE_URL=https://downloads.example.com/touchtype \
  touchtype-master:latest
```

For a non-Docker deployment, set `DOWNLOAD_BASE_URL` in the service environment before starting the app. The browser automatically highlights the likely download for the visitor’s operating system; visitors can choose another architecture when needed. If this variable is absent or invalid, buttons stay safely disabled instead of linking to a broken file.

For Docker Desktop on Windows, allow its WSL 2 integration if prompted. On a ChromeOS laptop, Docker is optional and may require Linux virtualisation support; the Go method above is usually simpler.

## 4. Build or use a binary

Build a binary for the current computer. On macOS, Linux, or ChromeOS with GNU Make installed:

```bash
make build
```

Run it:

```bash
# macOS/Linux/ChromeOS
./build/touchtype-master

# Windows PowerShell
.\build\touchtype-master.exe
```

On Windows PowerShell, where `make` is not normally installed, build directly with Go:

```powershell
New-Item -ItemType Directory -Force build | Out-Null
go build -o .\build\touchtype-master.exe .
.\build\touchtype-master.exe
```

To create all release binaries (Windows amd64, Linux amd64, macOS Intel, and macOS Apple Silicon), run:

```bash
make build-all
```

The outputs are written to `dist/`. macOS may require you to permit a downloaded unsigned binary in **System Settings → Privacy & Security**.

## How to play

1. Click **Start**.
2. Type the highlighted next character in the lesson. The falling prompt is the same character.
3. Keep the prompt from reaching the keyboard to preserve your combo and score. A missed prompt resets the combo and applies a small score penalty.
4. Accuracy is based on typed keys. Finish a lesson to advance automatically.
5. Use **Reset** to begin a new run from lesson one. Background music is off by default; use the music button to opt in.

Later lessons intentionally include spaces, capitals, punctuation, and symbols. Use a physical keyboard and the Shift key where needed; the on-screen keyboard is a visual guide, not a replacement for every modified key.

## Development and checks

```bash
make test
make vet
make build
```

On Windows without GNU Make, use the equivalent Go commands:

```powershell
go test ./...
go vet ./...
go build .
```

For hot reload, install Air and start development mode:

```bash
go install github.com/air-verse/air@latest
make dev
```

## Troubleshooting

| Problem | Fix |
| --- | --- |
| `go: command not found` | Reinstall Go, reopen the terminal, and run `go version`. |
| Browser cannot open the page | Confirm the terminal still shows the server running; then visit the exact port printed by the app. |
| Port 8080 is already in use | Run with `PORT=3000 go run .` on macOS/Linux/ChromeOS, or `$env:PORT=3000; go run .` in PowerShell. |
| Docker cannot connect | Start Docker Desktop (Windows/macOS) or the Docker service (Linux), then run `docker version`. |
| Rootless Docker cannot connect | Run `export DOCKER_HOST=unix://$XDG_RUNTIME_DIR/docker.sock`, then confirm `systemctl --user status docker`. |
| A symbol is marked wrong | Check keyboard layout and use Shift when the lesson displays an uppercase or shifted symbol. |

## License

Distributed under the [MIT License](LICENSE).

