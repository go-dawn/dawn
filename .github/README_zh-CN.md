# Dawn
<p align="center">
  <a href="https://pkg.go.dev/github.com/go-dawn/dawn?tab=doc">
    <img src="https://img.shields.io/badge/%F0%9F%93%9A%20godoc-pkg-00ACD7.svg?color=00ACD7&style=flat">
  </a>
  <a href="https://goreportcard.com/report/github.com/go-dawn/dawn">
    <img src="https://img.shields.io/badge/%F0%9F%93%9D%20goreport-A%2B-75C46B">
  </a>

  <a href="https://codecov.io/gh/go-dawn/dawn">
    <img src="https://codecov.io/gh/go-dawn/dawn/branch/main/graph/badge.svg?token=3VA39G2KNI"/>
  </a>

  <a href="https://github.com/go-dawn/dawn/actions?query=workflow%3ASecurity">
    <img src="https://img.shields.io/github/workflow/status/go-dawn/dawn/Security?label=%F0%9F%94%91%20gosec&style=flat&color=75C46B">
  </a>
  <a href="https://github.com/go-dawn/dawn/actions?query=workflow%3ATest">
    <img src="https://img.shields.io/github/workflow/status/go-dawn/dawn/Test?label=%F0%9F%A7%AA%20tests&style=flat&color=75C46B">
  </a>
  <a>
    <img src="https://counter.gofiber.io/badge/go-dawn/dawn">
  </a>
</p>

`Dawn`是一个基于[fiber](https://github.com/gofiber/fiber)的，个性化的，轻量的，提供了快速开发能力的`web`框架。它提供了配置、日志、`fiber`扩展、`gorm`扩展、事件系统等基础服务。

`Dawn`的核心理念是模块化。高层的业务模块可以调用低层的基础模块，例如数据库等等。遵循`DDD`的思想，每一个模块对应一个领域，且都可以轻易地转换为微服务。

每个模块都需要实现自己的`Init`，`Boot`这两个核心方法，然后注册到`Sloop`中。一般业务模块需要实现其`RegisterRoutes`方法，用于注册路由，提供`http`服务。

模块的封装本着不重复造轮子的原则，直接提供依赖库其原本的结构和方法。

目前用到的库有
- [klog](https://github.com/kubernetes/klog)
- [viper](https://github.com/spf13/viper)
- [godotenv](https://github.com/joho/godotenv)
- [fiber](https://github.com/gofiber/fiber)
- [gorm](https://github.com/go-gorm/gorm)
- [go-redis](https://github.com/go-redis/redis)
- [validator](https://github.com/go-playground/validator)

# 注意
**本项目还在开发中，请勿在生产环境中使用。**

# 为什么是dawn
这是为了致敬海贼王第一集——`Romance Dawn`。让我们向着浪漫扬帆起航。
