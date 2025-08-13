# Event

轻量级 Go 事件调度器，用于实现事件驱动编程。支持上下文传递、通配符监听器、结构化日志和错误处理。

## 特性

- 事件分发：发布和订阅领域事件
- 多监听器：一次订阅多个监听器
- 通配符监听器：单个监听器订阅所有事件
- 上下文支持：在事件系统中传递上下文
- 结构化日志：内置日志系统，支持上下文和键值对格式
- 同步/异步处理：选择阻塞或非阻塞事件处理
- 线程安全：并发访问支持
- 错误处理：捕获并记录事件处理器中的错误

## 安装

```shell
go get -u github.com/kevinyan815/event
```

## 快速开始

### 1. 定义事件

```go
type UserCreated struct {
    UserId   int
    UserName string
}

func (*UserCreated) EventName() string {
    return "UserCreated"
}

func (e *UserCreated) EntityID() string {
    return fmt.Sprintf("%d", e.UserId)
}

func (e *UserCreated) EntityType() string {
    return "User"
}
```

### 2. 创建监听器

```go
type UserCreatedListener struct{}

func (*UserCreatedListener) EventHandler() event.Handler {
    return func(e *event.Event) error {
        userData, ok := e.EventData().(*UserCreated)
        if !ok {
            return fmt.Errorf("cannot convert event data")
        }
        
        fmt.Printf("User created: %s (ID: %d)\n", userData.UserName, userData.UserId)
        return nil
    }
}

func (*UserCreatedListener) AsyncProcess() bool {
    return false  // 同步处理
}
```

### 3. 创建通配符监听器

```go
type LogEventListener struct{}

func (*LogEventListener) EventHandler() event.Handler {
    return func(e *event.Event) error {
        fmt.Printf("Event logged: %s (ID: %s)\n", e.Name, e.ID)
        return nil
    }
}

func (*LogEventListener) AsyncProcess() bool {
    return true  // 异步处理
}
```

### 4. 订阅事件

```go
func init() {
    dispatcher := event.Dispatcher()
    
    // 订阅特定事件的多个监听器
    dispatcher.Subscribe(
        event.NewEvent(&UserCreated{}),
        &UserCreatedListener{},
        &UserCreatedErrListener{},
    )
    
    // 订阅通配符监听器
    dispatcher.SubscribeWildcard(&LogEventListener{})
}
```

### 5. 分发事件

```go
func CreateUser(username string) {
    // 创建带有追踪信息的上下文
    ctx := context.WithValue(context.Background(), "trace", "trace-123")
    
    // 创建并分发带有上下文的事件
    event.Dispatcher().Dispatch(
        event.NewEventWithContext(ctx, &UserCreated{
            UserId:   1,
            UserName: username,
        })
    )
}
```

### 6. 错误处理

```go
func (*UserCreatedErrListener) EventHandler() event.Handler {
    return func(e *event.Event) error {
        // 返回错误将被自动记录
        return fmt.Errorf("处理用户创建事件时出错")
    }
}
```

### 7. 自定义日志
自定义日志需要实现 event.ILogger接口
```go
// ILogger
type ILogger interface {
	Debug(ctx context.Context, msg string, params ...LogParam)
	Info(ctx context.Context, msg string, params ...LogParam)
	Warn(ctx context.Context, msg string, params ...LogParam)
	Error(ctx context.Context, msg string, params ...LogParam)
}
```
```go
// 创建自定义日志实现
type MyCustomLogger struct {
    logger *log.Logger //  项目自己的Logger实现
}

func (l *MyCustomLogger) Info(ctx context.Context, msg string, params ...event.LogParam) {
    l.logger.WithContext(ctx).Printf("[INFO] %s, %s", msg, formatParams(params))
}
// 实现所有方法
// .....

// 设置自定义日志
event.SetLogger(&MyCustomLogger{logger: logger})

// 解析日志参数示例
func formatParams(params []LogParam) string {
	if len(params) == 0 {
		return ""
	}

	var parts []string
	for _, p := range params {
		parts = append(parts, fmt.Sprintf("%q, %v", p.Key, p.Value))
	}
	return ", " + strings.Join(parts, ", ")
}
```

## 完整示例

完整示例请查看 [examples 目录](https://github.com/kevinyan815/event/tree/master/examples)。

## 许可证

本项目采用 Apache License 2.0 许可证 - 详见 [LICENSE](https://github.com/kevinyan815/event/blob/master/LICENSE) 文件。
