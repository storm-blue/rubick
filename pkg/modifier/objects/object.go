package objects

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/storm-blue/rubick/pkg/common"
	"gopkg.in/yaml.v2"
	"regexp"
	"strconv"
	"strings"
)

const metadataKey = ".__object.__metadata"

func NewObject() StructuredObject {
	return _object{metadataKey: _metadata{}}
}

func FromJSON(jsonStr string) (StructuredObject, error) {
	result := map[interface{}]interface{}{}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, err
	}
	return FromMap(result), nil
}

func FromYAML(yamlStr string) (StructuredObject, error) {
	result := map[interface{}]interface{}{}
	if err := yaml.Unmarshal([]byte(yamlStr), result); err != nil {
		return nil, err
	}
	return FromMap(result), nil
}

func FromYAMLs(multiYaml string) ([]StructuredObject, error) {
	var result []StructuredObject
	decoder := yaml.NewDecoder(bytes.NewBuffer([]byte(multiYaml)))

	for {
		o := map[interface{}]interface{}{}
		if err := decoder.Decode(o); err != nil {
			if err.Error() == "EOF" {
				break
			}
		}
		result = append(result, FromMap(o))
	}

	return result, nil
}

func FromMap(m map[interface{}]interface{}) StructuredObject {
	m[metadataKey] = _metadata{}
	return _object(m)
}

func ToJSON(object StructuredObject) (string, error) {
	return object.ToJSON()
}

func ToYAML(object StructuredObject) (string, error) {
	return object.ToYAML()
}

func ToYAMLs(_objects []StructuredObject) (string, error) {
	buf := &bytes.Buffer{}
	encoder := yaml.NewEncoder(buf)
	for _, object := range _objects {
		o := object.(_object)
		metadata := object.Metadata()

		// remove metadata before marshal
		delete(o, metadataKey)

		err := encoder.Encode(o)

		// restore metadata after marshal
		o[metadataKey] = metadata

		if err != nil {
			return "", err
		}
	}
	return buf.String(), nil
}

type StructuredObject interface {
	Get(string) (interface{}, error)
	GetObject(string) (StructuredObject, error)
	GetArray(string) ([]interface{}, error)
	GetObjects(string) ([]StructuredObject, error)
	GetInt(string) (int, error)
	GetInt32(string) (int32, error)
	GetInt64(string) (int64, error)
	GetFloat32(string) (float32, error)
	GetFloat64(string) (float64, error)
	GetString(string) (string, error)
	Metadata() Metadata

	Delete(string) error
	Set(key string, value interface{}) error

	Exist(string) bool
	Len() int
	ToMap() map[interface{}]interface{}
	ToYAML() (string, error)
	ToJSON() (string, error)

	Clone() (StructuredObject, error)
}

type _object map[interface{}]interface{}

func (o _object) Metadata() Metadata {
	metadata, ok := o[metadataKey]
	if !ok {
		metadata = _metadata{}
		o[metadataKey] = metadata
	}
	return metadata.(Metadata)
}

func (o _object) Len() int {
	return len(o)
}

func (o _object) Clone() (StructuredObject, error) {
	yamlStr, err := o.ToYAML()
	if err != nil {
		return nil, err
	}
	return FromYAML(yamlStr)
}

func (o _object) ToMap() map[interface{}]interface{} {
	return o
}

func (o _object) GetObject(s string) (StructuredObject, error) {
	object_, err := o.Get(s)
	if err != nil {
		return nil, err
	}

	if object_ == nil {
		return nil, nil
	}

	switch object := object_.(type) {
	case _object:
		return object, nil
	case map[interface{}]interface{}:
		return _object(object), nil
	default:
		return nil, fmt.Errorf("GetObject error: value is not object, key: %v", s)
	}
}

func (o _object) GetObjects(s string) ([]StructuredObject, error) {
	array, err := o.GetArray(s)
	if err != nil {
		return nil, err
	}

	var result []StructuredObject
	for _, object_ := range array {
		if object, ok := object_.(map[interface{}]interface{}); ok {
			result = append(result, _object(object))
		} else {
			return nil, fmt.Errorf("GetObjects error: value is not objects, key: %v", s)
		}
	}
	return result, nil
}

func (o _object) GetArray(s string) ([]interface{}, error) {
	object_, err := o.Get(s)
	if err != nil {
		return nil, err
	}

	if object_ == nil {
		return nil, nil
	}

	object, ok := object_.([]interface{})
	if ok {
		return object, nil
	} else {
		return nil, fmt.Errorf("GetArray error: value is not array, key: %v", s)
	}
}

func (o _object) GetInt(s string) (int, error) {
	object_, err := o.Get(s)
	if err != nil {
		return 0, err
	}

	if object_ == nil {
		return 0, fmt.Errorf("GetInt error: value is nil, key: %v", s)
	}

	switch object := object_.(type) {
	case int:
		return object, nil
	case float64:
		return int(object), nil
	case int32:
		return int(object), nil
	case int64:
		return int(object), nil
	case float32:
		return int(object), nil
	default:
		return 0, fmt.Errorf("GetInt error: value is not int, key: %v", s)
	}
}

func (o _object) GetInt32(s string) (int32, error) {
	object_, err := o.Get(s)
	if err != nil {
		return 0, err
	}

	if object_ == nil {
		return 0, fmt.Errorf("GetInt32 error: value is nil, key: %v", s)
	}

	switch object := object_.(type) {
	case int:
		return int32(object), nil
	case float64:
		return int32(object), nil
	case int32:
		return object, nil
	case int64:
		return int32(object), nil
	case float32:
		return int32(object), nil
	default:
		return 0, fmt.Errorf("GetInt32 error: value is not int32, key: %v", s)
	}
}

func (o _object) GetInt64(s string) (int64, error) {
	object_, err := o.Get(s)
	if err != nil {
		return 0, err
	}

	if object_ == nil {
		return 0, fmt.Errorf("GetInt64 error: value is nil, key: %v", s)
	}

	switch object := object_.(type) {
	case int:
		return int64(object), nil
	case float64:
		return int64(object), nil
	case int32:
		return int64(object), nil
	case int64:
		return object, nil
	case float32:
		return int64(object), nil
	default:
		return 0, fmt.Errorf("GetInt64 error: value is not int64, key: %v", s)
	}
}

func (o _object) GetFloat32(s string) (float32, error) {
	object_, err := o.Get(s)
	if err != nil {
		return 0, err
	}

	if object_ == nil {
		return 0, fmt.Errorf("GetFloat32 error: value is nil, key: %v", s)
	}

	switch object := object_.(type) {
	case float64:
		return float32(object), nil
	case int:
		return float32(object), nil
	case int32:
		return float32(object), nil
	case int64:
		return float32(object), nil
	case float32:
		return object, nil
	default:
		return 0, fmt.Errorf("GetFloat32 error: value is not float32, key: %v", s)
	}
}

func (o _object) GetFloat64(s string) (float64, error) {
	object_, err := o.Get(s)
	if err != nil {
		return 0, err
	}

	if object_ == nil {
		return 0, fmt.Errorf("GetFloat64 error: value is nil, key: %v", s)
	}

	switch object := object_.(type) {
	case float64:
		return object, nil
	case int:
		return float64(object), nil
	case int32:
		return float64(object), nil
	case int64:
		return float64(object), nil
	case float32:
		return float64(object), nil
	default:
		return 0, fmt.Errorf("GetFloat64 error: value is not float64, key: %v", s)
	}
}

func (o _object) GetString(s string) (string, error) {
	object_, err := o.Get(s)
	if err != nil {
		return "", err
	}

	if object_ == nil {
		return "", fmt.Errorf("GetString error: value is nil, key: %v", s)
	}

	switch object := object_.(type) {
	case string:
		return object, nil
	default:
		return "", fmt.Errorf("GetString error: value is not string, key: %v", s)
	}
}

func (o _object) ToYAML() (string, error) {
	metadata := o.Metadata()

	// remove metadata before marshal
	delete(o, metadataKey)

	yamlBytes, err := yaml.Marshal(o)

	// restore metadata after marshal
	o[metadataKey] = metadata

	if err != nil {
		return "", err
	}
	return string(yamlBytes), nil
}

func (o _object) ToJSON() (string, error) {
	metadata := o.Metadata()

	// remove metadata before marshal
	delete(o, metadataKey)
	bs, err := json.Marshal(o)

	// restore metadata after marshal
	o[metadataKey] = metadata

	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (o _object) Delete(key string) error {
	return deleteObject(o, key)
}

func (o _object) Exist(key string) bool {
	return existObject(o, key)
}

func (o _object) Get(key string) (interface{}, error) {
	return getObject(o, key)
}

func (o _object) Set(key string, value interface{}) error {
	return setObject(o, key, value)
}

func getObject(objects map[interface{}]interface{}, fullKey string) (interface{}, error) {
	key, restKey, err := ParseNextSegment(fullKey)
	if err != nil {
		return nil, err
	}

	if restKey == "" {
		if key.isArray {
			arrayObject, ok := objects[key.key]
			if !ok {
				return nil, nil
			}
			array, ok := arrayObject.([]interface{})
			if !ok {
				return nil, fmt.Errorf("ojbect is not array: %v", arrayObject)
			}

			return getElement(array, key.index)
		} else {
			return objects[key.key], nil
		}
	} else { // has rest key
		if key.isArray {
			arrayObject, ok := objects[key.key]
			if !ok {
				return nil, nil
			}
			array, ok := arrayObject.([]interface{})
			if !ok {
				return nil, fmt.Errorf("ojbect is not array: %v", arrayObject)
			}

			_subObject, err := getElement(array, key.index)
			if err != nil {
				return nil, err
			}

			subObject, ok := _subObject.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("ojbect is not map: %v", _subObject)
			}

			return getObject(subObject, restKey)
		} else {
			if subObject, ok := objects[key.key]; ok {
				subObjects, ok := subObject.(map[interface{}]interface{})
				if !ok {
					return nil, fmt.Errorf("ojbect is not map: %v", subObject)
				}
				return getObject(subObjects, restKey)
			}
		}
	}

	return nil, nil
}

func existObject(objects map[interface{}]interface{}, fullKey string) bool {
	key, restKey, err := ParseNextSegment(fullKey)
	if err != nil {
		return false
	}

	if restKey == "" {
		if key.isArray {
			arrayObject, ok := objects[key.key]
			if !ok {
				return false
			}
			array, ok := arrayObject.([]interface{})
			if !ok {
				return false
			}
			return existElement(array, key.index)
		} else {
			_, ok := objects[key.key]
			return ok
		}
	} else { // has rest key
		if key.isArray {
			arrayObject, ok := objects[key.key]
			if !ok {
				return false
			}
			array, ok := arrayObject.([]interface{})
			if !ok {
				return false
			}

			_subObject, err := getElementForExist(array, key.index)
			if err != nil {
				return false
			}

			subObject, ok := _subObject.(map[interface{}]interface{})
			if !ok {
				return false
			}

			return existObject(subObject, restKey)
		} else {
			if subObject, ok := objects[key.key]; ok {
				subObjects, ok := subObject.(map[interface{}]interface{})
				if !ok {
					return false
				}
				return existObject(subObjects, restKey)
			} else {
				return false
			}
		}
	}
}

func getElementForExist(slice []interface{}, index YamlIndex) (interface{}, error) {
	switch index.indexType {
	case IndexNormal:
		if index.index < 0 || index.index > len(slice)-1 {
			return nil, fmt.Errorf("get elements for exsit by index error: out of range: %v", index.index)
		}
		return slice[index.index], nil
	case IndexMax:
		if len(slice) == 0 {
			return nil, nil
		}
		return slice[len(slice)-1], nil
	case IndexAppend:
		fallthrough
	case IndexLoop:
		return nil, fmt.Errorf("get elements for exsit by index error: unsupported index type: %v", index.indexType)
	default:
		return nil, fmt.Errorf("get elements by index error: unkown index type: %v", index.indexType)
	}
}

func deleteObject(objects _object, fullKey string) error {
	key, restKey, err := ParseNextSegment(fullKey)
	if err != nil {
		return err
	}

	if restKey == "" {
		if key.isArray {
			arrayObject, ok := objects[key.key]
			if !ok {
				return fmt.Errorf("can not find key: %v", fullKey)
			}
			array, ok := arrayObject.([]interface{})
			if !ok {
				return fmt.Errorf("ojbect is not array: %v", arrayObject)
			}

			objects[key.key], err = deleteElement(array, key.index)
			if err != nil {
				return err
			}
		} else {
			delete(objects, key.key)
		}
	} else { // has rest key
		if key.isArray {
			arrayObject, ok := objects[key.key]
			if !ok {
				return fmt.Errorf("can not find key: %v", fullKey)
			}
			array, ok := arrayObject.([]interface{})
			if !ok {
				return fmt.Errorf("ojbect is not array: %v", arrayObject)
			}

			_subObjects, err := getElementForDelete(array, key.index)
			if err != nil {
				return err
			}

			for _, _subObject := range _subObjects {
				if _subObject != nil {
					subObject, ok := _subObject.(map[interface{}]interface{})
					if !ok {
						return fmt.Errorf("ojbect is not map: %v", _subObjects)
					}

					if err := deleteObject(subObject, restKey); err != nil {
						return err
					}
				}
			}
		} else {
			if subObject, ok := objects[key.key]; ok {
				subObjects, ok := subObject.(map[interface{}]interface{})
				if !ok {
					return fmt.Errorf("ojbect is not map: %v", subObject)
				}
				if err := deleteObject(subObjects, restKey); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func setObject(objects map[interface{}]interface{}, fullKey string, value interface{}) error {
	if _, ok := value.(_object); ok {
		value = value.(_object).ToMap()
	}

	key, restKey, err := ParseNextSegment(fullKey)
	if err != nil {
		return err
	}

	if restKey == "" {
		if key.isArray {
			var array []interface{}
			arrayObject, ok := objects[key.key]
			if ok {
				array, ok = arrayObject.([]interface{})
				if !ok {
					return fmt.Errorf("ojbect is not array: %v", arrayObject)
				}
			} else {
				array = []interface{}{}
				objects[key.key] = array
			}

			array, err := setElements(array, key.index, value)
			if err != nil {
				return err
			}
			objects[key.key] = array
		} else {
			objects[key.key] = value
		}
	} else { // has rest key
		if key.isArray {
			var array []interface{}
			arrayObject, ok := objects[key.key]
			if ok {
				array, ok = arrayObject.([]interface{})
				if !ok {
					return fmt.Errorf("ojbect is not array: %v", arrayObject)
				}
			} else {
				array = []interface{}{}
				objects[key.key] = array
			}

			subObjects_, array, err := getElementForSet(array, key.index)
			if err != nil {
				return err
			}

			// getElementForSet maybe change the slice, set again
			objects[key.key] = array

			for _, subObject_ := range subObjects_ {
				if subObject_ != nil {
					subObject, ok := subObject_.(map[interface{}]interface{})
					if !ok {
						return fmt.Errorf("ojbect is not map: %v", subObject)
					}
					if err := setObject(subObject, restKey, value); err != nil {
						return err
					}
				}
			}
		} else {
			var subObjects map[interface{}]interface{}
			subObject, ok := objects[key.key]

			if ok {
				subObjects, ok = subObject.(map[interface{}]interface{})
				if !ok {
					return fmt.Errorf("ojbect is not map: %v", subObject)
				}
			} else {
				subObjects = map[interface{}]interface{}{}
				objects[key.key] = subObjects
			}

			if err := setObject(subObjects, restKey, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func getElement(slice []interface{}, index YamlIndex) (interface{}, error) {
	switch index.indexType {
	case IndexNormal:
		if index.index < 0 || index.index > len(slice)-1 {
			return nil, fmt.Errorf("get elements by index error: out of range: %v", index.index)
		}
		return slice[index.index], nil
	case IndexMax:
		if len(slice) == 0 {
			return nil, nil
		}
		return slice[len(slice)-1], nil
	case IndexAppend:
		return nil, fmt.Errorf("get elements by index error: unsupported index type: %v", index.indexType)
	case IndexLoop:
		return slice, nil
	default:
		return nil, fmt.Errorf("get elements by index error: unkown index type: %v", index.indexType)
	}
}

func existElement(slice []interface{}, index YamlIndex) bool {
	switch index.indexType {
	case IndexNormal:
		if index.index < 0 || index.index > len(slice)-1 {
			return false
		}
		return true
	case IndexMax:
		if len(slice) == 0 {
			return false
		}
		return true
	case IndexAppend:
		return false
	case IndexLoop:
		return true
	default:
		return false
	}
}

func getElementForDelete(slice []interface{}, index YamlIndex) ([]interface{}, error) {
	switch index.indexType {
	case IndexNormal:
		if index.index < 0 || index.index > len(slice)-1 {
			return nil, fmt.Errorf("get elements by index for delete error: out of range: %v", index.index)
		}
		return []interface{}{slice[index.index]}, nil
	case IndexMax:
		if len(slice) == 0 {
			return nil, nil
		}
		return []interface{}{slice[len(slice)-1]}, nil
	case IndexAppend:
		return nil, fmt.Errorf("get elements by index for delete error: unsupported index type: %v", index.indexType)
	case IndexSearch:
		for _, e := range slice {
			e_, ok := e.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("delete elements by index error: %v, index: %v", slice, index)
			}
			if lenientEqual(e_[index.key], index.value) {
				return []interface{}{e}, nil
			}
		}
		return nil, nil
	case IndexLoop:
		return slice, nil
	default:
		return nil, fmt.Errorf("get elements by index for delete error: unkown index type: %v", index.indexType)
	}
}

func deleteElement(slice []interface{}, index YamlIndex) ([]interface{}, error) {
	var result []interface{}

	// delete all elements
	if index.indexType == IndexLoop {
		return result, nil
	}

	for i, e := range slice {
		switch index.indexType {
		case IndexNormal:
			if i != index.index {
				result = append(result, e)
			}
		case IndexMax: // delete last elements
			if i < len(slice) {
				result = append(result, e)
			}
		case IndexSearch:
			e_, ok := e.(map[interface{}]interface{})
			if !ok {
				return nil, fmt.Errorf("delete elements by index error: %v, index: %v", slice, index)
			}
			if !lenientEqual(e_[index.key], index.value) {
				result = append(result, e)
			}
		case IndexAppend: // append(++) operation is not allowed in deletion
			return nil, fmt.Errorf("delete elements by index error: unsupported index type: %v", index.indexType)
		default:
			return nil, fmt.Errorf("delete elements by index error: unkown index type: %v", index.indexType)
		}
	}
	return result, nil
}

func getElementForSet(slice []interface{}, index YamlIndex) ([]interface{}, []interface{}, error) {
	switch index.indexType {
	case IndexNormal:
		if index.index < 0 || index.index > len(slice) {
			return nil, slice, fmt.Errorf("get elements by index for set error: out of range: %v, index: %v", slice, index.index)
		} else if index.index == len(slice) {
			e := map[interface{}]interface{}{}
			return []interface{}{e}, append(slice, e), nil
		} else {
			return []interface{}{slice[index.index]}, slice, nil
		}
	case IndexMax:
		if len(slice) == 0 {
			return nil, slice, nil
		}
		return []interface{}{slice[len(slice)-1]}, slice, nil
	case IndexAppend:
		e := map[interface{}]interface{}{}
		return []interface{}{e}, append(slice, e), nil
	case IndexLoop:
		return slice, slice, nil
	case IndexSearch:
		for _, e := range slice {
			e_, ok := e.(map[interface{}]interface{})
			if !ok {
				return nil, slice, fmt.Errorf("get elements by index for set error: %v, index: %v", slice, index)
			}
			if lenientEqual(e_[index.key], index.value) {
				return []interface{}{e}, slice, nil
			}
		}
		return nil, slice, nil
	default:
		return nil, slice, fmt.Errorf("get elements by index for set error: unkown index type: %v", index.indexType)
	}
}

func setElements(slice []interface{}, index YamlIndex, value interface{}) ([]interface{}, error) {
	switch index.indexType {
	case IndexNormal:
		if index.index > len(slice) || index.index < 0 {
			return slice, fmt.Errorf("set elements by index error: out of range: %v, index: %v", slice, index.index)
		} else if index.index == len(slice) {
			return append(slice, value), nil
		} else {
			slice[index.index] = value
			return slice, nil
		}
	case IndexMax:
		if len(slice) == 0 {
			return slice, nil
		}
		slice[len(slice)-1] = value
		return slice, nil
	case IndexAppend:
		return append(slice, value), nil
	case IndexSearch:
		var results []interface{}
		for _, e := range slice {
			e_, ok := e.(map[interface{}]interface{})
			if !ok {
				return slice, fmt.Errorf("set elements by index error: element is not map: %v", e)
			}
			if lenientEqual(e_[index.key], index.value) {
				results = append(results, value)
			} else {
				results = append(results, e)
			}
		}
		return results, nil
	default:
		return slice, fmt.Errorf("set elements by index error: unkown index type: %v", index)
	}
}

// ParseNextSegment from key
// Not support array like: a[1][2]
func ParseNextSegment(key string) (KeySegment, string, error) {

	// need parse '()'
	if strings.HasPrefix(key, "(") {
		_, rightIndex := common.FindFirstParenthesesPair(key)
		if rightIndex == -1 {
			return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
		}

		head := key[:rightIndex+1]
		head = common.UnwrapIfNeeded(head)
		if !isValidateKeySegment(head) {
			return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
		}

		tail := key[rightIndex+1:]
		if tail == "" {
			return KeySegment{key: head}, "", nil
		}

		if strings.HasPrefix(tail, ".") {
			restKey := strings.TrimPrefix(tail, ".")
			if restKey == "" {
				return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
			}
			return KeySegment{key: head}, restKey, nil
		} else if strings.HasPrefix(tail, "[") {
			_, rightIndexOfTail := common.FindFirstBracketsPair(tail)
			if rightIndexOfTail == -1 {
				return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
			}
			headOfTail := tail[:rightIndexOfTail+1]
			headOfTail = common.UnwrapBracketsIfNeeded(headOfTail)
			tailOfTail := tail[rightIndexOfTail+1:]

			if headOfTail == "" {
				return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
			}
			index, err := parseIndex(headOfTail)
			if err != nil {
				return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
			}

			if tailOfTail == "" {
				return KeySegment{key: head, isArray: true, index: index}, "", nil
			}

			if !strings.HasPrefix(tailOfTail, ".") || tailOfTail == "." {
				return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
			}
			return KeySegment{key: head, isArray: true, index: index}, strings.TrimPrefix(tailOfTail, "."), nil
		} else {
			return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
		}
	} else
	// do not need to parse '()'
	{
		firstLeftBracketIndex := strings.Index(key, "[")
		firstDotIndex := strings.Index(key, ".")
		// '[]' has higher priority than '.' , parse '[]'
		if firstLeftBracketIndex != -1 && firstLeftBracketIndex < firstDotIndex {
			leftBracketsIndex, rightBracketsIndex := common.FindFirstBracketsPair(key)
			if leftBracketsIndex == -1 || rightBracketsIndex == -1 {
				return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
			}

			head := key[:rightBracketsIndex+1]
			tail := key[rightBracketsIndex+1:]

			keyPartOfHead := head[:leftBracketsIndex]
			if !isValidateKeySegment(keyPartOfHead) {
				return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
			}

			indexPartOfHead := head[leftBracketsIndex:]
			indexPartOfHead = common.UnwrapBracketsIfNeeded(indexPartOfHead)
			index, err := parseIndex(indexPartOfHead)
			if err != nil {
				return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
			}

			if tail == "" {
				return KeySegment{isArray: true, key: keyPartOfHead, index: index}, "", nil
			}

			if !strings.HasPrefix(tail, ".") || tail == "." {
				return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
			}

			tail = strings.TrimPrefix(tail, ".")
			return KeySegment{isArray: true, key: keyPartOfHead, index: index}, tail, nil
		} else
		// parse by separator: '.'
		{
			headAndTailParts := strings.SplitN(key, ".", 2)
			if len(headAndTailParts) == 2 && headAndTailParts[1] == "" {
				return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
			}

			head := headAndTailParts[0]
			leftBracketsIndex, rightBracketsIndex := common.FindFirstBracketsPair(head)
			if leftBracketsIndex != -1 && rightBracketsIndex != -1 {
				if rightBracketsIndex != len(head)-1 {
					return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
				}
				keyPartOfHead := head[:leftBracketsIndex]
				if !isValidateKeySegment(keyPartOfHead) {
					return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
				}

				indexPartOfHead := head[leftBracketsIndex:]
				indexPartOfHead = common.UnwrapBracketsIfNeeded(indexPartOfHead)

				index, err := parseIndex(indexPartOfHead)
				if err != nil {
					return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
				}
				if len(headAndTailParts) == 1 {
					return KeySegment{key: keyPartOfHead, isArray: true, index: index}, "", nil
				}
				return KeySegment{key: keyPartOfHead, isArray: true, index: index}, headAndTailParts[1], nil
			} else {
				if !isValidateKeySegment(head) {
					return KeySegment{}, "", fmt.Errorf("parse next key error: %v", key)
				}
				if len(headAndTailParts) == 1 {
					return KeySegment{key: head}, "", nil
				}
				return KeySegment{key: head}, headAndTailParts[1], nil
			}
		}
	}
}

var keySegmentRegex = regexp.MustCompile("^[a-zA-Z0-9/_.-]+$")

func isValidateKeySegment(key string) bool {
	return keySegmentRegex.MatchString(key)
}

func IsValidKey(key string) bool {
	if key == "" {
		return false
	}
	_, rest, err := ParseNextSegment(key)
	if err != nil {
		return false
	}
	if rest == "" {
		return true
	} else {
		return IsValidKey(rest)
	}
}

func parseIndex(index string) (YamlIndex, error) {
	if index == "+" {
		return YamlIndex{
			indexType: IndexMax,
		}, nil
	}

	if index == "++" {
		return YamlIndex{
			indexType: IndexAppend,
		}, nil
	}

	if index == "*" {
		return YamlIndex{
			indexType: IndexLoop,
		}, nil
	}

	if strings.Contains(index, "=") {
		strs := strings.SplitN(index, "=", 2)
		return YamlIndex{
			indexType: IndexSearch,
			key:       strs[0],
			value:     strs[1],
		}, nil
	}

	numberIndex, err := strconv.Atoi(index)
	if err != nil {
		return YamlIndex{}, err
	}

	return YamlIndex{
		indexType: IndexNormal,
		index:     numberIndex,
	}, nil
}

type KeySegment struct {
	isArray bool
	key     string
	index   YamlIndex
}

type IndexType string

const (
	IndexNormal = "index-normal"
	IndexMax    = "index-max"
	IndexAppend = "index-append"
	IndexSearch = "index-search"
	IndexLoop   = "index-loop"
)

type YamlIndex struct {
	indexType IndexType
	index     int
	key       string
	value     string
}

func lenientEqual(v interface{}, s string) bool {
	switch v.(type) {
	case string:
		return s == v.(string)
	case int:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return false
		}
		return f == float64(v.(int))
	case int64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return false
		}
		return f == float64(v.(int64))
	case float32:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return false
		}
		return f == float64(v.(float32))
	case float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return false
		}
		return f == v.(float64)
	case bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return false
		}
		return b == v.(bool)
	default:
		return fmt.Sprint(v) == s
	}
}
