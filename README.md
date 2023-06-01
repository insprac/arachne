# Arachne

Arachne is a web crawler implemented in Go, capable of crawling multiple URLs
and saving the contents as Markdown files. This tool is perfect for those
looking to build a local repository of web pages for offline reading or for
performing subsequent data analysis.

## Features

- Scrape any webpage and save the content as Markdown files.
- It allows you to specify starting URLs and the allowed URL prefixes for
crawling.
- Utilizes concurrency for faster scraping.
- Automatically handles the conversion of HTML content to Markdown.
- Respects and parses relative and absolute URLs.
- Removes unnecessary elements from the webpage like scripts, styles, headers,
footers, navigation, and sidebars to focus on the main content.

## Usage

After cloning this repository and navigating into the project directory, you
can run the program using:

```bash
go run main.go \
    --start "https://www.example.com,https://www.another-example.com" \
    --allow-prefix "https://www.example.com,https://www.another-example.com/some/path" \
    --out "/path/to/output"
```

The program accepts the following command-line arguments:

- `--start`: A comma-separated list of URLs from which to start crawling.
- `--allow-prefix`: A comma-separated list of URL prefixes which the crawler is allowed to visit.
- `--out`: The directory where the Markdown files will be saved.

## Note

Please use this tool responsibly, respecting the `robots.txt` files and not
overloading servers. Be mindful of your network usage as well, as web scraping
can consume a lot of bandwidth depending on the size of the websites you are
scraping.

## Contribution

This is an open-source project. Feel free to fork and make any contributions
you feel will enhance the functionality of this web scraper.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
# arachne
