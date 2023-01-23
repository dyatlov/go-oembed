package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/dyatlov/go-oembed/oembed"
)

func main() {
	ctx := context.Background()

	data, err := ioutil.ReadFile("../providers.json")

	if err != nil {
		panic(err)
	}

	oe := oembed.NewOembed()
	oe.ParseProviders(bytes.NewReader(data))

	extras := url.Values{"autoplay": []string{"1"}}

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
			info, err := item.FetchOembedCtx(ctx, oembed.Options{URL: url, AcceptLanguage: "en-us", ExtraOpts: extras})
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
