package getopt

import (
	"fmt"
	"reflect"
	"testing"
)

func testCaseOne(t *testing.T, osArgs []string, expectArr []string) {
	resultArr := []string{}

	optHandler := func(option *Option, nametype NameType, value string) error {
		switch nametype {
		case ShortName:
			resultArr = append(resultArr, fmt.Sprintf("-%c", option.ShortName))
		case LongName:
			resultArr = append(resultArr, fmt.Sprintf("--%s", option.LongName))
		}
		if option.HasArg != NoArgument {
			resultArr = append(resultArr, fmt.Sprintf("{%s}", value))
		}
		return nil
	}

	getopt := &Getopt{
		Options: []Option{
			{'x', NoLongName, NoArgument, optHandler},
			{NoShortName, "xyz", NoArgument, optHandler},
			{'h', "help", NoArgument, optHandler},
			{'V', "version", NoArgument, optHandler},
			{'a', "caa", NoArgument, optHandler},
			{'b', "cba", RequiredArgument, optHandler},
		},
	}

	if err := getopt.Parse(osArgs); err != nil {
		t.Fatal(err)
	}

	resultArr = append(resultArr, "--")

	for _, arg := range getopt.Args() {
		resultArr = append(resultArr, fmt.Sprintf("{%s}", arg))
	}

	if !reflect.DeepEqual(expectArr, resultArr) {
		t.Fatalf("unexpected: %#v", resultArr)
	}
}

func testCaseAbbrev(t *testing.T, osArgs []string, expectArr []string) {
	resultArr := []string{}

	optHandler := func(option *Option, nametype NameType, value string) error {
		switch nametype {
		case ShortName:
			resultArr = append(resultArr, fmt.Sprintf("-%c", option.ShortName))
		case LongName:
			resultArr = append(resultArr, fmt.Sprintf("--%s", option.LongName))
		}
		if option.HasArg != NoArgument {
			resultArr = append(resultArr, fmt.Sprintf("{%s}", value))
		}
		return nil
	}

	getopt := &Getopt{
		AllowAbbrev: true,
		Options: []Option{
			{'h', "help", NoArgument, optHandler},
			{'V', "version", NoArgument, optHandler},
			{'a', "daa", NoArgument, optHandler},
			{'b', "cba", RequiredArgument, optHandler},
		},
	}

	if err := getopt.Parse(osArgs); err != nil {
		t.Fatal(err)
	}

	resultArr = append(resultArr, "--")

	for _, arg := range getopt.Args() {
		resultArr = append(resultArr, fmt.Sprintf("{%s}", arg))
	}

	if !reflect.DeepEqual(expectArr, resultArr) {
		t.Fatalf("unexpected: %#v", resultArr)
	}
}

func testCaseAlternative(t *testing.T, osArgs []string, expectArr []string) {
	resultArr := []string{}

	optHandler := func(option *Option, nametype NameType, value string) error {
		switch nametype {
		case ShortName:
			resultArr = append(resultArr, fmt.Sprintf("-%c", option.ShortName))
		case LongName:
			resultArr = append(resultArr, fmt.Sprintf("--%s", option.LongName))
		}
		if option.HasArg != NoArgument {
			resultArr = append(resultArr, fmt.Sprintf("{%s}", value))
		}
		return nil
	}

	getopt := &Getopt{
		AllowAlternative: true,
		Options: []Option{
			{'h', "help", NoArgument, optHandler},
			{'V', "version", NoArgument, optHandler},
			{'a', "daa", NoArgument, optHandler},
			{'b', "cba", RequiredArgument, optHandler},
		},
	}

	if err := getopt.Parse(osArgs); err != nil {
		t.Fatal(err)
	}

	resultArr = append(resultArr, "--")

	for _, arg := range getopt.Args() {
		resultArr = append(resultArr, fmt.Sprintf("{%s}", arg))
	}

	if !reflect.DeepEqual(expectArr, resultArr) {
		t.Fatalf("unexpected: %#v", resultArr)
	}
}

func testCaseOptional(t *testing.T, osArgs []string, expectArr []string) {
	resultArr := []string{}

	optHandler := func(option *Option, nametype NameType, value string) error {
		switch nametype {
		case ShortName:
			resultArr = append(resultArr, fmt.Sprintf("-%c", option.ShortName))
		case LongName:
			resultArr = append(resultArr, fmt.Sprintf("--%s", option.LongName))
		}
		if option.HasArg != NoArgument {
			resultArr = append(resultArr, fmt.Sprintf("{%s}", value))
		}
		return nil
	}

	getopt := &Getopt{
		AllowAbbrev: true,
		Options: []Option{
			{'h', "help", NoArgument, optHandler},
			{'V', "version", NoArgument, optHandler},
			{'a', "daa", NoArgument, optHandler},
			{'b', "cba", OptionalArgument, optHandler},
		},
	}

	if err := getopt.Parse(osArgs); err != nil {
		t.Fatal(err)
	}

	resultArr = append(resultArr, "--")

	for _, arg := range getopt.Args() {
		resultArr = append(resultArr, fmt.Sprintf("{%s}", arg))
	}

	if !reflect.DeepEqual(expectArr, resultArr) {
		t.Fatalf("unexpected: %#v", resultArr)
	}
}

func TestShortOption(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-h", "-V"},
		[]string{"-h", "-V", "--"},
	)
}

func TestShortOptionWithInlineArgument(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-abcd"},
		[]string{"-a", "-b", "{cd}", "--"},
	)
}

func TestShortOptionWithStandaloneArgument(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-ab", "cd"},
		[]string{"-a", "-b", "{cd}", "--"},
	)
}

func TestShortOptionWithEqualArgument(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-ab=cd"},
		[]string{"-a", "-b", "{cd}", "--"},
	)
}

func TestShortOptionWithEqualEmptyArgument(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-ab="},
		[]string{"-a", "-b", "{}", "--"},
	)
}

func TestShortOptionWithArgumentTwo(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-a", "-b", "-c"},
		[]string{"-a", "-b", "{-c}", "--"},
	)
}

func TestParameterBetweenShortOptions(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-a", "XXX", "-bcd", "ZZZ"},
		[]string{"-a", "-b", "{cd}", "--", "{XXX}", "{ZZZ}"},
	)
}

func TestAllParameters(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "--", "-a", "XXX", "-bcd"},
		[]string{"--", "{-a}", "{XXX}", "{-bcd}"},
	)
}

func TestHalfParameters(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-a", "-h", "--", "XXX", "-bcd"},
		[]string{"-a", "-h", "--", "{XXX}", "{-bcd}"},
	)
}

func TestLongOptionsWithArgument(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "--caa", "--cba", "XXX"},
		[]string{"--caa", "--cba", "{XXX}", "--"},
	)
}

func TestLongOptionsWithEqualArgument(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "--caa", "--cba=XXX"},
		[]string{"--caa", "--cba", "{XXX}", "--"},
	)
}

func TestLongOptionsWithEqualEmpty(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "--caa", "--cba="},
		[]string{"--caa", "--cba", "{}", "--"},
	)
}

func TestAbbrevLongOptions(t *testing.T) {
	testCaseAbbrev(
		t,
		[]string{"PROG", "--d", "--c=XXX"},
		[]string{"--daa", "--cba", "{XXX}", "--"},
	)
}

func TestAbbrevLongOptionsWithEqualArgument(t *testing.T) {
	testCaseAbbrev(
		t,
		[]string{"PROG", "--daa", "--cb=XXX"},
		[]string{"--daa", "--cba", "{XXX}", "--"},
	)
}

func TestAlternativeLongOptionsWithArgument(t *testing.T) {
	testCaseAlternative(
		t,
		[]string{"PROG", "--daa", "-cba", "XXX"},
		[]string{"--daa", "--cba", "{XXX}", "--"},
	)
}

func TestAlternativeLongOptionsWithEqualArgument(t *testing.T) {
	testCaseAlternative(
		t,
		[]string{"PROG", "-cba=XXX"},
		[]string{"--cba", "{XXX}", "--"},
	)
}

func TestOptionalLongOption(t *testing.T) {
	testCaseOptional(
		t,
		[]string{"PROG", "--cba"},
		[]string{"--cba", "{}", "--"},
	)
}

func TestOptionalShortOption(t *testing.T) {
	testCaseOptional(
		t,
		[]string{"PROG", "-b"},
		[]string{"-b", "{}", "--"},
	)
}

func TestOptionalShortOptionTwo(t *testing.T) {
	testCaseOptional(
		t,
		[]string{"PROG", "-b", "-v"},
		[]string{"-b", "{-v}", "--"},
	)
}
