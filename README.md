<div align=center>
<img src="https://testingcf.jsdelivr.net/gh/4444TENSEI/CDN/img/avatar/kaze/png/rounded.png" alt="KazeCryptoAPI" width="200" height="200"/>
<h1>KazeFrame</h1>
<h3>迷你、轻盈的Gin+Gorm+Redis开发脚手架。</h3>
<p>使用经典MVC模式，实现登陆、注册、邮件发送等多个在线项目开发必备基础接口</p>
<p>包含JWT鉴权、频繁请求限制中间件，以及一些实用、通用工具</p>
<p>以及本地运行日志/日志切割、详细的站点访问日志记录</p>
<p>封装数据CRUD操作已经全部封装好，极易拓展新增业务</p>
<p>
    <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
    <img src="https://img.shields.io/badge/redis-%23DD0031.svg?&style=for-the-badge&logo=redis&logoColor=white" />
</p>
</div>

> # 简介

> 相信很多时候，你都需要一次又一次的重新开始，那么由KazeFrame来为你提供，相信这个项目会是一个适合你的出发点。

1. 在线项目必备接口：登录、注册、接口请求日志等等
2. 邮件服务：注册验证码、密码重置验证码发送，邮件HTML模板
3. 中间件：权限认证中间件、基于IP统计的访问频率限制中间件、站点访问捕捉中间件
4. 通用工具：JWT生成、密码散列加密、RSA非对称加密工具
5. 用户在线状态：通过访问Token+刷新Token机制，实现用户在线状态监控以及令牌有效期控制
6. 数据库访问层合理封装（`dao`包）：没有过度抽象，为你的新增业务接口，快速实现CRUD操作
7. 学习：源码注释极为详细
8. 详细的在线接口文档：[点击查看](https://apifox.com/apidoc/shared-dd0fa318-5437-499c-9227-956e9003fabc)

![](https://testingcf.jsdelivr.net/gh/4444TENSEI/CDN/img/server/readme/KazeFrame/001.webp)

![](https://testingcf.jsdelivr.net/gh/4444TENSEI/CDN/img/server/readme/KazeFrame/002.webp)

![](https://testingcf.jsdelivr.net/gh/4444TENSEI/CDN/img/server/readme/KazeFrame/003.webp)

> # 开始

1. ### 拉取项目，进入到项目根目录，启动你的编辑器

   ```
   git clone https://github.com/4444TENSEI/KazeFrame.git
   ```

2. ### 修改配置文件`./static/server/config/config.yml`

   > 自行去除`.example`后缀，配置说明内部有详细注释。*项目默认使用MySQL进行初始化，得益于Gorm你可以选择其他多种数据库，如需其他数据库可以自行前往./internal/config/db.go修改依赖*

3. ### 运行您所配置好的数据库、Redis

4. ### 启动服务

   ```
   go run .
   ```

   > 推荐使用`Air`进行项目开发（项目运行时自动监测代码文件改动，自动重新编译进行热重载），另外项目启动时将会在项目目录下创建`运行日志.log`，能够用于排查部署到服务器后的运行错误

   > # 目录结构

```
KazeFrame
├─ main.go  //主程序
├─ go.mod   //模块配置
├─ go.sum   //模块依赖
├─ .air.toml    //air配置
├─ script   //各类便捷脚本
│  ├─ pack_dist_linux.ps1   //Linux可执行文件打包脚本
│  ├─ pack_dist_win.ps1 //Windows可执行文件打包脚本
│  └─ pack_source_code.ps1  //源码备份打包脚本
├─ pkg  //公用包
│  └─ util  //通用工具函数
│     ├─ bcrypt.go  //用于数据库用户密码散列加密
│     ├─ rsa.go //rsa密钥生成与验证
│     ├─ jwt.go //jwt生成与验证
│     └─ response_map.go    //Gin快捷json响应封装工具函数
└─ internal //内部模块
   ├─ service   //内部服务模块
   │  ├─ email_service.go   //通用邮件服务
   │  └─ user_service.go    //通用用户服务
   ├─ router    //路由模块
   │  └─ router.go    //定义接口路由
   ├─ model //结构模型
   │  ├─ table_basic.go //存放数据库自动创建的表结构
   │  ├─ user_request.go    //用户请求结构体
   │  └─ user_response.go   //用户响应结构体
   ├─ middleware    //中间件
   │  ├─ auth.go     //jwt权限验证中间件
   │  ├─ cors.go     //跨域处理中间件
   │  ├─ request_limit.go    //请求频率限制中间件
   │  └─ request_log.go //请求日志中间件(记录到数据库表)
   ├─ dao   //通用数据库访问层，高度封装CRUD操作，在这里拓展新的数据库表仓库方法
   │  ├─ init.go    //初始化数据库，并在这里为你的新接口、新数据表操作增加仓库初始化方法
   │  ├─ create.go  //创建数据
   │  ├─ delete.go  //删除数据
   │  └─ update.go  //更新数据
   │  ├─ read.go    //查询数据
   ├─ config    //配置模块
   │  ├─ config.go   //config.yml配置结构定义
   │  ├─ db.go   //数据库配置
   │  ├─ init.go    //负责项目整体运行初始化
   │  ├─ logger.go    //zap日志/切割配置
   │  ├─ redis.go    //redis配置
   │  └─ seed.go    //初始化表数据
   ├─ cache //缓存模块
   │  ├─ create_captcha.go   //创建验证码
   │  ├─ key_map.go //redis键名映射表
   │  └─ set_status.go   //设置在线状态
   └─ api   //接口模块
      ├─ user   //用户接口
      │  ├─ count_online.go //统计在线人数
      │  ├─ delete_user.go  //删除用户
      │  ├─ find_all_user.go    //查询所有用户
      │  ├─ find_user_exact.go  //精确查询用户接口
      │  ├─ find_user_fuzzy.go  //模糊查询用户接口
      │  ├─ get_online_list.go  //获取在线用户列表
      │  ├─ get_profile.go  //获取用户信息
      │  ├─ keep_online.go   //使用刷新令牌刷新访问令牌维持用户在线
      │  ├─ login.go     //登录
      │  ├─ logout.go    //注销
      │  ├─ register.go //注册
      │  ├─ reset_password.go    //重置密码
      │  └─ update_profile.go    //更新个人信息
      ├─ email //邮件发送接口
      │  └─ user_captcha.go
      └─ monitor //接口访问日志、Redis缓存键操作接口
         ├─ redis_clear.go   //清除Redis缓存接口、如接口请求频繁、邮件验证码
         ├─ db_log_clear.go   //清除数据库中的"接口请求日志"表数据
         └─ db_log_search.go  //查询数据库中的"接口请求日志"表数据
```

> # **支持项目**

- #### ⭐等待您小小的一个`Star`！

- #### ⭐如果有任何Bug / 更多的建议，欢迎你随时提出issue！

- #### ⭐联络作者：2479757568@qq.com
