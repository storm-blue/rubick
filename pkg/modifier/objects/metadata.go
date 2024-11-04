package objects

const removedKey = "__metadata.__removed"
const _true = "true"
const _false = "false"

type Metadata interface {
	Removed() bool
	MarkRemoved(bool)
	Set(key string, value string)
	Get(key string) string
}

type _metadata map[string]string

func (m _metadata) Removed() bool {
	return m[removedKey] == _true
}

func (m _metadata) MarkRemoved(b bool) {
	if b {
		m[removedKey] = _true
	} else {
		m[removedKey] = _false
	}
}

func (m _metadata) Set(key string, value string) {
	m[key] = value
}

func (m _metadata) Get(key string) string {
	return m[key]
}
