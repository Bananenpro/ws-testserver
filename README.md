# ws-testserver

## Description

A simple websocket server, which enables you to manually send and receive any text.

## Setup

### Prerequisites

- [Go](https://go.dev/) 1.17+

Clone the repository into a local directory:

```sh
git clone https://github.com/Bananenpro/ws-testserver.git
cd ws-testserver
```

Compile and run the server:

```sh
go run ./cmd/ws-server/main.go
```

Attach a control client:

```sh
go run ./cmd/ws-server-attach/main.go <client-id>
```

## License

Copyright (c) 2022 Julian Hofmann

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
