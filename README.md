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
  <a href="https://github.com/go-dawn/dawn/blob/master/.github/README_zh-CN.md">
      <img height="20px" src="https://img.shields.io/badge/CN-flag.svg?color=555555&style=flat&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAxMjAwIDgwMCIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiPg0KPHBhdGggZmlsbD0iI2RlMjkxMCIgZD0ibTAsMGgxMjAwdjgwMGgtMTIwMHoiLz4NCjxwYXRoIGZpbGw9IiNmZmRlMDAiIGQ9Im0tMTYuNTc5Niw5OS42MDA3bDIuMzY4Ni04LjEwMzItNi45NTMtNC43ODgzIDguNDM4Ni0uMjUxNCAyLjQwNTMtOC4wOTI0IDIuODQ2Nyw3Ljk0NzkgOC40Mzk2LS4yMTMxLTYuNjc5Miw1LjE2MzQgMi44MTA2LDcuOTYwNy02Ljk3NDctNC43NTY3LTYuNzAyNSw1LjEzMzF6IiB0cmFuc2Zvcm09Im1hdHJpeCg5LjkzMzUyIC4yNzc0NyAtLjI3NzQ3IDkuOTMzNTIgMzI0LjI5MjUgLTY5NS4yNDE1KSIvPg0KPHBhdGggZmlsbD0iI2ZmZGUwMCIgaWQ9InN0YXIiIGQ9Im0zNjUuODU1MiwzMzIuNjg5NWwyOC4zMDY4LDExLjM3NTcgMTkuNjcyMi0yMy4zMTcxLTIuMDcxNiwzMC40MzY3IDI4LjI1NDksMTEuNTA0LTI5LjU4NzIsNy40MzUyLTIuMjA5NywzMC40MjY5LTE2LjIxNDItMjUuODQxNS0yOS42MjA2LDcuMzAwOSAxOS41NjYyLTIzLjQwNjEtMTYuMDk2OC0yNS45MTQ4eiIvPg0KPGcgZmlsbD0iI2ZmZGUwMCI+DQo8cGF0aCBkPSJtNTE5LjA3NzksMTc5LjMxMjlsLTMwLjA1MzQtNS4yNDE4LTE0LjM5NDUsMjYuODk3Ni00LjMwMTctMzAuMjAyMy0zMC4wMjkzLTUuMzc4MSAyNy4zOTQ4LTEzLjQyNDItNC4xNjQ3LTMwLjIyMTUgMjEuMjMyNiwyMS45MDU3IDI3LjQ1NTQtMTMuMjk5OC0xNC4yNzIzLDI2Ljk2MjcgMjEuMTMzMSwyMi4wMDE3eiIvPg0KPHBhdGggZD0ibTQ1NS4yNTkyLDMxNS45Nzk1bDkuMzczNC0yOS4wMzE0LTI0LjYzMjUtMTcuOTk3OCAzMC41MDctLjA1NjYgOS41MDUtMjguOTg4NiA5LjQ4MSwyOC45OTY0IDMwLjUwNywuMDgxOC0yNC42NDc0LDE3Ljk3NzQgOS4zNDkzLDI5LjAzOTItMjQuNzE0LTE3Ljg4NTgtMjQuNzI4OCwxNy44NjUzeiIvPg0KPC9nPg0KPHVzZSB4bGluazpocmVmPSIjc3RhciIgdHJhbnNmb3JtPSJtYXRyaXgoLjk5ODYzIC4wNTIzNCAtLjA1MjM0IC45OTg2MyAxOS40MDAwNSAtMzAwLjUzNjgxKSIvPg0KPC9zdmc+DQo=">
    </a>
</p>

`Dawn` is an opinionated `web` framework that provides rapid development capabilities which on top of [fiber](https://github.com/gofiber/fiber). It provides basic services such as configuration, logging, `fiber` extension, `gorm` extension, and event system. 

The core idea of ​​Dawn is modularity. High-level business modules can invoke low-level modules, such as databases, cache and so on. Following the idea of ​​`DDD`, each module corresponds to a domain and can be easily converted into microservices.

Each module needs to implement its own two core methods of `Init` and `Boot`, and then register it in `Sloop`. General business modules need to implement its `RegisterRoutes` method to register routes and provide `http` services.

The modules should be based on the principle of not recreating the wheel, and directly provides the original structure and method of the dependent library.

The libraries currently used are
- [klog](https://github.com/kubernetes/klog)
- [viper](https://github.com/spf13/viper)
- [godotenv](https://github.com/joho/godotenv)
- [fiber](https://github.com/gofiber/fiber)
- [gorm](https://github.com/go-gorm/gorm)
- [go-redis](https://github.com/go-redis/redis)
- [validator](https://github.com/go-playground/validator)

# Notice
**This project is still under development, please do not use it in a production environment.**

# Why dawn?
Tribute to the first episode of one piece romance dawn. Let us set sail towards romance with the sloop.
