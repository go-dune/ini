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
	"fmt"
	"sort"
)

func ExampleConfig_SetValue() {
	c := NewConfig()
	c.SetValue("Section From SetValue", "Label 0", "Value 0")
	fmt.Println(c.Value("Section From SetValue", "Label 0"))
	// Output: Value 0
}

func ExampleConfig_SetMap() {
	c := NewConfig()
	c.SetMap(map[string]map[string]string{
		"Section 0": {
			"Label 0": "Value 0",
		},
	})

	fmt.Println(c.Value("Section 0", "Label 0"))
	// Output: Value 0
}

func ExampleConfig_Sections() {
	c := NewConfig()
	c.SetValue("Section 0", "Label 0", "Value 0")
	c.SetValue("Section 1", "Label 1", "Value 1")
	sections := c.Sections()
	sort.Strings(sections)
	fmt.Printf("%q", sections)
	// Output: ["Section 0" "Section 1"]
}
