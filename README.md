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

```json
{
  "botKey": "<INSERT REALLY LONG BOT KEY HERE>",
  "wolframID": "<WOLFRAM ALPHA REQUEST KEY>",
  "status": "<BOT STATUS FOR DISCORD>"
}
```


### Running

```
cd $GOPATH/src/github.com/foxeh/gofox
go build
./gofox
```

