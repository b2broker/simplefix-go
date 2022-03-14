package fix

import (
	"fmt"
	"strconv"
	"time"
)

// Value is an interface implementing all basic methods required to process field values of FIX messages.
type Value interface {
	// ToBytes returns a byte representation of a field value.
	ToBytes() []byte

	// FromBytes parses values stored in a byte array.
	FromBytes([]byte) error

	// Value returns a field value.
	Value() interface{}

	// String returns a string representation of a value.
	String() string

	// IsNull is used to check whether a field value is empty.
	IsNull() bool

	// Set replaces a specified field value with a value of a corresponding type.
	Set(d interface{}) error
}

// Raw is a structure that is used to convert the message data to a byte array.
type Raw struct {
	value []byte
}

// NewRaw creates a new instance of a Raw object.
func NewRaw(v []byte) *Raw {
	return &Raw{
		value: v,
	}
}

func (v *Raw) ToBytes() []byte {
	return v.value
}

func (v *Raw) FromBytes(d []byte) (err error) {
	v.value = d
	return nil
}

func (v *Raw) IsNull() bool {
	return v.value == nil
}

func (v *Raw) Value() interface{} {
	return v.value
}

// Set parses and assigns the field value stored as a byte array.
func (v *Raw) Set(d interface{}) error {
	if res, ok := d.([]byte); ok {
		v.value = res
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "Byte Array")
}

func (v *Raw) String() string {
	return string(v.value)
}

// String is a structure used for converting string values.
type String struct {
	value string
	valid bool
}

func NewString(v string) *String {
	return &String{value: v, valid: true}
}

// Set parses and assigns the field value stored as a string.
func (v *String) Set(d interface{}) error {
	if d == nil {
		v.valid = false
		return nil
	}

	if res, ok := d.(string); ok {
		v.value = res
		v.valid = true
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "String")
}

func (v *String) ToBytes() []byte {
	if !v.valid || v.value == "" {
		return nil
	}
	return []byte(v.value)
}

func (v *String) IsNull() bool {
	return !v.valid
}

func (v *String) Value() interface{} {
	return v.value
}

func (v *String) FromBytes(d []byte) (err error) {
	if d == nil {
		v.valid = false
		return nil
	}

	v.valid = true
	v.value = string(d)

	return nil
}

func (v *String) String() string {
	return v.value
}

// Int is a structure used for converting integer values.
type Int struct {
	value int
	valid bool
}

func NewInt(value int) *Int {
	return &Int{value: value, valid: true}
}

func (v *Int) IsNull() bool {
	return !v.valid
}

// Set parses and assigns the field value stored as an integer number.
func (v *Int) Set(d interface{}) error {
	if d == nil {
		v.valid = false
		return nil
	}

	if res, ok := d.(int); ok {
		v.value = res
		v.valid = true
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "Int")
}

func (v *Int) Value() interface{} {
	return v.value
}
func (v *Int) String() string {
	return strconv.Itoa(v.value)
}

func (v *Int) FromBytes(d []byte) (err error) {
	if d == nil {
		v.valid = false
		return nil
	}

	v.valid = true
	v.value, err = strconv.Atoi(string(d))

	return err
}

func (v *Int) ToBytes() []byte {
	if !v.valid {
		return nil
	}
	return []byte(strconv.Itoa(v.value))
}

// Uint is a structure used for converting values to the uint64 type.
type Uint struct {
	value uint64
	valid bool
}

func NewUint(value uint64) *Uint {
	return &Uint{value: value}
}

// Set parses and assigns the field value stored as a uint64 number.
func (v *Uint) Set(d interface{}) error {
	if d == nil {
		v.valid = false
		return nil
	}

	if res, ok := d.(uint64); ok {
		v.value = res
		v.valid = true
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "Uint")
}

func (v *Uint) IsNull() bool {
	return !v.valid
}

func (v *Uint) FromBytes(d []byte) (err error) {
	if d == nil {
		v.valid = false
		return nil
	}

	v.valid = true
	v.value, err = strconv.ParseUint(string(d), 10, 64)

	return err
}

func (v *Uint) Value() interface{} {
	return v.value
}

func (v *Uint) String() string {
	return fmt.Sprintf("%d", v.value)
}

func (v *Uint) ToBytes() []byte {
	if !v.valid {
		return nil
	}
	return []byte(strconv.FormatUint(v.value, 10))
}

// Float is a structure used for converting values to the float64 type.
type Float struct {
	value float64
	valid bool
}

func NewFloat(value float64) *Float {
	return &Float{value: value}
}

func (v *Float) IsNull() bool {
	return !v.valid
}

func (v *Float) Value() interface{} {
	return v.value
}

func (v *Float) FromBytes(d []byte) (err error) {
	if d == nil {
		v.valid = false
		return nil
	}

	v.valid = true
	v.value, err = strconv.ParseFloat(string(d), 64)

	return err
}

func (v *Float) ToBytes() []byte {
	if !v.valid {
		return nil
	}
	return []byte(strconv.FormatFloat(v.value, 'f', -1, 64))
}

func (v *Float) String() string {
	return fmt.Sprintf("%f", v.value)
}

// Set parses and assigns the field value stored as a float64 number.
func (v *Float) Set(d interface{}) error {
	if d == nil {
		v.valid = false
		return nil
	}

	if res, ok := d.(float64); ok {
		v.value = res
		v.valid = true
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "Float")
}

// Time is a structure used for converting date-time values.
type Time struct {
	value time.Time
	valid bool
}

func NewTime(value time.Time) *Time {
	return &Time{value: value, valid: true}
}

// Set parses and assigns the field value stored in the date-time format.
func (v *Time) Set(d interface{}) error {
	if d == nil {
		v.valid = false
		return nil
	}

	if res, ok := d.(time.Time); ok {
		v.value = res
		v.valid = true
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "Date-Time")
}

func (v *Time) IsNull() bool {
	return !v.valid
}

func (v *Time) Value() interface{} {
	return v.value
}

func (v *Time) ToBytes() []byte {
	if !v.valid {
		return nil
	}
	return []byte(v.value.Format(TimeLayout)) // TODO: set layout outside.
}

func (v *Time) FromBytes(d []byte) (err error) {
	if d == nil {
		v.valid = false
		return nil
	}

	v.valid = true
	v.value, err = time.Parse(TimeLayout, string(d))

	return err
}

func (v *Time) String() string {
	return v.value.Format(TimeLayout)
}

const (
	True  = "Y"
	False = "N"
)

// Bool is a structure used for converting Boolean values.
type Bool struct {
	value bool
	valid bool
}

func (v *Bool) ToBytes() []byte {
	if !v.valid {
		return nil
	}

	if v.value {
		return []byte(True)
	}
	return []byte(False)
}

func (v *Bool) FromBytes(d []byte) error {
	if d == nil {
		v.valid = false
		return nil
	}

	v.valid = true
	v.value = string(d) == True

	return nil
}

func (v *Bool) Value() interface{} {
	return v.value
}

func (v *Bool) String() string {
	if !v.valid {
		return ""
	}

	if v.value {
		return True
	}
	return False
}

func (v *Bool) IsNull() bool {
	return !v.valid
}

// Set parses and assigns the field value stored in the Boolean format.
func (v *Bool) Set(d interface{}) error {
	if d == nil {
		v.valid = false
		return nil
	}

	if res, ok := d.(bool); ok {
		v.value = res
		v.valid = true
		return nil
	}

	return fmt.Errorf("could not convert %s to %s", d, "Boolean")
}
