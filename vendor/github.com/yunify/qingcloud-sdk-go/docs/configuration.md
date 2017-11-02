# Configuration Guide

## Summary

This SDK uses a structure called "Config" to store and manage configuration, read ["config/config.go"](https://github.com/yunify/qingcloud-sdk-go/blob/master/config/config.go) file's comments of public functions for more information.

Except for AccessKeyID and SecretAccessKey, you can also configure the API servers for private cloud usage scenario. All available configureable items are list in default configuration file.

___Default Configuration File:___

``` yaml
# QingCloud services configuration

qy_access_key_id: 'ACCESS_KEY_ID'
qy_secret_access_key: 'SECRET_ACCESS_KEY'

host: 'api.qingcloud.com'
port: 443
protocol: 'https'
uri: '/iaas'
connection_retries: 3

# Valid log levels are "debug", "info", "warn", "error", and "fatal".
log_level: 'warn'

```

## Usage

Just create a config structure instance with your API AccessKey, and initialize services that you need using Init() function of target service.

### Code Snippet

Create default configuration

``` go
defaultConfig, _ := config.NewDefault()
```

Create configuration from AccessKey

``` go
configuration, _ := config.New("ACCESS_KEY_ID", "SECRET_ACCESS_KEY")

anotherConfiguration := config.NewDefault()
anotherConfiguration.AccessKeyID = "ACCESS_KEY_ID"
anotherConfiguration.SecretAccessKey = "SECRET_ACCESS_KEY"
```

Load user configuration

``` go
userConfig, _ := config.NewDefault().LoadUserConfig()
```

Load configuration from config file

``` go
configFromFile, _ := config.NewDefault().LoadConfigFromFilepath("PATH/TO/FILE")
```

Change API server

``` go
moreConfiguration, _ := config.NewDefault()

moreConfiguration.Protocol = "https",
moreConfiguration.Host = "api.private.com",
moreConfiguration.Port = 4433,
moreConfiguration.URI = "/iaas",
```
