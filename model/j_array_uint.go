package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JsonArrayUint []uint

func (o JsonArrayUint) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *JsonArrayUint) Scan(input interface{}) (err error) {

	switch v := input.(type) {
	case []byte:
		return json.Unmarshal(v, o)
	case string:
		return json.Unmarshal([]byte(v), o)
	default:
		err = fmt.Errorf("unexpected type %T in JsonArrayUint", v)
	}
	return err
}
func (m JsonArrayUint) MarshalJSONArrayUint() ([]byte, error) {
	return json.Marshal(m)
}

func (m *JsonArrayUint) UnmarshalJSONArrayUint(data []byte) error {
	return json.Unmarshal(data, m)
}
func (j JsonArrayUint) IsNull() bool {
	return len(j) == 0
}

func (j JsonArrayUint) Equals(j1 JsonArrayUint) bool {
	t1 := []uint(j)
	t2 := []uint(j1)
	if len(t1) != len(t2) {
		return false
	}
	for k, vv := range t1 {
		if t2[k] != vv {
			return false
		}
	}
	return true
}

func (j JsonArrayUint) AppendOrRemove(e uint) (isExist bool, a JsonArrayUint) {
	temp := []uint{}
	for _, v := range j {
		if v == e {
			isExist = true
			continue
		}
		temp = append(temp, v)
	}
	//不存在 添加e
	if !isExist {
		temp = append(temp, e)
	}
	a = JsonArrayUint(temp)
	return isExist, a
}
