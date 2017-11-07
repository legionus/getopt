// Package getopt provides simple GNU-like parser for the command-line arguments.
//
// A simple short option is a '-' followed by a short option character. If the option
// has a required argument, it may be written directly after the option character
// or as the next parameter (i.e. separated by whitespace on the command line). If the option
// has an optional argument, it must be written directly after the option character if present.
//
// It is possible to specify several short options after one '-', as long as all (except possibly
// the last) do not have required or optional arguments.
//
// A long option begins with '--' followed by the long option name. If the option has a required
// argument, it may be written directly after the long option name, separated by '=', or as the next
// argument (i.e. separated by whitespace on the command line). If the option has an optional argument,
// it must be written directly after the long option name, separated by '=', if present.
// Long options may be abbreviated, as long as the abbreviation is not ambiguous.
//
// Each parameter not starting with a '-', and not a required argument of a previous option,
// is a non-option parameter.  Each parameter after a '--' parameter is always interpreted
// as a non-option parameter.
//
package getopt

import (
	"fmt"
	"strings"
)

const (
	NoShortName = 0x0
	NoLongName  = ""
)

// ArgumentType
type ArgumentType int

const (
	// NoArgument means the option does not take an argument.
	NoArgument ArgumentType = iota
	// RequiredArgument means the option requires an argument.
	RequiredArgument
	// OptionalArgument means the option takes an optional argument.
	OptionalArgument
)

type NameType int

const (
	_ NameType = iota
	ShortName
	LongName
)

// Option describes command-line option, his short and long form.
type Option struct {
	// ShortName specifies short form of the option. If there is no such form, it should be NoShortName.
	ShortName byte
	// LongName specifies long form of the option. If there is no such form, it should be NoLongName.
	LongName string
	// HasArg describes the need to have the argument. Option may not require additional arguments (NoArgument),
	// Option may require an additional argument (RequiredArgument) or the argument may be optional (OptionalArgument).
	HasArg ArgumentType
	// Handler specifies the handler that will be called if the option is specified on the command line.
	Handler OptionFunc
}

type OptionFunc func(*Option, NameType, string) error

type Getopt struct {
	// AllowAlternative allows long options to start with a single `-'. See (getopt -a).
	AllowAlternative bool
	// AllowAbbrev allows long options be abbreviated, as long as the abbreviation is not ambiguous.
	AllowAbbrev bool
	// Options describes short and long options.
	Options []Option
	args    []string
}

func (g Getopt) getShortOption(c byte, options []Option) (*Option, error) {
	for _, option := range options {
		if option.ShortName == NoShortName {
			continue
		}
		if c == option.ShortName {
			return &option, nil
		}
	}
	if g.AllowAlternative {
		return nil, nil
	}
	return nil, fmt.Errorf("invalid option -- '%c'", c)
}

func (g Getopt) getLongOption(name string, options []Option) (*Option, error) {
	var ret *Option
	for i := range options {
		option := options[i]
		if option.LongName == NoLongName {
			continue
		}
		if g.AllowAbbrev {
			if strings.HasPrefix(option.LongName, name) {
				if ret != nil {
					return nil, fmt.Errorf("option '--%s' is ambiguous; possibilities: '--%s' '--%s'", name, ret.LongName, option.LongName)
				}
				ret = &option
			}
			continue
		}
		if name == option.LongName {
			ret = &option
			break
		}
	}
	if ret != nil {
		return ret, nil
	}
	return nil, fmt.Errorf("unrecognized option -- '%s'", name)
}

func (g Getopt) splitArg(s string) (int, string, string) {
	i := strings.IndexByte(s, '=')
	if i > 0 {
		return i, s[0:i], s[i+1:]
	}
	return 0, s, ""
}

func (g Getopt) Args() []string {
	return g.args
}

func (g *Getopt) Parse(args []string) error {
	optind := 1
	for ; optind < len(args); optind++ {
		if args[optind] == "--" {
			break
		}

		var argType NameType
		if len(args[optind]) > 1 {
			if args[optind][0] == '-' {
				argType = ShortName
				if args[optind][1] == '-' {
					argType = LongName
				}
			}
		}

		var (
			v  string
			eq int
		)
		if argType == ShortName {
			var (
				option *Option
				err    error
				i      int
				n      rune
			)
			eq, args[optind], v = g.splitArg(args[optind])
			for i, n = range args[optind][1:] {
				option, err = g.getShortOption(byte(n), g.Options)
				if err != nil {
					return err
				} else if option == nil {
					argType = LongName
					args[optind] = "-" + args[optind]
					if eq > 0 {
						args[optind] += "=" + v
					}
					goto longArg
				}
				if option.HasArg != NoArgument {
					i++
					break
				}
				if err := option.Handler(option, ShortName, ""); err != nil {
					return err
				}
				option = nil
			}

			if option == nil {
				continue
			}

			if option.HasArg != NoArgument {
				if eq > 0 {
					if err := option.Handler(option, ShortName, v); err != nil {
						return err
					}
					continue
				}
				if i < len(args[optind][1:]) {
					if err := option.Handler(option, ShortName, args[optind][i+1:]); err != nil {
						return err
					}
					continue
				}
				if optind+1 < len(args) {
					optind++
					if err := option.Handler(option, ShortName, args[optind]); err != nil {
						return err
					}
					continue
				}
			}

			if option.HasArg == RequiredArgument {
				return fmt.Errorf("option requires an argument -- '%c'", option.ShortName)
			}
			if err := option.Handler(option, ShortName, ""); err != nil {
				return err
			}
			continue
		}
	longArg:
		if argType == LongName {
			eq, args[optind], v = g.splitArg(args[optind])

			option, err := g.getLongOption(args[optind][2:], g.Options)
			if err != nil {
				return err
			}

			if option.HasArg != NoArgument {
				if eq > 0 {
					if err := option.Handler(option, LongName, v); err != nil {
						return err
					}
					continue
				}
				if optind+1 < len(args) {
					optind++
					if err := option.Handler(option, LongName, args[optind]); err != nil {
						return err
					}
					continue
				}
			}

			if option.HasArg == RequiredArgument {
				return fmt.Errorf("option '--%s' requires an argument", option.LongName)
			}
			if err := option.Handler(option, LongName, ""); err != nil {
				return err
			}
			continue
		}

		g.args = append(g.args, args[optind])
	}

	for optind++; optind < len(args); optind++ {
		g.args = append(g.args, args[optind])
	}

	return nil
}
