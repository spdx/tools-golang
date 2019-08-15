package rdf2v1

import (
	"regexp"
	"strings"
)

// interface for elements
type Value interface {
	V() string
}

// string
type ValueStr struct {
	Val string
}

// ValueStr from string
func Str(v string) ValueStr { return ValueStr{v} }

// string from ValueStr
func (v ValueStr) V() string { return v.Val }

// func (v ValueList) L() []string {

func (v ValueStr) Equal(w ValueStr) bool { return v.Val == w.Val }

// bool
type ValueBool struct {
	Val bool
}

func Bool(v bool) ValueBool { return ValueBool{v} }
func (v ValueBool) Value() string {
	if v.Val {
		return "true"
	}
	return "false"
}

type ValueCreator struct {
	val   string
	what  string
	name  string
	email string
}

// Get the original value of this ValueCreator
func (c ValueCreator) V() string { return c.val }

// Get the `what` part from the format `what: name (email)`.
func (c ValueCreator) What() string { return c.what }

// Get the `name` part from the format `what: name (email)`
func (c ValueCreator) Name() string { return c.name }

// Get the `email` part from the format `what: name (email)`
func (c ValueCreator) Email() string { return c.email }

// parses and populates values of value creator
func (c *ValueCreator) SetValue(v string) {
	c.val = v
	RegexCreator := regexp.MustCompile("^([^:]*):([^\\(]*)(\\((.*)\\))?$")
	match := RegexCreator.FindStringSubmatch(v)
	if len(match) == 5 {
		c.what = strings.TrimSpace(match[1])
		c.name = strings.TrimSpace(match[2])
		c.email = strings.TrimSpace(match[4])
	}
}

// Create and populate a new ValueCreator.
func ValueCreatorNew(val string) ValueCreator {
	var valuecreator ValueCreator
	(&valuecreator).SetValue(val)
	return valuecreator
}

type ValueDate struct {
	Val string
}

func (d ValueDate) ValDate() string    { return d.Val }
func (d *ValueDate) SetValue(v string) { d.Val = v }

// New ValueDate.
func ValueDateNew(val string) ValueDate {
	var valuedate ValueDate
	(&valuedate).SetValue(val)
	return valuedate
}

func ValueList(list []ValueStr) []string {
	var str []string
	for _, v := range list {
		str = append(str, v.Val)
	}
	return str
}
func ValueStrList(list []string) []ValueStr {
	var str []ValueStr
	for _, v := range list {
		str = append(str, Str(v))
	}
	return str
}
