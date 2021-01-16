# Schematic

A terraform style IaC (Infrastructure as Code) manager of services

## Declaration Language

### Basic Types (`B`)

* `String`
* `Integer`
* `Float`
* `Boolean`

```JSON
"StringValue" = "some string"
"IntegerValue" = 5821
"FloatValue" = 63.64f
"BooleanValue" = false
```

---

### Expressions (`E`)

```Typescript
<String> = <B | V<T>>
```

---

### Variables (`V<T>`)

```HCL
variable "someVar" {
    "value" = <T>
}
```

*or*

```Typescript
variable "anotherVar" = <B>
```

---

### Complex Types (`C<T>`)

* `Block<T>`
* `Array<T>`

#### Block

A composite object of fields and values:

```HCL
"String" {
    ..."String" = <B | C<T> | V>
}
```

#### Array

A list of values or composites

```HCL
"String" = [
    ...<B | C<T> | V>
]
```

---

## Schema Strucure

There are three top level components to a schema:

* Instance
* Data
* Capture

Instances are paramterized and filled instantiations of a template, data are a source that can be referenced and templates are skeletons that describe the structure of a custom instance.

---
  
### Instance

| Field        	| Type           	    |
|--------------	|----------------------	|
| name         	| `String`         	    |
| type         	| `String`         	    |
| hasDependency | `Array<String>` |
| count        	| `Integer`  	    |
| structure    	| `Block<InstBody>` 	    |

```HCL
instance "String<Name>" "String<Type>" {
    "hasDependency" = [
        "String<INSTANCE NAME>"
    ]
    "count" = Integer
    "structure" {
        <InstBody>
    }
}
```

Example of an instance

```HCL
variable "testcapsule_clsid" {
    "value" = 85831
}

instance "capsule::config" "test_capsule" {
    "hasDependency" = []
    "count" = 2
    structure {
        "inbuilt" = true
        "containerId" = "testContainer"
        "config" {
            "pidsMax" = 20
            "memMax" = 4096
            "netClsId" = var.testcapsule_clsid
            "terminateOnClose" = true
        }
    }
}
```

---

#### InstBody

There are two variants of the `InstBody`, one as a custom defintion via schema defintion and the second with inbuilt service support. The services are as follows:

* `capsule::config`
* `capsule::boxfile`
* `aws::s3`
* `aws::ec2`
* `aws::route53`
* `docker::config`
* `docker::dockerfile`

A custom schema is defined as a set of `"String" = <B | C<T> | V>`

The schema for the inbuilt support is:

| Field | Type |
|-------|------|
| inbuilt | `Boolean` |
| fieldSet | `{...<B> | C<T> | V>}` |

*Capsule (Config)*

```HCL
structure {
    "inbuilt" = true
    "containerId" = <String | V<String>>
    "config" {
        "pidsMax" = <Integer | V<Integer>>
        "memMax" = <Integer | V<Integer>>
        "netClsId" = <Integer | V<Integer>>
        "terminateOnClose" = <Boolean | V<Boolean>>
    }
}
```

*Capsule (BoxFile)*

```HCL
structure {
    "inbuilt" = true
    "boxFileLocation" = <String | V<String>>
}
```

*AWS (S3)*

```HCL
structure {
    "inbuilt" = true
    "instanceType" = "s3"
    "bucketName" = <String | V<String>>
    "arn" = <String | V<String>>
    "accessPolicy" {
        "public" = <"ALL" | "NONE" | "GROUP">
        "private" = <"ALL" | "NONE" | "GROUP" | "SELF">
    }
}
```

*Docker (Dockerfile)*

```HCL
structure {
    "inbuilt" = true
    "dockerfileLocation" = <String | V<String>>
}
```

An example of a Capsule structure is

```HCL
variable "testcon_clsid" {
    "value" = 85831
}

...

structure {
    "inbuilt" = true
    "containerId" = "testContainer"
    "config" {
        "pidsMax" = 20
        "memMax" = 4096
        "netClsId" = var.testcon_clsid
        "terminateOnClose" = true
    }
}
```

---

### Data

Data declarations are "intakes" for data for a specific existing source such as a running container or active S3 bucket

```HCL
data "<"file" | "service">" "String<NAME>" {
    "reference" = String<'L' | 'W'>::<String | V<String>
    "schema" {...<B | C<T> | V>}
}
```

```typescript
data.<String<TYPE>>.<String<NAME>>.<String<ATTRIUTE>>[.<String<ATTRIUTE>>]
```

```HCL
data "service" "service_manager_container" {
    "reference" = "L::/etc/capsule/containers/service_manager_container"
    "schema" {
        "containerId" = <String | V<String>>
        "proc" {
            "pidsMax" = <Integer | V<Integer>>
            "memMax" = <Integer | V<Integer>>
            "terminateOnClose" = <Boolean | V<Boolean>>
        }
        "network" {
            "netClsId" = <Integer | V<Integer>>
        }
    }
}

instance "test_inst_type" "example_inst" {
    "someProperty" = data.service.service_manager_container.network.netClsId
    ...
}
```

---

### Capture

Captures define the source code for instances

```HCL
capture "String<NAME>" {
    "source" = <String | V<String>>
    "hasDependency" = [
        ...<String | V<String>>
    ]
    "handler" = <String | V<String>>
}
```

For example

```HCL
variable "service_manager_net_id" {
    value = 2586425
}

capture "capsule" {
    "source" = "github.com/EngineersBox/Capsule"
    "hasDependency" = [
        "github.com/google/golang"
    ]
    "handler" = "capsule"
}

instance "capsule" "test_continer" {
    "hasDependency" = []
    "count" = 2
    structure {
        "inbuilt" = true
        "containerId" = "service_id_573"
        "config" {
            "pidsMax" = 20
            "memMax" = 4096
            "netClsId" = var.service_manager_net_id
            "terminateOnClose" = true
        }
    }
}
```
