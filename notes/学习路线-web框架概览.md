## 一 学习路线——Web框架概览

Web框架 -> Server定义 -> 静态路由匹配 -> 通配符匹配 -> 参数路径
-> 正则匹配 -> Context 处理输入输出 -> AOP方案 -> 可观测性 Middleware
-> 支持路由化 Middleware -> 页面渲染 -> 文件处理与静态资源 -> Session -> 用户服务 Web 接口

对于一个 Web 框架来说，至少要提供三个抽象：
- 代表服务器的抽象，这里我们称之为 Server
- 代表上下文的抽象，这里我们称之为 Context
- 路由树

### 1.1 Web核心—— Server

对于一个 Web 框架来说，我们首先要有一个整体代表服务器的抽象，也就是 Server。

Server 从特性上来说，至少要提供三部分功能：
- 生命周期控制：即启动、关闭。如果在后期，我们还要考虑增加生命周期回调特性。
- 路由注册接口：提供路由注册功能
- 作为 http 包 到 Web 框架的桥梁

Server —— http.Handler 接口

http 包暴露了一个接口，Handler，他是我们引入自定义 Web 框架相关的连接点。

### 1.2 Server —— 接口定义

Server 定义版本一：只组合 http.Handler

优点：
- 用户在使用的时候只需要调用 http.ListenAndServe 就可以
- 和 HTTPS 协议无缝衔接
- 极简设计

缺点：
- 难以控制生命周期，并且在控制生命周期的时候增加回调支持
- 缺乏控制力，如果将来希望支持优雅退出的功能，将难以支持

Server 定义版本二：组合 http.Handler 并且增加 Start 方法。

优点：
- Server 即可以当成普通的 http.Handler 来使用，又可以作为一个独立的实体，拥有自己的管理生命周期的能力
- 完全的控制

缺点：
- 如果用户不希望使用 ListenAndServTLS，那么 Server 需要提供 HTTPS 的支持

版本一和版本二都直接耦合了 Go 自带的 http 包，如果我们希望换为 fasthttp 或者类似的 http 包，则会非常困难

### 1.3 Server —— 注册路由 API 设计

大体上有两类方法：
- 针对任意方法的：如 Gin 和 Iris 的Handle 方法、Echo 的 Add 方法
- 针对不同 HTTP 方法的：如 Get、POST、Delete，这一类方法基本上都是委托给前一类方法

所以实际上，核心方法只需要有一个，例如 Handle。其他方法都建立在这上面。

### 1.4 Server —— ServeHTTP 方法

ServeHTTP 方法是作为 http 包与 Web 框架的关联点，需要在 ServeHTTP 内部，执行：
- 构建起 Web 框架的上下文
- 查找路由树，并执行命中路由的代码

