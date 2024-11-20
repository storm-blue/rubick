# Rubick

****
rubick是一个处理yaml文件的工具。起因是私有化部署的场景，需要做k8s集群资源的搬迁工作。
但是往往直接导出的yaml文件无法直接使用，需要进行一定的清洗工作。最初尝试写代码去清洗，
但不具备可复用性，于是写了这个工具。

## 工具使用

包含三个子工具：

export: 导出工具，将k8s资源导出为文件

modify: 修改工具，可以通过脚本修改目标yaml文件

exec: 根据提供的配置文件进行k8s yaml导出并清洗

```
[__kubeconfig__]
/root/.kube/config

[deployment]
test/*
prod/*
*/app-a

[service]
test/*
prod/*
*/app-a

[__scripts__]
# common scripts
DELETE(metadata.annotations.(kubectl.kubernetes.io/last-applied-configuration))
DELETE(metadata.creationTimestamp)
DELETE(metadata.resourceVersion)
DELETE(metadata.uid)
DELETE(status)

# deployment scripts
IF VALUE_OF(kind) == "Deployment" THEN PRINT(metadata.namespace)
IF VALUE_OF(kind) == "Deployment" && EXISTS(metadata.labels.(github.com/app)) THEN SET(metadata.name, VALUE_OF(metadata.labels.(github.com/app)))
IF VALUE_OF(kind) == "Deployment" && NOT_EXISTS(metadata.labels.(github.com/app)) THEN TRIM_PREFIX(metadata.name, "test-")


# service scripts
IF VALUE_OF(kind) == "Service" THEN DELETE(spec.clusterIP)
IF VALUE_OF(kind) == "Service" THEN DELETE(spec.clusterIPs)
IF VALUE_OF(kind) == "Service" && LENGTH_OF(spec.ports) > 1 THEN PRINT(metadata.name)
IF VALUE_OF(kind) == "Service" && LENGTH_OF(spec.ports) > 1 THEN REMOVE()
IF VALUE_OF(kind) == "Service" THEN SET(spec.ports[0].port, 80)
```

执行: ``rubick -h``可以查看提示

## 脚本语法

基本语法：

IF [ condition ] THEN [ action ]

condition是条件判断，action是要执行的动作，例如：

```
IF VALUE_OF(metadata.name) == "app-a" THEN SET(metadata.namespace, "test")
```

当然也可以没有condition,直接执行动作：

```
DELETE(metadata.namespace)
```

condition也支持嵌套：

```
IF (VALUE_OF(metadata.name) == "app" && EXISTS(metadata.labels.(github.com/app))) || LENGTH_OF(spec.ports) > 1 THEN ...
```

### 支持的condition方法

**VALUE_OF**

```
IF VALUE_OF(metadata.name) == "app" THEN ...
```

**LENGTH_OF**

```
IF LENGTH_OF(metadata.labels) < 5 THEN ...
```

**EXISTS**

```
IF EXISTS(metadata.labels.app-name) THEN ...
```

**NOT_EXISTS**

```
IF NOT_EXISTS(metadata.labels.app-name) THEN ...
```

**HAS_PREFIX**

```
IF HAS_PREFIX(metadata.labels.app-name, "dev-") THEN ...
```

**HAS_SUFFIX**

```
IF HAS_SUFFIX(metadata.labels.app-name, "-app") THEN ...
```

### 支持的action方法

**DELETE**

满足条件则删除对应的key：

```
IF ... THEN DELETE(metadata.labels.app-name)
```

**SET**

满足条件则设置目标的值：

```
IF ... THEN SET(metadata.labels.app-name, "app1")
```

也支持和VALUE_OF配合使用：

```
IF ... THEN SET(metadata.labels.app-name, VALUE_OF(metadata.name))
```

**REPLACE_PART**

满足条件则对目标值进行字符串替换操作：

```
IF ... THEN REPLACE_PART(metadata.labels.app-name, "app-a", "app-b")
```

也支持和VALUE_OF配合使用：

```
IF ... THEN SET(metadata.labels.app-name, VALUE_OF(metadata.name))
```

**TRIM_PREFIX**

满足条件则对目标值进行移除前缀操作：

```
IF ... THEN TRIM_PREFIX(metadata.labels.app-name, "dev-")
```

也支持和VALUE_OF配合使用：

```
IF ... THEN TRIM_PREFIX(metadata.labels.app-name, VALUE_OF(metadata.namespace))
```

**TRIM_SUFFIX**

满足条件则对目标值进行移除后缀操作：

```
IF ... THEN TRIM_SUFFIX(metadata.labels.app-name, "-app")
```

也支持和VALUE_OF配合使用：

```
IF ... THEN TRIM_SUFFIX(metadata.labels.app-name, VALUE_OF(metadata.namespace))
```

**PRINT**

会在控制台打印对应的值：

```
IF ... THEN PRINT(metadata.labels.app-name)
```

**REMOVE**

满足条件则移除当前YAML：

```
IF ... THEN REMOVE()
```

## Object Key语法

脚本很多地方需要指定一个key，这个key指向YAML文件对象的某个位置。例如下面的yaml:

```
a:
  b: v1
  c:
  - d: v2
    e: v3
  - d: v4
    e: v5
  - d: v6
    e: v7
```

其中：

```
VALUE_OF(a.b) = v1

VALUE_OF(a.c[0].d) = v2

VALUE_OF(a.c[1].e) = v5

LENGTH_OF(a.c) = 3
```

对于数组，还支持一些额外的操作：

```
# 将数组c所有的子对象的属性"d"设置为“v100”
SET(a.c[*].d, "v100")

# 将数组c最后一个子对象的属性"d"设置为“v100”
SET(a.c[+].d, "v100")

# 删除数组c最后一个子对象
DELETE(a.c[+])

# 将数组增加一个子对象，并且把它的属性"d"设置为“v100”
SET(a.c[++].d, "v100")

# 寻找数组c中属性“d”的值为"v4"的子对象，并且将它的属性“d”设置为“v100”
SET(a.c[d=v4].d, "v100")
```
