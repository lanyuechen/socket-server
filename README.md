### 运行
```
go run socket.go
```

### 编译
```
go build socket.go 

//交叉编译
GOOS=linux GOARCH=amd64 go build socket.go
```

|OS|ARCH|OS version|
|---|---|---|
| linux | 386 / amd64 / arm | >= Linux 2.6 |
| darwin | 386 / amd64 | OS X (Snow Leopard + Lion) |
| freebsd | 386 / amd64 | >= FreeBSD 7 |
| windows | 386 / amd64 | >= Windows 2000 |

### 启动（默认监听3000端口）
```
./socket 
```

### 使用方法(JS)
```
//创建socket
const s = new WebSocket('ws://localhost:3000/socket');
s.onmessage = function(res) {		//接收数据
	console.log(res.data);
}

//另外一个环境下创建其他socket
const s2 = new WebSocket('ws://localhost:3000/socket');
s2.send('something to send');
```