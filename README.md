GO Oembed
=========

[![GoDoc](https://godoc.org/github.com/dyatlov/go-oembed/oembed?status.svg)](https://godoc.org/github.com/dyatlov/go-oembed/oembed)

Go Oembed provides methods to retrieve oEmbed data from known providers.
The provider list can be fetched from this repository or from [oembed.com/providers.json](https://oembed.com/providers.json)


Install
-------

`go get github.com/dyatlov/go-oembed/oembed`


Example
-------

```go
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dyatlov/go-oembed/oembed"
)

func main() {
	data, err := ioutil.ReadFile("../providers.json")

	if err != nil {
		panic(err)
	}

	oe := oembed.NewOembed()
	oe.ParseProviders(bytes.NewReader(data))

	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter url: ")
		url, _ := reader.ReadString('\n')

		url = strings.Trim(url, "\r\n")

		if url == "" {
			break
		}

		item := oe.FindItem(url)

		if item != nil {
			info, err := item.FetchOembed(oembed.Options{URL: url})
			if err != nil {
				fmt.Printf("An error occured: %s\n", err.Error())
			} else {
				if info.Status >= 300 {
					fmt.Printf("Response status code is: %d\n", info.Status)
				} else {
					fmt.Printf("Oembed info:\n%s\n", info)
				}
			}
		} else {
			fmt.Println("nothing found :(")
		}

	}
}
```
