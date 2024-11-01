package action

import (
	"github.com/storm-blue/rubick/pkg/log"
	"github.com/storm-blue/rubick/pkg/modifier/objects"
)

func NewContext(logKeys []string) Context {
	return &actionContext{
		logKeys: logKeys,
	}
}

type Context interface {
	Log(object objects.StructuredObject, action Action, err error)
	Logs() []*Log
}

type actionContext struct {
	logKeys []string
	logs    []*Log
}

func (c *actionContext) Logs() []*Log {
	return c.logs
}

func (c *actionContext) Log(object objects.StructuredObject, action Action, err error) {
	m := map[string]interface{}{}

	for _, key := range c.logKeys {
		if v, err := object.Get(key); err != nil {
			log.Errorf("log action error: %v", err)
		} else {
			m[key] = v
		}
	}
	c.logs = append(c.logs, &Log{
		logKeysMap: m,
		action:     action,
		err:        err,
	})
}

type Log struct {
	logKeysMap map[string]interface{}
	action     Action
	err        error
}

func (l *Log) GetKey(key string) interface{} {
	return l.logKeysMap[key]
}

func (l *Log) Action() Action {
	return l.action
}

func (l *Log) Err() error {
	return l.err
}
