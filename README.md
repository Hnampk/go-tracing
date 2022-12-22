# go-tracing

serves logs and prometheus metrics to help trace a function execution at runtime.

## fabric-flogging

`fabric-flogging` is cloned from flogging of Hyperledger Fabric at https://github.com/hyperledger/fabric/tree/master/common/flogging

## how to use

1. Define tracing

```go
var myTracer  = tracing.MustGetTracer("example_file")
```

2. Count on function call

call `StartFunction()` at the beginning of the function call. For example

```go
func myfunction(){
    myTracer.StartFunction("my_trace_number")
    // function logic
}
```

```log
2022-12-22 13:43:09.339 +07 [example_file] myfunction -> INFO 026 [my_trace_number] StartFunction at 1671691389339
```

3. Count on function end

call `EndFunction()` at the end of the function call. For example

```go
func myfunction(){

    // function logic
    myTracer.EndFunction("my_trace_number")
}
```

```log
2022-12-22 13:43:09.339 +07 [example_file] myfunction -> INFO 026 [my_trace_number] EndFunction at 1671691389339
```

4. Trace function duration

-   option 1: call `EndFunctionWithDurationSince()` at the end of the function call. For example

```go
func myfunction(){
    startTime := myTracer.StartFunction("my_trace_number")
    // function logic
    myTracer.EndFunctionWithDurationSince("my_trace_number", startTime)
}
```

-   option 2: call `EndFunctionWithDurationSince()` with `defer` at the beginning of function call. For example

```go
func myfunction(){
    defer myTracer.EndFunctionWithDurationSince("my_trace_number", time.Now())
    // function logic
}
```

both options will have the same logging like this

```log
2022-12-22 13:43:09.891 +07 [example_file] myfunction -> INFO 02f [my_trace_number] EndFunction at 1671691389891, duration=551ms
```

## metrics

check `URL:PORT/metrics` to find prometheus metrics
