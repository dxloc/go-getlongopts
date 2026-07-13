package getlongopts

import (
	"cmp"
	"fmt"
	"os"
	"strings"

	"github.com/dxloc/gosort"
)

const (
	ArgTypeNone = iota
	ArgTypeDir
	ArgTypeFile
	ArgTypeOther
)

type OptSetFunc func(value string)
type OptTableEntry struct {
	typ   int
	setFn OptSetFunc
}

type LongOption struct {
	Long        string
	Short       string
	ArgType     int
	SetFn       OptSetFunc
	Description string
}

type Parser struct {
	opts             []LongOption
	optTable         map[string]OptTableEntry
	maxLongOptLength int
}

var ErrUnknownOption = fmt.Errorf("unknown option")

// Compare two 'LongOption's
func (o LongOption) Compare(a, b LongOption) int {
	if a.Short != b.Short {
		return cmp.Compare(a.Short, b.Short)
	} else {
		return cmp.Compare(a.Long, b.Long)
	}
}

// Parse command line options and return the remaining arguments.
func (p *Parser) Parse(args []string) ([]string, error) {
	var key string

	argc := len(args)
	i := 1

	for ; i < argc; i++ {
		argv := args[i]
		if argv[0] != '-' {
			break
		}
		if argv[1] == '-' {
			key = argv[2:]
		} else {
			key = argv[1:]
		}
		if ent, ok := p.optTable[key]; ok {
			if ent.typ == ArgTypeNone {
				if ent.setFn != nil {
					ent.setFn("")
				}
			} else {
				i++
				if i >= argc {
					return nil, ErrUnknownOption
				}
				argv = args[i]
				if ent.setFn != nil {
					ent.setFn(argv)
				}
			}
		} else {
			return nil, ErrUnknownOption
		}
	}
	if i < argc {
		return args[i:], nil
	} else {
		return nil, nil
	}
}

// Built-in bash-completion generator
func (p *Parser) BashCompletion(progname string) string {
	var s strings.Builder
	var c strings.Builder

	ss := strings.Split(progname, "/")
	progname = ss[len(ss)-1]

	s.WriteString("#!/bin/bash\n")
	s.WriteString("_")
	s.WriteString(progname)
	s.WriteString("() {\n")
	s.WriteString("    local cur prev words cword\n")
	s.WriteString("    _init_completion || return\n")
	s.WriteString("    case ${prev} in\n")
	for _, opt := range p.opts {
		if opt.Short == "" && opt.Long == "" {
			continue
		}
		if opt.Long != "" {
			c.WriteString(" --")
			c.WriteString(opt.Long)
		} else if opt.Short != "" {
			c.WriteString(" -")
			c.WriteString(opt.Short)
		}
		switch opt.ArgType {
		case ArgTypeDir:
			s.WriteString("        ")
			if opt.Short != "" {
				s.WriteString("-")
				s.WriteString(opt.Short)
				if opt.Long != "" {
					s.WriteString("|")
				}
			}
			if opt.Long != "" {
				s.WriteString("--")
				s.WriteString(opt.Long)
			}
			s.WriteString(")\n")
			s.WriteString("            COMPREPLY=( $(compgen -d -- ${cur}) )\n")
			s.WriteString("            ;;\n")
		case ArgTypeFile:
			s.WriteString("        ")
			if opt.Short != "" {
				s.WriteString("-")
				s.WriteString(opt.Short)
				if opt.Long != "" {
					s.WriteString("|")
				}
			}
			if opt.Long != "" {
				s.WriteString("--")
				s.WriteString(opt.Long)
			}
			s.WriteString(")\n")
			s.WriteString("            COMPREPLY=( $(compgen -f -- ${cur}) )\n")
			s.WriteString("            ;;\n")
		}
	}
	s.WriteString("        *)\n")
	s.WriteString("            COMPREPLY=( $(compgen -W '")
	s.WriteString(c.String())
	s.WriteString("' -- ${cur}) )\n")
	s.WriteString("            ;;\n")
	s.WriteString("    esac\n")
	s.WriteString("} &&\n")
	s.WriteString("    complete -o filenames -F _")
	s.WriteString(progname)
	s.WriteString(" ")
	s.WriteString(progname)
	s.WriteString("\n")
	s.WriteString("# ex: ts=4 sw=4 et filetype=sh")

	return s.String()
}

// Built-in help function
func (p *Parser) Usage() string {
	var s strings.Builder
	s.WriteString("Usage:\n")
	s.WriteString("  ")
	s.WriteString(os.Args[0])
	s.WriteString(" [options]\n")
	s.WriteString("\n")
	s.WriteString("Options:\n")
	for _, opt := range p.opts {
		if opt.Short == "" && opt.Long == "" {
			continue
		}
		s.WriteString("  ")
		if opt.Short != "" {
			s.WriteString("-")
			s.WriteString(opt.Short)
			if opt.Long != "" {
				s.WriteString("|")
			}
		} else {
			s.WriteString("   ")
		}
		if opt.Long != "" {
			s.WriteString("--")
			s.WriteString(opt.Long)
		}

		longLen := len(opt.Long)
		for longLen < p.maxLongOptLength+8 {
			s.WriteString(" ")
			longLen++
		}

		if opt.Description != "" {
			s.WriteString(opt.Description)
		}
		s.WriteString("\n")
	}
	return s.String()
}

// Initialize and returns a new Parser with the given options from the 'opts' parameter.
//
// The new Parser is also initialized with the built-in bash-completion generator
// (-b|--bash-completion) and (-h|--help) help options. Users can override the built-in
// options by their own (-b), (--bash-completion), (-h), (--help) options.
func NewParser(opts []LongOption) *Parser {
	p := Parser{optTable: make(map[string]OptTableEntry)}
	bashOpts := LongOption{
		ArgType: ArgTypeNone,
		Short:   "b",
		Long:    "bash-completion",
		SetFn: func(value string) {
			fmt.Println(p.BashCompletion(os.Args[0]))
			os.Exit(0)
		},
		Description: "Print bash completion script and exit",
	}
	helpOpts := LongOption{
		ArgType: ArgTypeNone,
		Short:   "h",
		Long:    "help",
		SetFn: func(value string) {
			fmt.Println(p.Usage())
			os.Exit(0)
		},
		Description: "Print this message and exit",
	}
	for _, opt := range opts {
		if opt.Short == "" && opt.Long == "" {
			continue
		}
		if opt.Short == "b" {
			bashOpts.Short = ""
		}
		if opt.Short == "h" {
			helpOpts.Short = ""
		}
		if opt.Long == "bash-completion" {
			bashOpts.Long = ""
		}
		if opt.Long == "help" {
			helpOpts.Long = ""
		}
		if opt.Short != "" {
			p.optTable[opt.Short] = OptTableEntry{
				typ:   opt.ArgType,
				setFn: opt.SetFn,
			}
		}
		if opt.Long != "" {
			p.optTable[opt.Long] = OptTableEntry{
				typ:   opt.ArgType,
				setFn: opt.SetFn,
			}
		}
		p.opts = append(p.opts, opt)
		if len(opt.Long) > p.maxLongOptLength {
			p.maxLongOptLength = len(opt.Long)
		}
	}
	if bashOpts.Short != "" {
		p.optTable[bashOpts.Short] = OptTableEntry{
			typ:   bashOpts.ArgType,
			setFn: bashOpts.SetFn,
		}
	}
	if bashOpts.Long != "" {
		p.optTable[bashOpts.Long] = OptTableEntry{
			typ:   bashOpts.ArgType,
			setFn: bashOpts.SetFn,
		}
	}
	if helpOpts.Short != "" {
		p.optTable[helpOpts.Short] = OptTableEntry{
			typ:   helpOpts.ArgType,
			setFn: helpOpts.SetFn,
		}
	}
	if helpOpts.Long != "" {
		p.optTable[helpOpts.Long] = OptTableEntry{
			typ:   helpOpts.ArgType,
			setFn: helpOpts.SetFn,
		}
	}
	p.opts = append(opts, []LongOption{bashOpts, helpOpts}...)
	s := gosort.NewSorter[LongOption](0)
	s.Sort(p.opts, 0)
	return &p
}
