![Logo](https://raw.githubusercontent.com/sawyersteven/reach/master/img/Logo_wide.png)

# Reach
Reach is a quick and simple CLI application for testing availability of remote http resources.

### Installation
Download the appropriate release for your operating system and copy the binary to a directory that is in your PATH, or add reach to your PATH in order to access it from a terminal.


### Usage
    reach [OPTIONS] URL

    Options:
    -c, --nocolor               Print output without colors
    --maxredirects=REDIRECTS    Maximum redirects to follow [default: 20]
    --timeout=SECONDS           HTTP request timeout in seconds [default: 15]
    --help                      Display this help message
    --version                   Display version and license info


Reach will attempt to reach the given URL and will communicate the connection status back to the terminal. Simple color-coded status results are easy to read but may be turned off if your terminal does not support ANSI colors.

### Examples
![Examples](https://raw.githubusercontent.com/sawyersteven/reach/master/img/Examples.png)
