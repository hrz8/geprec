# geprec

## setup buf

```bash
buf config init
```

## setup module proto

```bash
buf lint
buf build
```

## generate

```bash
buf generate # add plugin
```

## deps

```bash
buf dep update
```

## in action

### http
```bash
curl -X POST http://localhost:3008/v1/greeter/hello -d '{"name": "Eji"}'
```

### grpc
```bash
grpcurl -plaintext -d '{"name": "Eji"}' localhost:3009 greeter.v1.GreeterService/SayHello
```
