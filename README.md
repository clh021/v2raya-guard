# v2raya-guard
v2raya 的 自动化守护

## 使用

编译，根据实际情况配置好环境变量，运行
```bash
username=username \
password=password \
serverbaseurl="http://192.168.1.2:2017/api/" \
cronExp="30 3,12,20 * * *" \
v2raya-guard
```

## Docker部署
