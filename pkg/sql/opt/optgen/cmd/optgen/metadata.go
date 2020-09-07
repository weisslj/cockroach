// Copyright 2018 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"fmt"
	"strings"

	"github.com/weisslj/cockroach/pkg/sql/opt/optgen/lang"
)

// metadata generates and stores the mapping from Optgen language expressions to
// the Go types that the code generators use to implement those expressions. In
// addition, when generating strongly-typed Go code, it is often necessary to
// insert casts, struct referencing (&), and pointer dereferences (*). Since the
// Optgen code generators do not have access to regular Go metadata for these
// types, aspects of it must be replicated here.
type metadata struct {
	compiled  *lang.CompiledExpr
	types     map[string]*typeDef
	exprTypes map[lang.Expr]*typeDef
}

// typeDef describes the name and characteristics of each Go type used by the
// code generators. This is used by the Optgen code generators to generate
// correct strongly-typed code.
type typeDef struct {
	// name is the correctly qualified Go name of the type, as it should appear
	// in the current package. For example, if the current package is "memo":
	//
	//   RelExpr
	//   opt.Expr
	//   FiltersExpr
	//   ScanLimit
	//   *tree.Subquery
	//   types.T
	//
	name string

	// fullName is the fully qualified Go name of the type, irrespective of the
	// current package. For example, if the current package is "memo":
	//
	//   memo.RelExpr
	//   opt.Expr
	//   memo.FiltersExpr
	//   memo.ScanLimit
	//   *tree.Subquery
	//   types.T
	//
	fullName string

	// friendlyName is a human-friendly name for the type that is a simple
	// identifier having no special characters like "." or "*". It is used in .opt
	// files to refer to types without needing package names and type modifiers
	// like "*". It's also used to generate methods based on the type, like Intern
	// methods, where special characters are not allowed. For example:
	//
	//   RelExpr
	//   Expr
	//   FiltersExpr
	//   ScanLimit
	//   Subquery
	//   DatumType
	//
	friendlyName string

	// isExpr is true if the type is either one of the expression-related
	// interfaces (opt.Expr, memo.RelExpr, opt.ScalarExpr), or one of the auto-
	// generated expression types (memo.ScanExpr, memo.AndExpr, etc.).
	isExpr bool

	// isPointer is true if the type is a Go pointer or interface type (types.T,
	// memo.RelExpr, *tree.Subquery, etc.).
	isPointer bool

	// usePointerIntern is true if the type should be treated as a pointer during
	// interning, meaning that an instance will be hashed by address rather than
	// by value.
	usePointerIntern bool

	// passByVal is true if the type should be passed by value to custom functions,
	// as well as stored by value in structs and variables.
	passByVal bool

	// isGenerated is true if the type's definition was auto-generated by Optgen,
	// vs. manually defined.
	isGenerated bool

	// listItemType links to the type of items in the list, if this type is a
	// a list type (e.g. memo.FiltersExpr). If this is not a list type, then
	// listItemType is nil.
	listItemType *typeDef
}

// isListType is true if this type is represented as a Go slice. For example:
//
//   type FiltersExpr []FiltersItem
//
func (t *typeDef) isListType() bool {
	return t.listItemType != nil
}

// asParam returns the Go type used to pass this Optgen type around as a
// parameter. For example:
//
//   func SomeFunc(input memo.RelExpr, filters memo.FiltersExpr)
//   func SomeFunc(scanPrivate *ScanPrivate)
//
func (t *typeDef) asParam() string {
	// If the type should not be passed by value, then pass it as a pointer.
	if t.passByVal {
		return t.name
	}
	return fmt.Sprintf("*%s", t.name)
}

// newMetadata creates a new instance of metadata from the compiled expression.
// The pkg parameter is used to correctly qualify type names. For example, if
// pkg is "memo", then:
//
//   memo.RelExpr     => RelExpr
//   opt.ScalarExpr   => opt.ScalarExpr
//   memo.ScanPrivate => ScanPrivate
//
func newMetadata(compiled *lang.CompiledExpr, pkg string) *metadata {
	md := &metadata{
		compiled:  compiled,
		types:     make(map[string]*typeDef),
		exprTypes: make(map[lang.Expr]*typeDef),
	}

	// Add all types used in Optgen defines here.
	md.types = map[string]*typeDef{
		"RelExpr":        {fullName: "memo.RelExpr", isExpr: true, isPointer: true},
		"Expr":           {fullName: "opt.Expr", isExpr: true, isPointer: true},
		"ScalarExpr":     {fullName: "opt.ScalarExpr", isExpr: true, isPointer: true},
		"Operator":       {fullName: "opt.Operator", passByVal: true},
		"ColumnID":       {fullName: "opt.ColumnID", passByVal: true},
		"ColSet":         {fullName: "opt.ColSet", passByVal: true},
		"ColList":        {fullName: "opt.ColList", passByVal: true},
		"TableID":        {fullName: "opt.TableID", passByVal: true},
		"SchemaID":       {fullName: "opt.SchemaID", passByVal: true},
		"SequenceID":     {fullName: "opt.SequenceID", passByVal: true},
		"ValuesID":       {fullName: "opt.ValuesID", passByVal: true},
		"Ordering":       {fullName: "opt.Ordering", passByVal: true},
		"OrderingChoice": {fullName: "physical.OrderingChoice", passByVal: true},
		"TupleOrdinal":   {fullName: "memo.TupleOrdinal", passByVal: true},
		"ScanLimit":      {fullName: "memo.ScanLimit", passByVal: true},
		"ScanFlags":      {fullName: "memo.ScanFlags", passByVal: true},
		"JoinFlags":      {fullName: "memo.JoinFlags", passByVal: true},
		"ExplainOptions": {fullName: "tree.ExplainOptions", passByVal: true},
		"StatementType":  {fullName: "tree.StatementType", passByVal: true},
		"ShowTraceType":  {fullName: "tree.ShowTraceType", passByVal: true},
		"bool":           {fullName: "bool", passByVal: true},
		"int":            {fullName: "int", passByVal: true},
		"string":         {fullName: "string", passByVal: true},
		"DatumType":      {fullName: "types.T", isPointer: true},
		"ColType":        {fullName: "coltypes.T", isPointer: true},
		"Datum":          {fullName: "tree.Datum", isPointer: true},
		"TypedExpr":      {fullName: "tree.TypedExpr", isPointer: true},
		"Subquery":       {fullName: "*tree.Subquery", isPointer: true, usePointerIntern: true},
		"CreateTable":    {fullName: "*tree.CreateTable", isPointer: true, usePointerIntern: true},
		"Constraint":     {fullName: "*constraint.Constraint", isPointer: true, usePointerIntern: true},
		"FuncProps":      {fullName: "*tree.FunctionProperties", isPointer: true, usePointerIntern: true},
		"FuncOverload":   {fullName: "*tree.Overload", isPointer: true, usePointerIntern: true},
		"PhysProps":      {fullName: "*physical.Required", isPointer: true},
		"Presentation":   {fullName: "physical.Presentation", passByVal: true},
		"RelProps":       {fullName: "props.Relational"},
		"ScalarProps":    {fullName: "props.Scalar"},
	}

	// Add types of generated op and private structs.
	for _, define := range compiled.Defines {
		// Derive friendly name of type.
		var friendlyName string
		if define.Tags.Contains("ListItem") || define.Tags.Contains("Private") {
			friendlyName = string(define.Name)
		} else {
			friendlyName = fmt.Sprintf("%sExpr", define.Name)
		}

		fullName := fmt.Sprintf("memo.%s", friendlyName)
		typ := &typeDef{fullName: fullName, isGenerated: true}
		if !define.Tags.Contains("Private") {
			typ.isExpr = true
		}
		if define.Tags.Contains("List") {
			typ.passByVal = true
		}

		md.types[friendlyName] = typ
		md.exprTypes[define] = typ
	}

	// 1. Associate each DefineField with its type.
	// 2. Link list types to the types of their list items. A list item type has
	//    the same name as its list parent + the "Item" prefix.
	for _, define := range compiled.Defines {
		// Associate each DefineField with its type.
		for _, field := range define.Fields {
			md.exprTypes[field] = md.lookupType(string(field.Type))
		}

		if define.Tags.Contains("List") {
			listTyp := md.typeOf(define)
			itemTyp := md.lookupType(fmt.Sprintf("%sItem", define.Name))
			if itemTyp != nil {
				listTyp.listItemType = itemTyp
			} else {
				listTyp.listItemType = md.lookupType("ScalarExpr")
			}
		}
	}

	// Now walk each type and fill in any remaining fields.
	for friendlyName, typ := range md.types {
		typ.friendlyName = friendlyName

		// Remove package prefix from types in the same package. For now, don't
		// handle pointer types, since at this time, none are needed in the packages
		// into which files are generated.
		if strings.HasPrefix(typ.fullName, pkg+".") {
			typ.name = typ.fullName[len(pkg)+1:]
		} else {
			typ.name = typ.fullName
		}

		// If type is a pointer/interface, then it should always be passed byref.
		if typ.isPointer {
			typ.passByVal = true
		}
	}

	return md
}

// typeOf returns a type definition for a *lang.Define or *lang.DefineField
// Optgen expression.
func (m *metadata) typeOf(e lang.Expr) *typeDef {
	return m.exprTypes[e]
}

// lookupType returns the type definition with the given friendly name (e.g.
// RelExpr, DatumType, ScanExpr, etc.).
func (m *metadata) lookupType(friendlyName string) *typeDef {
	return m.types[friendlyName]
}

// fieldName maps the Optgen field name to the corresponding Go field name. In
// particular, fields named "_" are mapped to the name of a Go embedded field,
// which is equal to the field's type name:
//
//   define Scan {
//     _ ScanPrivate
//   }
//
// gets compiled into:
//
//   type ScanExpr struct {
//	   ScanPrivate
//     ...
//   }
//
// Note that the field's type name is always a simple alphanumeric identifier
// with no package specified (that's only specified in the fullName field of the
// typeDef).
func (m *metadata) fieldName(field *lang.DefineFieldExpr) string {
	if field.Name == "_" {
		return string(field.Type)
	}
	return string(field.Name)
}

// childFields returns the set of fields for an operator define expression that
// are considered children of that operator. Private (non-expression) fields are
// omitted from the result. For example, for the Project operator:
//
//   define Project {
//     Input       RelExpr
//     Projections ProjectionsExpr
//     Passthrough ColSet
//  }
//
// The Input and Projections fields are children, but the Passthrough field is
// a private field and will not be returned.
func (m *metadata) childFields(define *lang.DefineExpr) lang.DefineFieldsExpr {
	// Skip until non-expression field is found.
	n := 0
	for _, field := range define.Fields {
		typ := m.typeOf(field)
		if !typ.isExpr {
			break
		}
		n++
	}
	return define.Fields[:n]
}

// privateField returns the private field for an operator define expression, if
// one exists. For example, for the Project operator:
//
//   define Project {
//     Input       RelExpr
//     Projections ProjectionsExpr
//     Passthrough ColSet
//  }
//
// The Passthrough field is the private field. If no private field exists for
// the operator, then privateField returns nil.
func (m *metadata) privateField(define *lang.DefineExpr) *lang.DefineFieldExpr {
	// Skip until non-expression field is found.
	n := 0
	for _, field := range define.Fields {
		typ := m.typeOf(field)
		if !typ.isExpr {
			return define.Fields[n]
		}
		n++
	}
	return nil
}

// fieldLoadPrefix returns "&" if the address of the given field should be taken
// when loading that field from an instance in order to pass it elsewhere (like
// to a function). For example:
//
//   f.ConstructUnion(union.Left, union.Right, &union.SetPrivate)
//
// The Left and Right fields are passed by value, but the SetPrivate is passed
// by reference.
func (m *metadata) fieldLoadPrefix(field *lang.DefineFieldExpr) string {
	if !m.typeOf(field).passByVal {
		return "&"
	}
	return ""
}

// fieldStorePrefix is the inverse of fieldLoadPrefix, used when a field value
// is stored to an instance:
//
//   union.Left = left
//   union.Right = right
//   union.SetPrivate = *setPrivate
//
// Since SetPrivate values are passed by reference, they must be dereferenced
// before copying them to a target field.
func (m *metadata) fieldStorePrefix(field *lang.DefineFieldExpr) string {
	if !m.typeOf(field).passByVal {
		return "*"
	}
	return ""
}
