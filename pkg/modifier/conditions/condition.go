package conditions

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
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

type Condition interface {
	Calculate(objects.StructuredObject) (bool, error)
	String() string
	And(condition Condition) Condition
	Or(condition Condition) Condition
}

type CombinationCondition struct {
	left     Condition
	right    Condition
	operator int
}

func (c *CombinationCondition) Calculate(object objects.StructuredObject) (bool, error) {
	leftResult, err := c.left.Calculate(object)
	if err != nil {
		return false, err
	}
	rightResult, err := c.right.Calculate(object)
	if err != nil {
		return false, err
	}

	switch c.operator {
	case And:
		return leftResult && rightResult, nil
	case Or:
		return leftResult || rightResult, nil
	default:
		return false, fmt.Errorf("calculate error: unsupported operator: %v", c.operator)
	}
}

func (c *CombinationCondition) String() string {
	var leftString, rightString string
	if _, ok := c.left.(*CombinationCondition); ok {
		leftString = c.left.String()
	} else {
		leftString = fmt.Sprintf("(%v)", c.left.String())
	}

	rightString = c.right.String()

	switch c.operator {
	case And:
		return fmt.Sprintf("%v && (%v)", leftString, rightString)
	case Or:
		return fmt.Sprintf("%v || (%v)", leftString, rightString)
	default:
		return fmt.Sprintf("Format error: unsupported operator: %v", c.operator)
	}
}

func (c *CombinationCondition) And(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: And,
	}
}

func (c *CombinationCondition) Or(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: Or,
	}
}

func New() CreateCondition_Start {
	return CreateCondition_Start{}
}

type CreateCondition_Start struct{}

func (c CreateCondition_Start) ValueOf(key string) CreateCondition_ValueOf {
	return CreateCondition_ValueOf{
		condition: &valueOfCondition{
			operator: None,
			key:      key,
		},
	}
}

func (c CreateCondition_Start) Not(condition Condition) Condition {
	return &notCondition{
		condition: condition,
	}
}

func (c CreateCondition_Start) Exists(key string) Condition {
	return &existsCondition{
		key: key,
	}
}

func (c CreateCondition_Start) LengthOf(key string) CreateCondition_LengthOf {
	return CreateCondition_LengthOf{
		condition: &lengthOfCondition{
			operator: None,
			key:      key,
		},
	}
}

type CreateCondition_ValueOf struct {
	condition *valueOfCondition
}

func (c CreateCondition_ValueOf) GreaterThan(value interface{}) Condition {
	c.condition.operator = GreaterThan
	c.condition.value = value
	return c.condition
}

func (c CreateCondition_ValueOf) GreaterThanOrEqual(value interface{}) Condition {
	c.condition.operator = GreaterThanOrEqual
	c.condition.value = value
	return c.condition
}

func (c CreateCondition_ValueOf) LesserThan(value interface{}) Condition {
	c.condition.operator = LesserThan
	c.condition.value = value
	return c.condition
}

func (c CreateCondition_ValueOf) LesserThanOrEqual(value interface{}) Condition {
	c.condition.operator = LesserThanOrEqual
	c.condition.value = value
	return c.condition
}

func (c CreateCondition_ValueOf) EqualTo(value interface{}) Condition {
	c.condition.operator = EqualTo
	c.condition.value = value
	return c.condition
}

func (c CreateCondition_ValueOf) NotEqual(value interface{}) Condition {
	c.condition.operator = NotEqual
	c.condition.value = value
	return c.condition
}

type CreateCondition_LengthOf struct {
	condition *lengthOfCondition
}

func (c CreateCondition_LengthOf) GreaterThan(value int) Condition {
	c.condition.operator = GreaterThan
	c.condition.value = value
	return c.condition
}

func (c CreateCondition_LengthOf) GreaterThanOrEqual(value int) Condition {
	c.condition.operator = GreaterThanOrEqual
	c.condition.value = value
	return c.condition
}

func (c CreateCondition_LengthOf) LesserThan(value int) Condition {
	c.condition.operator = LesserThan
	c.condition.value = value
	return c.condition
}

func (c CreateCondition_LengthOf) LesserThanOrEqual(value int) Condition {
	c.condition.operator = LesserThanOrEqual
	c.condition.value = value
	return c.condition
}

func (c CreateCondition_LengthOf) EqualTo(value int) Condition {
	c.condition.operator = EqualTo
	c.condition.value = value
	return c.condition
}

func (c CreateCondition_LengthOf) NotEqual(value int) Condition {
	c.condition.operator = NotEqual
	c.condition.value = value
	return c.condition
}

type valueOfCondition struct {
	operator int
	key      string
	value    interface{}
}

func (s *valueOfCondition) And(condition Condition) Condition {
	return &CombinationCondition{
		left:     s,
		right:    condition,
		operator: And,
	}
}

func (s *valueOfCondition) Or(condition Condition) Condition {
	return &CombinationCondition{
		left:     s,
		right:    condition,
		operator: Or,
	}
}

func (s *valueOfCondition) Calculate(object objects.StructuredObject) (bool, error) {
	v, err := object.Get(s.key)
	if err != nil {
		return false, err
	}

	if v == nil {
		return false, nil
	}

	return calculate(s.operator, v, s.value)
}

func (s *valueOfCondition) String() string {
	operatorString := ""
	switch s.operator {
	case GreaterThan:
		operatorString = ">"
	case GreaterThanOrEqual:
		operatorString = ">="
	case LesserThan:
		operatorString = "<"
	case LesserThanOrEqual:
		operatorString = "<="
	case EqualTo:
		operatorString = "=="
	case NotEqual:
		operatorString = "!="
	default:
	}

	valueString := ""
	switch tv := s.value.(type) {
	case int8:
		valueString = strconv.Itoa(int(tv))
	case int32:
		valueString = strconv.Itoa(int(tv))
	case int64:
		valueString = strconv.Itoa(int(tv))
	case int:
		valueString = strconv.Itoa(tv)
	case uint:
		valueString = strconv.Itoa(int(tv))
	case float32:
		valueString = strconv.FormatFloat(float64(tv), 'g', 20, 64)
	case float64:
		valueString = strconv.FormatFloat(tv, 'g', 20, 64)
	case string:
		valueString = tv
	case []interface{}:
		valueString = "<array>"
	case map[interface{}]interface{}:
		valueString = "<object>"
	default:
		valueString = "<object>"
	}

	return fmt.Sprintf("%v %v %v", s.key, operatorString, valueString)
}

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
	default:
		return false, fmt.Errorf("calculate error: unsupported type: %v", reflect.TypeOf(objectValue))
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

type notCondition struct {
	condition Condition
}

func (c *notCondition) And(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: And,
	}
}

func (c *notCondition) Or(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: Or,
	}
}

func (c *notCondition) Calculate(object objects.StructuredObject) (bool, error) {
	result, err := c.condition.Calculate(object)
	if err != nil {
		return false, err
	}
	return !result, nil
}

func (c *notCondition) String() string {
	return fmt.Sprintf("!(%v)", c.condition.String())
}

type existsCondition struct {
	key string
}

func (e *existsCondition) And(condition Condition) Condition {
	return &CombinationCondition{
		left:     e,
		right:    condition,
		operator: And,
	}
}

func (e *existsCondition) Or(condition Condition) Condition {
	return &CombinationCondition{
		left:     e,
		right:    condition,
		operator: Or,
	}
}

func (e *existsCondition) Calculate(object objects.StructuredObject) (bool, error) {
	result, err := object.Get(e.key)
	if err != nil {
		return false, err
	}
	return result != nil, nil
}

func (e *existsCondition) String() string {
	return fmt.Sprintf("EXISTS(%v)", e.key)
}

type lengthOfCondition struct {
	operator int
	key      string
	value    int
}

func (c *lengthOfCondition) And(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: And,
	}
}

func (c *lengthOfCondition) Or(condition Condition) Condition {
	return &CombinationCondition{
		left:     c,
		right:    condition,
		operator: Or,
	}
}

func (c *lengthOfCondition) Calculate(object objects.StructuredObject) (bool, error) {
	v, err := object.Get(c.key)
	if err != nil {
		return false, err
	}

	if v == nil {
		return false, nil
	}

	size := 0
	switch v_ := v.(type) {
	case []interface{}:
		size = len(v_)
	case map[interface{}]interface{}:
		size = len(v_)
	default:
		return false, fmt.Errorf("calculate error: unsupported type: %v", reflect.TypeOf(v))
	}

	return calculateNumber(c.operator, float64(size), c.value)
}

func (c *lengthOfCondition) String() string {
	operatorString := ""
	switch c.operator {
	case GreaterThan:
		operatorString = ">"
	case GreaterThanOrEqual:
		operatorString = ">="
	case LesserThan:
		operatorString = "<"
	case LesserThanOrEqual:
		operatorString = "<="
	case EqualTo:
		operatorString = "=="
	case NotEqual:
		operatorString = "!="
	default:
	}

	return fmt.Sprintf("%v %v %v", c.key, operatorString, c.value)
}
