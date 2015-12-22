# dto-layer-generator

[![Join the chat at https://gitter.im/francoishill/dto-layer-generator](https://badges.gitter.im/francoishill/dto-layer-generator.svg)](https://gitter.im/francoishill/dto-layer-generator?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
Generate the DTO (data transfer object) layer for your code - language agnostic by using plugins

## WORK IN PROGRESS!!!

**Both the code and the documentation is a work in progress.**

This is still a work in progress. The aim of this package is to have a DTO-definition that can be used to generate code for specific languages (supported by using plugins).

See the example YAML:
```
name: Employee
url: /employees
fields:
- name: Id
  type: int64
- name: Name
  type: string
- name: EmployeeNumber
  type: string
- name: Employer
  type: object
  fields:
  - name: Id
    type: int64
  - name: Name
    type: string
  - name: Department
    type: string
- name: Projects
  type: objectarray
  fields:
  - name: Id
    type: int64
  - name: Name
    type: string
```

This would generate the go struct:
```
type Employee struct {
	Id             int64
	Name           string
	EmployeeNumber string
	Employer       struct {
		Id         int64
		Name       string
		Department string
	}
	Projects []struct {
		Id   int64
		Name string
	}
}
```

The long-term plan is to generate code for the definition the same as protocol buffers. But instead actually be "human-readable". So you will be able to define this DTO and have for instance the java/golang backend and javascript front-end code generated to do the send/receive marshaling.

Because it will generate "type-safe" code and not use reflection, we will also get compile-time safety.