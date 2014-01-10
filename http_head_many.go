package main

import (
  "fmt"
  "os"
  "bufio"
  "net/http"
  "time"
)

type HttpResponse struct {
  url      string
  response *http.Response
  err      error
}

func asyncHttpGets(urls []string) []*HttpResponse {
  ch := make(chan *HttpResponse)
  responses := []*HttpResponse{}
  for _, url := range urls {
      go func(url string) {
          resp, err := http.Head(url)
          if err!=nil {
            fmt.Printf("Error when fetching %s: %s\n", url, err)
          }
          ch <- &HttpResponse{url, resp, err}
      }(url)
  }

  for {
      select {
      case r := <-ch:
          responses = append(responses, r)
          if len(responses) == len(urls) {
              return responses
          }
      case <-time.After(100 * time.Millisecond):
          fmt.Printf(".")
      }
  }
  return responses
}

func main() {
  var urls []string;

  if len(os.Args)<2 {
     fmt.Printf("usage: %s urls.txt #one URL on the line\n", os.Args[0])
     return
  }
  file, err := os.Open(os.Args[1])
  if err != nil {
    fmt.Printf("Error accessing the file with URLs: %s", err)
    return
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    urls = append(urls, scanner.Text())
  }
  if len(urls) < 1 {
    fmt.Printf("Error reading URLs: %s", scanner.Err())
    return
  }

  results := asyncHttpGets(urls)

  ofile, err := os.Create(os.Args[1] + ".csv")
  if err!=nil {
    fmt.Printf("Error opening file for results: %s", err)
    return
  }

  w := bufio.NewWriter(ofile)
  for _, result := range results {
    if result.err == nil {
      cl := result.response.Header["Content-Length"]
      if len(cl) > 0 {
        fmt.Fprintf(w, "%s, %s\n", result.url, cl[0])
      }
    }
  }
  w.Flush()
  ofile.Close()
  fmt.Println("done");
}

