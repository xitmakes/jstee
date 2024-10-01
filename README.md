
```markdown
# jstee

`jstee` is a Go tool designed for bug bounty hunters and security researchers. It scrapes JavaScript file links from a list of hosts and saves them to a file.

## Features
- Extracts JavaScript file URLs (`<script>` tags) from web pages.
- Supports relative and absolute URLs.
- Saves the list of JavaScript files into a `js.txt` file for further analysis.

## Installation

You can install `jstee` using the following command:

```bash
go install -v github.com/xitmakes/cmd/jstee@latest
```

Replace `your-username` with your actual GitHub username.

## Usage

1. Prepare a `hosts.txt` file with a list of hosts (one per line). For example:

   ```txt
   https://example.com
   https://another-site.org
   ```

2. Run the tool using the following command:

   ```bash
   jstee -f hosts.txt
   ```

3. The tool will scrape the JavaScript file URLs and save them in a file called `js.txt`.

### Example

If `hosts.txt` contains:

```txt
https://example.com
https://test.com
```

Running the command:

```bash
jstee -f hosts.txt
```

Will produce a `js.txt` file containing the scraped JavaScript URLs:

```txt
https://example.com/assets/js/main.js
https://test.com/scripts/app.js
```

## Contributing

Feel free to submit issues or pull requests if you'd like to contribute to this project.

## License

This project is licensed under the MIT License.
```

