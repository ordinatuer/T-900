package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "encoding/csv"
  "os"
  "strconv"
  "time"
)

func main() {
	// auth()
	// return

  house := 23593 // ID Дома
  cookieValue := "0fa95ee88f919a2026a31a2a7111b094" // строка Cookie из браузера

  // перебор квартир от roomMax до roomMin
  roomMax := 25
  roomMin := 1

  cookieTmpl := "ToddVoight=%s"
  url := "https://t900.skynet.ru/genisys/t800/orders/check?bl_home_ID=%d&location_flat=%d&check_flat=no"
  method := "GET"
  skynetStatus := ""

  client := &http.Client {
  }

  fileName := fmt.Sprintf("%d.csv", house)
  file, ferr := os.Create(fileName)
  if ferr != nil {
  	fmt.Println("Error file read", ferr)
  	return
  }

  w := csv.NewWriter(file)
  defer w.Flush()

  _url := ""
  cookie := fmt.Sprintf(cookieTmpl, cookieValue)
  _t := time.Now()

  for i:=roomMax;roomMin<=i;i-- {
  	  _url = fmt.Sprintf(url, house, i)

	  req, err := http.NewRequest(method, _url, nil)

	  if err != nil {
	    fmt.Println(err)
	    return
	  }
	  req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"107\", \"Chromium\";v=\"107\", \"Not=A?Brand\";v=\"24\"")
	  req.Header.Add("sknt-project-name", "t900.sknt.ru")
	  req.Header.Add("x-environment", "javascript")
	  req.Header.Add("sec-ch-ua-mobile", "?0")
	  req.Header.Add("sknt-client-req", "UmoVwG5gOSO0T3DtGrQ3h")
	  req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	  req.Header.Add("accept", "version=0.1")
	  req.Header.Add("x-requested-with", "XMLHttpRequest")
	  req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	  req.Header.Add("Cookie", cookie)

	  res, err := client.Do(req)
	  if err != nil {
	    fmt.Println(err)
	    return
	  }
	  defer res.Body.Close()

	  body, err := ioutil.ReadAll(res.Body)
	  if err != nil {
	    fmt.Println(err)
	    return
	  }

	  if 401 == res.StatusCode {
	  	fmt.Println("Auth error. Maybe smtg wrong with cookie")
	  	return
	  }
	  if 200 != res.StatusCode {
	  	fmt.Println("Request error with status %d")
	  	return
	  }

	  // fmt.Println(string(body))

	  var dat map[string]interface{}
	  er := json.Unmarshal(body, &dat)
	  if er != nil {
	  	fmt.Println(err)
	  	return
	  }

	  if dat["result"] == "error" {
	  	skynetStatus = "Skynet"
	  } else if dat["result"] == "ok" {
	  	skynetStatus = ""
	  } else {
	  	skynetStatus = " -- "
	  }
	  row := []string{strconv.Itoa(i), skynetStatus}

	  csverr := w.Write(row)
	  if csverr != nil {
	  	fmt.Println("Error file writing", csverr)
	  	return
	  }

	  fmt.Println(res.StatusCode, dat["result"])
  }

  fmt.Println(time.Since(_t))
}

// func auth() {
// 	fmt.Println("Auth")
// }