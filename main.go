package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/insprac/arachne/crawler"
)

var outputDir = ""

func main() {
	startUrlsArg := flag.String("start", "", "Comma-separated list of URLs to start crawling")
	allowPrefixesArg := flag.String("allow-prefix", "", "Comma-separated list of allowed URL prefixes")
	out := flag.String("out", "", "Output directory")

	flag.Parse()

	if *startUrlsArg == "" {
		log.Fatalf("--start parameter must be provided")
	}

	if *allowPrefixesArg == "" {
		log.Fatalf("--allow-prefix parameter must be provided")
	}

	if *out == "" {
		log.Fatalf("--out parameter must be provided")
	}

	startUrls := strings.Split(*startUrlsArg, ",")
	allowPrefixes := strings.Split(*allowPrefixesArg, ",")
	outputDir = *out

	err := crawler.Crawl(startUrls, allowPrefixes, saveDocument)
	if err != nil {
		log.Fatalf("Error during crawling: %v\n", err)
	}
}

func saveDocument(doc crawler.Document) error {
	parsedUrl, err := url.Parse(doc.Uri)
	if err != nil {
		return fmt.Errorf("Error parsing URL: %v", err)
	}

	dirPath := filepath.Join(outputDir, parsedUrl.Host, filepath.Dir(parsedUrl.Path))

	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("Error creating directory: %v", err)
	}

	filePath := filepath.Join(dirPath, filepath.Base(parsedUrl.Path)+".md")
	err = ioutil.WriteFile(filePath, []byte(doc.Body), 0644)
	if err != nil {
		return fmt.Errorf("An error occurred: %v", err)
	}

	return nil
}
