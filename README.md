# GameDB

GameDB is a powerful command-line tool designed to scrape and manage repack game data from various online sources. With support for multiple data sources and the ability to provide a RESTful API.

## Features

- **Data Sources**:
  - DODI
  - Fitgirl
  - FreeGOG
  - KaOSKrew
  - OnlineFix
  - Xatab

- **Database**:
  - Stores game data in MongoDB
  - Supports Redis for caching to improve performance

- **RESTful API**:
  - Provides an API for external access to the game data

- **Command-Line Interface**:
  - Comprehensive set of commands for managing game data

## Usage

GameDB provides a variety of commands to manage game data. Below are some examples:

- **Crawl Data**:
    ```sh
    gamedb crawl -p <platform> -a
    ```
    This command will scrape game data from the specified platform and add it to the database.

- **Start Server**:
    ```sh
    gamedb server -a <addr>
    ```
    This command will start the RESTful API server at the specified address.

Read `internal/cmd` for more details.

## Configuration

Edit the `config.json` file to set up your environment:

- **MongoDB**: Configure connection details for the MongoDB database.
- **Redis**: Optionally configure Redis for caching.
- **Other Settings**: Adjust other settings as needed for your deployment.

Read `internal/config/config.go` for more details.

## Example Workflow

1. **Clone and configure the project**:
    ```sh
    git clone https://github.com/nitezs/GameDB.git
    cd GameDB
    mv config.example.json config.json
    nano config.json  # Edit the configuration file
    ```

2. **Run the crawler**:
    ```sh
    go run . crawl -p fitgirl -a
    ```

3. **Start the server**:
    ```sh
    go run . server -a :8080
    ```

4. **Access the API**:
    You can now access the API at `http://localhost:8080`.

## Routes

- GET /raw/:id - Get raw game data
- GET /game/search - Search for game infos 
- GET /game/:id - Get game info
- GET /game/name/:name - Get game info by name
- GET /ranking/:type - Get game ranking, type can be top, week-top, best-of-the-year, most-played

## License

This project is licensed under the GNU General Public License v3.0 License.