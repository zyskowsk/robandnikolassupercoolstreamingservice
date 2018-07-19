package main

import {
    "fmt"
    "net/http"
}


func main() {
    lc := make(chan string, 10)
    exit := make(chan bool)

    go func () {
        for {
            var line string
            fmt.Scanf("%s", &line)
            lc <- line

            if line == "exit" {
                exit <- true
            }
        }
    }()

    go func() {
        for {
            fmt.Printf(<-lc)
        }
    }()

    <-exit
}
