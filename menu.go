package exoskeleton

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"
)

// MenuRelativeTo constructs a menu of Commands with usage strings relative to a given Module.
func (e *Entrypoint) MenuRelativeTo(c Commands, m Module) Menu {
	usage := Usage(m) + " <command> [<args>]"

	var items MenuItems

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

		if summary, err := cache.Read(cmd); err != nil {
			e.onError(err)
		} else if summary != "" {
			heading := e.menuHeadingFor(m, cmd)
			items = append(items, &MenuItem{Name: name, Summary: summary, Heading: heading})
		}
	}

	width := items.MaxWidth()

	byHeading := make(map[string]MenuItems)
	var orderedHeadings []string
	for _, menuItem := range items {
		menuItem.Width = width
		if _, present := byHeading[menuItem.Heading]; !present {
			orderedHeadings = append(orderedHeadings, menuItem.Heading)
		}
		byHeading[menuItem.Heading] = append(byHeading[menuItem.Heading], menuItem)
	}

	var sections MenuSections
	for _, heading := range orderedHeadings {
		menuItems := byHeading[heading]
		if len(menuItems) > 0 {
			sort.Sort(menuItems)
			sections = append(sections, MenuSection{heading, menuItems})
		}
	}

	trailer := fmt.Sprintf(
		"Run \033[96m%s help %s\033[0m to print information on a specific command.",
		Usage(e),
		strings.TrimLeft(UsageRelativeTo(m, e)+" <command>", " "),
	)

	return Menu{Usage: usage, Sections: sections, Trailer: trailer}
}

type Menu struct {
	Usage    string
	Sections MenuSections
	Trailer  string
}

func (m Menu) String() string {
	return fmt.Sprintf("USAGE\n   %s\n\n%s\n\n%s", m.Usage, m.Sections, m.Trailer)
}

type MenuSections []MenuSection

func (m MenuSections) String() string {
	var s []string
	for _, section := range m {
		s = append(s, section.String())
	}
	return strings.Join(s, "\n\n")
}

type MenuSection struct {
	Heading   string
	MenuItems MenuItems
}

func (section MenuSection) String() string {
	return fmt.Sprintf("\033[1m%s\033[0m\n   %s", section.Heading, section.MenuItems)
}

type MenuItems []*MenuItem

// Implement sort.Interface so that MenuItems can be sorted by Name
func (m MenuItems) Len() int           { return len(m) }
func (m MenuItems) Less(i, j int) bool { return m[i].Name < m[j].Name }
func (m MenuItems) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

func (m MenuItems) MaxWidth() (longestCommand int) {
	for _, menuItem := range m {
		if len(menuItem.Name) > longestCommand {
			longestCommand = len(menuItem.Name)
		}
	}
	return
}

func (m MenuItems) String() string {
	var s []string
	for _, mi := range m {
		s = append(s, mi.String())
	}
	return strings.Join(s, "\n   ")
}

type MenuItem struct {
	Name    string
	Summary string
	Heading string
	Width   int
}

func (mi *MenuItem) String() string {
	return fmt.Sprintf("%*s  %s", -mi.Width, mi.Name, mi.Summary)
}
