package getopt

import (
	"fmt"
	"reflect"
	"testing"
)

func testCaseOne(t *testing.T, osArgs []string, expectArr []string) {
	resultArr := []string{}

	optHandler := func(option *GetoptOption, nametype GetoptNameType, value string) error {
		switch nametype {
		case GetoptShortName:
			resultArr = append(resultArr, fmt.Sprintf("-%c", option.ShortName))
		case GetoptLongName:
			resultArr = append(resultArr, fmt.Sprintf("--%s", option.LongName))
		}
		if option.HasArg != GetoptNoArgument {
			resultArr = append(resultArr, fmt.Sprintf("{%s}", value))
		}
		return nil
	}

	getopt := &Getopt{
		Options: []GetoptOption{
			{'x', GetoptNoLongName, GetoptNoArgument, optHandler},
			{GetoptNoShortName, "xyz", GetoptNoArgument, optHandler},
			{'h', "help", GetoptNoArgument, optHandler},
			{'V', "version", GetoptNoArgument, optHandler},
			{'a', "caa", GetoptNoArgument, optHandler},
			{'b', "cba", GetoptRequiredArgument, optHandler},
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

	optHandler := func(option *GetoptOption, nametype GetoptNameType, value string) error {
		switch nametype {
		case GetoptShortName:
			resultArr = append(resultArr, fmt.Sprintf("-%c", option.ShortName))
		case GetoptLongName:
			resultArr = append(resultArr, fmt.Sprintf("--%s", option.LongName))
		}
		if option.HasArg != GetoptNoArgument {
			resultArr = append(resultArr, fmt.Sprintf("{%s}", value))
		}
		return nil
	}

	getopt := &Getopt{
		AllowAbbrev: true,
		Options: []GetoptOption{
			{'h', "help", GetoptNoArgument, optHandler},
			{'V', "version", GetoptNoArgument, optHandler},
			{'a', "daa", GetoptNoArgument, optHandler},
			{'b', "cba", GetoptRequiredArgument, optHandler},
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

	optHandler := func(option *GetoptOption, nametype GetoptNameType, value string) error {
		switch nametype {
		case GetoptShortName:
			resultArr = append(resultArr, fmt.Sprintf("-%c", option.ShortName))
		case GetoptLongName:
			resultArr = append(resultArr, fmt.Sprintf("--%s", option.LongName))
		}
		if option.HasArg != GetoptNoArgument {
			resultArr = append(resultArr, fmt.Sprintf("{%s}", value))
		}
		return nil
	}

	getopt := &Getopt{
		AllowAlternative: true,
		Options: []GetoptOption{
			{'h', "help", GetoptNoArgument, optHandler},
			{'V', "version", GetoptNoArgument, optHandler},
			{'a', "daa", GetoptNoArgument, optHandler},
			{'b', "cba", GetoptRequiredArgument, optHandler},
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

	optHandler := func(option *GetoptOption, nametype GetoptNameType, value string) error {
		switch nametype {
		case GetoptShortName:
			resultArr = append(resultArr, fmt.Sprintf("-%c", option.ShortName))
		case GetoptLongName:
			resultArr = append(resultArr, fmt.Sprintf("--%s", option.LongName))
		}
		if option.HasArg != GetoptNoArgument {
			resultArr = append(resultArr, fmt.Sprintf("{%s}", value))
		}
		return nil
	}

	getopt := &Getopt{
		AllowAbbrev: true,
		Options: []GetoptOption{
			{'h', "help", GetoptNoArgument, optHandler},
			{'V', "version", GetoptNoArgument, optHandler},
			{'a', "daa", GetoptNoArgument, optHandler},
			{'b', "cba", GetoptOptionalArgument, optHandler},
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

func TestGetoptShortOption(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-h", "-V"},
		[]string{"-h", "-V", "--"},
	)
}

func TestGetoptShortOptionWithInlineArgument(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-abcd"},
		[]string{"-a", "-b", "{cd}", "--"},
	)
}

func TestGetoptShortOptionWithStandaloneArgument(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-ab", "cd"},
		[]string{"-a", "-b", "{cd}", "--"},
	)
}

func TestGetoptShortOptionWithEqualArgument(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-ab=cd"},
		[]string{"-a", "-b", "{cd}", "--"},
	)
}

func TestGetoptShortOptionWithEqualEmptyArgument(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-ab="},
		[]string{"-a", "-b", "{}", "--"},
	)
}

func TestGetoptShortOptionWithArgumentTwo(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-a", "-b", "-c"},
		[]string{"-a", "-b", "{-c}", "--"},
	)
}

func TestGetoptParameterBetweenShortOptions(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-a", "XXX", "-bcd", "ZZZ"},
		[]string{"-a", "-b", "{cd}", "--", "{XXX}", "{ZZZ}"},
	)
}

func TestGetoptAllParameters(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "--", "-a", "XXX", "-bcd"},
		[]string{"--", "{-a}", "{XXX}", "{-bcd}"},
	)
}

func TestGetoptHalfParameters(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "-a", "-h", "--", "XXX", "-bcd"},
		[]string{"-a", "-h", "--", "{XXX}", "{-bcd}"},
	)
}

func TestGetoptLongOptionsWithArgument(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "--caa", "--cba", "XXX"},
		[]string{"--caa", "--cba", "{XXX}", "--"},
	)
}

func TestGetoptLongOptionsWithEqualArgument(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "--caa", "--cba=XXX"},
		[]string{"--caa", "--cba", "{XXX}", "--"},
	)
}

func TestGetoptLongOptionsWithEqualEmpty(t *testing.T) {
	testCaseOne(
		t,
		[]string{"PROG", "--caa", "--cba="},
		[]string{"--caa", "--cba", "{}", "--"},
	)
}

func TestGetoptAbbrevLongOptions(t *testing.T) {
	testCaseAbbrev(
		t,
		[]string{"PROG", "--d", "--c=XXX"},
		[]string{"--daa", "--cba", "{XXX}", "--"},
	)
}

func TestGetoptAbbrevLongOptionsWithEqualArgument(t *testing.T) {
	testCaseAbbrev(
		t,
		[]string{"PROG", "--daa", "--cb=XXX"},
		[]string{"--daa", "--cba", "{XXX}", "--"},
	)
}

func TestGetoptAlternativeLongOptionsWithArgument(t *testing.T) {
	testCaseAlternative(
		t,
		[]string{"PROG", "--daa", "-cba", "XXX"},
		[]string{"--daa", "--cba", "{XXX}", "--"},
	)
}

func TestGetoptAlternativeLongOptionsWithEqualArgument(t *testing.T) {
	testCaseAlternative(
		t,
		[]string{"PROG", "-cba=XXX"},
		[]string{"--cba", "{XXX}", "--"},
	)
}

func TestGetoptOptionalLongOption(t *testing.T) {
	testCaseOptional(
		t,
		[]string{"PROG", "--cba"},
		[]string{"--cba", "{}", "--"},
	)
}

func TestGetoptOptionalShortOption(t *testing.T) {
	testCaseOptional(
		t,
		[]string{"PROG", "-b"},
		[]string{"-b", "{}", "--"},
	)
}

func TestGetoptOptionalShortOptionTwo(t *testing.T) {
	testCaseOptional(
		t,
		[]string{"PROG", "-b", "-v"},
		[]string{"-b", "{-v}", "--"},
	)
}
