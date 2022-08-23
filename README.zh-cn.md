## pixiv
go 实现的pixiv爬虫

## Before
先初始化您的module，然后在module中使用 `go install github.com/yfaimisaka/pixiv@latest` 来获取它

您需要在项目根目录中添加一个 `config.yaml` 来配置代理，下载等

例子：
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

upload:  # 现在上传只支持minio
  use: true 
  endPoint: xxxxxxxxxxxxx:9000
  accessKeyID: xxxxxxxxxx
  secretAccessKey: xxxxxxxxxx
  useSSL: false
  bucketName: pixiv
```
### 用法
#### Single(workId string) 
> 使用 workid 获取单张图片

样例：
```go
import "github.com/yfaimisaka/pixiv"

pixiv.Single("90938571").Download() // 下载图片，保存到你配置的download.single.path
pixiv.Single("90938571").Upload() // 上传到你的 minio 服务器
```

#### KeyWord(word string)
> 按关键字获取图片

样例：
```go
import "github.com/yfaimisaka/pixiv"

pixiv.KeyWord("灼眼のシャナ").Num(10).Download() //在关键词上下载10张关键字为 灼眼のシャナ的图片
pixiv.KeyWord("灼眼のシャナ").Num(10).Upload() // 上传10张关键词图片到你的minio服务器
```
#### Rank()
> 获得今天的排名

样例：
```go
import "github.com/yfaimisaka/pixiv"

pixiv.Rank().Num(10).Download() // 下载今天排名的前10张图片
pixiv.Rank().Num(10).Upload() // 上传排名前 10 的图片
```

#### Tag(tagName string)
> 获取特定标签的图片

样例：

```go
import "github.com/yfaimisaka/pixiv"

pixiv.Tag("绯弾のアリア").Num(10).Download() //下载10张标签为 绯弾のアリア 的图片
pixiv.Tag("绯弾のアリア").Num(10).Upload() // 上传10张标签图片到你的minio服务器
```

#### User(userId string)
> 根据userid获取特定画师的作品

样例：

```go
import "github.com/yfaimisaka/pixiv"

pixiv.User("104180").Num(10).Download() // 下载userid=104180的10张图片
pixiv.User("104180").Num(10).Upload() // 上传10张关键词图片到你的minio服务器
```
