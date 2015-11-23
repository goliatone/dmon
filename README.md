

## Development
Install dependencies.

TCP server:
```
"github.com/firstrow/tcp_server"
```

Binary data builder:
```
$ go get -u github.com/jteeuwen/go-bindata/...
```

```
go-bindata -o bin.go bin/check
```

Change the package of the generated asset to `data`, and move to a data directory. If we change the contents of bin dir, we should recompile.


Build:
```
$ go build -o dmon && chmod 775
```

Run:

```
$ ./dmon
```

Open a `netcat` connection and send a container id:

```
$ nc localhost 9386
menagerie
OK
```
It will return **KO** if the container **menagerie** was not found, and **OK** it container was found.

We can pass a container id or a container name.


```
$ nc localhost 9386
ed8bc298f850
OK
```

```
$ docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS                    NAMES
ed8bc298f850        b8e8e082f4e5        "node app.js"       About an hour ago   Up About an hour    0.0.0.0:1337->1337/tcp   menagerie
```

## Deployment
To run `docker` without `sudo` command.

Create the `docker` group and add your user:
```
sudo usermod -a -G docker ubuntu
```
Run `newgrp docker` to recognize the new group or log out and log in to have the change applied to groups.

Add upstart `dmon.conf` file. Read more on the [upstart cookbook][uc]

Check if conf file is ok:

Add logrotate conf file, `/etc/logrotate.d/dmon`:

```
/var/log/dmon.log {
    weekly
    size 20M
    missingok
    rotate 52
    notifempty
    nocreate
}
```

Ensure your port is reachable from outside of the box, if you are on an EC2 instance add a rule to the security group of the machine.
`dmon` listends on port **9386** by default:

* port: 9386
* protocol: tcp
* Source: 0.0.0.0/0

## Development
If you need to build the binary on Ubuntu:


```
# apt-get install golang
```

Set `GOPATH`:
```
$ export GOPATH=/opt/GO
```

You might need to `apt-get install git` and some other dependencies.

```
$ git clone https://github.com/goliatone/dmon.git
```

<!--
https://github.com/go-godo/godo
s3cmd put dmon s3://vecna.io

http://upstart.ubuntu.com/getting-started.html
-->

[uc]: http://upstart.ubuntu.com/cookbook/#pre-start-example-debian-and-ubuntu-specific



->
mkdir /opt/dmon/

sudo install

```shell
#!/bin/sh

wget wget https://s3.amazonaws.com/com.goliatone.dmon/ubuntu/dmon
wget wget https://s3.amazonaws.com/com.goliatone.dmon/ubuntu/dmon.conf
wget wget https://s3.amazonaws.com/com.goliatone.dmon/ubuntu/dmon.log

mv dmon /opt/dmon
cp dmon.conf /etc/init/dmon.conf
```
