## Overview
Image service is one of services for social network. From original image (http link), it resizes to small images with predefine sizes, then stores in local server and upload to AWS S3 storage. This service also has open api for other services also.

Input is array of http images, technically, for resize and upload images:
+ New goroutine for resize and store local per image.
+ New goroutine for upload image 
+ Using buffer channel, wait group to communicate between goroutine

## System requirement
- Centos 7
- Go 1.13+
- [libvips](https://libvips.github.io/libvips/install.html)  8.8.0 (require for resize image)
## Build
#### 1.Locally
Go to root folder project build
```sh
go build -o resizeimage
```
Run
```sh
./resizeimage
```
#### 2.Server
Build on server Centos
```sh
PKG_CONFIG_PATH=/usr/local/lib/pkgconfig go build -o resizeimage
```
or setup config path and build like local.
```sh
export PKG_CONFIG_PATH=/usr/local/lib/pkgconfig
```
### Run service on server
+ Setup enviroment (golang and lipvips)
  Note: Setup libvips for centos7, please run script: ./setup-libvips.sh
+ Run LD_LIBRARY_PATH=/usr/local/lib ./resizeimage in detach mode

### Errors
When use s3 storage, may be have error 
    ```
     failed: to s3://bucketname/folder/ A client error    (RequestTimeTooSkewed) occurred when calling the UploadPart operation: The difference between the request time and the current time is too large
     ```
```sh
yum install ntp
```
Edit file /etc/ntp.conf and change to 
```sh
  server 0.amazon.pool.ntp.org iburst
  server 1.amazon.pool.ntp.org iburst
  server 2.amazon.pool.ntp.org iburst
  server 3.amazon.pool.ntp.org iburst
```
```sh
systemctl restart ntpd
```