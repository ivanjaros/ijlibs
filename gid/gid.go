package gid

// Package GID provides UUID, or Global ID, generator.
// Instead of using different implementations of UUID v4, like github.com/google/uuid,
// we have one single implementation to be used throughout the whole code base.
// Additionally, instead of UUID v4, this package uses the ULID implementation - https://github.com/ulid/spec.
// This has advantages over UUID in that it is shorter by 10 characters, from 36 to 26,
// it uses only alpha-numeric characters without dashes, it is time k-sortable,
// and most importantly - it is case insensitive.
// Due to the nature of the implementation, this package implements sync.Mutex to prevent data race issues.
// Performance of mutex vs sync.Pool is described in here https://github.com/oklog/ulid/issues/33#issuecomment-455659488.
// ULID is also perfectly suitable for high throughput(like api gateway) due to its significantly low
// collision rate of 1,208,925,819,614,629,174,706,176 / 1 ms.
// Note: ULID is case insensitive and produces upper-cased format by default. To preserve uniformity of gids
// throughout the system and because gids might get confused with other strings since they do no look like uuids
// we always enforce upper-case formatting. Mixed-case or lower-case strings are still valid but upper-case should
// be always enforced so storage engines will not confuse aBc and ABC.
// ULIDS are essentially [16]byte so storing them as such can save 10 bytes in comparison to their string versions.
// Due to the nature of this library and its use, in case of failure it will always panics instead of returning an error.

import (
	"github.com/oklog/ulid"
	"io"
	"math/rand"
	"sync"
	"time"
)

const (
	ByteLength   = 16
	StringLength = 26
)

var (
	f Factory
)

func init() {
	f = NewFactory(nil)
}

type Factory interface {
	New() (string, error)
	MustNew() string
}

// Creates new GID factory.
// If r is nil, rand.New(rand.NewSource(time.Now().UnixNano())) will be used by default.
func NewFactory(r *rand.Rand) Factory {
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	return &defaultFactory{
		entropy: ulid.Monotonic(r, 0),
	}
}

type defaultFactory struct {
	mx      sync.Mutex
	entropy io.Reader
}

func (f *defaultFactory) New() (string, error) {
	f.mx.Lock()
	id, err := ulid.New(ulid.Timestamp(time.Now()), f.entropy)
	f.mx.Unlock()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func (f *defaultFactory) MustNew() string {
	id, err := f.New()
	if err != nil {
		panic(err)
	}
	return id
}

// Generates new GID, panics in case of failure.
func New() string {
	return f.MustNew()
}

// Validates the id as true GID.
func IsValid(id string) bool {
	_, e := ulid.ParseStrict(id)
	return e == nil
}

// Validates the id as true GID in byte form.
func IsValidBytes(id []byte) bool {
	var typed ulid.ULID
	if err := typed.UnmarshalBinary(id); err != nil {
		return false
	}
	if _, err := ulid.ParseStrict(typed.String()); err != nil {
		return false
	}
	return true
}

// Valid GID is [16]byte array but for ease of use we return byte slice.
// Panics if provided id is not a GID.
func ToBytes(id string) []byte {
	parsed := ulid.MustParseStrict(id)
	return parsed[:]
}

// Valid GID is [16]byte array but for ease of use we return byte slice.
func TryToBytes(id string) ([]byte, bool) {
	parsed, err := ulid.ParseStrict(id)
	return parsed[:], err == nil
}

// Valid GID is [16]byte array but for ease of use we accept byte slice.
// Panics if provided id is not a GID.
func FromBytes(id []byte) string {
	var typed ulid.ULID
	// this only copies the bytes into the ulid and nothing more
	if err := typed.UnmarshalBinary(id); err != nil {
		panic(err)
	}
	// we have to validate the ulid from provided id before returning GID
	if _, err := ulid.ParseStrict(typed.String()); err != nil {
		panic(err)
	}
	return typed.String()
}

// Valid GID is [16]byte array but for ease of use we accept byte slice.
func TryFromBytes(id []byte) (string, bool) {
	var typed ulid.ULID
	// this only copies the bytes into the ulid and nothing more
	if err := typed.UnmarshalBinary(id); err != nil {
		return "", false
	}
	// we have to validate the ulid from provided id before returning GID
	if _, err := ulid.ParseStrict(typed.String()); err != nil {
		return "", false
	}
	return typed.String(), true
}

func CreatedAt(id string) time.Time {
	parsed := ulid.MustParseStrict(id)
	return time.UnixMilli(int64(parsed.Time()))
}

func CreatedAtBytes(id []byte) time.Time {
	var typed ulid.ULID
	// this only copies the bytes into the ulid and nothing more
	if err := typed.UnmarshalBinary(id); err != nil {
		return time.Unix(0, 0)
	}
	// we have to validate the ulid from provided id before returning GID
	if _, err := ulid.ParseStrict(typed.String()); err != nil {
		return time.Unix(0, 0)
	}
	return time.UnixMilli(int64(typed.Time()))
}

func StringsToArgs(ids []string) []interface{} {
	args := make([]interface{}, len(ids))
	for k, v := range ids {
		args[k] = ToBytes(v)
	}
	return args
}

func BinsToStrings(bins [][]byte) []string {
	out := make([]string, len(bins))
	for k, v := range bins {
		out[k] = FromBytes(v)
	}
	return out
}
