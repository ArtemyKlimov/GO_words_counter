package main

import (
    "fmt"
    "bufio"
    "net/http"
    "os"
    "sync"
    "runtime"
    "io/ioutil"
    "log"
    "regexp"
    //"time"
)

func goCounter(url string) (int, error) {
    resp, err := http.Get(url)
    if err != nil {
        //fmt.Println("Some problem occured while opening URL: ", url)
        return -1, fmt.Errorf("could not get %s: %v", url, err)
    }
    defer resp.Body.Close()
    html, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        //fmt.Println("Some problem occured while reading URL: ", url)
        return -1, fmt.Errorf("can not read html from %s: %v", url, err)
    }
    regExp:= regexp.MustCompile("Go")
    matches := regExp.FindAllStringIndex(string(html), -1)
    fmt.Println("Count for", url, ":", len(matches))
    return len(matches), nil
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    scanner := bufio.NewScanner(os.Stdin)
    tasks := make(chan string)
    index := 1
    go func() {
        for scanner.Scan() {
            tasks <- scanner.Text()
            index++
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
    //time.Sleep(time.Millisecond * 1000)
    fmt.Println("Total:",total)
}
