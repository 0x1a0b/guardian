package helpers

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

//IsMatch checks for matching string
func IsMatch(pattern string, str string) (bool, error) {
	return regexp.MatchString(pattern, str)
}

//UnEscapeRawValue unescapes raw query
func UnEscapeRawValue(rawQuery string) string {
	rawQuery = strings.Replace(rawQuery, "%%", "%25%", -1)
	rawQuery = strings.Replace(rawQuery, "%'", "%25'", -1)
	rawQuery = strings.Replace(rawQuery, `%"`, `%25"`, -1)
	re := regexp.MustCompile(`%$`)
	rawQuery = re.ReplaceAllString(rawQuery, `%25`)
	decodeQuery, _ := url.QueryUnescape(rawQuery)

	decodeQuery = PreProcessString(decodeQuery)
	//fmt.Println("UnEscapeRawValue decodeQuery", decodeQuery)
	return decodeQuery
}

// PreProcessString ...
func PreProcessString(value string) string {
	value2 := strings.Replace(value, `'`, ``, -1)
	value2 = strings.Replace(value2, `"`, ``, -1)
	value2 = strings.Replace(value2, `+`, ` `, -1)
	value2 = strings.Replace(value2, `/**/`, ` `, -1)
	return value2
}

//HeadersToString ...
func HeadersToString(header http.Header) (res string) {
	for name, values := range header {
		for _, value := range values {
			res += fmt.Sprintf("%s: %s ", name, value)
		}
	}
	return
}

//GetHeadersNames Gets the header name
func GetHeadersNames(header http.Header) (res []string) {
	for name := range header {
		res = append(res, name)
	}
	return
}

//CookiesToString ...
func CookiesToString(cookie []*http.Cookie) (res string) {
	for _, values := range cookie {
		res += fmt.Sprintf("%s=%s ", values.Name, values.Value)
	}
	return
}

//GetCookiesNames ...
func GetCookiesNames(cookie []*http.Cookie) (res []string) {
	for _, values := range cookie {
		res = append(res, values.Name)
	}
	return
}

//CalcTime ...
func CalcTime(start time.Time, end time.Time) int64 {
	return end.Sub(start).Nanoseconds() / int64(time.Millisecond)
}

//CalcTimeNow ...
func CalcTimeNow(end time.Time) int64 {
	return time.Since(end).Nanoseconds() / int64(time.Millisecond)
}

//StringContains searches given string in a string slice
func StringContains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
