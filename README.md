# DTO Layer Generator

[![Join the chat at https://gitter.im/zero-boilerplate/dto-layer-generator](https://badges.gitter.im/zero-boilerplate/dto-layer-generator.svg)](https://gitter.im/zero-boilerplate/dto-layer-generator?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
Generate the DTO (data transfer object) layer for your code - language agnostic by using plugins

# Support
[![Support via Gratipay](https://cdn.rawgit.com/gratipay/gratipay-badge/2.3.0/dist/gratipay.png)](https://gratipay.com/~FrancoisHill/)


```
WORK IN PROGRESS!!!

The long-term plan is to generate code for the definition the same as protocol buffers. But instead actually be "human-readable". So you will be able to define this DTO and have for instance the java/golang backend and javascript front-end code generated to do the send/receive marshaling.

Because it will generate "type-safe" code and not use reflection, we will also get compile-time safety.
```


# Documentation

## Placeholders

In the `example/out` dir you will see an example of how the placeholders are used in your project folder. Note that the format of placeholders are usually inside a commented line, like **//{{BEGIN `EMPLOYEE_PLACEHOLDER`}}** and has another line after that for ending **//{{END `EMPLOYEE_PLACEHOLDER`}}**. The reason for two lines is that after this generator "injected" the generated code it can still find where to replace/overwrite that section the next time it generates it.

## Example

Refer to the `example` folder in this git repo. You will see `simple_example.yml` file that is a demo setup of how your "DTO definition/setup" will look. The file `nested_example.yml` is one with nested objects in your DTO.

Windows:
```
go get github.com/zero-boilerplate/dto-layer-generator
cd "%gopath%/src/github.com/zero-boilerplate/dto-layer-generator"
"%GOPATH%/bin/dto-layer-generator.exe" "example/simple_example.yml"
:: "%GOPATH%/bin/dto-layer-generator.exe" "example/nested_example.yml"
```

Linux:
```
go get github.com/zero-boilerplate/dto-layer-generator
cd "$gopath/src/github.com/zero-boilerplate/dto-layer-generator"
"$GOPATH/bin/dto-layer-generator.exe" "example/simple_example.yml"
# "$GOPATH/bin/dto-layer-generator.exe" "example/nested_example.yml"
```
Now have look at the above mentioned files (in the **Placeholders** section) to see their structs/classes generated from the input `...example.yml` file.

## Supported types

Supported types (using the builtin types of golang, newly added plugins should have a convert map like in the Java plugin where it calls `ConvertTypeName`):
- `string`
- `bool`
- `byte`
- `float32`
- `float64`
- `int`
- `int8`
- `int16`
- `int32`
- `int64`
- `uint`
- `uint8`
- `uint16`
- `uint32`
- `uint64`

# Contributing

## Adding an open-source plugin

It is always the simplest to just look at an example. See the java example in dir `plugins/java.go`. But it boils down to that you should implement the interface `plugins/Plugin` and call the `RegisterPlugin()` method with the alias and your new plugin's instance. This registration is best done in the plugin's own file with the `init()` func - like in the java example.

After updating code locally you can just call the `go install` command in the `%gopath%/src/github.com/zero-boilerplate/dto-layer-generator` dir to regenerate the binary into the `%GOPATH%/bin` dir.