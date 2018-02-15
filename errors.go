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

import "fmt"

// SyntaxError represents a parsing error.
type SyntaxError struct {
	path   string // File path.
	lineNr uint   // Line number.
	msg    string // Error description.
}

// Error formats the error to a human readable sentence.
func (e *SyntaxError) Error() string {
	if e.path == "" {
		return fmt.Sprintf("%d: %s", e.lineNr, e.msg)
	}

	return fmt.Sprintf("%s:%d: %s", e.path, e.lineNr, e.msg)
}
