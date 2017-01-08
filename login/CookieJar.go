package login

import (
	"bytes"
	"fmt"
	"net/http"
)

type CookieJar struct {
	Cookies []*http.Cookie
}

func (jar *CookieJar) Update(cookies []*http.Cookie) {
	for _, newCookie := range cookies {
		found := false
		for i, oldCookie := range jar.Cookies {
			if newCookie.Name == oldCookie.Name {
				found = true
				jar.Cookies[i] = newCookie
				break
			}
		}
		if !found {
			jar.Cookies = append(jar.Cookies, newCookie)
		}
	}
}

func (jar *CookieJar) Encode() string {
	var buf bytes.Buffer
	for _, cookie := range jar.Cookies {
		buf.WriteString(fmt.Sprintf("%s=%s; ", cookie.Name, cookie.Value))
	}
	return buf.String()
}
func (jar *CookieJar) EncodeOne(x int) string {
	var buf bytes.Buffer
	if jar.Cookies[x] != nil {
		buf.WriteString(fmt.Sprintf("%s=%s; ", jar.Cookies[x].Name, jar.Cookies[x].Value))
	} else {
		buf.WriteString(fmt.Sprintf("%s=%s", "Cookie", ""))
	}

	return buf.String()
}
