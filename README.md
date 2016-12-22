# www
bran.land server code

## Build/deploy instructions

On the build machine:
```shell
go generate ./data
go build
scp www bran.land:www
```

On the server:
```shell
sudo chown www:www www
sudo setcap 'cap_net_bind_service=+ep' www
sudo mv www /var/lib/www/www
sudo systemctl restart www
```
