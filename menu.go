package exoskeleton

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"
)

func (e *Entrypoint) buildMenu(c Commands, m Module) menu {
	usage := Usage(m) + " <command> [<args>]"

	var items menuItems

	seen := make(map[string]bool)
	cache := &summaryCache{Path: e.cachePath, onError: e.onError}

	for _, cmd := range c {
		name := UsageRelativeTo(cmd, m)
		if _, ok := cmd.(Module); ok {
			name += ":"
		}

		if seen[name] {
			continue
		} else {
			seen[name] = true
		}

		if summary := cache.Read(cmd); summary != "" {
			heading := e.menuHeadingFor(m, cmd)
			items = append(items, &menuItem{Name: name, Summary: summary, Heading: heading})
		}
	}

	width := items.MaxWidth()

	byHeading := make(map[string]menuItems)
	var orderedHeadings []string
	for _, menuItem := range items {
		menuItem.Width = width
		if _, present := byHeading[menuItem.Heading]; !present {
			orderedHeadings = append(orderedHeadings, menuItem.Heading)
		}
		byHeading[menuItem.Heading] = append(byHeading[menuItem.Heading], menuItem)
	}

	var sections menuSections
	for _, heading := range orderedHeadings {
		menuItems := byHeading[heading]
		if len(menuItems) > 0 {
			sort.Sort(menuItems)
			sections = append(sections, menuSection{heading, menuItems})
		}
	}

	trailer := fmt.Sprintf(
		"Run \033[96m%s help %s\033[0m to print information on a specific command.",
		Usage(e),
		strings.TrimLeft(UsageRelativeTo(m, e)+" <command>", " "),
	)

	return menu{Usage: usage, Sections: sections, Trailer: trailer}
}

type menu struct {
	Usage    string
	Sections menuSections
	Trailer  string
}

func (m menu) String() string {
	return fmt.Sprintf("USAGE\n   %s\n\n%s\n\n%s", m.Usage, m.Sections, m.Trailer)
}

type menuSections []menuSection

func (m menuSections) String() string {
	var s []string
	for _, section := range m {
		s = append(s, section.String())
	}
	return strings.Join(s, "\n\n")
}

type menuSection struct {
	Heading   string
	MenuItems menuItems
}

func (section menuSection) String() string {
	return fmt.Sprintf("\033[1m%s\033[0m\n   %s", section.Heading, section.MenuItems)
}

type menuItems []*menuItem

// Implement sort.Interface so that MenuItems can be sorted by Name
func (m menuItems) Len() int           { return len(m) }
func (m menuItems) Less(i, j int) bool { return m[i].Name < m[j].Name }
func (m menuItems) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

func (m menuItems) MaxWidth() (longestCommand int) {
	for _, menuItem := range m {
		if len(menuItem.Name) > longestCommand {
			longestCommand = len(menuItem.Name)
		}
	}
	return
}

func (m menuItems) String() string {
	var s []string
	for _, mi := range m {
		s = append(s, mi.String())
	}
	return strings.Join(s, "\n   ")
}

type menuItem struct {
	Name    string
	Summary string
	Heading string
	Width   int
}

func (mi *menuItem) String() string {
	return fmt.Sprintf("%*s  %s", -mi.Width, mi.Name, mi.Summary)
}
