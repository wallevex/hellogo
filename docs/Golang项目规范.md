# Golang项目规范

## Golang项目结构

1. 参考[https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)

## MySQL规范

1. 不要使用外键约束
2. 不要直接进行物理删除，要使用delete_at进行软删除，然后再统一回收

