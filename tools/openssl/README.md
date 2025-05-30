# OpenSSL自行生成私钥和证书
[Shing Light Productions OpenSSL下载](https://slproweb.com/products/Win32OpenSSL.html)

```shell
# 1)生成私钥（RSA 2048位），按照提示填写私钥加密密码
openssl genrsa -des3 -out server.key 2048
 
# 2)根据openssl.cnf配置文件生成csr文件
openssl req -new -key server.key -out server.csr -config openssl.cnf
 
# 3)备份加密后的私钥，然后将私钥解密
cp server.key server.key.org 
openssl rsa -in server.key.org -out server.key
 
# 4)生成crt文件，有效期1年（365天）
openssl x509 -req -in server.csr -signkey server.key -out server.crt -days 365 -extensions v3_req -extfile openssl.cnf
```
