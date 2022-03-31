package rtconfig

import (
	"errors"
	"time"

	"github.com/spf13/cast"
)

// Config errors.
var (
	ErrEmptyKey    = errors.New("key is empty")
	ErrNilVariable = errors.New("variable is nil")
)

// Value describes a config value.
type Value interface {
	IsNil() bool
	String() string
	Bool() bool
	Int() int
	Int32() int32
	Int64() int64
	Uint() uint
	Uint32() uint32
	Uint64() uint64
	Float64() float64
	Time() time.Time
	Duration() time.Duration
	IntSlice() []int
	StringSlice() []string
	StringMap() map[string]interface{}
	StringMapString() map[string]string
}

// WatcherCallback is a callback function for a variable watcher.
type WatcherCallback func(oldValue, newValue Value)

type value struct {
	value interface{}
}

func (v value) IsNil() bool                        { return v.value == nil }
func (v value) String() string                     { return cast.ToString(v.value) }
func (v value) Bool() bool                         { return cast.ToBool(v.value) }
func (v value) Int() int                           { return cast.ToInt(v.value) }
func (v value) Int32() int32                       { return cast.ToInt32(v.value) }
func (v value) Int64() int64                       { return cast.ToInt64(v.value) }
func (v value) Uint() uint                         { return cast.ToUint(v.value) }
func (v value) Uint32() uint32                     { return cast.ToUint32(v.value) }
func (v value) Uint64() uint64                     { return cast.ToUint64(v.value) }
func (v value) Float64() float64                   { return cast.ToFloat64(v.value) }
func (v value) Time() time.Time                    { return cast.ToTime(v.value) }
func (v value) Duration() time.Duration            { return cast.ToDuration(v.value) }
func (v value) IntSlice() []int                    { return cast.ToIntSlice(v.value) }
func (v value) StringSlice() []string              { return cast.ToStringSlice(v.value) }
func (v value) StringMap() map[string]interface{}  { return cast.ToStringMap(v.value) }
func (v value) StringMapString() map[string]string { return cast.ToStringMapString(v.value) }
