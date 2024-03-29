// Copyright 2019 leiysky
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ast

import (
	"fmt"
	"io"
	"strings"
)

// RestoreContext is used for restore Cypher
type RestoreContext struct {
	w io.Writer
}

// NewRestoreContext will create a RestoreContext
func NewRestoreContext(w io.Writer) *RestoreContext {
	return &RestoreContext{
		w: w,
	}
}

// Write writes plain text into io.Writer.
func (c *RestoreContext) Write(v ...interface{}) {
	fmt.Fprint(c.w, v...)
}

// WriteIdent writes an identifier with backquote into io.Writer.
func (c *RestoreContext) WriteIdent(s string) {
	fmt.Fprint(c.w, "`", s, "`")
}

// WriteKeyword writes s with full upper case into io.Writer.
func (c *RestoreContext) WriteKeyword(s string) {
	fmt.Fprint(c.w, strings.ToUpper(s))
}

// WriteString writes single-quoted string literal into io.Writer.
func (c *RestoreContext) WriteString(s string) {
	fmt.Fprint(c.w, "'", s, "'")
}

// Writef writes values with format into io.Writer.
func (c *RestoreContext) Writef(format string, v ...interface{}) {
	fmt.Fprintf(c.w, format, v...)
}
