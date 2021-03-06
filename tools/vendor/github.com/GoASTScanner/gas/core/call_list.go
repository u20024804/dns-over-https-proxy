//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"go/ast"
)

type set map[string]bool

/// CallList is used to check for usage of specific packages
/// and functions.
type CallList map[string]set

/// NewCallList creates a new empty CallList
func NewCallList() CallList {
	return make(CallList)
}

/// AddAll will add several calls to the call list at once
func (c CallList) AddAll(selector string, idents ...string) {
	for _, ident := range idents {
		c.Add(selector, ident)
	}
}

/// Add a selector and call to the call list
func (c CallList) Add(selector, ident string) {
	if _, ok := c[selector]; !ok {
		c[selector] = make(set)
	}
	c[selector][ident] = true
}

/// Contains returns true if the package and function are
/// members of this call list.
func (c CallList) Contains(selector, ident string) bool {
	if idents, ok := c[selector]; ok {
		_, found := idents[ident]
		return found
	}
	return false
}

/// ContainsCallExpr resolves the call expression name and type
/// or package and determines if it exists within the CallList
func (c CallList) ContainsCallExpr(n ast.Node, ctx *Context) bool {
	selector, ident, err := GetCallInfo(n, ctx)
	if err != nil {
		return false
	}
	// Try direct resolution
	if c.Contains(selector, ident) {
		return true
	}

	// Also support explicit path
	if path, ok := GetImportPath(selector, ctx); ok {
		return c.Contains(path, ident)
	}
	return false
}
