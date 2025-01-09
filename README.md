# News Aggregator

News Aggregator is a Go-based application designed to collect and display news articles from various sources, along with a Telegram bot integration for delivering news directly to users.

## Features

- Fetches news articles from multiple sources defined in a JSON configuration file.
- Serves the aggregated news through a web interface.
- Telegram bot for delivering news updates directly to users.
- Dockerized for easy deployment.

## Prerequisites

- [Go](https://golang.org/doc/install) 1.16 or higher
- [Docker](https://docs.docker.com/get-docker/) (optional, for containerized deployment)
- A Telegram bot token (to enable Telegram bot functionality).

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/kolosiv/news-aggr.git
   cd news-aggr
   ```

2. **Build the application:**

   ```bash
   go build -o news-aggr
   ```

3. **Run the application:**

   ```bash
   ./news-aggr
   ```

   The application will start serving at `http://localhost:8080`.

## Telegram Bot Setup

To use the Telegram bot functionality:

1. Create a Telegram bot using [BotFather](https://core.telegram.org/bots#botfather) and obtain the bot token.
2. Add the bot token to your environment variables or configuration file:
   ```bash
   export TELEGRAM_BOT_TOKEN="your-telegram-bot-token"
   ```
3. Run the application, and the bot will be ready to interact with users.

   Users can start the bot and receive news updates directly within Telegram by subscribing to specific categories or keywords.

## Configuration

The news sources are defined in the `sources.json` file. Each entry should include the source's name and URL. Modify this file to add or remove news sources as needed.

## Docker Deployment

1. **Build the Docker image:**

   ```bash
   docker build -t news-aggr .
   ```

2. **Run the Docker container:**

   ```bash
   docker run -p 8080:8080 -e TELEGRAM_BOT_TOKEN="your-telegram-bot-token" news-aggr
   ```

   The application will be accessible at `http://localhost:8080`, and the Telegram bot will be operational.

## License

This project is licensed under the MIT License.
