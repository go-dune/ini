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

// Package ini provides routines to parse configuration files encoded in a format
// akin to Microsoft Windows INI files.
package ini

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Compiled regular expressions.
var (
	// Match comments.
	reComment = regexp.MustCompile(`[;|#]+(.*)$`)
	// Match sections.
	reSection = regexp.MustCompile(`\[([^]]+)]`)
	// Match append instructions.
	reAppend = regexp.MustCompile(`([^+]+)\+=(.*)`)
	// Match assignment instructions.
	reAssign = regexp.MustCompile(`([^=]+)=(.*)`)
)

// Parser is an INI format parser.
type Parser struct {
	Config       *Config         // Configuration instance.
	curSection   string          // Section being parsed.
	curLabel     string          // Label being parsed.
	fileStack    []string        // File stack, top is file being parsed.
	lineNrStack  []uint          // Line number stack.
	visitedFiles map[string]bool // Set of visited files.
}

// NewParser creates a new instance of Parser.
func NewParser(c *Config) *Parser {
	p := new(Parser)
	p.visitedFiles = make(map[string]bool)
	if c == nil {
		p.Config = NewConfig()
	} else {
		p.Config = c
	}

	return p
}

// Parse parses an INI format stream.
func (p *Parser) Parse(reader io.Reader) error {
	return p.parseReader(reader, "")
}

// ParseFile parses an INI format file.
func (p *Parser) ParseFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return p.parseReader(file, path)
}

func (p *Parser) parseReader(reader io.Reader, path string) error {
	err := p.pushFile(path)
	if err != nil {
		return err
	}

	defer p.popFile()

	bio := bufio.NewReader(reader)

	eof := false
	for !eof {
		line, err := bio.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				eof = true
			}
		}

		err = p.handleLine(line)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) pushFile(path string) error {
	absPath, _ := filepath.Abs(path)

	_, exists := p.visitedFiles[absPath]
	if exists {
		return &SyntaxError{p.curFile(), p.curLineNr(), "include loop"}
	}

	log.Printf("parsing %v\n", path)
	p.visitedFiles[absPath] = true
	p.fileStack = append(p.fileStack, path)
	p.lineNrStack = append(p.lineNrStack, 0)
	return nil
}

func (p *Parser) popFile() {
	stackLen := len(p.fileStack)
	if stackLen > 0 {
		p.fileStack = p.fileStack[:stackLen-1]
		p.lineNrStack = p.lineNrStack[:stackLen-1]
	}
}

func (p *Parser) curFile() string {
	return p.fileStack[len(p.fileStack)-1]
}

func (p *Parser) incLineNr() {
	p.lineNrStack[len(p.lineNrStack)-1]++
}

func (p *Parser) curLineNr() uint {
	return p.lineNrStack[len(p.lineNrStack)-1]
}

func removeComments(line string) string {
	return strings.TrimSpace(reComment.ReplaceAllString(line, ""))
}

func readSectionName(line string) (bool, string) {
	matches := reSection.FindStringSubmatch(line)
	if len(matches) == 2 {
		return true, strings.TrimSpace(matches[1])
	}

	return false, ""
}

func readLabelValue(re *regexp.Regexp, line string) (bool, string, string) {
	matches := re.FindStringSubmatch(line)
	if len(matches) == 3 {
		return true, strings.TrimSpace(matches[1]), strings.TrimSpace(matches[2])
	}

	return false, "", ""
}

func readAppend(line string) (bool, string, string) {
	return readLabelValue(reAppend, line)
}

func readAssign(line string) (bool, string, string) {
	return readLabelValue(reAssign, line)
}

func (p *Parser) insertValue(section string, label string, value string, append bool) error {
	if section == "" {
		return &SyntaxError{p.curFile(), p.curLineNr(), "empty section name"}
	}

	if label == "" {
		return &SyntaxError{p.curFile(), p.curLineNr(), "empty label"}
	}

	if append {
		p.Config.AppendValue(section, label, value, " ")
	} else {
		p.Config.SetValue(section, label, value)
	}

	return nil
}

func (p *Parser) resolveIncludeFile(path string) string {
	curFolder := filepath.Dir(p.curFile())
	return filepath.Join(curFolder, strings.TrimSpace(path))
}

func (p *Parser) handleInclude(line string) error {
	incPath := p.resolveIncludeFile(strings.TrimPrefix(line, "Include "))

	file, err := os.Open(incPath)
	if err != nil {
		return nil
	}
	defer file.Close()

	return p.parseReader(file, incPath)
}

func (p *Parser) handleRequire(line string) error {
	incPath := p.resolveIncludeFile(strings.TrimPrefix(line, "Require "))
	return p.ParseFile(incPath)
}

func (p *Parser) setCurSection(section string) error {
	if section == "" {
		return &SyntaxError{p.curFile(), p.curLineNr(), "empty section name"}
	}

	p.curSection = section
	return nil
}

func (p *Parser) handleLine(line string) error {
	p.incLineNr()

	// Remove comments and clean string.
	cleanLine := removeComments(strings.TrimSpace(line))
	if cleanLine == "" {
		return nil
	}

	// Section.
	secRv, secName := readSectionName(cleanLine)
	if secRv {
		if strings.HasPrefix(secName, "Require ") {
			return p.handleRequire(secName)
		} else if strings.HasPrefix(secName, "Include ") {
			return p.handleInclude(secName)
		} else {
			return p.setCurSection(secName)
		}
	}

	// Append operator.
	apRv, apLabel, apValue := readAppend(cleanLine)
	if apRv {
		p.curLabel = apLabel
		return p.insertValue(p.curSection, apLabel, apValue, true)
	}

	// Assign operator.
	asRv, asLabel, asValue := readAssign(cleanLine)
	if asRv {
		p.curLabel = asLabel
		return p.insertValue(p.curSection, asLabel, asValue, false)
	}

	// Multi-line value.
	return p.insertValue(p.curSection, p.curLabel, asValue, true)
}
