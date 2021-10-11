#!/bin/bash

set -e -o pipefail

echo "" > test.log

function run() {
  echo '-------------------------------------------'
  echo "Test => $1"
  eval "$2"
  echo '✅'
  echo '-------------------------------------------'
}

#run 'dns registry A' 'dig registry-1.docker.io A'
#run 'dns registry AAAA' 'dig registry-1.docker.io AAAA'
#run 'dns auth A' 'dig auth.docker.io A'
#run 'dns auth AAAA' 'dig auth.docker.io AAAA'
#run 'dns cloudflare A' 'dig production.cloudflare.docker.com A'
#run 'dns cloudflare AAAA' 'dig production.cloudflare.docker.com AAAA'

run 'ipv4 registry curl' 'curl -4svo /dev/null https://registry-1.docker.io/'
run 'ipv4 auth curl' 'curl -4svo /dev/null https://auth.docker.io/'
run 'ipv4 cloudflare  curl' 'curl -4svo /dev/null https://production.cloudflare.docker.com/'
run 'ipv4 cloudflare trace' 'curl -4sv http://production.cloudflare.docker.com/cdn-cgi/trace'

#run 'ipv6 registry curl' 'curl -6svo /dev/null https://registry-1.docker.io/'
run 'ipv6 auth curl' 'curl -6svo /dev/null https://auth.docker.io/'
run 'ipv6 cloudflare  curl' 'curl -6svo /dev/null https://production.cloudflare.docker.com/'
#run 'ipv6 cloudflare trace' 'curl -6sv http://production.cloudflare.docker.com/cdn-cgi/trace'

curl -svo /dev/null https://registry-1.docker.io/ -w "\nContent Type: %{content_type} \nHTTP Code: %{http_code} \nHTTP Connect:%{http_connect} \nNumber Connects: %{num_connects} \nNumber Redirects: %{num_redirects} \nRedirect URL: %{redirect_url} \nSize Download: %{size_download} \nSize Upload: %{size_upload} \nSSL Verify: %{ssl_verify_result} \nTime Handshake: %{time_appconnect} \nTime Connect: %{time_connect} \nName Lookup Time: %{time_namelookup} \nTime Pretransfer: %{time_pretransfer} \nTime Redirect: %{time_redirect} \nTime Start Transfer: %{time_starttransfer} \nTime Total: %{time_total} \nEffective URL: %{url_effective}\n" 2>&1

# TODO: follow this url
https://support.cloudflare.com/hc/en-us/articles/203118044-Gathering-information-for-troubleshooting-sites#h_8c9c815c-0933-49c0-ac00-b700700efce7
