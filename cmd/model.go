package cmd

import (
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pixambi/gbrain/internal/db"
)

const (
	projectsView = iota
	projectView
	projectTitleView
	nodeView
	nodeTitleView
	nodeContentView
)

type model struct {
	state            uint
	db               db.Db
	projects         []db.Project
	nodes            []db.Node
	textArea         textarea.Model
	textInput        textinput.Model
	currentNode      db.Node
	currentProject   db.Project
	projectListIndex int
	nodeListIndex    int
	width            int
	height           int
	err              error
}

func NewApp(db db.Db) model {
	if err := db.Init(); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	projects, err := db.GetProjects()
	if err != nil {
		log.Fatalf("Error getting projects: %v", err)
	}

	ti := textinput.New()
	ti.Placeholder = "Enter title..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 30

	ta := textarea.New()
	ta.Placeholder = "Enter content..."
	ta.Focus()
	ta.SetWidth(80)
	ta.SetHeight(20)

	return model{
		state:            projectsView,
		db:               db,
		projects:         projects,
		textArea:         ta,
		textInput:        ti,
		projectListIndex: 0,
		nodeListIndex:    0,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textArea.SetWidth(msg.Width - 4)
		m.textArea.SetHeight(msg.Height - 8)
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		switch m.state {
		case projectsView:
			switch key {
			case "q", "ctrl+c", "esc":
				return m, tea.Quit

			case "n":
				m.textInput.Reset()
				m.textInput.Focus()
				m.state = projectTitleView
				return m, textinput.Blink

			case "j", "down":
				if len(m.projects) > 0 && m.projectListIndex < len(m.projects)-1 {
					m.projectListIndex++
				}

			case "k", "up":
				if m.projectListIndex > 0 {
					m.projectListIndex--
				}

			case "enter":
				if len(m.projects) > 0 {
					// Load nodes for this project
					m.currentProject = m.projects[m.projectListIndex]
					nodes, err := m.db.GetNodesByProjectID(m.currentProject.ID)
					if err != nil {
						m.err = err
						return m, nil
					}
					m.nodes = nodes
					m.nodeListIndex = 0
					m.state = projectView
				}
			}

		case projectTitleView:
			switch key {
			case "esc":
				m.state = projectsView
				return m, nil

			case "enter":
				if m.textInput.Value() != "" {
					project := db.Project{
						Name: m.textInput.Value(),
					}
					if err := m.db.AddProject(project); err != nil {
						m.err = err
						return m, nil
					}

					projects, err := m.db.GetProjects()
					if err != nil {
						m.err = err
						return m, nil
					}
					m.projects = projects
					m.state = projectsView
				}
			}

			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)

		case projectView:
			switch key {
			case "esc", "q":
				m.state = projectsView
				return m, nil

			case "n":
				m.textInput.Reset()
				m.textInput.Focus()
				m.currentNode = db.Node{ProjectID: m.currentProject.ID}
				m.state = nodeTitleView
				return m, textinput.Blink

			case "j", "down":
				if len(m.nodes) > 0 && m.nodeListIndex < len(m.nodes)-1 {
					m.nodeListIndex++
				}

			case "k", "up":
				if m.nodeListIndex > 0 {
					m.nodeListIndex--
				}

			case "enter":
				if len(m.nodes) > 0 {
					m.currentNode = m.nodes[m.nodeListIndex]
					m.state = nodeView
				}
			}

		case nodeTitleView:
			switch key {
			case "esc":
				m.state = projectView
				return m, nil

			case "enter":
				if m.textInput.Value() != "" {
					m.currentNode.Title = m.textInput.Value()
					m.textArea.Reset()
					m.textArea.Focus()
					m.state = nodeContentView
					return m, textarea.Blink
				}
			}

			m.textInput, cmd = m.textInput.Update(msg)
			cmds = append(cmds, cmd)

		case nodeContentView:
			switch key {
			case "esc":
				m.state = nodeTitleView
				m.textInput.SetValue(m.currentNode.Title)
				m.textInput.Focus()
				return m, textinput.Blink

			case "ctrl+s":
				m.currentNode.Content = m.textArea.Value()
				if err := m.db.AddNode(m.currentNode); err != nil {
					m.err = err
					return m, nil
				}

				nodes, err := m.db.GetNodesByProjectID(m.currentProject.ID)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.nodes = nodes
				m.state = projectView
				return m, nil
			}

			m.textArea, cmd = m.textArea.Update(msg)
			cmds = append(cmds, cmd)

		case nodeView:
			switch key {
			case "esc", "q":
				m.state = projectView
				return m, nil

			case "e":
				m.textInput.SetValue(m.currentNode.Title)
				m.textInput.Focus()
				m.state = nodeTitleView
				return m, textinput.Blink
			}
		}
	}

	return m, tea.Batch(cmds...)
}
