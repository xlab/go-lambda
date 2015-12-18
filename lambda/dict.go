package lambda

import "errors"

// HAS NOT BEEN AUTOGENERATED on Fri Dec 18 19:05:46 2015 MSK
// I wish this has been autogenerated.

var ErrNoValue = errors.New("no value")

type Dict struct {
	m map[string]interface{}
}

func MakeDict(size int) *Dict {
	return &Dict{
		m: make(map[string]interface{}, size),
	}
}

func (d Dict) Int(key string) (v int, err error) {
	if vv, ok := d.m[key]; ok {
		if v, ok := vv.(int); ok {
			return v, nil
		}
	}
	err = ErrNoValue
	return
}

func (d Dict) SetInt(key string, v int) error {
	d.m[key] = v
	return nil
}

func (d Dict) Int64(key string) (v int64, err error) {
	if vv, ok := d.m[key]; ok {
		if v, ok := vv.(int64); ok {
			return v, nil
		}
	}
	err = ErrNoValue
	return
}

func (d Dict) SetInt64(key string, v int64) error {
	d.m[key] = v
	return nil
}

func (d Dict) Bool(key string) (v bool, err error) {
	if vv, ok := d.m[key]; ok {
		if v, ok := vv.(bool); ok {
			return v, nil
		}
	}
	err = ErrNoValue
	return
}

func (d Dict) SetBool(key string, v bool) error {
	d.m[key] = v
	return nil
}

func (d Dict) String(key string) (v string, err error) {
	if vv, ok := d.m[key]; ok {
		if v, ok := vv.(string); ok {
			return v, nil
		}
	}
	err = ErrNoValue
	return
}

func (d Dict) SetString(key, v string) error {
	d.m[key] = v
	return nil
}

func (d Dict) Float64(key string) (v float64, err error) {
	if vv, ok := d.m[key]; ok {
		if v, ok := vv.(float64); ok {
			return v, nil
		}
	}
	err = ErrNoValue
	return
}

func (d Dict) SetFloat64(key string, v float64) error {
	d.m[key] = v
	return nil
}