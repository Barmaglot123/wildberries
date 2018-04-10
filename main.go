package main

import (
    "fmt"
    "bufio"
    "os"
    "net/http"
    "io/ioutil"
    "net/url"
    "strings"
    "log"
)

func main() {
    var t int = 0
    total := make(chan int)
    defer close(total)
    waiter := make(chan int)
    defer close(waiter)
    finished := make(chan bool)
    defer close(finished)
    const k = 5

    go readFromInput(total, finished, waiter, k)

    for {
        select {
        case c := <-total:
            t += c
        case <-finished:
            fmt.Println("Total: ", t)
            return
        }
    }
}

func readFromInput (total chan int, finished chan bool, waiter chan int, k int) {
    s := bufio.NewScanner(os.Stdin)
    routinesCount := 0

    for s.Scan() {
        for {
            if k <= routinesCount{
                w := <- waiter
                routinesCount -= w
            } else {
                break
            }
        }

        u := s.Text()
        routinesCount += 1
        _, err := url.ParseRequestURI(u)

        if err != nil {
            panic("Unvalid URL")
        }
        go countEntries(u, total, waiter)
    }

    if s.Err() != nil {
        panic(s.Err())
    }

    for {
        w := <- waiter
        routinesCount -= w

        if routinesCount == 0{
           finished <- true
           return
       }
    }
}

func countEntries(url string, total, waiter chan int) {
    b, err := sendGetRequest(url)

    if err != nil { log.Println(err) }

    count := strings.Count(string(b), "Go")
    fmt.Println("Count for", url, ":", count)

    total <- count
    waiter <- 1
}

func sendGetRequest(url string) ([]byte, error) {
    req, _ := http.NewRequest("GET", url, nil)
    client := &http.Client{}
    res, err := client.Do(req)
    defer res.Body.Close()
    if err != nil { return nil, err }
    body, _ := ioutil.ReadAll(res.Body)

    return body, err
}