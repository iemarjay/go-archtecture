package cache

import "time"

type cache interface {
	Set(key string, value string, expiry time.Duration) error
	Get(key string) (string, error)
	Forget(key string) error
}

type Data struct {
	value string
	ttl   time.Time
}

type Map struct {
	values map[string]*Data
}

func NewMap(values map[string]*Data) *Map {
	return &Map{values: values}
}

func (m *Map) Set(key string, value string, expiry time.Duration) error {
	ttl := time.Now().Add(expiry)
	m.values[key] = &Data{value: value, ttl: ttl}
	return nil
}

func (m *Map) Get(key string) (string, error) {
	if data, exists := m.values[key]; exists && data.ttl.After(time.Now()) {
		return data.value, nil
	}

	return "", nil
}

func (m *Map) Forget(key string) error {
	delete(m.values, key)
	return nil
}
