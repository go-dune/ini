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

import "sync"

// Config organizes configuration values in sections. Each section may
// contain an arbitrary number of unique labels with associated values.
// Instances of Config are thread-safe.
type Config struct {
	cfg  map[string]map[string]string
	lock sync.RWMutex
}

// NewConfig creates a new instance of Config.
func NewConfig() *Config {
	c := new(Config)
	c.cfg = make(map[string]map[string]string)
	return c
}

func cloneMap(m map[string]string) map[string]string {
	newM := make(map[string]string)
	for key, value := range m {
		newM[key] = value
	}

	return newM
}

func cloneTable(t map[string]map[string]string) map[string]map[string]string {
	newT := make(map[string]map[string]string)

	for key, value := range t {
		newT[key] = cloneMap(value)
	}

	return newT
}

func (c *Config) setValueNoLock(section string, label string, value string) {
	_, exists := c.cfg[section]
	if !exists {
		c.cfg[section] = make(map[string]string)
	}

	c.cfg[section][label] = value
}

// Map returns a copy of the configuration contents.
func (c *Config) Map() map[string]map[string]string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return cloneTable(c.cfg)
}

// SetMap replaces the current configuration with a copy of map m.
func (c *Config) SetMap(m map[string]map[string]string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cfg = cloneTable(m)
}

// Sections returns an unordered slice of all section names.
func (c *Config) Sections() []string {
	var sections []string

	c.lock.RLock()
	defer c.lock.RUnlock()

	for section := range c.cfg {
		sections = append(sections, section)
	}

	return sections
}

// SetSection replaces the current contents of section s with a copy of map m.
func (c *Config) SetSection(s string, m map[string]string) {
	newM := cloneMap(m)

	c.lock.Lock()
	defer c.lock.Unlock()

	c.cfg[s] = newM
}

// Labels returns an unordered slice of all labels of a section s.
func (c *Config) Labels(s string) []string {
	var labels []string

	c.lock.RLock()
	defer c.lock.RUnlock()

	for label := range c.cfg[s] {
		labels = append(labels, label)
	}

	return labels
}

// Value retrieves the value of label l of section s.
func (c *Config) Value(s string, l string) string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.cfg[s][l]
}

// SetValue assigns value to label l of section s.
func (c *Config) SetValue(s string, l string, value string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.setValueNoLock(s, l, value)
}

// AppendValue appends value to label l of section s. The new value will be the concatenation of
// the old value, sep, and value. If no contents exist this function behaves like SetValue.
func (c *Config) AppendValue(s string, l string, value string, sep string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	curValue := c.cfg[s][l]
	if curValue == "" {
		c.setValueNoLock(s, l, value)
	} else {
		c.setValueNoLock(s, l, curValue+sep+value)
	}
}
