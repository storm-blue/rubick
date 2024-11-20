package conditions

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

func (c CreateCondition_Start) LengthOf(key string) CreateCondition_LengthOf {
	return CreateCondition_LengthOf{
		condition: &lengthOfCondition{
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

func (c CreateCondition_Start) HasPrefix(key string, prefix string) Condition {
	return &hasPrefixCondition{
		key:    key,
		prefix: prefix,
	}
}

func (c CreateCondition_Start) HasSuffix(key string, suffix string) Condition {
	return &hasSuffixCondition{
		key:    key,
		suffix: suffix,
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
