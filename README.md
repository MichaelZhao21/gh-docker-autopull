# Github Docker Autopull

This script will listen on a port for a Github repo webhook call. If the specified branch gets updated, the repo gets cloned and a Docker build will be run. 

## Obviously still WIP

TODO:

1. Server that listens for a webhook
2. Check for correct repo and update event
3. Clone repo and run docker build
4. Add support for multiple repos/webhooks
5. Daemonize process

