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

- **Matt** is a fast, keyboard-driven, modern terminal file manager featuring a **Matt Black** design aesthetic. Built with Go and Bubble Tea, it offers a 3-column pane system, async disk analyzer, live fuzzy finder, file operations, bookmarks, command history, syntax-highlighted previews, and a premium dark UI with scroll indicators, zebra-striping, and breadcrumb navigation.
- ## Table of Contents
- [Key Features](#key-features)
  - [Premium UI](#premium-ui)
  - [File Management & Bulk Operations](#file-management--bulk-operations)
  - [Archive Management](#archive-management)
  - [Navigation & Search](#navigation--search)
  - [Async Disk Analyzer](#async-disk-analyzer)
  - [Integrated Terminal & Openers](#integrated-terminal--openers)
  - [Preview & Inspection](#preview--inspection)
- [Installation](#installation)
  - [Quick One-Liner](#quick-one-liner)
  - [Clone & Build](#clone--build)
  - [Updating](#updating)
- [Keyboard Controls](#keyboard-controls)
  - [Navigation](#navigation-1)
  - [File Operations](#file-operations-1)
  - [Tools & Navigation](#tools--navigation-1)
- [Configuration](#configuration)
  - [Data Files](#data-files)
- [License](#license)


<div align=center>
  <table>
      <tr>
       <td>
      <img width="600" height="600" alt="mattt" src="https://github.com/user-attachments/assets/bdc8921c-0dd8-4e6d-9db5-d5bcccb974f3" />
       </td>
        <td>
      <img width="600" height="600" alt="matt" src="https://github.com/user-attachments/assets/afb6c3ae-f9f7-4bfa-89af-4e8c89efc1a4" /
      </td>
      </tr>
  </table>
</div>

---

## Key Features

### Premium UI
- **Matt Black Aesthetic** - Custom dark palette with onyx backgrounds, charcoal borders, and high-contrast grey accents
- **Solid Takeover Background** - Full terminal takeover layout prevents transparent patches and preserves background color consistency
- **3-Column Split View** - Navigation tree, explorer, and preview/metadata inspector
- **Breadcrumb Path Bar** - Interactive breadcrumb navigation in the header
- **Scroll Indicators** - Visual arrows with item counts when content overflows
- **Zebra-Striping** - Alternating row backgrounds for better readability
- **Selection Accent Bar** - Accent indicator on the active selection
- **Focus Indicators** - Active / inactive pane markers
- **Mode Indicator** - Shows current mode (NORMAL / FILTER / COMMAND)

### File Management & Bulk Operations
- **Multi-Selection** - Press Space to select multiple files/directories, marked with a checkmark
- **Bulk Operations** - Perform bulk copy (c), cut (x), paste (p), and delete (d) on all selected files simultaneously
- **Create Files & Folders** - n for new file, N for new folder
- **Rename** - m to move/rename with inline prompt
- **Symlink Resolution** - Shows symlink targets
- **Action Undo** - Press u to revert the last copy, paste, delete, or rename file action

### Archive Management
- **Zip Compression** - Press z to compress selected files/folders into a .zip archive
- **Archive Extraction** - Press Z on an archive (.zip, .tar.gz, .tar, .tgz) to extract it inline instantly

### Navigation & Search
- **Live Fuzzy Search** - Press / for instant fuzzy filtering
- **Git Status Indicators** - Color-coded file list reflecting Git status:
  - [M] (Yellow) - Modified files
  - [A] (Green) - Staged/added files
  - [?] (Gray) - Untracked files
- **Bookmarks** - B to bookmark, b to jump to saved locations
- **Quick Jump** - g to first item, G to last item
- **Vim-Style Navigation** - h/j/k/l or arrow keys

### Async Disk Analyzer
- **Non-Blocking Analysis** - Alt+D runs disk scan in background (no more UI freezing!)
- **Emoji-Free Layout** - Modern text-based tabular layout (DIR/FILE labels)
- **Visual Progress Bars** - Percentage breakdown
- **Sort Toggle** - Press s to toggle between size and name sorting
- **Full Scrolling** - Navigate large directory analyses

### Integrated Terminal & Openers
- **System Opener** - Press o or O on any file to open it with the default system association handler
- **Shell Command Execution** - Run any command from the bottom pane
- **Command History** - Up/Down arrow keys cycle through previous commands
- **Smart cd** - Syncs all 3 panes when changing directories
- **Config Aliases** - Custom shortcuts like ll, findbig, sysinfo
- **Permission Prompts** - Safe confirmation for destructive operations

### Preview & Inspection
- **Rich Previews** - Formatted previews for specific extensions:
  - **Markdown (.md)** - Styled headers, styled bullet points, and code blocks
  - **CSV (.csv)** - Structured, aligned, zebra-striped data table grid with column dividers
- **Syntax Highlighting** - Powered by Chroma with line numbers
- **ASCII Image Preview** - Renders image formats as ASCII art
- **Binary Hex View** - Hex dump for binary files
- **File Metadata Inspector** - Permissions, owner, size, timestamps

---

## Installation

### Quick One-Liner

#### Linux / macOS
```bash
curl -fsSL https://raw.githubusercontent.com/Chintanpatel24/Matt/main/install.sh | bash
```

Or using wget:
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
| `Space` | Toggle multi-select checkbox on current item |
| `n` | Create new file (enter name in prompt) |
| `N` | Create new folder (enter name in prompt) |
| `m` | Rename/move selected item |
| `c` | Copy selected item(s) to clipboard |
| `x` | Cut selected item(s) to clipboard |
| `p` | Paste clipboard item(s) into current directory |
| `d` | Delete selected item(s) (with confirmation) |
| `u` | Revert/Undo the last file action (copy, paste, move, delete) |
| `z` | Compress selected item(s) into archive.zip |
| `Z` | Extract selected archive inline |
| `o` / `O` | Open file in system default application handler |

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
    "border_active": "#71717a",
    "text_primary": "#f8fafc",
    "text_muted": "#a1a1aa",
    "directory": "#38bdf8",
    "executable": "#34d399",
    "selection": "#3f3f46",
    "accent": "#a1a1aa",
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

---

## License

This project is licensed under the [MIT License](LICENSE).
