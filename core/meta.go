package core

type MetaDataKey string

// `MetaData` stores map of `interface{}` with
// `MetaDataKey`-type keys
type MetaData struct {
	meta map[MetaDataKey]interface{}
}

// Constructor allocates new `MetaData` struct
func NewMetaData() *MetaData {
	return &MetaData{
		meta: make(map[MetaDataKey]interface{}),
	}
}

// Deletes entry by key
func (m *MetaData) Del(key MetaDataKey) bool {
	if _, ok := m.meta[key]; !ok {
		return false
	} else {
		delete(m.meta, key)
		return true
	}
}

func (m *MetaData) SetBool(key MetaDataKey, val bool) bool {
	_, ok := m.meta[key]
	m.meta[key] = val
	return !ok
}

func (m *MetaData) GetBool(key MetaDataKey) (bool, bool) {
	var typeVal bool

	if val, ok := m.meta[key]; !ok {
		return typeVal, false
	} else {
		if typeVal, ok := val.(bool); !ok {
			return typeVal, false
		} else {
			return typeVal, true
		}
	}
}

func (m *MetaData) SetString(key MetaDataKey, val string) bool {
	_, ok := m.meta[key]
	m.meta[key] = val
	return !ok
}

func (m *MetaData) GetString(key MetaDataKey) (string, bool) {
	var typeVal string

	if val, ok := m.meta[key]; !ok {
		return typeVal, false
	} else {
		if typeVal, ok := val.(string); !ok {
			return typeVal, false
		} else {
			return typeVal, true
		}
	}
}

func (m *MetaData) SetInt(key MetaDataKey, val int) bool {
	_, ok := m.meta[key]
	m.meta[key] = val
	return !ok
}

func (m *MetaData) GetInt(key MetaDataKey) (int, bool) {
	var typeVal int

	if val, ok := m.meta[key]; !ok {
		return typeVal, false
	} else {
		if typeVal, ok := val.(int); !ok {
			return typeVal, false
		} else {
			return typeVal, true
		}
	}
}

func (m *MetaData) SetFloat(key MetaDataKey, val float64) bool {
	_, ok := m.meta[key]
	m.meta[key] = val
	return !ok
}

func (m *MetaData) GetFloat(key MetaDataKey) (float64, bool) {
	var typeVal float64

	if val, ok := m.meta[key]; !ok {
		return typeVal, false
	} else {
		if typeVal, ok := val.(float64); !ok {
			return typeVal, false
		} else {
			return typeVal, true
		}
	}
}
