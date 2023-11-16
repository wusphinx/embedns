# embedns
一个轻量型支持API动态更新的DNS Server(coredns+embed etcd)

# example 
1. 启动服务
```
go run main.go
```

2. 写入A记录
```
etcdctl put /coredns/helloworld/x1 '{"host":"192.168.1.8","ttl":60}'
```

3. 域名解析验证
```
❯ nslookup helloworld localhost
Server:		localhost
Address:	::1#53

Name:	helloworld
Address: 192.168.1.8
```