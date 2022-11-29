package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const house = 17783 // ID Дома
// перебор квартир от roomMax до roomMin
const roomMax = 25
const roomMin = 15

const login = "bednyh.d"  // логин
const password = "SqguXZ" // пароль

const cookieName = "ToddVoight"

func main() {
	if roomMax <= roomMin {
		fmt.Println("Wrong flats range")
		return
	}

	cookieTmpl := "%s=%s"
	url := "https://t900.skynet.ru/genisys/t800/orders/check?bl_home_ID=%d&location_flat=%d&check_flat=no"
	method := "GET"
	skynetStatus := ""

	client := &http.Client{}
	// получить значение куки через страницу входа
	cookieValue := auth()

	fileName := fmt.Sprintf("%d.csv", house)
	file, ferr := os.Create(fileName)
	if ferr != nil {
		fmt.Println("Error file read", ferr)
		return
	}

	w := csv.NewWriter(file)
	defer w.Flush()

	cookie := fmt.Sprintf(cookieTmpl, cookieName, cookieValue)
	skynetOn := 0  // счётчик подключенных
	skynetOff := 0 // счётчик НЕподключенных
	skynetAll := 0 // счётчик всех из диапазона
	_t := time.Now()

	// перебор квартир из указанного  диапазона
	for i := roomMax; roomMin <= i; i-- {
		_url := fmt.Sprintf(url, house, i)

		req, err := http.NewRequest(method, _url, nil)

		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\"")
		req.Header.Add("sknt-project-name", "t900.sknt.ru")
		req.Header.Add("x-environment", "javascript")
		req.Header.Add("sec-ch-ua-mobile", "?0")
		req.Header.Add("sknt-client-req", "UmoVwG5gOSO0T3DtGrQ3h") // ??
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
		req.Header.Add("accept", "version=0.1")
		req.Header.Add("x-requested-with", "XMLHttpRequest")
		req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
		req.Header.Add("Cookie", cookie)

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if res.StatusCode == 401 {
			fmt.Println("Auth error. Maybe smtg wrong with cookie")
			return
		}
		if res.StatusCode != 200 {
			fmt.Println("Request error with status %d")
			return
		}

		var dat map[string]interface{}
		er := json.Unmarshal(body, &dat)
		if er != nil {
			fmt.Println(err)
			continue
		}

		if dat["result"] == "error" {
			skynetStatus = "Skynet"
			skynetOn++
		} else if dat["result"] == "ok" {
			skynetStatus = ""
			skynetOff++
		} else {
			skynetStatus = " -- "
		}
		row := []string{strconv.Itoa(i), skynetStatus}

		csverr := w.Write(row)
		if csverr != nil {
			fmt.Println("Error file writing", csverr)
			return
		}

		fmt.Println(i, dat["result"])

		skynetAll = i
	}

	fmt.Println(time.Since(_t))
	fmt.Printf("Save in file %s", fileName)
	fmt.Println(percents(skynetOn, skynetOff, skynetAll))
}

func auth() string {
	url := "https://t900.skynet.ru/genisys/auth/"
	method := "POST"
	authRequestData := `{"login":"%s","password":"%s"}`

	payload := strings.NewReader(fmt.Sprintf(authRequestData, login, password))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		panic(err)
	}
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\"")
	req.Header.Add("sknt-project-name", "t900.sknt.ru")
	req.Header.Add("x-environment", "javascript")
	req.Header.Add("sknt-client-req", "cAuC5YcMs5GduxaZhGiMo") // ??
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "version=0.1")
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("Cookie", "metrika_enabled=1")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Auth failed")
	}
	return getCookieValueByName(res.Cookies(), "ToddVoight")
}
