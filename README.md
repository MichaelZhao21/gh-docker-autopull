# Github Docker Autopull

This script will listen on a port for a Github repo webhook call. If the specified branch gets updated, the repo gets cloned and a Docker build will be run.

Note that the use case of this project is pretty specific. I have a single cloud server (hosted on a [Digital Ocean Droplet](https://www.digitalocean.com/products/droplets)) that has a running Docker daemon on it.

## Setup

As listed above, you will need to set up a droplet or your own cloud server. Your server will need to have [NGINX](https://www.nginx.com/) installed. To set up the [reverse proxy](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/) for your instance of `gh-autopull`, you will need to add the following specification to `sites-enabled`:

```nginx
server {
        server_name <EXTERNAL_ENDPOINT_URL>;

        location / {
                proxy_pass http://<SERVER_IP>:<PORT>;
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection 'upgrade';
                proxy_set_header Host $host;
                proxy_cache_bypass $http_upgrade;
                proxy_set_header X-Real-IP $remote_addr;
        }

        listen 80;
}
```

The fields you will need to fill in:

- `<EXTERNAL_ENDPOINT_URL>` = url that Github will call from the webhook
- `<SERVER_IP>` = Internal IP address of the server (I think u can use the loopback address but not sure)
- `<PORT>` = port that `gh-autopull` script is running on

You may also want to set up a certificate [using certbot]()

# TOOOOODOOOOOOO LIST

TODO:

2. Check for correct repo and update event
3. Clone repo and run docker build
4. Add support for multiple repos/webhooks
5. Daemonize process
6. [Extra] Header hash check for security
