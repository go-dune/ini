//***************************************************************************
// Copyright 2018 OceanScan - Marine Systems & Technology, Lda.             *
//***************************************************************************
// Licensed under the Apache License, Version 2.0 (the "License");          *
// you may not use this file except in compliance with the License.         *
// You may obtain a copy of the License at                                  *
//                                                                          *
// http://www.apache.org/licenses/LICENSE-2.0                               *
//                                                                          *
// Unless required by applicable law or agreed to in writing, software      *
// distributed under the License is distributed on an "AS IS" BASIS,        *
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. *
// See the License for the specific language governing permissions and      *
// limitations under the License.                                           *
//***************************************************************************
// Author: Ricardo Martins                                                  *
//***************************************************************************

package ini

import (
	"reflect"
	"strings"
	"testing"
)

func TestConstructor(t *testing.T) {
	c := NewConfig()
	p := NewParser(c)

	if c != p.Config {
		t.Errorf("expected: %q, actual: %q", p.Config, c)
	}
}

func TestRemoveComments(t *testing.T) {
	// Leading ; comment.
	input := "; Leading comment"
	expected := ""
	actual := removeComments(input)
	if expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}

	// Leading # comment.
	input = "# Leading comment"
	expected = ""
	actual = removeComments(input)
	if expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}

	// Embedded comment.
	input = "Some text # Leading comment"
	expected = "Some text"
	actual = removeComments(input)
	if expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}

	// Comment within comment.
	input = "Some text # Leading comment ; Another comment"
	expected = "Some text"
	actual = removeComments(input)
	if expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}

	// Multiple comment characters.
	input = "########################################################################"
	expected = ""
	actual = removeComments(input)
	if expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}
}

func TestReadSectionName(t *testing.T) {
	// Regular section name.
	input := "[Section]"
	expected := "Section"
	rv, actual := readSectionName(input)
	if !rv || actual != expected {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}

	// Section name with spaces.
	input = "[  Section A   ]"
	expected = "Section A"
	rv, actual = readSectionName(input)
	if !rv || actual != expected {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}

	// Empty section name.
	input = "[ ]"
	expected = ""
	rv, actual = readSectionName(input)
	if !rv || actual != expected {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}
}

func TestReadAssign(t *testing.T) {
	var tests = []struct {
		in    string
		rv    bool
		label string
		value string
	}{
		{"Label 1=More Text", true, "Label 1", "More Text"},
		{"Label 1 = More Text", true, "Label 1", "More Text"},
		{"Label 1     =       More Text", true, "Label 1", "More Text"},
		{"Label 1     = ", true, "Label 1", ""},
		{"Label 1 += More Text", true, "Label 1 +", "More Text"},
	}

	for idx, tt := range tests {
		rv, label, value := readAssign(tt.in)
		if rv != tt.rv || label != tt.label || value != tt.value {
			t.Errorf("idx: %d, expected: %t, %q, %q, actual: %t, %q, %q",
				idx, tt.rv, tt.label, tt.value, rv, label, value)
		}
	}
}

func TestReadAppend(t *testing.T) {
	var tests = []struct {
		in    string
		rv    bool
		label string
		value string
	}{
		{"Label 1+=More Text", true, "Label 1", "More Text"},
		{"Label 1 += More Text", true, "Label 1", "More Text"},
		{"Label 1 = More Text", false, "", ""},
	}

	for idx, tt := range tests {
		rv, label, value := readAppend(tt.in)
		if rv != tt.rv || label != tt.label || value != tt.value {
			t.Errorf("idx: %d, expected: %t, %q, %q, actual: %t, %q, %q",
				idx, tt.rv, tt.label, tt.value, rv, label, value)
		}
	}
}

func TestParseEmptyLabel(t *testing.T) {
	p := NewParser(nil)

	err := p.ParseFile("testdata/invalid_label.ini")
	if err == nil {
		t.Errorf("expected error")
	}

	actual := err.Error()
	expected := "testdata/invalid_label.ini:2: empty label"
	if actual != expected {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}
}

func TestParseEmptySection(t *testing.T) {
	p := NewParser(nil)
	input := "[]\n" +
		"Label 1 += B"

	err := p.Parse(strings.NewReader(input))
	if err == nil {
		t.Errorf("expected error")
	}

	actual := err.Error()
	expected := "1: empty section name"
	if actual != expected {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}
}

func TestParseDoubleAppend(t *testing.T) {
	p := NewParser(nil)
	input := "[Section]\n" +
		"Label 1 += A\n" +
		"Label 1 += B"

	expected := map[string]map[string]string{
		"Section": {"Label 1": "A B"},
	}

	p.Parse(strings.NewReader(input))

	actual := p.Config.Map()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}
}

func TestInvalidSection(t *testing.T) {
	p := NewParser(nil)
	input := "[ ]"

	err := p.Parse(strings.NewReader(input))
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestSimple(t *testing.T) {
	p := NewParser(nil)
	input := "# Some comment\n" +
		"[Section A]\n" +
		"Label A = Value A\n" +
		"\n" +
		"[Section A]\n" +
		"Label A += B C D"

	err := p.Parse(strings.NewReader(input))
	if err != nil {
		t.Errorf("unexpected error")
	}

	expected := "Value A B C D"
	actual := p.Config.Value("Section A", "Label A")
	if actual != expected {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}
}

func TestParseFileError(t *testing.T) {
	p := NewParser(nil)
	err := p.ParseFile("testdata/__no_such_file")
	if err == nil {
		t.Errorf("an error was expected, got none")
	}
}

func TestParseFileInclude(t *testing.T) {
	p := NewParser(nil)
	err := p.ParseFile("testdata/valid00.ini")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := map[string]map[string]string{
		"valid00": {
			"valid00 L0": "valid00 V0",
			"valid00 L1": "valid00 V1",
		},
		"include00": {
			"include00 L0": "include00 V0",
			"include00 L1": "include00 V1",
		},
	}

	actual := p.Config.Map()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nexpected: %q\nactual: %q", expected, actual)
	}
}

func TestParseFileIncludeIgnore(t *testing.T) {
	p := NewParser(nil)
	err := p.ParseFile("testdata/include_ignore.ini")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := map[string]map[string]string{
		"valid00": {
			"valid00 L0": "valid00 V0",
			"valid00 L1": "valid00 V1",
		},
	}

	actual := p.Config.Map()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nexpected: %q\nactual: %q", expected, actual)
	}
}

func TestParseFileRequire(t *testing.T) {
	p := NewParser(nil)
	err := p.ParseFile("testdata/valid01.ini")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := map[string]map[string]string{
		"valid01": {
			"valid01 L0": "valid01 V0",
			"valid01 L1": "valid01 V1",
		},
		"include01": {
			"include01 L0": "include01 V0",
			"include01 L1": "include01 V1",
		},
	}

	actual := p.Config.Map()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nexpected: %q\nactual: %q", expected, actual)
	}
}

func TestRequireLoop(t *testing.T) {
	p := NewParser(nil)
	err := p.ParseFile("testdata/require_loop.ini")
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestIncludeLoop(t *testing.T) {
	p := NewParser(nil)
	err := p.ParseFile("testdata/include_loop.ini")
	if err == nil {
		t.Errorf("expected error")
	}
}
