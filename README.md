# 阿里云CDN自动HTTPS证书更新系统
## 介绍
此项目的出发点在于宝塔，宝塔会为网站定时更新证书。因此，我们可以设置一个定时任务来运行此程序，负责把宝塔更新写入好的证书和密钥文件上传到阿里云的CAD。

## 国内
处于国内加速的目的，api的endpoint基本都选择杭州节点（cn-hangzhou），你若有需要也可以按需秀嘎哎。

## 协议
本程序由[MIT](./LICENSE)协议发布。

## 宝塔SSL文件位置
证书夹：/www/server/panel/vhost/ssl
注意：证书夹里面似乎有冲突。例如域名通配符`*.example.com`和普通域名`example.com`会占用同一目录，导致覆盖。因此最好使用下面的Let's Encrypt文件夹。

Let's Encrypt：/www/server/panel/vhost/letsencrypt
