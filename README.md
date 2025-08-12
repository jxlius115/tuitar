# tuitar

A terminal-based guitar tablature editor built with Go and Bubble Tea.

## Features

- ✨ **Intuitive Terminal Interface**: Vim-like keyboard navigation
- 📝 **Real-time Tab Editing**: Create and edit guitar tabs with instant feedback
- 🎵 **MIDI Playback**: Play your tabs with visual highlighting (basic implementation)
- 💾 **Local Storage**: SQLite-based tab management
- 🔍 **Tab Browser**: Browse, search, and organize your tabs
- ⌨️  **Keyboard-driven**: Efficient workflows without mouse dependency
- 🌍 **Cross-platform**: Works on Linux, macOS, and Windows

## Installation

```bash
# Install from source
git clone https://github.com/Cod-e-Codes/tuitar
cd tuitar
go build -o tuitar

# Or install directly
go install github.com/Cod-e-Codes/tuitar@latest
```

## Usage

```bash
# Start the application
./tuitar

# The application will create a tabs.db SQLite database in the current directory
```

## Key Bindings

### Global
- `q` / `Ctrl+C` - Quit application
- `?` - Toggle help
- `Tab` - Switch between browser and editor
- `Ctrl+N` - Create new tab
- `Ctrl+S` - Save current tab

### Browser Mode
- `j` / `↓` - Move down
- `k` / `↑` - Move up
- `Enter` - Edit selected tab

### Editor Mode
- `h` / `←` - Move cursor left
- `j` / `↓` - Move cursor down (to next string)
- `k` / `↑` - Move cursor up (to previous string)
- `l` / `→` - Move cursor right
- `i` - Enter insert mode
- `Esc` - Return to normal mode
- `0-9` - Insert fret number (in insert mode)
- `-` - Insert rest (in insert mode)
- `x` - Delete character at cursor
- `Space` - Play/pause tab (basic implementation)

## Project Structure

The application follows a clean architecture pattern:

- `models/` - Core data structures and business logic
- `storage/` - Data persistence layer (SQLite)
- `ui/` - Bubble Tea UI components and views  
- `midi/` - MIDI playback functionality (basic implementation)
- `parser/` - Import/export parsers for different tab formats

## Building from Source

```bash
# Clone the repository
git clone https://github.com/Cod-e-Codes/tuitar
cd tuitar

# Install dependencies
go mod tidy

# Build
go build -o tuitar

# Run
./tuitar
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [SQLite](https://github.com/mattn/go-sqlite3) - Local database

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT License - see LICENSE file for details

## Roadmap

- [ ] Full MIDI playback with actual audio output
- [ ] Guitar Pro file import/export
- [ ] Advanced tab notation (bends, slides, etc.)
- [ ] Multi-instrument support
- [ ] Plugin system for extensibility
- [ ] Network sync capabilities
