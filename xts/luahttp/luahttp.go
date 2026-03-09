package luahttp

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	lua "github.com/speedata/go-lua"
)

const luaHTTPResponseTypeName = "http.response"

type httpResponse struct {
	res      *http.Response
	body     string
	bodySize int
}

var client = &http.Client{}

func httpGet(l *lua.State) int {
	return doRequestAndPush(l, "GET", 1, 2)
}

func httpDelete(l *lua.State) int {
	return doRequestAndPush(l, "DELETE", 1, 2)
}

func httpHead(l *lua.State) int {
	return doRequestAndPush(l, "HEAD", 1, 2)
}

func httpPatch(l *lua.State) int {
	return doRequestAndPush(l, "PATCH", 1, 2)
}

func httpPost(l *lua.State) int {
	return doRequestAndPush(l, "POST", 1, 2)
}

func httpPut(l *lua.State) int {
	return doRequestAndPush(l, "PUT", 1, 2)
}

func httpRequest(l *lua.State) int {
	method, _ := l.ToString(1)
	return doRequestAndPush(l, strings.ToUpper(method), 2, 3)
}

func doRequestAndPush(l *lua.State, method string, urlIdx, optIdx int) int {
	url, _ := l.ToString(urlIdx)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		l.PushNil()
		l.PushString(err.Error())
		return 2
	}

	// Parse options table if present
	if l.IsTable(optIdx) {
		// headers
		l.Field(optIdx, "headers")
		if l.IsTable(-1) {
			l.PushNil()
			for l.Next(-2) {
				key, _ := l.ToString(-2)
				value, _ := l.ToString(-1)
				req.Header.Set(key, value)
				l.Pop(1)
			}
		}
		l.Pop(1)

		// cookies
		l.Field(optIdx, "cookies")
		if l.IsTable(-1) {
			l.PushNil()
			for l.Next(-2) {
				key, _ := l.ToString(-2)
				value, _ := l.ToString(-1)
				req.AddCookie(&http.Cookie{Name: key, Value: value})
				l.Pop(1)
			}
		}
		l.Pop(1)

		// query
		l.Field(optIdx, "query")
		if l.IsString(-1) {
			q, _ := l.ToString(-1)
			req.URL.RawQuery = q
		}
		l.Pop(1)

		// body
		l.Field(optIdx, "body")
		if l.IsString(-1) {
			body, _ := l.ToString(-1)
			req.ContentLength = int64(len(body))
			req.Body = io.NopCloser(strings.NewReader(body))
		}
		l.Pop(1)

		// timeout
		l.Field(optIdx, "timeout")
		if !l.IsNoneOrNil(-1) {
			var duration time.Duration
			if n, ok := l.ToNumber(-1); ok {
				duration = time.Second * time.Duration(int(n))
			} else if s, ok := l.ToString(-1); ok {
				duration, err = time.ParseDuration(s)
				if err != nil {
					l.Pop(1)
					l.PushNil()
					l.PushString(err.Error())
					return 2
				}
			}
			if duration > 0 {
				ctx, cancel := context.WithTimeout(req.Context(), duration)
				defer cancel()
				req = req.WithContext(ctx)
			}
		}
		l.Pop(1)

		// auth
		l.Field(optIdx, "auth")
		if l.IsTable(-1) {
			l.Field(-1, "user")
			l.Field(-2, "pass")
			if l.IsString(-2) && l.IsString(-1) {
				user, _ := l.ToString(-2)
				pass, _ := l.ToString(-1)
				req.SetBasicAuth(user, pass)
			} else {
				l.Pop(3) // pass, user, auth
				l.PushNil()
				l.PushString("auth table must contain non-nil user and pass fields")
				return 2
			}
			l.Pop(2) // pass, user
		}
		l.Pop(1) // auth
	}

	res, err := client.Do(req)
	if err != nil {
		l.PushNil()
		l.PushString(err.Error())
		return 2
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		l.PushNil()
		l.PushString(err.Error())
		return 2
	}

	pushHTTPResponse(l, res, string(body), len(body))
	return 1
}

func pushHTTPResponse(l *lua.State, res *http.Response, body string, bodySize int) {
	if lua.NewMetaTable(l, luaHTTPResponseTypeName) {
		l.PushGoFunction(httpResponseIndex)
		l.SetField(-2, "__index")
	}
	l.Pop(1) // pop metatable

	l.PushUserData(&httpResponse{
		res:      res,
		body:     body,
		bodySize: bodySize,
	})
	lua.SetMetaTableNamed(l, luaHTTPResponseTypeName)
}

func checkHTTPResponse(l *lua.State) *httpResponse {
	ud := l.ToUserData(1)
	if v, ok := ud.(*httpResponse); ok {
		return v
	}
	lua.ArgumentError(l, 1, "http.response expected")
	return nil
}

func httpResponseIndex(l *lua.State) int {
	res := checkHTTPResponse(l)
	field := lua.CheckString(l, 2)

	switch field {
	case "headers":
		l.NewTable()
		for key := range res.res.Header {
			l.PushString(res.res.Header.Get(key))
			l.SetField(-2, key)
		}
		return 1
	case "cookies":
		l.NewTable()
		for _, cookie := range res.res.Cookies() {
			l.PushString(cookie.Value)
			l.SetField(-2, cookie.Name)
		}
		return 1
	case "status_code":
		l.PushInteger(res.res.StatusCode)
		return 1
	case "url":
		l.PushString(res.res.Request.URL.String())
		return 1
	case "body":
		l.PushString(res.body)
		return 1
	case "body_size":
		l.PushInteger(res.bodySize)
		return 1
	}
	return 0
}

// Open starts the http lua module
func Open(l *lua.State) int {
	lua.NewLibrary(l, []lua.RegistryFunction{
		{"get", httpGet},
		{"delete", httpDelete},
		{"head", httpHead},
		{"patch", httpPatch},
		{"post", httpPost},
		{"put", httpPut},
		{"request", httpRequest},
	})
	return 1
}
