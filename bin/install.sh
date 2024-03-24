cd /usr/bin \
&& rm -f peurl \
&& curl -s https://api.github.com/repos/jackgpalfrey/peurl-cli/releases/latest \
| grep "browser_download_url" \
| cut -d : -f 2,3 \
| tr -d \" \
| wget -qi - \
&& chmod +x peurl
