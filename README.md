# Better-Brook

Brook is a cross-platform network tool designed for developers. Better-Brook on the other hand is a fork of Brook that is designed for experienced users. The aim is to provide a more user-friendly interface and a way to easily manage multiple Brook instances and clients.

## Installation

At the moment I'm not providing any pre-built binaries (coming soon). You can however build the project yourself, it's pretty straightforward and fast.

First, you need to have the Go compiler installed. You can download it and install it easily on Windows, Linux and macOS, just follow the instructions on the [official website](https://go.dev/doc/install).

Once you have Go installed, you can build the project as easily as running the following commands:

```bash
git clone https://github.com/Klaushayan/better-brook
cd better-brook/cli/brook
go build .
```

And a binary named `brook` will be created in the `cli/brook` directory.

## Server

Both the server and the client work the same way they do in the original Brook. The only difference is that a few options have been added.

```
brook server -l :9999 -p hello
```

## Client

[GUI Client](https://txthinking.github.io/brook/)

> replace 1.2.3.4 with your server IP

-   brook server: `1.2.3.4:9999`
-   password: `hello`

[CLI Client](https://txthinking.github.io/brook/)

> create socks5://127.0.0.1:1080

`brook client -s 1.2.3.4:9999 -p hello`

## Documentation

For more features, please check the [**Documentation**](https://txthinking.github.io/brook/).
-   [Telegram](https://t.me/brookgroup)
