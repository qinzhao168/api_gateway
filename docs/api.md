# <span id="1">目录</span>

* **[目录](#1)**
* **[pass-api协议](#2)**
* **[版本](#3)**
* **[提示](#4)**
* **[更改](#5)**
* **[传输](#6)**
* **[格式](#7)**
* **[规范](#8)**
    - [规范](#8.1)
    - [请求](#8.2)
    - [响应](#8.3)
    - [错误码](#8.4)
* **[协议](#9)**
* **[API-GATEWAY模块](#10)**
    - [部署应用](#10.1)
    - [删除应用](#10.2)
    - [更新应用配置 包括弹性伸缩  容器个数  启动  停止  重新部署](#10.3)
    - [获取应用状态](#10.4)
    - [获取应用中的所以容器](#10.5)

# <span id="2">pass-api协议</span>

## <span id="3">版本</span>
---

**v1.0**

## <span id="4">提示</span>
---

本文为markdown格式文本，可使用beyond compare或类似工具比较版本间的修改。
改动时请拉取最新代码进行改动(推荐)，或者在code.yunfancdn.cn对应文件下进行编辑。
不要用空格缩进，而应该用tab缩进。

## <span id="5">更改</span>
---
- 2016/10/8, 黄佳, 1.0
  * 协议模板创建


## <span id="6">传输</span>
---

> 使用HTTP作为传输层; 

> 使用UTF-8编码; 

## <span id="7">格式</span>
---

> 请求使用原始的HTTP格式；

> 响应使用JSON封装，详情见下面响应说明；

> 时间格式采用如下形式：[beginTime, endTime];

消息格式为Json,
参考：http://www.json.org/json-zh.html

## <span id="8">规范</span>
---

### <span id="8.1">规范</span>

> 大体上符合REST风格；

> URL都采用单数，复数的情况使用路径文件夹形式，例如POST BaseURI/supplier/, 注意最后的'/'表示文件夹；

> 命名采用小写开头，驼峰格式，例如supplierId;

### <span id="8.2">请求</span>

> GET: 用于读取信息，参数在query中，成功返回200；幂等；

> POST: 主要用于创建，也可以用于更改，参数在body中，成功返回201；非幂等；

> PUT: 用于更改已有资源，参数和POST一样，成功返回201；非幂等；

> DELETE：用于删除资源，成功返回204；非幂等；

### <span id="8.3">响应</span>

- 格式如下：

```text
{
    "api": "1.0",
    "status": 200,
    "err": "OK",
    "data": {
        "totalSize": 200
    }
}
```

- 空数组：
 "data": []

- 空对象：
 "data": {}

### <span id="8.4">错误码</span>

- 200 OK - [GET]：服务器成功返回用户请求的数据，该操作是幂等的（Idempotent）。
- 201 CREATED - [POST/PUT/PATCH]：用户新建或修改数据成功。
- 204 NO CONTENT - [DELETE]：用户删除数据成功。
- 400 INVALID REQUEST - [POST/PUT/PATCH]：用户发出的请求有错误，服务器没有进行新建或修改数据的操作，该操作是幂等的。。
- 404 NOT FOUND - [*]：用户发出的请求针对的是不存在的记录，服务器没有进行操作，该操作是幂等的。
- 500 INTERNAL SERVER ERROR - [*]：服务器发生错误，用户将无法判断发出的请求是否成功。

### <span id="8.5">请求地址</span>

> RootURI: http://120.24.7.22:9105/

> ApiURI: RootURI/api/v1/

## <span id="9">协议</span>
---

- build组件api

## <span id="10"build组件模块</span>
---

### <span id="10.1">构建应用</span>

构建应用。

URI: ApiURI/build

Method: POST

**请求**

- JSON:

```text
{
    "name": "app-test",
    "version": "1.0",
    "base": "192.168.3.205:5000/golang:pro",
    "image": "",
    "tarball": "",
    "repo": "192.168.3.205:5000",
    "status": "",
    "type": "application",
    "userid": 3,
    "username": "huangjia",
}
```

**响应**

- HTTP Status: 201;
- JSON:

```text
{
    "api": "1.0", 
    "status": "201", 
    "err": "OK", 
    "data": {
        "success": true, 
        "reason": "build success"
    }
}
```

### <span id="10.2">查询全部应用镜像</span>

查询全部应用镜像。

URI: ApiURI/list

Method: GET

**请求**

- ApiURI/list 
- 说明:

**响应**

- HTTP Status: 200
- JSON:

```text
{
    "api": "1.0", 
    "status": "200", 
    "data": {
        "data": [
            {
                "id": 1, 
                "name": "test", 
                "version": "1.0", 
                "base": "test:1.0", 
                "image": "test:pro", 
                "tarball": "", 
                "repo": "192.168.1.123", 
                "status": "build success", 
                "type": "application", 
                "userid": 2, 
                "username": "huangjia", 
                "createat": "2016-08-29T17:14:03+08:00"
            }, 
            {
                "id": 2, 
                "name": "test1", 
                "version": "1.1", 
                "base": "test:1.1", 
                "image": "test1:test", 
                "tarball": "", 
                "repo": "192.168.1.123", 
                "status": "build success", 
                "type": "application", 
                "userid": 3, 
                "username": "lili", 
                "createat": "2016-08-29T17:36:19+08:00"
            }
        ], 
        "total": 3
    }
}
```

### <span id="10.3">上传文件</span>

上传文件。

URI: ApiURI/upload

Method: POST

**请求**

- ApiURI/upload
- 说明：需要web前端用form表单或者ajax提交上传


**响应**

- HTTP Status: 200;
- JSON:

```text
{
    "api": "1.0", 
    "status": 201, 
    "err": "OK", 
    "data": {
        "reason": "upload success", 
        "success": true
    }
}
```

### <span id="10.4">输出docker构建应用的日志</span>

输出docker构建应用的日志。

URI: ApiURI/logswrite

Method: GET

**请求**

- ApiURI/logswrite
- 说明：需要web前端通过ajax或者其他方式建立websocket的client端


**响应**

- HTTP Status: 200
- client端动态的获取日志信息

### <span id="10.5">获取github，gitlab项目列表</span>

获取github，gitlab项目列表。

URI: ApiURI/projects

Method: GET

**请求**

- ApiURI/projects?login=name&password=psword&type=gitlab
- 说明：login github或者gitlab用户名 password github或gitlab的用户密码 type 代码库的类型，目前支持github  gitlab两种

**响应**

- HTTP Status: 200;
- JSON:

```text
{
  "api": "1.0",
  "status": 200,
  "err": "OK",
  "data": [
    {
      "archived": false,
      "avatar_url": null,
      "created_at": "2016-08-25T06:51:18.339Z",
      "creator_id": 43,
      "default_branch": "master",
      "description": "存放文档",
      "http_url_to_repo": "http://code.yunfancdn.cn/huangjia/document.git",
      "id": 80,
      "issues_enabled": true,
      "last_activity_at": "2016-09-01T11:59:24.978Z",
      "merge_requests_enabled": true,
      "name": "document",
      "name_with_namespace": "huangjia / document",
      "namespace": {
        "avatar": null,
        "created_at": "2016-08-25T06:27:27.235Z",
        "description": "",
        "id": 46,
        "name": "huangjia",
        "owner_id": 43,
        "path": "huangjia",
        "updated_at": "2016-08-25T06:27:27.235Z"
      },
      "owner": {
        "avatar_url": "http://www.gravatar.com/avatar/b3a1d7010a5c0683166415f78736fe1a?s=40&d=identicon",
        "id": 43,
        "name": "huangjia",
        "state": "active",
        "username": "huangjia"
      },
      "path": "document",
      "path_with_namespace": "huangjia/document",
      "public": false,
      "snippets_enabled": false,
      "ssh_url_to_repo": "git@code.yunfancdn.cn:huangjia/document.git",
      "tag_list": [],
      "visibility_level": 0,
      "web_url": "http://code.yunfancdn.cn/huangjia/document",
      "wiki_enabled": true
    },
    {
      "archived": false,
      "avatar_url": null,
      "created_at": "2016-08-25T06:35:13.038Z",
      "creator_id": 43,
      "default_branch": "master",
      "description": "build is yunfan pass platform 's module，use it to build application image",
      "http_url_to_repo": "http://code.yunfancdn.cn/huangjia/build.git",
      "id": 79,
      "issues_enabled": true,
      "last_activity_at": "2016-08-29T11:03:04.978Z",
      "merge_requests_enabled": true,
      "name": "build",
      "name_with_namespace": "huangjia / build",
      "namespace": {
        "avatar": null,
        "created_at": "2016-08-25T06:27:27.235Z",
        "description": "",
        "id": 46,
        "name": "huangjia",
        "owner_id": 43,
        "path": "huangjia",
        "updated_at": "2016-08-25T06:27:27.235Z"
      },
      "owner": {
        "avatar_url": "http://www.gravatar.com/avatar/b3a1d7010a5c0683166415f78736fe1a?s=40&d=identicon",
        "id": 43,
        "name": "huangjia",
        "state": "active",
        "username": "huangjia"
      },
      "path": "build",
      "path_with_namespace": "huangjia/build",
      "public": true,
      "snippets_enabled": false,
      "ssh_url_to_repo": "git@code.yunfancdn.cn:huangjia/build.git",
      "tag_list": [],
      "visibility_level": 20,
      "web_url": "http://code.yunfancdn.cn/huangjia/build",
      "wiki_enabled": true
    }
  ]
}
```