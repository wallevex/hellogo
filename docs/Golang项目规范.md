# Golang项目规范

## Golang项目结构和编码规范

1. 项目结构规范参考 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
2. 高频JSON序列化反序列化使用 [easyjson](https://github.com/mailru/easyjson)
3. 配置文件使用TOML
4. 判断型函数统一命名为IsCorrectXXX
5. 描述性的修饰词放前面，数量/聚合/统计词放后面
6. 包名统一用小写格式，文件夹和文件名最好也全是小写（太长可以用下划线）
7. logger.Named()用点号分隔，比如logger.Named(repo.Student)、logger.Named(repo.Teacher)
8. 结构化日志的key用下划线格式
9. 只在有上下文信息的地方打错误日志（比如client id，trace id，请求路径），像repo这些底层函数不要打错误日志，而是包装日志信息返回出来
10. debug日志每个地方都可以适当添加，利于定位
11. 单元测试要用testify
12. 统一时区。要么全部存储为CST本地时区时间，要么全部存储为UTC国际标准时区。

## API设计规范

1. 参考 [Google API设计指南](https://cloud.google.com/apis/design?hl=zh-cn) 和 [API Improvement Proposals ](https://aip.bybutter.com/general)
2. API路径用小写和中划线，Params参数和JSON Body里的参数统一用小驼峰格式

## MySQL规范

1. MySQL表名建议使用单数形式，理由参考 [Table Naming Dilemma: Singular vs. Plural Names](https://stackoverflow.com/questions/338156/table-naming-dilemma-singular-vs-plural-names)
2. 不要使用外键约束
3. 不要直接进行物理删除，要使用delete_at进行软删除，然后再统一回收
4. 使用sqlx+squirrel，尽量不要用orm
5. 如果字段较多可以用orm生成表结构

## 部署规范
1. Debug端口统一在HTTP端口上加上10000
