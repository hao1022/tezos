package common

import (
    "fmt"
    "bytes"
    "io/ioutil"
    "net/http"
)

func Get(host string, route string) []byte {
    url := fmt.Sprintf("http://%s%s", host, route)
    resp, _ := http.Get(url)
    if resp == nil {
	fmt.Println("%s get nil response", url)
        return nil
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}

func Post(host string, route string, data string) []byte {
    data := []byte(data)
    url := fmt.Sprintf("http://%s%s", host, target)
    resp, _ := http.Post(url, "application/json", bytes.NewBuffer(data))
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    return body
}
