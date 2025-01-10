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