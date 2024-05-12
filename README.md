# simple-demo

## 抖音项目服务端简单示例

具体功能内容参考飞书说明文档

工程无其他依赖，直接编译运行即可

```shell
go build && ./simple-demo
```

### 功能说明

接口功能不完善，仅作为示例

* 用户登录数据保存在内存中，单次运行过程中有效
* 视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可

### 测试

test 目录下为不同场景的功能测试case，可用于验证功能实现正确性

其中 common.go 中的 _serverAddr_ 为服务部署的地址，默认为本机地址，可以根据实际情况修改

测试数据写在 demo_data.go 中，用于列表接口的 mock 测试

# 更新日志
## v1.0&v1.1
技术栈：gin,gorm,mysql,viper(用于配置文件),bcrypt(密码加密，本来用MD5加盐加密，后面想到还要存盐值不安全，便改为bcrypt，而且查阅资料得该加密更安全),go-jwt，ffmpeg(获取视频封面)  
由于是初代版本，总体写下来的感觉大部分都是  
1.(最主要)gorm对数据库的增删改查  
2.gin框架相关编写设计jwt的中间键鉴权，跨域CORS中间键编写，以及获取合法token情况下的user信息  

## 数据库表格设计
**索引是根据查询的情况创建的，对于查询频率高的创建索引，需要留意联合索引创建时字段的顺序，满足最左匹配，防止失效**
### 1.用户个人信息表格
| 列名          | 数据类型 | 约束            | 索引 | 备注 |
| ------------- | -------- | --------------- | ---- | ---- |
| Id            | int64    |      主键           |      |  用户id    |
| CreatedAt     | time.Time|       无          |      |   创建时间   |
| UpdatedAt     | time.Time|        无         |      |   更新时间   |
| Name          | string   |       unique          |  idx_users_name  |  用户名    |
| Avatar        | string   |        无         |      |   头像url   |
| BackgroundImage | string |         无        |      |   主页图片url   |
| Signature     | string   |      无          |      |   个性签名   |
| FollowCount   | int64    |        无         |    | 关注数量 |
| FollowerCount | int64    |      无           |    | 粉丝数 |
| Password      | string   | size:255,not null |   |   密码   |
| WorkCount     | int      | default:0       |      |     作品数量 |
| TotalFavorited | int     | default:0       |      |   获赞数   |
| FavoriteCount | int      | default:0       |      |    喜欢数  |

### 2.视频信息表格
| 列名          | 数据类型 | 约束            | 索引          | 备注        |
| ------------- | -------- | --------------- | ------------- | ----------- |
| Id            | int64    | 主键            |               | 视频ID      |
| CreatedAt     | time.Time| 无              | idx_created_at | 创建时间    |
| AuthorID      | int64    | 外键(User)      | idx_author_id | 作者ID      |
| Title         | string   | 无              |               | 视频标题    |
| PlayUrl       | string   | 无              |               | 播放地址    |
| CoverUrl      | string   | 无              |               | 封面地址    |
| FavoriteCount | int64    | 默认值为0       |               | 被喜欢次数  |
| CommentCount  | int64    | 默认值为0       |               | 评论次数    |


### 3.关注信息表格
| 列名            | 数据类型 | 约束           | 索引               | 备注         |
| --------------- | -------- | -------------- | ------------------ | ------------ |
| Id              | int64    | 主键           |                    | 关注ID        |
| UserId          | int64    | 外键(User)     | idx_follow_user_id | 用户ID        |
| FollowerUserId  | int64    | 外键(User)     | idx_follower_user_id | 粉丝用户ID   |
| CreatedAt       | time.Time| 无             |    | 创建时间     |

### 4.点赞（喜欢）表格
| 列名        | 数据类型 | 约束           | 索引               | 备注          |
| ----------- | -------- | -------------- | ------------------ | ------------- |
| Id          | int64    | 主键           |                    | 收藏ID        |
| UserId      | int64    | 外键(User)     | idx_favorite_user_id | 用户ID        |
| VideoId     | int64    | 外键(Video)    | idx_favorite_video_id | 视频ID      |
| IsFavorite  | bool     | 无             |                    | 是否喜欢      |
| CreatedAt   | time.Time| 无             |     | 创建时间      |

### 5.评论表格
| 列名        | 数据类型 | 约束           | 索引                | 备注       |
| ----------- | -------- | -------------- | ------------------- | ---------- |
| Id          | int64    | 主键           |                     | 评论ID     |
| VideoId     | int64    | 外键(Video)    | idx_comment_video_id | 视频ID     |
| UserId      | int64    | 外键(User)     | idx_comment_user_id | 用户ID     |
| Content     | string   | 无             |                     | 评论内容   |
| CreatedAt   | time.Time| 无             |     | 创建时间   |

### 6.私信表格
| 列名          | 数据类型 | 约束              | 索引                  | 备注      |
| ------------- | -------- | ----------------- | --------------------- | --------- |
| Id            | int64    | 主键              |                       | 消息ID    |
| ToUserId      | int64    | 无                | idx_to_from_created   | 收信用户ID|
| FromUserId    | int64    | 无                | idx_to_from_created   | 发信用户ID|
| Content       | string   | 无                |                       | 消息内容  |
| CreatedAt     | time.Time| 无                | idx_to_from_created   | 创建时间  |

## 接口实现
### 1./douyin/feed/ - 视频流接口（难点）
1.遇到的第一个问题就是时间的问题，需要留意时间戳是毫秒还是秒（服务端需要返回视频流的最早时间，下次请求会带上这个时间，从而实现不断追溯之前的视频）  
2.遇到的第二个问题就是是否关注了视频的作者（查询用户-粉丝表）以及判断是否点赞，这个就涉及到数据库表格的设计，需要有个用户-视频点赞表  

### 2./douyin/user/register(login)/ - 用户注册(登录)接口
1.比较常规，正常写就可以，只是注册密码需要用bcrypt加密，登录时候需要请求的password和查询数据库获得用户的加密后的密码同样用bcrypt的compare函数进行比对  

### 3. /douyin/user/ - 用户信息
1.很常规，正常响应即可

### 4./douyin/publish/action/ - 视频投稿（难点）
1.首先弄清楚思路，投稿的视频服务器需要有地方存，因此需要将post的视频存在某个文件夹里，并且还需要设置为静态资源，输入url能够访问到的（因为视频流获取视频就是靠url获取）  
2.使用ffmpeg获取视频第一帧作为封面，同理也需要存在某个文件夹里面作为静态资源  
3.**难点！！！** 由于需要对两个数据库（用户表对作品数量++，视频表添加视频信息）,因此使用事务来保证数据安全，一旦出错立即回滚  

### 5./douyin/publish/list/ - 发布列表
1.常规操作，正常读取数据库  

### 6./douyin/favorite/action/ - 赞操作
1.**难点** 首先需要对点赞数据库查找是否有这一条点赞记录，若没有则添加，而且还要在视频用户第三方表上（记录用户和视频是否点赞的表格）**但是！！！** 我后面重新思考优化代码的时候发现这个表格根本就不需要，有favorite就够了，只需要再添加一个is_favorite的字段就可以了。好的话题扯回来，点赞需要对多个表格进行修改，因此用到了事务操作，一旦出错就回滚  
2.**难点** 关于高并发场景下需要对视频的获赞数量字段以及视频作者的获赞数量字段进行更改的时候需要加悲观锁，为了防止并发问题导致读写不正确。  
3.**难点**边界处理，重复点赞、重复取消点赞的问题  

### 7./douyin/favorite/list/ - 喜欢列表
很常规

### 8./douyin/comment/action/ - 评论操作
1.**难点**同样也需要对评论数量上锁，开启事务

### 9./douyin/comment/list/ - 视频评论列表
很常规

### 10./douyin/relation/action/ - 关系操作
1.**难点**同样要对被关注的人的粉丝数量上锁操作，开启事务   
2.**难点**边界情况处理，比如不能重复关注，重复取关  

### 11./douyin/relatioin/follow/list/ - 用户关注列表
1.**坑点**由于发起请求的人不一定是被看关注列表的人，因此需要在获取列表的时候还要逐个判断请求人和关注列表的人是否关注

### 12./douyin/relation/follower/list/ - 用户粉丝列表
**同上坑点**

### 13./douyin/message/action/ - 消息操作
很常规，正常上传聊天内容于数据库

### 14./douyin/message/chat/ - 聊天记录
**没能的解决问题：** 由于前端那边点击发送后也会有聊天气泡弹出，导致我获得聊天记录的时候会重复获得这个最新发出去的消息气泡，我想了好久觉得只有后端应该是解决不了这个问题的    
**个人认为解决方案：** 前端这边点击发送可以不弹出气泡，后端这边在遇到发送消息的请求的时候可以直接调用获得历史信息的接口，完成即使反馈。  

## 优化空间：
1.登录注册相关用户数据存redis，聊天内容存redis，评论，点赞相关对数据库频繁操作也可以用redis优化  
2.可以对视频进行压缩，**看看能不能解决每次单个视频请求时间过长的问题**
