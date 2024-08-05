# Cloud Bench Checker

Connect to multiple clouds such as public cloud or cloud native via public APIs, and perform security baseline checks according to benchmark recommendations.

## Feature
* :white_check_mark: Support for multiple clouds with parallel execution
* :white_check_mark: Support for switching from various authorization profiles
* :white_check_mark: Flexible baseline configuration in [YAML](https://yaml.org/) format
* :white_check_mark: Flexible configuration to extract required data from cloud response with the support of [JSONPath](https://goessner.net/articles/JsonPath/)
* :white_check_mark: Flexible result validation with the support of [JSON Schema](https://json-schema.org/)

## SECURITY DISCLAIMER
**ALWAYS** use the *READONLY* cloud authorizations (ak/sk/ClusterRole/etc...) to be configured in the project, and **NEVER** trust any rule provided by others, even if it is cloned or downloaded from this site.

## Quick start
### Install
1. Clone this repo.
1. Build as following:
```sh
go build ./bin/cmd/main.go
``` 

### Prepare cloud auth config
To conform to file of `baseline.tmpl.conf`, authorization information should be stored in environment variables.
An easy way to do this is by creating a file similar to this:
```
TENCENTCLOUD_SECRET_ID=xxx
TENCENTCLOUD_SECRET_KEY=xxx
TENCENTCLOUD_REGION=xxx
ALIBABA_CLOUD_ACCESS_KEY_ID=xxx
ALIBABA_CLOUD_ACCESS_KEY_SECRET=xxx
ALIBABA_CLOUD_REGION=xxx
AZURE_CLIENT_ID=xxx
AZURE_TENANT_ID=xxx
AZURE_CLIENT_SECRET=xxx
AZURE_SUBSCRIPTION_ID=xxx
```
And then export the file as environment variables using one of the following commands:

<details><summary>under linux</summary>

```sh
export $(cat ./env.txt)
```
</details>

<details><summary>under Windows with Powershell</summary>

```powershell
(Get-Content .\env.txt).ForEach({ $name, $value = $_ -Split "="; Set-Item -Path "env:$name" -Value $value })
```
</details>

### Run
To perform baseline checks with tag `test` in the file of `baseline.tmpl.conf`:
```sh
./main -t test -c ./template/baseline.tmpl.conf
```

## Further guide
Please see [documentation](doc).

## Roadmap
- [x] Framework
    - [x] listor
    - [x] checker
    - [x] baseline
    - [x] auth controller
- [ ] Connector
    - [ ] cloud connector
        - [x] tencent cloud
            - [x] tencent cos
        - [x] aliyun cloud
            - [x] aliyun oss
        - [x] k8s
            - [ ] version validator
        - [ ] aws
        - [x] azure ( :warning: beta version)
        - [ ] maybe openstack?
        - [ ] support of multiple region
    - [ ] cross platform connector
        - [ ] api connector
- [ ] Versioning and compatibility for config file
- [ ] Interaction
    - [x] command tool
    - [ ] api
    - [ ] webui
- [ ] Tool
    - [x] baseline config manager: [project](example/baseline_manager)
    - [ ] building support
    - [ ] dockerize support
- [ ] Doc
    - [ ] usage
        - [x] [auth file format](doc/Auth.md)
    - [x] [baseline config format](doc/Baseline.md)
    - [ ] api
