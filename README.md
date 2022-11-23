# Creek

[Brook](https://github.com/txthinking/brook) is a cross-platform network tool designed for developers. Creek on the other hand is based on brook (fully compatible), but it offers more features. The aim is to provide a more user-friendly interface and a way to easily manage multiple Brook instances and clients.

## Installation

At the moment I'm not providing any pre-built binaries (coming soon). You can however build the project yourself, it's pretty straightforward and fast.

First, you need to have the Go compiler installed. You can download it and install it easily on Windows, Linux and macOS, just follow the instructions on the [official website](https://go.dev/doc/install).

Once you have Go installed, you can build the project as easily as running the following commands:

```bash
git clone https://github.com/Klaushayan/creek
cd creek/cli/creek
go build .
```

And a binary named `creek` will be created in the `cli/creek` directory.

## Server

Both the server and the client work the same way they do in the original Brook. The only difference is that a few options have been added.

```
creek server -l :9999 -p hello
```

## Client

Creek servers are completely compatible with Brook clients.

[GUI Client](https://txthinking.github.io/brook/)

> replace 1.2.3.4 with your server IP

-   brook server: `1.2.3.4:9999`
-   password: `hello`

[CLI Client](https://txthinking.github.io/brook/)

> create socks5://127.0.0.1:1080

`creek client -s 1.2.3.4:9999 -p hello`

## Roadmap

-   [x] Limit the number of IP addresses that can connect to the server
-   [x] Block IP addresses
-   [ ] Logging into a file
-   [ ] An API for managing the server
-   [ ] Session management to manage multiple instances of the server
-   [ ] Web panel for easier configuration, management and monitoring

## Documentation

For more features, please check the [**Documentation**](https://txthinking.github.io/brook/).
-   [Telegram](https://t.me/brookgroup)
