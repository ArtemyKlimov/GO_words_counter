package main

import (
    "fmt"
    "bufio"
    "net/http"
    "os"
    "sync"
    "runtime"
    "strings"
    "io/ioutil"
    "log"
    //_ "net/http/pprof"
)

func goCounter(url string) (int, error) {
    resp, err := http.Get(url)
    if err != nil {
        return 0, fmt.Errorf("could not get %s: %v", url, err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        if resp.StatusCode == http.StatusTooManyRequests {
            return 0, fmt.Errorf("You are being rate limited:")
        }
        return 0, fmt.Errorf("bad response from server: %s", resp.Status)
    }
    html, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return 0, fmt.Errorf("can not read html from %s: %v", url, err)
    }
    result := strings.Count(string(html), "Go")
    fmt.Printf("Count for %s: %d\n", url, result)
    return result, nil
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    scanner := bufio.NewScanner(os.Stdin)
    tasks := make(chan string, 10)
    go func() {
        for scanner.Scan() {
            if scanner.Text() == "" {
                break
            }
            tasks <- scanner.Text()
        }
        if err := scanner.Err(); err != nil {
            fmt.Fprintln(os.Stderr, "reading standard input:", err)
        }
        close(tasks)
    }()

    results := make(chan int, 10)
    var wg sync.WaitGroup
    wg.Add(5)
    go func() {
        wg.Wait()
        close(results)
    }()

    for i := 0; i < 5; i++ {
        go func() {
            defer wg.Done()
            for t := range tasks {
                res, err := goCounter(t)
                if err != nil {
                    log.Printf("error ocured: %v, %v", t, err)
                    continue
                }
                results <-res
            }
        }()
    }
    total := 0
    for r := range results {
        total += r
    }
    fmt.Println("Total:",total)
}
