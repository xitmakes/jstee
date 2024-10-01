package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

func main() {
	// Define the command-line flag for the input file
	inputFile := flag.String("f", "hosts.txt", "File containing the list of hosts to scrape JS files from")
	flag.Parse()

	// Open the file containing the list of hosts
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Println("Error opening", *inputFile, ":", err)
		return
	}
	defer file.Close()

	// Use a scanner to read hosts from the file line by line
	var hostList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			hostList = append(hostList, line)
		}
	}

	// Check for scanner error
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading hosts:", err)
		return
	}

	// Regular expression to find <script src="..."> tags
	re := regexp.MustCompile(`<script[^>]+src=["']([^"']+)["']`)

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create or truncate js.txt file to save JS file URLs
	outputFile, err := os.Create("js.txt")
	if err != nil {
		fmt.Println("Error creating js.txt:", err)
		return
	}
	defer outputFile.Close()

	// Loop through hosts and fetch the HTML to extract JS file links
	for _, host := range hostList {
		wg.Add(1)
		go func(h string) {
			defer wg.Done()

			// Ensure the URL starts with http/https
			if !strings.HasPrefix(h, "http://") && !strings.HasPrefix(h, "https://") {
				h = "http://" + h
			}

			// Parse the base URL for handling relative URLs
			baseURL, err := url.Parse(h)
			if err != nil {
				fmt.Println("Error parsing URL:", h, "-", err)
				return
			}

			// Make HTTP GET request to fetch HTML
			resp, err := http.Get(h)
			if err != nil {
				fmt.Println("Error fetching:", h, "-", err)
				return
			}
			defer resp.Body.Close()

			// Check if the response is successful
			if resp.StatusCode != http.StatusOK {
				fmt.Println("Failed to fetch:", h, "Status code:", resp.StatusCode)
				return
			}

			// Read the HTML content
			bodyScanner := bufio.NewScanner(resp.Body)
			var bodyContent string
			for bodyScanner.Scan() {
				bodyContent += bodyScanner.Text()
			}

			// Extract JS file URLs using the regex
			matches := re.FindAllStringSubmatch(bodyContent, -1)

			// Write each found JS link to the output file
			for _, match := range matches {
				if len(match) > 1 {
					jsLink := match[1]

					// Handle relative URLs by converting them to absolute
					parsedJSLink, err := url.Parse(jsLink)
					if err != nil {
						fmt.Println("Error parsing JS link:", jsLink, "-", err)
						continue
					}

					if !parsedJSLink.IsAbs() {
						// Resolve relative URL to an absolute URL using the base URL
						jsLink = baseURL.ResolveReference(parsedJSLink).String()
					}

					// Write the JS link to the file
					if _, err := outputFile.WriteString(jsLink + "\n"); err != nil {
						fmt.Println("Error writing to js.txt:", err)
					} else {
						fmt.Printf("Saved JS link from %s: %s\n", h, jsLink)
					}
				}
			}
		}(host)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("Finished scraping JavaScript file links with jstree.")
}


