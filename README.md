# gitstats

gitstats is a command-line tool for analyzing Git activity in a locally cloned repository.

## Installation

1. Ensure you have Go installed on your system.

2. Clone the repository:
   ```
   git clone https://github.com/taylormichaelhall/gitstats.git
   ```
3. Navigate to the project directory:

   ```
   cd gitstats
   ```

4. Install dependencies:

   ```
   go mod tidy
   ```

5. Build the project:
   ```
   go install
   ```

## Usage

### Contributor Analysis

```
./gitstats contributors /path/to/repo
```

This command displays a bar graph of contributor statistics and allows you to view detailed information for individual contributors.

### File Change Frequency

```
./gitstats files /path/to/repo
```

This command shows a chart of the most frequently changed files in the repository.

To ignore specific files or directories:

```
./gitstats files /path/to/repo -i package.json -i build/
```

### Help

For more information on available commands and options:

```
./gitstats help
```
