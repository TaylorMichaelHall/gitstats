# gitstats

gitstats is a command-line tool for analyzing Git activity in a locally cloned repository.

## Installation

1. Ensure you have Go installed on your system.

2. Clone the repository:
   ```shell
   git clone https://github.com/taylormichaelhall/gitstats.git
   ```
3. Navigate to the project directory:

   ```shell
   cd gitstats
   ```

4. Install dependencies:

   ```shell
   go mod tidy
   ```

5. Build the project:
   ```shell
   go build -o gitstats
   ```

Place the binary somewhere in your PATH or run it from the project directory.

## Usage

### Contributor Analysis

```shell
gitstats contributors /path/to/repo
```

This command displays a bar graph of contributor statistics and allows you to view detailed information for individual contributors.

### File Change Frequency

```shell
gitstats files /path/to/repo
```

This command shows a chart of the most frequently changed files in the repository.

To ignore specific files or directories:

```shell
gitstats files /path/to/repo -i package.json -i build/
```

### Help

For more information on available commands and options:

```
gitstats help
```
