package main

import (
	"fmt"
	"net/http"
	"time"
	"math/rand"
)

func getCookieValueByName(cookies []*http.Cookie, name string) string {
	result := ""

	for _, cookie := range cookies {
		if cookie.Name == name {
			result = cookie.Value
			break
		}
	}
	if result == "" {
		panic("Smtng wrong with cookie")
	}
	return result
}

func percents(on int, off int, all int) string {
	sum := on + off

	percentOn := int(float32(on) / float32(sum) * 100.00)
	percentOff := int(100 - percentOn)

	return fmt.Sprintf(`
Total:             %d flats
Already connected: %d flats - %d percent
Not connected    : %d flats - %d percent
`,
		sum, on, percentOn, off, percentOff)
}

func waitLtl() {
	// задержка
	rand.Seed(time.Now().UnixNano())

	n := rand.Intn(1000)
	n = n + 500
	time.Sleep(time.Duration(n) * time.Millisecond)
}
