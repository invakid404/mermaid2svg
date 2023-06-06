apt update
apt install -y firefox-esr jq

json=$(curl -s https://api.github.com/repos/mozilla/geckodriver/releases/latest)
url=$(printf "%s" "$json" | jq -r '.assets[].browser_download_url | select(contains("linux64") and endswith("gz"))')
curl -s -L "$url" | tar -xz
chmod +x geckodriver
mv geckodriver /usr/bin/
