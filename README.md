### 生成激活码
```
code_generator 设备ID
```

### 自启动
- mihomo
> /etc/init.d/mymihomo
```
#!/bin/sh /etc/rc.common

START=99
USE_PROCD=1


start_service() {
        procd_open_instance
        procd_set_param command /usr/bin/mihomo
        procd_append_param command -d /etc/mihomo
        procd_set_param respawn
        procd_close_instance
}
```
- mylinux
>  /etc/init.d/mylinux
```
#!/bin/sh /etc/rc.common

START=99
USE_PROCD=1


start_service() {
        procd_open_instance
        procd_set_param command /usr/bin/mylinux
        procd_set_param respawn
        procd_close_instance
}
```