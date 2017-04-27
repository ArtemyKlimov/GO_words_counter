package main

import (
  "fmt"
  "net/http"
  "bufio"
  "os"
  "regexp"
  "io/ioutil"
  "time"
)

func worker(id int, jobs<-chan string, results chan<-int) {
  t0 := time.Now()
  for url := range jobs {
    resp, err := http.Get(url)
    if err != nil {
      fmt.Println("problem while opening url", url)
      results<-0
      continue  //if can not get http request - skip and go furher
    }
    defer resp.Body.Close()
    html, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      continue
    }
    regExp:= regexp.MustCompile("Go")
    matches := regExp.FindAllStringIndex(string(html), -1)
    t1 := time.Now()
    fmt.Println("Count for", url, ":", len(matches), "Elapsed time:", t1.Sub(t0),  "works id", id)
    results<-len(matches)
  }
}

func main(){
  scanner := bufio.NewScanner(os.Stdin)
  jobs := make(chan string, 100) //
  results := make(chan int, 100)
  t0 := time.Now()
  for w:= 0; w<5; w++{
    go worker(w, jobs, results)
  }
  var task int = 0
  res := 0
  for scanner.Scan() {
      jobs <- scanner.Text()
      task ++
  }
  close(jobs)
  for a := 1; a <= task; a++ {
    res+=<-results
  }
  close(results)
  t2 := time.Now()
  fmt.Println("Total:",res, "Elapsed total time:", t2.Sub(t0) );
}

//echo -e 'https://golang.org/help\nhttps://golang.org/project\nhttps://golang.org/pkg\nhttps://golang.org/dl\nhttps://gobyexample.ru\nhttp://golang-book.ru\nhttps://golang.org\nhttps://godoc.org\nlll\nhttps://gobyexample.com\nhttps://golang.org/pkg' | go run r4.go
