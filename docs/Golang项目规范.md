# Golang项目规范

## Golang项目结构

1. 项目结构规范参考 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
2. 高频JSON序列化反序列化使用 [easyjson](https://github.com/mailru/easyjson)

## API设计规范

1. 参考 [Google API设计指南](https://cloud.google.com/apis/design?hl=zh-cn) 和 [API Improvement Proposals ](https://aip.bybutter.com/general)

## MySQL规范

1. MySQL表名建议使用单数形式，理由参考 [Table Naming Dilemma: Singular vs. Plural Names](https://stackoverflow.com/questions/338156/table-naming-dilemma-singular-vs-plural-names)
1. 不要使用外键约束
2. 不要直接进行物理删除，要使用delete_at进行软删除，然后再统一回收
