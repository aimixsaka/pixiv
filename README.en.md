## pixiv
**A go pixiv spider**

### Before
Use `go mod init xxx` to init your module, then use `go install github.com/yfaimisaka/pixiv@latest` to add it to your module.

You need to add a `config.yaml` in your project root to config `proxy`, `download` ...

Example:
```yaml
proxy:
  host: 127.0.0.1
  port: 10809
    
download:
  single: 
    path: "C:\\xxxx\\pixiv\\sigle"
  keyword:
    path: "C:\\xxxx\\pixiv\\keyname"
  rank:
    path: "C:\\xxxx\\pixiv\\rank"
  tag:
    path: "C:\\xxxx\\pixiv\\tag"
  user:
    path: "C:\\xxxx\\pixiv\\user"

upload:  # now only support minio
  use: true 
  endPoint: xxxxxxxxxxxxx:9000
  accessKeyID: xxxxxxxxxx
  secretAccessKey: xxxxxxxxxx
  useSSL: false
  bucketName: pixiv
```
### Usage
#### Single() 
> get single picture with workid

Sample:
```go
import "github.com/yfaimisaka/pixiv"

pixiv.Single("90938571").Download() // download a picture, save it to your download.single.path
pixiv.Single("90938571").Upload() // upload to your minio server
```
#### KeyWord(word string)
> get pictures by keyword

Sample:
```go
import "github.com/yfaimisaka/pixiv"

pixiv.KeyWord("灼眼のシャナ").Num(10).Download() // download 10 pictures on keyword=灼眼のシャナ
pixiv.KeyWord("灼眼のシャナ").Num(10).Upload() // upload 10 keyword pictures to your minio server
```

#### Rank()
> get today's rank

Sample:
```go
import "github.com/yfaimisaka/pixiv"

pixiv.Rank().Num(10).Download() // download 10 pictures from today's rank
pixiv.Rank().Num(10).Upload() // upload top 10 pictures from rank
```

#### Tag(tagName string)
> get pictures on specific tag

Sample:
```go
import "github.com/yfaimisaka/pixiv"

pixiv.Tag("緋弾のアリア").Num(10).Download() // download 10 pictures on tag=緋弾のアリア
pixiv.Tag("緋弾のアリア").Num(10).Upload() // upload 10 keyword pictures to your minio server
```

#### User(userId string)
> get pictures on userId


Sample:
```go
import "github.com/yfaimisaka/pixiv"

pixiv.User("104180").Num(10).Download() // download 10 pictures on userid=104180
pixiv.User("104180").Num(10).Upload() // upload 10 keyword pictures to your minio server
```
