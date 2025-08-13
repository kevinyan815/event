# Event

A lightweight event dispatcher for Go projects to implement event-driven programming. Supports context propagation, wildcard listeners, structured logging, and error handling.

[中文介绍](https://github.com/kevinyan815/event/blob/master/README.cn.md)

## Features

- Event Dispatching: Publish and subscribe to domain events
- Multiple Listeners: Subscribe multiple listeners to an event at once
- Wildcard Listeners: Subscribe to all events with a single listener
- Context Support: Propagate context through the event system
- Structured Logging: Built-in logging system with context and key-value formatting
- Synchronous/Asynchronous Processing: Choose between blocking and non-blocking handlers
- Thread Safety: Concurrent access support
- Error Handling: Capture and log errors from event handlers

## Installation

```shell
go get -u github.com/kevinyan815/event
```

## Quick Start

### 1. Define Events

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

### 2. Create Listeners

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
    return false  // Synchronous processing
}
```

### 3. Create Wildcard Listeners

```go
type LogEventListener struct{}

func (*LogEventListener) EventHandler() event.Handler {
    return func(e *event.Event) error {
        fmt.Printf("Event logged: %s (ID: %s)\n", e.Name, e.ID)
        return nil
    }
}

func (*LogEventListener) AsyncProcess() bool {
    return true  // Asynchronous processing
}
```

### 4. Subscribe to Events

```go
func init() {
    dispatcher := event.Dispatcher()
    
    // Subscribe multiple listeners to a specific event
    dispatcher.Subscribe(
        event.NewEvent(&UserCreated{}),
        &UserCreatedListener{},
        &UserCreatedErrListener{},
    )
    
    // Subscribe wildcard listener
    dispatcher.SubscribeWildcard(&LogEventListener{})
}
```

### 5. Dispatch Events

```go
func CreateUser(username string) {
    // Create context with tracing info
    ctx := context.WithValue(context.Background(), "trace", "trace-123")
    
    // Create and dispatch event with context
    event.Dispatcher().Dispatch(
        event.NewEventWithContext(ctx, &UserCreated{
            UserId:   1,
            UserName: username,
        })
    )
}
```

### 6. Error Handling

```go
func (*UserCreatedErrListener) EventHandler() event.Handler {
    return func(e *event.Event) error {
        // Returned errors will be automatically logged
        return fmt.Errorf("error processing user creation event")
    }
}
```

### 7. Custom event logger
Customized evemt logger should implement interface event.ILogger 
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
type MyCustomLogger struct {
    logger *log.Logger //  your project's own logger
}

func (l *MyCustomLogger) Info(ctx context.Context, msg string, params ...event.LogParam) {
    l.logger.WithContext(ctx).Printf("[INFO] %s, %s", msg, formatParams(params))
}
// Implement all methods in event.ILogger
// .....

// Set event logger
event.SetLogger(&MyCustomLogger{logger: logger})

// Example for parse log params
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
## Complete Examples

Find complete working examples in the [examples directory](https://github.com/kevinyan815/event/tree/master/examples).

## License

This project is licensed under the Apache License 2.0 License - see the [LICENSE](https://github.com/kevinyan815/event/blob/master/LICENSE) file for details.
