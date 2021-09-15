package schema

// Code generated by go-ipld-prime gengo.  DO NOT EDIT.

import (
	"github.com/ipld/go-ipld-prime/datamodel"
)

var _ datamodel.Node = nil // suppress errors when this dependency is not referenced
// Type is a struct embeding a NodePrototype/Type for every Node implementation in this package.
// One of its major uses is to start the construction of a value.
// You can use it like this:
//
// 		schema.Type.YourTypeName.NewBuilder().BeginMap() //...
//
// and:
//
// 		schema.Type.OtherTypeName.NewBuilder().AssignString("x") // ...
//
var Type typeSlab

type typeSlab struct {
	Advertisement            _Advertisement__Prototype
	Advertisement__Repr      _Advertisement__ReprPrototype
	Any                      _Any__Prototype
	Any__Repr                _Any__ReprPrototype
	Bool                     _Bool__Prototype
	Bool__Repr               _Bool__ReprPrototype
	Bytes                    _Bytes__Prototype
	Bytes__Repr              _Bytes__ReprPrototype
	EntryChunk               _EntryChunk__Prototype
	EntryChunk__Repr         _EntryChunk__ReprPrototype
	Float                    _Float__Prototype
	Float__Repr              _Float__ReprPrototype
	Int                      _Int__Prototype
	Int__Repr                _Int__ReprPrototype
	Link                     _Link__Prototype
	Link__Repr               _Link__ReprPrototype
	Link_Advertisement       _Link_Advertisement__Prototype
	Link_Advertisement__Repr _Link_Advertisement__ReprPrototype
	Link_EntryChunk          _Link_EntryChunk__Prototype
	Link_EntryChunk__Repr    _Link_EntryChunk__ReprPrototype
	List                     _List__Prototype
	List__Repr               _List__ReprPrototype
	List_Bytes               _List_Bytes__Prototype
	List_Bytes__Repr         _List_Bytes__ReprPrototype
	List_String              _List_String__Prototype
	List_String__Repr        _List_String__ReprPrototype
	Map                      _Map__Prototype
	Map__Repr                _Map__ReprPrototype
	String                   _String__Prototype
	String__Repr             _String__ReprPrototype
}

// --- type definitions follow ---

// Advertisement matches the IPLD Schema type "Advertisement".  It has Struct type-kind, and may be interrogated like map kind.
type Advertisement = *_Advertisement
type _Advertisement struct {
	PreviousID _Link_Advertisement__Maybe
	Provider   _String
	Signature  _Bytes
	Entries    _Link
	Metadata   _Bytes
	IsRm       _Bool
}

// Any matches the IPLD Schema type "Any".
// Any has Union typekind, which means its data model behaviors are that of a map kind.
type Any = *_Any
type _Any struct {
	x _Any__iface
}
type _Any__iface interface {
	_Any__member()
}

func (_Bool) _Any__member()   {}
func (_Int) _Any__member()    {}
func (_Float) _Any__member()  {}
func (_String) _Any__member() {}
func (_Bytes) _Any__member()  {}
func (_Map) _Any__member()    {}
func (_List) _Any__member()   {}
func (_Link) _Any__member()   {}

// Bool matches the IPLD Schema type "Bool".  It has bool kind.
type Bool = *_Bool
type _Bool struct{ x bool }

// Bytes matches the IPLD Schema type "Bytes".  It has bytes kind.
type Bytes = *_Bytes
type _Bytes struct{ x []byte }

// EntryChunk matches the IPLD Schema type "EntryChunk".  It has Struct type-kind, and may be interrogated like map kind.
type EntryChunk = *_EntryChunk
type _EntryChunk struct {
	Entries _List_Bytes
	Next    _Link_EntryChunk__Maybe
}

// Float matches the IPLD Schema type "Float".  It has float kind.
type Float = *_Float
type _Float struct{ x float64 }

// Int matches the IPLD Schema type "Int".  It has int kind.
type Int = *_Int
type _Int struct{ x int64 }

// Link matches the IPLD Schema type "Link".  It has link kind.
type Link = *_Link
type _Link struct{ x datamodel.Link }

// Link_Advertisement matches the IPLD Schema type "Link_Advertisement".  It has link kind.
type Link_Advertisement = *_Link_Advertisement
type _Link_Advertisement struct{ x datamodel.Link }

// Link_EntryChunk matches the IPLD Schema type "Link_EntryChunk".  It has link kind.
type Link_EntryChunk = *_Link_EntryChunk
type _Link_EntryChunk struct{ x datamodel.Link }

// List matches the IPLD Schema type "List".  It has list kind.
type List = *_List
type _List struct {
	x []_Any__Maybe
}

// List_Bytes matches the IPLD Schema type "List_Bytes".  It has list kind.
type List_Bytes = *_List_Bytes
type _List_Bytes struct {
	x []_Bytes
}

// List_String matches the IPLD Schema type "List_String".  It has list kind.
type List_String = *_List_String
type _List_String struct {
	x []_String
}

// Map matches the IPLD Schema type "Map".  It has map kind.
type Map = *_Map
type _Map struct {
	m map[_String]MaybeAny
	t []_Map__entry
}
type _Map__entry struct {
	k _String
	v _Any__Maybe
}

// String matches the IPLD Schema type "String".  It has string kind.
type String = *_String
type _String struct{ x string }
