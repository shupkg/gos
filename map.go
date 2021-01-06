package gos

import (
	"encoding/json"
	"reflect"
	"time"
)

type M map[string]interface{}

func (m M) Set(key string, val interface{}) M {
	if val != nil {
		switch v := val.(type) {
		case string:
			if val == "" {
				return m
			}
		case int, int8, int16, int32, int64:
			if v := reflect.ValueOf(val).Int(); v == 0 {
				return m
			}
		case uint, uint8, uint16, uint32, uint64:
			if v := reflect.ValueOf(val).Uint(); v == 0 {
				return m
			}
		case float64, float32:
			if v := reflect.ValueOf(val).Float(); v == 0 {
				return m
			}
		case time.Time:
			if v.IsZero() || v.Before(time.Unix(0, 0)) {
				return m
			}
		}
		m[key] = val
	}
	return m
}

func (m M) Get(key string) (interface{}, bool) {
	val, find := m[key]
	return val, find
}

func (m M) Del(key string) M {
	delete(m, key)
	return m
}

func (m M) Merge(maps ...M) M {
	for _, mp := range maps {
		for key, val := range mp {
			m.Set(key, val)
		}
	}
	return m
}

func (m M) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

type Map struct {
	M
}

func (m *Map) Set(key string, val interface{}) *Map {
	m.M = m.M.Set(key, val)
	return m
}

func (m *Map) Del(key string) *Map {
	m.M = m.M.Del(key)
	return m
}

func (m *Map) Merge(maps ...M) *Map {
	m.M = m.M.Merge(maps...)
	return m
}

func (m Map) String() string {
	v, _ := json.MarshalIndent(m.M, "", "  ")
	return string(v)
}

func (m Map) MarshalBinary() ([]byte, error) {
	return json.Marshal(m.M)
}

func (m Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.M)
}
