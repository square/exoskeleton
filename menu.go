package exoskeleton

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"
)

const menuTemplate = "\033[1m" + `USAGE` + "\033[0m" + `
   {{.Usage}}

{{- range .Sections}}

` + "\033[1m" + `{{.Heading}}` + "\033[0m" + `
{{- range .MenuItems}}
   {{rpad .Name .Width}}  {{.Summary}}
{{- end}}
{{- end}}

Run ` + "\033[96m" + `{{.HelpUsage}} <command>` + "\033[0m" + ` to print information on a specific command.`

var templateFuncs = template.FuncMap{
	"rpad": func(s string, padding int) string { return fmt.Sprintf("%*s", -padding, s) },
}

// SummaryFunc is a function that is expected to return the heading
type SummaryFunc func(Command) (string, error)

type menuOptions struct {
	Depth      int
	HeadingFor MenuHeadingForFunc
	SummaryFor SummaryFunc
	Template   *template.Template
}

// Menu is the data passed to MenuOptions.Template when it is executed.
type Menu struct {
	Usage     string
	HelpUsage string
	Sections  MenuSections
}

type MenuSections []MenuSection

type MenuSection struct {
	Heading   string
	MenuItems MenuItems
}

type MenuItems []*MenuItem

// implement sort.Interface so that MenuItems can be sorted by Name

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

type MenuItem struct {
	Name    string
	Summary string
	Heading string
	Width   int
}

// menuFor renders a menu of commands for a Module.
func menuFor(m Module, opts *menuOptions) (string, []error) {
	if opts.Template == nil {
		opts.Template = template.Must(template.New("menu").Funcs(templateFuncs).Parse(menuTemplate))
	}

	menu, errs := buildMenu(m, opts)
	b := new(bytes.Buffer)
	if err := opts.Template.Execute(b, menu); err != nil {
		panic(err)
	}
	return b.String(), errs
}

// buildMenu constructs a Menu of Commands with their short summary strings for a given Module.
func buildMenu(m Module, opts *menuOptions) (*Menu, []error) {
	if opts.SummaryFor == nil {
		opts.SummaryFor = func(cmd Command) (string, error) { return cmd.Summary() }
	}

	if opts.HeadingFor == nil {
		opts.HeadingFor = func(Module, Command) string { return "COMMANDS" }
	}

	c, err := m.Subcommands()
	if err != nil {
		return &Menu{}, []error{err}
	}

	c, errs := expandModules(c, opts.Depth)
	var items MenuItems
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

	return &Menu{
		Usage:     Usage(m) + " <command> [<args>]",
		Sections:  sections,
		HelpUsage: helpUsage(m),
	}, errs
}

func expandModules(cmds Commands, depth int) (Commands, []error) {
	all := Commands{}
	var errs []error

	for _, cmd := range cmds {
		if m, ok := cmd.(Module); ok && depth != 0 {
			subcmds, err := m.Subcommands()
			if err != nil {
				errs = append(errs, err)
			}
			expandedSubcmds, expansionErrs := expandModules(subcmds, depth-1)
			all = append(all, expandedSubcmds...)
			errs = append(errs, expansionErrs...)
		} else {
			all = append(all, cmd)
		}
	}

	return all, errs
}

func helpUsage(m Module) string {
	args := argsRelativeTo(m, nil)
	return strings.Join(append([]string{args[0], "help"}, args[1:]...), " ")
}
