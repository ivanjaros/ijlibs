package utils

// @todo: these are the former "utils" for backwards compatibility.
// when refactoring is done, move them where needed and remove unused ones.

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/rand"
	"net/mail"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// this allows us to override the time during testing
var TimeNow = time.Now

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func InArray(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

func IntInArray(needle int, haystack []int) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

func UintInArray(needle uint, haystack []uint) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

func LengthRange(text string, min int, max int) bool {
	ln := len(text)
	return ln >= min && ln <= max
}

func ArrayUnique(arr []string) []string {
	if len(arr) < 2 {
		return arr
	}

	sort.Strings(arr)
	j := 0
	for i := 1; i < len(arr); i++ {
		if arr[j] == arr[i] {
			continue
		}
		j++
		arr[j] = arr[i]
	}
	return arr[:j+1]
}

func ArrayIntersect(a []string, b []string) []string {
	c := make([]string, 0, len(a))
	for _, v := range a {
		if InArray(v, b) {
			c = append(c, v)
		}
	}
	return c
}

func ArrayDiff(a []string, b []string) []string {
	c := make([]string, 0, len(a))
	for _, v := range a {
		if InArray(v, b) == false {
			c = append(c, v)
		}
	}
	return c
}

func NumRange(num int, min int, max int) bool {
	return num >= min && num <= max
}

func Hash_sha1(value string) string {
	if value == "" {
		return ""
	}
	hasher := sha1.New()
	hasher.Write([]byte(value))
	return hex.EncodeToString(hasher.Sum(nil))
}

func BackTrace(wrap int) (string, string, int) {
	function, file, line, _ := runtime.Caller(1 + wrap)

	i := strings.LastIndex(file, "/")
	if i != -1 {
		file = file[i+1:]
	}

	return file, runtime.FuncForPC(function).Name(), line
}

func MaxUint(uints []uint) uint {
	var max uint
	for _, num := range uints {
		if num > max {
			max = num
		}
	}
	return max
}

func CleanEmptyStrings(in []string) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func IsEmail(email string) bool {
	if _, ok := mail.ParseAddress("Foo <" + email + ">"); ok != nil {
		return false
	}
	return true
}

func FormatTimeStamp(ts int64, emptyOnZero ...bool) string {
	if len(emptyOnZero) > 0 && emptyOnZero[0] == true && ts == 0 {
		return ""
	}
	return time.Unix(ts, 0).Format(time.RFC3339)
}

func GetTimeStampDateRange(ts int64) (int64, int64) {
	y, m, d := time.Unix(ts, 0).Date()
	from := time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix()
	to := time.Date(y, m, d, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1).Unix() - 1
	return from, to
}

// expected format is "%d-%d" as min-max, for example: 2-16
func ParseDomainLength(l string) (min int, max int, err error) {
	parts := strings.Split(l, "-")
	if len(parts) != 2 {
		err = errors.New("invalid domain length format")
		return
	}

	min, err = strconv.Atoi(parts[0])
	if err != nil {
		return
	}

	max, err = strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	if min > max {
		err = errors.New("invalid domain length range")
	}

	return
}

// parses TLD periods, expected format is "N,N,N.."
func ParsePeriods(p string) ([]int, error) {
	parts := strings.Split(p, ",")
	periods := make([]int, len(parts))
	for k := range parts {
		n, err := strconv.Atoi(parts[k])
		if err != nil {
			return periods, err
		}
		periods[k] = n
	}
	return periods, nil
}

func RandomString(size int) string {
	chars := []rune("23456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")
	l := len(chars)
	uid := make([]rune, size)
	for k := range uid {
		uid[k] = chars[r.Intn(l)]
	}
	return string(uid)
}

type NopWriter struct{}

func (w *NopWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func UintEquals(a, b []uint) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		var found bool
		for j := range b {
			if a[i] == b[j] {
				found = true
				break
			}
		}
		if found == false {
			return false
		}
	}

	return true
}

func StrRev(in string) string {
	r := []rune(in)
	var res []rune
	for i := len(r) - 1; i >= 0; i-- {
		res = append(res, r[i])
	}
	return string(res)
}
