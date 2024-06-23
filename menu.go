package exoskeleton

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"
)

type SummaryFunc func(Command) (string, error)

type buildMenuOptions struct {
	HeadingFor MenuHeadingForFunc
	SummaryFor SummaryFunc
}

func buildMenu(c Commands, m Module, opts *buildMenuOptions) (menu, []error) {
	if opts.SummaryFor == nil {
		opts.SummaryFor = func(cmd Command) (string, error) { return cmd.Summary() }
	}

	if opts.HeadingFor == nil {
		opts.HeadingFor = func(Module, Command) string { return "COMMANDS" }
	}

	usage := Usage(m) + " <command> [<args>]"

	var items menuItems
	var errs []error

	seen := make(map[string]bool)

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

		if summary, err := opts.SummaryFor(cmd); err != nil {
			errs = append(errs, err)
		} else if summary != "" {
			heading := opts.HeadingFor(m, cmd)
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

	a := argsRelativeTo(m, nil)
	helpUsage := strings.Join(append([]string{a[0], "help"}, a[1:]...), " ")

	return menu{Usage: usage, Sections: sections, HelpUsage: helpUsage}, errs
}

type menu struct {
	Usage     string
	HelpUsage string
	Sections  menuSections
}

func (m menu) String() string {
	return fmt.Sprintf("USAGE\n   %s\n\n%s\n\nRun \033[96m%s <command>\033[0m to print information on a specific command.", m.Usage, m.Sections, m.HelpUsage)
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
