<div align=center>
<img width="371" height="228" alt="matt" src="https://github.com/user-attachments/assets/805c3699-3c2b-4646-bc29-4155ebb512b4" />
<!-- 
  ███▄ ▄███▓ ▄▄▄     ▄▄▄█████▓▄▄▄█████▓
  ▓██▒▀█▀ ██▒▒████▄   ▓  ██▒ ▓▒▓  ██▒ ▓▒
  ▓██    ▓██░▒██  ▀█▄ ▒ ▓██░ ▒░▒ ▓██░ ▒░
  ▒██    ▒██ ░██▄▄▄▄██░ ▓██▓ ░ ░ ▓██▓ ░ 
  ▒██▒   ░██▒ ▓█   ▓██▒ ▒██▒ ░   ▒██▒ ░ 
  ░ ▒░   ░  ░ ▒▒   ▓▒█░ ▒ ░░     ▒ ░░   
  ░  ░      ░  ▒   ▒▒ ░   ░        ░    
  ░      ░     ░   ▒    ░        ░      
         ░         ░  ░
  -->
</div>
<div align=center>
 
[![CI](https://github.com/Chintanpatel24/Matt/actions/workflows/ci.yml/badge.svg)](https://github.com/Chintanpatel24/Matt/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://go.dev)
[![GitHub Stars](https://img.shields.io/github/stars/Chintanpatel24/Matt?style=flat&color=eab308)](https://github.com/Chintanpatel24/Matt)
</div>

- **Matt** is a fast, keyboard-driven, modern terminal file manager featuring a **Matt Black** design aesthetic. Built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea), it offers a **3-column pane system**, **async disk analyzer**, **live fuzzy finder**, **file operations**, **bookmarks**, **command history**, **syntax-highlighted previews**, and a **premium dark UI** with scroll indicators, zebra-striping, and breadcrumb navigation.

<div align=center>
<img width="1000" alt="matttt" src="https://github.com/user-attachments/assets/36131f34-ebc1-4478-a347-75c983a51efd" />
</div>

##  Key Features

### Premium UI
- **Matt Black Aesthetic** - Custom dark palette with onyx backgrounds, charcoal borders, and lite grey accents
- **3-Column Split View** - Navigation tree, explorer, and preview/metadata inspector
- **Breadcrumb Path Bar** - Interactive breadcrumb navigation in the header (`~ ❯ Projects ❯ Matt`)
- **Scroll Indicators** - Visual `▲`/`▼` arrows with item counts when content overflows
- **Zebra-Striping** - Alternating row backgrounds for better readability
- **Selection Accent Bar** - `▌` accent indicator on the active selection
- **Focus Indicators** - `◉` active / `○` inactive pane markers
- **Mode Indicator** - Shows current mode (`NORMAL` / `FILTER` / `COMMAND`)

### File Management
- **Create Files & Folders** - `n` for new file, `N` for new folder
- **Copy / Cut / Paste** - `c` copy, `x` cut, `p` paste
- **Rename** - `m` to move/rename with inline prompt
- **Delete** - `d` with confirmation dialog
- **Symlink Resolution** - Shows symlink targets with `→` notation

### Navigation & Search
- **Live Fuzzy Search** - Press `/` for instant fuzzy filtering
- **Bookmarks** - `B` to bookmark, `b` to jump to saved locations
- **Quick Jump** - `g` to first item, `G` to last item
- **Vim-Style Navigation** - `h/j/k/l` or arrow keys

### Async Disk Analyzer
- **Non-Blocking Analysis** - `Alt+D` runs disk scan in background (no more UI freezing!)
- **Visual Progress Bars** - `██████░░ 75%` percentage breakdown
- **Sort Toggle** - Press `s` to toggle between size and name sorting
- **Full Scrolling** - Navigate large directory analyses with `▲`/`▼` indicators

### Integrated Terminal
- **Shell Command Execution** - Run any command from the bottom pane
- **Command History** - `↑`/`↓` arrow keys cycle through previous commands
- **Smart `cd`** - Syncs all 3 panes when changing directories
- **Config Aliases** - Custom shortcuts like `ll`, `findbig`, `sysinfo`
- **Permission Prompts** - Safe confirmation for destructive operations

### Preview & Inspection
- **Syntax Highlighting** - Powered by Chroma with line numbers
- **ASCII Image Preview** - Renders `.png`, `.jpg`, `.gif` as ASCII art
- **Binary Hex View** - Hex dump for binary files
- **File Metadata Inspector** - Permissions, owner, size, timestamps

---

## Installation

### Quick One-Liner

#### Linux / macOS
```bash
curl -fsSL https://raw.githubusercontent.com/Chintanpatel24/Matt/main/install.sh | bash
```

Or using `wget`:
```bash
wget -qO- https://raw.githubusercontent.com/Chintanpatel24/Matt/main/install.sh | bash
```

#### Windows (PowerShell)
```powershell
iwr -useb https://raw.githubusercontent.com/Chintanpatel24/Matt/main/install.ps1 | iex
```

### Clone & Build
```bash
git clone https://github.com/Chintanpatel24/Matt.git
cd Matt
make build
make install
```

### Updating
```bash
curl -fsSL https://raw.githubusercontent.com/Chintanpatel24/Matt/main/update.sh | bash
```

---

## Keyboard Controls

### Navigation
| Key | Description |
|:---|:---|
| `↑` / `↓` (`k` / `j`) | Navigate files/directories in focused pane |
| `→` / `Enter` (`l`) | Open folder / expand directory |
| `←` (`h`) | Go up to parent directory |
| `Tab` / `Shift+Tab` | Cycle focus between panes & terminal |
| `g` | Jump to first item in list |
| `G` | Jump to last item in list |
| `/` | Live fuzzy search & filter |
| `.` | Toggle hidden files |
| `r` | Refresh directory view |

### File Operations
| Key | Description |
|:---|:---|
| `n` | Create new file (enter name in prompt) |
| `N` | Create new folder (enter name in prompt) |
| `m` | Rename/move selected item |
| `c` | Copy selected item to clipboard |
| `x` | Cut selected item to clipboard |
| `p` | Paste clipboard item into current directory |
| `d` | Delete selected item (with confirmation) |

### Tools & Navigation
| Key | Description |
|:---|:---|
| `Alt+D` | Toggle async Disk Space Analyzer |
| `b` | Open bookmarks list |
| `B` | Bookmark current directory |
| `:` | Focus bottom terminal / run shell commands |
| `↑` / `↓` (in terminal) | Browse command history |
| `s` (in analyzer) | Toggle sort by size/name |
| `Esc` | Close modal / unfocus mode |
| `q` / `Ctrl+C` | Quit Matt |

---

## Configuration

Matt reads configuration from `~/.config/matt/config.json`. Customize themes, aliases, and more:

```json
{
  "show_hidden": false,
  "aliases": {
    "ll": "ls -la",
    "findbig": "find . -type f -size +10M",
    "count": "find . -type f | wc -l",
    "sysinfo": "uname -a && uptime"
  },
  "theme": {
    "bg": "#09090b",
    "bg_surface": "#18181b",
    "bg_zebra": "#0f0f12",
    "border": "#27272a",
    "border_active": "#e4e4e7",
    "text_primary": "#f8fafc",
    "text_muted": "#a1a1aa",
    "directory": "#38bdf8",
    "executable": "#34d399",
    "selection": "#27272a",
    "accent": "#e4e4e7",
    "error": "#f87171",
    "warning": "#fbbf24",
    "success": "#4ade80"
  }
}
```

### Data Files
| File | Purpose |
|:---|:---|
| `~/.config/matt/config.json` | Main configuration & theme |
| `~/.config/matt/bookmarks.json` | Saved directory bookmarks |
| `~/.config/matt/history.json` | Command history (max 100 entries) |

<!--
## Installation Options

### Option 1: Quick One-Liner Commands

#### Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/Chintanpatel24/Matt/main/install.sh | bash
```

Or using `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/Chintanpatel24/Matt/main/install.sh | bash
```

#### Windows (PowerShell)

```powershell
iwr -useb https://raw.githubusercontent.com/Chintanpatel24/Matt/main/install.ps1 | iex
```

---

### Option 2: Clone & Install Manually

You can clone the repository and run `install.sh` or use `make`:

```bash
git clone https://github.com/Chintanpatel24/Matt.git
cd Matt
./install.sh
```

Or build with `make`:

```bash
make build
make install
```

---

### Updating Matt

To update Matt to the latest version at any time:

```bash
curl -fsSL https://raw.githubusercontent.com/Chintanpatel24/Matt/main/update.sh | bash
```

Or from inside your local cloned repository:

```bash
git pull
./install.sh
```

---

## Key Features

- **Sleek Matt Black Aesthetic**: Custom HSL dark palette (`#09090b` onyx background, `#18181b` pane surface, `#27272a` charcoal borders, matte gold accents).
- **Startup Session Screen**: On launch, directly choose between opening current directory (`./`) or starting a fresh home session (`~/`).
- **3-Column Split View**:
  - **Left Pane (Navigation Tree)**: Directory listing with live fuzzy filtering (`/`) & mouse click support.
  - **Center Pane (Explorer)**: Folder contents view with quick drill-down.
  - **Right Pane (Code & Image Inspector)**: Real-time syntax-highlighted code inspector, image pixel previewer (`.png`, `.jpg`, `.gif`, `.webp`), binary hex view, and directory overview.
- **Disk Space Analyzer (`Alt+D`)**: `ncdu`-style visual disk usage breakdown with percentage progress bars (`██████░░ 75%`).
- **Instant Fuzzy Search (`/`)**: Live fuzzy filtering across file lists as you type.
- **Config & Alias Extensions (`~/.config/matt/config.json`)**: Define custom command aliases (e.g. `ll`, `findbig`, `count`, `sysinfo`) and theme overrides.
- **Integrated Bottom Terminal**:
  - Direct shell command execution (`cd`, `touch`, `mkdir`, `rm`, `find`, `git`, `cat`, etc.).
  - Executing `cd <path>` syncs all 3 upper panes immediately to the new directory!
- **Interactive Permission Prompt System**: Safe confirmation dialogs for destructive operations (`rm`, `chmod`, `chown`) and feature approvals.

---

## Keyboard & Mouse Controls

| Key / Action | Description |
| :--- | :--- |
| `Up` / `Down` (`k` / `j`) | Navigate files/directories in focused pane |
| `Right` / `Enter` (`l`) | Open folder / expand directory |
| `Left` (`h`) | Go up to parent directory (`..`) |
| `Alt+D` | Toggle Disk Space Analyzer view (`ncdu` style) |
| `/` | Live fuzzy search & filter directory list |
| `Left Mouse Click` | Focus pane & select item directly |
| `Tab` / `Shift+Tab` | Cycle active focus between 3 upper panes & bottom terminal |
| `:` | Focus bottom terminal / run shell commands & extension aliases |
| `.` | Toggle hidden files |
| `r` | Refresh directory view |
| `d` | Delete highlighted file/folder (triggers permission dialog) |
| `Esc` | Unfocus mode / dismiss modal |
| `q` or `Ctrl+C` | Quit Matt |

---

## Configuration & Extension API

Matt automatically reads configuration from `~/.config/matt/config.json`. You can define custom shell command aliases and custom color schemes:

```json
{
  "show_hidden": false,
  "aliases": {
    "ll": "ls -la",
    "findbig": "find . -type f -size +10M",
    "count": "find . -type f | wc -l",
    "sysinfo": "uname -a && uptime"
  },
  "theme": {
    "bg": "#09090b",
    "bg_surface": "#18181b",
    "border": "#27272a",
    "border_active": "#eab308",
    "text_primary": "#f4f4f5",
    "accent": "#eab308"
  }
}
```
