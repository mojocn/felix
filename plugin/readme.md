

## 1. build plugin

```bash
cd .
go build -buildmode=plugin -o=plugin_doctor.so plugin_bad_docter.go
ll;
file plugin_doctor.so
```

## 2. build user plugin example

```bash
 go build use_plugin_example.go
```

## 3. run example with plugin

```go
./use_plugin_example 
```