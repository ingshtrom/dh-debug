# Docker Hub Debug (network connections)

This is a bunch of random things in order to debug network connections to Docker Hub (UI, Registry,
CDN for registry). 

[This](https://support.cloudflare.com/hc/en-us/articles/203118044-Gathering-information-for-troubleshooting-sites)
and [this](https://www.cloudflare.com/en-gb/learning/network-layer/what-is-mtr/) have been helpful
in influencing what types of tests go in here.

## Note on Sudo Usage
The MTR tests require sudo in order to use raw sockets. If you run `dh-debug` without
sudo, then the `traceroute` tests will give us some of that same information. **We highly recommend
using `sudo` in order to get the MTR tests in there as well as they are valuable.**

# Using dh-debug

Copy the `dh-config.json` file into the directory you want to run the `dh-debug` binary from. You
can modify as you wish, but the tests included in `dh-config.json` are what we recommend.

Run a debug test suite using `dh-debug` or `sudo dh-debug` (see the note above about using sudo). 

`dh-debug -s` can print some really basic, high-level metrics of tests run and their exit codes or
errors encountered. 

`dh-debug -p` will print out all of the test runs and their accompanying output for further
debugging.

# Build
`make build`

# Test
`make test`

