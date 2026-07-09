[![Deploy to AWS ECS](https://github.com/Foxeh/gofox/actions/workflows/yitties-aws.yml/badge.svg?branch=master)](https://github.com/Foxeh/gofox/actions/workflows/yitties-aws.yml/badge.svg?branch=master)

# GoFox
Discord bot written in Go.

### Installing
Use go get to download repository into your Go environment
```
go get github.com/foxeh/gofox
```
### Settings

* Rename or create new file

`config.example.json > config.json`

* Edit contents

Configuration is provided through environment variables. This is how the ECS deployment is
configured:

| Environment variable | Purpose |
| --- | --- |
| `BOT_KEY` | Discord bot token |
| `BOT_STATUS` | Discord status for the bot |

### Running

```
cd $GOPATH/src/github.com/foxeh/gofox
go build
./gofox
```