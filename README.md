# Web Scraping Project
This project is a web scraping application built using the Go programming language and the Gin framework. It leverages the [https://github.com/gocolly/colly](https://github.com/gocolly/colly) package to scrape web content.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

## Installation

1. Clone the repository:
    
    ```shell
    git clone git@github.com:jobayer12/ScrapifyGo.git
    ```

2. Change to the project directory:

   ```shell
    cd ScrapifyGo
   ```

3. Install the required dependencies:
   ```shell
    go mod tidy
   ```

4. Run the application:
    ```shell
    go run cmd/api/main.go
    ```
   or
    ```shell
    make run
   ```

## Usage

Once the application is running, you can access the API endpoints to perform various scraping tasks. The base URL for the API is http://localhost:8080.

## API Endpoints

### 1. Scrape Sitemap Data

- *Endpoint:* /api/v1/sitemap?url=http://example.com/sitemap.xml
- *Method:* GET
- *Description:* Scrapes the sitemap data from a given URL.

- *Response:*
    ```json
    {
    "data": [
        {
            "changefreq": "string",
            "lastmod": "string",
            "loc": "string",
            "priority": "string"
        }
    ],
    "error": "string",
    "status": 0
   }
   ```


### 2. Scrape URL from Website

- *Endpoint:* /api/v1/url?url=http://example.com/sitemap.xml
- *Method:* GET
- *Description:* Scrapes URLs from a given web page.

- *Response:*
  ```json
  {
    "data": [
        "string"
    ],
    "error": "string",
    "status": 0
  }
  ```


### 3. Scrape Emails from URLs

- *Endpoint:* /api/v1/email?url=http://example.com/sitemap.xml
- *Method:* GET
- *Description:* Scrapes email addresses from a list of URLs.

- *Response:*
  ```json
  {
    "data": [
        "string"
    ],
    "error": "string",
    "status": 0
  }
  ```

## Contributing
Contributions are welcome! Please submit a pull request or open an issue to discuss any changes.