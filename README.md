# docker log puller
Easily pull docker logs from docker contianers 

#### Purpose
I wanted an easy way for my docker host machines to pull logs from containers instead of having containers implement thier own log collection logic or sharing volumes between host and containers.

#### Usage
The [config.json](https://github.com/Sjeanpierre/docker_log_puller/blob/master/config.json) drives the behavior of this tool by specifying which docker containers to pull from, which logs to pull from the FS, and where to store them on the host machine.

Once the config.json is ready, simply place it in the same directory as the binary and run it.

#### Disclaimer
This method of getting log files from docker containers probably does not rank highly on anyones docker best practices list, but it scratched an itch for me and might help others.

it might be worth it to consider mounting a shared between the host and the containers for storing logs.
