# www
bran.land server code

## Build/deploy instructions

Requires Bazel to build.

On the build machine:
```shell
$ bazel build :www
$ scp www bran.land:www
```

On the server:
```shell
# chown www:www www
# setcap 'cap_net_bind_service=+ep' www
# mv www /home/www/www
# systemctl restart www
```
