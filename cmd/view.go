package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	appNameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Background(lipgloss.Color("0")).
			Bold(true).
			Padding(1, 2).
			Align(lipgloss.Center)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Bold(true).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().
			Padding(0, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("170")).
				Background(lipgloss.Color("237")).
				Bold(true).
				Padding(0, 2)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderBottom(true).
			Bold(true).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Padding(0, 1)

	linkStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Underline(true)

	selectedLinkStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Background(lipgloss.Color("237")).
				Underline(true).
				Bold(true)

	editModeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("202")).
			Bold(true).
			Padding(0, 1)
)

func (m model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	var s strings.Builder

	s.WriteString(appNameStyle.Render("Gbrain"))
	s.WriteString("\n\n")

	switch m.state {
	case projectsView:
		s.WriteString(titleStyle.Render("Projects"))
		s.WriteString("\n\n")

		if len(m.projects) == 0 {
			s.WriteString(infoStyle.Render("No projects yet. Press 'n' to create one."))
		} else {
			for i, project := range m.projects {
				style := itemStyle
				if i == m.projectListIndex {
					style = selectedItemStyle
				}
				s.WriteString(style.Render(project.Name))
				s.WriteString("\n")
			}
		}

		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render("j/k: navigate • n: new project • d: delete project • enter: open • q: quit"))

	case confirmDeleteProjectView:
		s.WriteString(warningStyle.Render("Delete Project"))
		s.WriteString("\n\n")
		s.WriteString(warningStyle.Render(fmt.Sprintf("Are you sure you want to delete '%s' and all its nodes?", m.currentProject.Name)))
		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render("y: confirm delete • n: cancel"))

	case projectTitleView:
		s.WriteString(titleStyle.Render("New Project"))
		s.WriteString("\n\n")
		s.WriteString(m.textInput.View())
		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render("enter: save • esc: cancel"))

	case projectView:
		s.WriteString(titleStyle.Render(fmt.Sprintf("Project: %s", m.currentProject.Name)))
		s.WriteString("\n\n")

		if len(m.nodes) == 0 {
			s.WriteString(infoStyle.Render("No nodes yet. Press 'n' to create one."))
		} else {
			for i, node := range m.nodes {
				style := itemStyle
				if i == m.nodeListIndex {
					style = selectedItemStyle
				}
				s.WriteString(style.Render(node.Title))
				s.WriteString("\n")
			}
		}

		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render("j/k: navigate • n: new node • d: delete node • enter: view • esc: back"))

	case confirmDeleteNodeView:
		s.WriteString(warningStyle.Render("Delete Node"))
		s.WriteString("\n\n")
		s.WriteString(warningStyle.Render(fmt.Sprintf("Are you sure you want to delete '%s'?", m.currentNode.Title)))
		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render("y: confirm delete • n: cancel"))

	case nodeTitleView:
		action := "New"
		if m.currentNode.ID != 0 {
			action = "Edit"
		}
		s.WriteString(titleStyle.Render(fmt.Sprintf("%s Node Title", action)))
		s.WriteString("\n\n")
		s.WriteString(m.textInput.View())
		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render("enter: continue to content • esc: cancel"))

	case nodeContentView:
		s.WriteString(titleStyle.Render(fmt.Sprintf("Node: %s", m.currentNode.Title)))
		s.WriteString("\n\n")
		s.WriteString(m.textArea.View())
		s.WriteString("\n\n")
		s.WriteString(infoStyle.Render("ctrl+s: save • esc: back to title"))
		s.WriteString("\n")
		s.WriteString(editModeStyle.Render("Note: Use [[Node Title]] to create links"))

	case nodeView:
		s.WriteString(titleStyle.Render(m.currentNode.Title))
		s.WriteString("\n\n")

		content := renderContent(m.currentNode.Content, m.links, m.currentLinkIndex)
		s.WriteString(content)

		s.WriteString("\n\n")

		if len(m.links) > 0 {
			s.WriteString(infoStyle.Render("tab: cycle links • enter: follow link • b: go back • e: edit • d: delete • esc: back"))
		} else {
			s.WriteString(infoStyle.Render("e: edit • d: delete • esc: back"))
		}
	}

	return s.String()
}
