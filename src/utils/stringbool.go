package utils

import (
	"strings"
)

type StringBool string

const enable StringBool = "enable"
const disable StringBool = "disable"
const enableBool StringBool = "true"
const disableBool StringBool = "false"

func (s *StringBool) check() bool {
	*s = StringBool(strings.ToLower(string(*s)))
	return *s == enable || *s == disable || *s == enableBool || *s == disableBool
}

func (s *StringBool) is(v StringBool, defaultVal ...bool) (res bool) {
	if !s.check() {
		if len(defaultVal) == 1 {
			res = defaultVal[0]
			return
		} else {
			return false
		}
	}

	return *s == v
}

func (s *StringBool) IsEnable(defaultVal ...bool) (res bool) {
	res = s.is(enable, defaultVal...) || s.is(enableBool, defaultVal...)
	return
}

func (s *StringBool) IsDisable(defaultVal ...bool) (res bool) {
	res = s.is(disable, defaultVal...) || s.is(disableBool, defaultVal...)
	return
}

func (s *StringBool) setDefault(v StringBool) {
	if !s.check() {
		*s = v
	}
}

func (s *StringBool) SetDefaultEnable() {
	s.setDefault(enable)
}

func (s *StringBool) SetDefaultDisable() {
	s.setDefault(disable)
}

func (s *StringBool) ToString() string {
	if s.IsEnable() {
		return string(enable)
	}
	return string(disable)
}

func (s *StringBool) ToStringDefaultEnable() string {
	if s.IsEnable(true) {
		return string(enable)
	}
	return string(disable)
}

func (s *StringBool) ToStringDefaultDisable() string {
	if s.IsEnable(false) {
		return string(enable)
	}
	return string(disable)
}

func (s *StringBool) ToBool(defaultVal ...bool) bool {
	return s.IsEnable(defaultVal...)
}
