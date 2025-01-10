# dnsgen - DNS Name Permutation Generator

`dnsgen` is a command-line tool that generates permutations of domain names. It's useful for security testing, reconnaissance, and discovering variations of a given domain. This tool is a Go implementation inspired by the original Python project [dnsgen](https://github.com/AlephNullSK/dnsgen).

## Features

*   Generates permutations of domain names based on common patterns:
    *   Inserting words at different positions
    *   Modifying numbers (incrementing/decrementing)
    *   Adding environment prefixes (e.g., `dev`, `staging`)
    *   Adding cloud provider terms (e.g., `api-aws`, `cdn-azure`)
    *   Adding region prefixes (e.g., `us-east`, `eu-west`)
    *   Adding microservice patterns (e.g., `auth-service`, `payment-api`)
    *   Adding common internal tooling prefixes (e.g., `jenkins`, `gitlab`)
    *   Adding common port numbers (e.g., `8080`, `8443`)
*   Supports custom wordlists to tailor the permutations.
*   Fast generation mode for a reduced set of common permutations.
*   Handles input from files or standard input (pipes).

## Installation

### Prerequisites

*   Go (version 1.20 or later)

### Go Install

You can directly install the `dnsgen` tool using the `go install` command:

```bash
go install github.com/behf/dnsgen@latest
```

### Build from Source

1. Clone the repository:

    ```bash
    git clone https://github.com/behf/dnsgen.git
    cd dnsgen
    ```
2. Build the executable:

    ```bash
    go build
    ```

    This will create an executable named `dnsgen` in the current directory.

## Usage

```bash
./dnsgen [OPTIONS] [input_file]
```

### Options

*   `-l`, `--wordlen`: Minimum length of custom words extracted from domains (default: 6).
*   `-w`, `--wordlist`: Path to a custom wordlist file.
*   `-f`, `--fast`: Use fast generation mode (fewer permutations).
*   `-o`, `--output`: Output file path. If not specified, results are printed to standard output.
*   `-v`, `--verbose`: Enable verbose logging (debug level).

### Input

`dnsgen` can read input domains from either:

1. **A file:** Provide the filename as the `input_file` argument. Each line in the file should contain a single domain name.

    ```bash
    ./dnsgen input_domains.txt
    ```
2. **Standard input (pipe):** Pipe the output of another command to `dnsgen`.

    ```bash
    echo "example.com" | ./dnsgen
    cat domains.txt | ./dnsgen -w custom_words.txt -v
    ```

### Examples

*   Generate permutations from a file named `domains.txt` and write the output to `variations.txt`:

    ```bash
    ./dnsgen -o variations.txt domains.txt
    ```

*   Use a custom wordlist and fast mode, taking input from a pipe:

    ```bash
    cat domains.txt | ./dnsgen -w my_wordlist.txt -f
    ```

*   Enable verbose logging to see detailed debug information:

    ```bash
    ./dnsgen -v input.txt
    ```

## Default Wordlist

The `words.txt` file included in the repository contains a list of common words used for generating permutations. You can customize this file or use the `-w` option to provide your own.

**Note:** The `words.txt` file should contain one word per line, and lines starting with `#` are considered comments and are ignored.

## Logging

`dnsgen` uses `logrus` for logging. The `-v` flag enables verbose logging (debug level), which provides more detailed information about the generation process.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.

## Acknowledgements

*   This project is inspired by the original Python `dnsgen` tool: [https://github.com/AlephNullSK/dnsgen](https://github.com/AlephNullSK/dnsgen)
*   Uses the `cobra` library for command-line argument parsing: [https://github.com/spf13/cobra](https://github.com/spf13/cobra)
*   Uses the `logrus` library for logging: [https://github.com/sirupsen/logrus](https://github.com/sirupsen/logrus)
*   Uses `golang.org/x/net/publicsuffix` for TLD extraction.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues to improve `dnsgen`.
