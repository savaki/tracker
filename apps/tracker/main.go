package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Ref struct {
	Value     string    `json:"value"`
	Cookie    string    `json:"cookie"`
	Ip        string    `json:"ip"`
	Timestamp time.Time `json:"timestamp"`
}

func CookieKey(key string) string {
	return "ref_" + key
}

func NewCookie(key, value string) http.Cookie {
	outkey := CookieKey(key)
	expires := time.Now().AddDate(5, 0, 0)
	domain := ".loyal3.com"
	maxAge := 86400 * 365 * 5
	// raw := fmt.Sprintf("%s=%s; domain=%s; maxAge=%d;", outkey, value, domain, maxAge)
	raw := fmt.Sprintf("%s=%s;", outkey, value)

	cookie := http.Cookie{
		Name:       outkey,
		Value:      value,
		Path:       "/",
		Domain:     domain,
		Expires:    expires,
		RawExpires: expires.Format(time.UnixDate),
		MaxAge:     maxAge,
		Secure:     false,
		HttpOnly:   true,
		Raw:        raw,
		Unparsed:   []string{raw},
	}

	return cookie
}

func main() {
	image, err := ioutil.ReadFile("cnbc.png")
	if err != nil {
		panic(err)
	}

	contentLength := strconv.Itoa(len(image))

	router := mux.NewRouter()
	router.HandleFunc("/pixel.png", func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		refValue := req.FormValue("ref")

		if c, _ := req.Cookie(CookieKey(refValue)); c == nil {
			ref := &Ref{
				Value:     refValue,
				Cookie:    "localhost:8080",
				Ip:        req.RemoteAddr,
				Timestamp: time.Now(),
			}

			cookie := NewCookie(ref.Value, ref.Cookie)
			w.Header().Set("Set-Cookie", cookie.Raw)

			data, _ := json.Marshal(ref)
			fmt.Println(string(data))
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", contentLength)
		w.WriteHeader(200)
		w.Write(image)
	})

	http.ListenAndServe(":8080", router)
}
