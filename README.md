# GBrain

GBrain is a terminal-based note-taking application that lets you organize your notes in a database of linked nodes. 

## Features

- **Project-based organization**: Group related notes into separate projects
- **Linked notes**: Create connections between notes using `[[WikiLink]]` syntax
- **Terminal UI**: keyboard-driven interface using [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Storage**: Your data is stored locally in a BoltDB database

## Key Bindings

### Global
- `q` or `Esc`: Go back or quit

### Projects View
- `j`/`down`: Navigate down
- `k`/`up`: Navigate up
- `n`: New project
- `d`: Delete project
- `Enter`: Open project

### Project View (Notes List)
- `j`/`down`: Navigate down
- `k`/`up`: Navigate up
- `n`: New note
- `d`: Delete note
- `Enter`: View note
- `Esc`: Back to projects

### Note View
- `Tab`: Cycle links
- `Enter`: Follow link
- `b`: Go back to previous note
- `e`: Edit note
- `d`: Delete note
- `Esc`: Back to notes list

### Editing
- `Enter`: Save title and continue to content
- `Ctrl+s`: Save note
- `Esc`: Cancel or go back

## License

This project is licensed under the GNU General Public License Version 3 - see the LICENSE file for details.

## Acknowledgements

GBrain uses the following libraries:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the terminal UI
- [Bubbles](https://github.com/charmbracelet/bubbles) for UI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling
- [BoltDB](https://github.com/etcd-io/bbolt) for data storage