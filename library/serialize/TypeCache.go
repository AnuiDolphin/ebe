package serialize

import (
	"ebe/types"
	"fmt"
	"reflect"
	"sync"
)

// typeCache provides concurrent-safe caching of reflection metadata
var typeCache = &TypeCache{
	typeToEBEType: &sync.Map{},
	structFields:  &sync.Map{},
	elementTypes:  &sync.Map{},
}

// TypeCache holds cached reflection metadata to avoid repeated expensive operations
type TypeCache struct {
	// typeToEBEType caches reflect.Type -> types.Types mappings
	typeToEBEType *sync.Map
	
	// structFields caches struct field metadata
	structFields *sync.Map
	
	// elementTypes caches array/slice element types
	elementTypes *sync.Map
}

// StructFieldInfo contains cached information about a struct field
type StructFieldInfo struct {
	Name     string
	Type     reflect.Type
	Index    int
	Exported bool
}

// StructInfo contains cached information about a struct type
type StructInfo struct {
	Fields []StructFieldInfo
	Empty  bool // True if struct has no exported fields
}

// GetEBEType returns the cached EBE type for a reflect.Type, computing and caching if not found
func (tc *TypeCache) GetEBEType(t reflect.Type) (types.Types, error) {
	// Try to get from cache first
	if cached, found := tc.typeToEBEType.Load(t); found {
		return cached.(types.Types), nil
	}
	
	// Compute the EBE type
	ebeType, err := computeEBEType(t)
	if err != nil {
		return 0, err
	}
	
	// Cache the result
	tc.typeToEBEType.Store(t, ebeType)
	return ebeType, nil
}

// GetStructInfo returns cached struct field information, computing and caching if not found
func (tc *TypeCache) GetStructInfo(t reflect.Type) (*StructInfo, error) {
	// Try to get from cache first
	if cached, found := tc.structFields.Load(t); found {
		return cached.(*StructInfo), nil
	}
	
	// Compute struct info
	structInfo, err := computeStructInfo(t)
	if err != nil {
		return nil, err
	}
	
	// Cache the result
	tc.structFields.Store(t, structInfo)
	return structInfo, nil
}

// GetElementType returns cached array/slice element type, computing and caching if not found
func (tc *TypeCache) GetElementType(t reflect.Type) (types.Types, error) {
	// Try to get from cache first
	if cached, found := tc.elementTypes.Load(t); found {
		return cached.(types.Types), nil
	}
	
	// Compute element type
	elemType, err := computeElementType(t)
	if err != nil {
		return 0, err
	}
	
	// Cache the result
	tc.elementTypes.Store(t, elemType)
	return elemType, nil
}

// computeEBEType computes the EBE type for a reflect.Type (internal function)
func computeEBEType(t reflect.Type) (types.Types, error) {
	switch t.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return types.UInt, nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return types.SInt, nil
	case reflect.Float32, reflect.Float64:
		return types.Float, nil
	case reflect.Bool:
		return types.Boolean, nil
	case reflect.String:
		return types.String, nil
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return types.Buffer, nil
		}
		return types.Array, nil
	case reflect.Array:
		return types.Array, nil
	case reflect.Map:
		return types.Map, nil
	case reflect.Struct:
		return types.Struct, nil
	case reflect.Chan:
		return 0, fmt.Errorf("channels not supported")
	default:
		return 0, fmt.Errorf("unsupported type: %v", t)
	}
}

// computeStructInfo computes struct field information (internal function)
func computeStructInfo(t reflect.Type) (*StructInfo, error) {
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct type, got %v", t.Kind())
	}
	
	info := &StructInfo{
		Fields: make([]StructFieldInfo, 0, t.NumField()),
		Empty:  true,
	}
	
	// Analyze each field
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		exported := field.PkgPath == "" // Exported field
		
		if exported {
			info.Empty = false
		}
		
		fieldInfo := StructFieldInfo{
			Name:     field.Name,
			Type:     field.Type,
			Index:    i,
			Exported: exported,
		}
		
		info.Fields = append(info.Fields, fieldInfo)
	}
	
	return info, nil
}

// computeElementType computes the EBE type for array/slice elements (internal function)
func computeElementType(t reflect.Type) (types.Types, error) {
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		return 0, fmt.Errorf("expected slice or array type, got %v", t.Kind())
	}
	
	elemType := t.Elem()
	
	// Apply original logic restrictions for array elements
	switch elemType.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return types.UInt, nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return types.SInt, nil
	case reflect.Float32, reflect.Float64:
		return types.Float, nil
	case reflect.Bool:
		return types.Boolean, nil
	case reflect.String:
		return types.String, nil
	case reflect.Slice:
		if elemType.Elem().Kind() == reflect.Uint8 {
			return types.Buffer, nil
		}
		return 0, fmt.Errorf("nested slices not yet supported")
	case reflect.Struct:
		return types.Struct, nil
	case reflect.Map:
		return types.Map, nil
	case reflect.Chan:
		return 0, fmt.Errorf("channels not supported")
	default:
		return 0, fmt.Errorf("unsupported type: %v", elemType)
	}
}

// ClearCache clears all cached data (useful for testing)
func (tc *TypeCache) ClearCache() {
	tc.typeToEBEType = &sync.Map{}
	tc.structFields = &sync.Map{}
	tc.elementTypes = &sync.Map{}
}