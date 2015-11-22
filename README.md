

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



https://github.com/go-godo/godo
