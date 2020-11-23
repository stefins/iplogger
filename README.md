# iplogger
A tiny Go application to log your public IP address every 20 minutes.

## Installation

```bash
go get -u github.com/iamstefin/iplogger
```
or download from releases

## Usage

```bash
iplogger &
disown
```
### To view the logged ip addresses
```bash
cat ~/.iplogger/log.txt
```

## LICENSE
MIT
