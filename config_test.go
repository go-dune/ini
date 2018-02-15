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
	"sort"
	"testing"
)

func TestSetGetMap(t *testing.T) {
	c := NewConfig()

	input := map[string]map[string]string{
		"S0": {"L0": "V0", "L1": "V1"},
		"S1": {"L0": "V0"},
	}
	expected := input

	c.SetMap(input)
	actual := c.Map()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestSetGetValue(t *testing.T) {
	c := NewConfig()
	actual := c.Value("A", "B")
	expected := ""
	if expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}

	c.SetValue("A", "B", "C")
	actual = c.Value("A", "B")
	expected = "C"
	if expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}
}

func TestAppendValue(t *testing.T) {
	c := NewConfig()
	c.AppendValue("S0", "L0", "A", "")
	actual := c.Value("S0", "L0")
	expected := "A"
	if expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}

	c.AppendValue("S0", "L0", "B", "")
	actual = c.Value("S0", "L0")
	expected = "AB"
	if expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}
}

func TestLabels(t *testing.T) {
	c := NewConfig()

	c.SetValue("S1", "L1", "V1")
	c.SetValue("S1", "L2", "V2")

	actual := c.Labels("S1")
	expected := []string{"L1", "L2"}
	sort.Strings(actual)

	if len(expected) != len(actual) {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}

	for i := 0; i < len(actual); i++ {
		if expected[i] != actual[i] {
			t.Errorf("%d | expected: %q, actual: %q", i, expected, actual)
		}
	}
}

func TestSetSection(t *testing.T) {
	c := NewConfig()

	c.SetSection("A", map[string]string{"B": "C", "D": "E"})

	actual := c.Value("A", "B")
	expected := "C"
	if expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}

	actual = c.Value("A", "D")
	expected = "E"
	if expected != actual {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}
}

func TestSections(t *testing.T) {
	c := NewConfig()

	c.SetValue("A", "B", "C")
	c.SetValue("1", "2", "3")

	actual := c.Sections()
	expected := []string{"1", "A"}
	sort.Strings(actual)

	if len(expected) != len(actual) {
		t.Errorf("expected: %q, actual: %q", expected, actual)
	}

	for i := 0; i < len(actual); i++ {
		if expected[i] != actual[i] {
			t.Errorf("%d | expected: %q, actual: %q", i, expected, actual)
		}
	}
}
