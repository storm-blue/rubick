package conditions

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"reflect"
	"strconv"
)

const (
	None = iota
	GreaterThan
	GreaterThanOrEqual
	LesserThan
	LesserThanOrEqual
	EqualTo
	NotEqual
	And
	Or
	Not
)

func calculate(operator int, objectValue interface{}, conditionValue interface{}) (bool, error) {
	switch v_ := objectValue.(type) {
	case int8:
		return calculateNumber(operator, float64(v_), conditionValue)
	case int32:
		return calculateNumber(operator, float64(v_), conditionValue)
	case int64:
		return calculateNumber(operator, float64(v_), conditionValue)
	case int:
		return calculateNumber(operator, float64(v_), conditionValue)
	case uint:
		return calculateNumber(operator, float64(v_), conditionValue)
	case float32:
		return calculateNumber(operator, float64(v_), conditionValue)
	case float64:
		return calculateNumber(operator, v_, conditionValue)
	case string:
		return calculateString(operator, v_, conditionValue)
	case []interface{}:
		return calculateObjectArray(operator, v_, conditionValue)
	case map[interface{}]interface{}:
		return calculateObjectMap(operator, v_, conditionValue)
	case nil:
		return calculateNilAndValue(operator, conditionValue)
	default:
		return false, fmt.Errorf("calculate error: unsupported type: %v", reflect.TypeOf(objectValue))
	}
}

func calculateNilAndValue(operator int, v2 interface{}) (bool, error) {
	switch operator {
	case GreaterThan:
		fallthrough
	case GreaterThanOrEqual:
		fallthrough
	case LesserThan:
		fallthrough
	case LesserThanOrEqual:
		return false, fmt.Errorf("calculateNilAndValue error: can not compare nil and %v ", reflect.TypeOf(v2))
	case EqualTo:
		return v2 == nil, nil
	case NotEqual:
		return v2 != nil, nil
	default:
		// Not possible
		return false, fmt.Errorf("calculate number error: unsupported operator: %v", operator)
	}
}

func calculateNumber(operator int, n1 float64, v2 interface{}) (bool, error) {
	n2, err := tryParseToNumber(v2)
	if err != nil {
		return false, err
	}
	switch operator {
	case GreaterThan:
		return n1 > n2, nil
	case GreaterThanOrEqual:
		return n1 >= n2, nil
	case LesserThan:
		return n1 < n2, nil
	case LesserThanOrEqual:
		return n1 <= n2, nil
	case EqualTo:
		return n1 == n2, nil
	case NotEqual:
		return n1 != n2, nil
	default:
		// Not possible
		return false, fmt.Errorf("calculate number error: unsupported operator: %v", operator)
	}
}

func calculateString(operator int, s1 string, v2 interface{}) (bool, error) {
	s2, err := tryParseToString(v2)
	if err != nil {
		return false, err
	}
	switch operator {
	case EqualTo:
		return s1 == s2, nil
	case NotEqual:
		return s1 != s2, nil
	default:
		// Not possible
		return false, fmt.Errorf("calculate string error: unsupported operator: %v", operator)
	}
}

func calculateObjectArray(operator int, a1 []interface{}, v2 interface{}) (bool, error) {
	a2, err := tryParseToArray(v2)
	if err != nil {
		return false, err
	}
	switch operator {
	case EqualTo:
		bs, err := yaml.Marshal(a1)
		if err != nil {
			return false, err
		}
		var a1_ []interface{}
		if err := yaml.Unmarshal(bs, &a1_); err != nil {
			return false, err
		}
		return reflect.DeepEqual(a1_, a2), nil
	case NotEqual:
		bs, err := yaml.Marshal(a1)
		if err != nil {
			return false, err
		}
		var a1_ []interface{}
		if err := yaml.Unmarshal(bs, &a1_); err != nil {
			return false, err
		}
		return !reflect.DeepEqual(a1_, a2), nil
	default:
		return false, fmt.Errorf("calculate array error: unsupported operator: %v", operator)
	}
}

func calculateObjectMap(operator int, m1 map[interface{}]interface{}, v2 interface{}) (bool, error) {
	m2, err := tryParseToMap(v2)
	if err != nil {
		return false, err
	}
	switch operator {
	case EqualTo:
		bs, err := yaml.Marshal(m1)
		if err != nil {
			return false, err
		}
		var m1_ map[interface{}]interface{}
		if err := yaml.Unmarshal(bs, &m1_); err != nil {
			return false, err
		}
		return reflect.DeepEqual(m1_, m2), nil
	case NotEqual:
		bs, err := yaml.Marshal(m1)
		if err != nil {
			return false, err
		}
		var m1_ map[interface{}]interface{}
		if err := yaml.Unmarshal(bs, &m1_); err != nil {
			return false, err
		}
		return !reflect.DeepEqual(m1_, m2), nil
	default:
		return false, fmt.Errorf("calculate map error: unsupported operator: %v", operator)
	}
}

func tryParseToNumber(v interface{}) (float64, error) {
	switch x := v.(type) {
	case int8:
		return float64(x), nil
	case int32:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case int:
		return float64(x), nil
	case uint:
		return float64(x), nil
	case float32:
		return float64(x), nil
	case float64:
		return x, nil
	case string:
		return strconv.ParseFloat(x, 64)
	default:
		return 0, fmt.Errorf("parse to number error: unsupported type: %v", reflect.TypeOf(v))
	}
}

func tryParseToString(v interface{}) (string, error) {
	switch x := v.(type) {
	case string:
		return x, nil
	default:
		return "", fmt.Errorf("parse to string error: unsupported type: %v", reflect.TypeOf(v))
	}
}

func tryParseToArray(v interface{}) ([]interface{}, error) {
	switch x := v.(type) {
	case string:
		var a []interface{}
		if err := json.Unmarshal([]byte(x), &a); err != nil {
			return nil, err
		}
		bs, err := yaml.Marshal(a)
		if err != nil {
			return nil, err
		}
		var a_ []interface{}
		return a_, yaml.Unmarshal(bs, &a_)
	case []interface{}:
		return x, nil
	default:
		return nil, fmt.Errorf("parse to array error: unsupported type: %v", reflect.TypeOf(v))
	}
}

func tryParseToMap(v interface{}) (map[interface{}]interface{}, error) {
	switch x := v.(type) {
	case string:
		m := map[string]interface{}{}
		if err := json.Unmarshal([]byte(x), &m); err != nil {
			return nil, err
		}
		bs, err := yaml.Marshal(m)
		if err != nil {
			return nil, err
		}
		m_ := map[interface{}]interface{}{}
		return m_, yaml.Unmarshal(bs, &m_)
	case map[interface{}]interface{}:
		return x, nil
	default:
		return nil, fmt.Errorf("parse to map error: unsupported type: %v", reflect.TypeOf(v))
	}
}
