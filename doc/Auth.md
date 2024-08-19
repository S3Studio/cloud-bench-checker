# Prepare Authorization Info for Clouds

## Location for authorization information
The authorization information is stored either in files or environment variables, depending on the baseline configuration file.

* If the value in `profile_name` for a cloud is `$ENV`, auth info is stored in default location, which usually refers to environment variables.
    * If the cloud type is 'k8s', the default location for auth info is `${HOME}/.kube/config`.

* Otherwise, the file named the same as value in `profile_name` is loaded from the `.auth` directory for the auth info.
    * [properties](https://docs.oracle.com/cd/E23095_01/Platform.93/ATGProgGuide/html/s0204propertiesfileformat01.html) format is used in many cases.
    * If the cloud type is 'k8s', the format of the file is the same as [kubeconfig](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/).
It is recommended to export from a valid kubeconfig using the following command:
```sh
kubectl config view --raw > ./.auth/file_name
```

## Available keys
The mentioned keys are applicable for both environment variables and files in properties format.

### Tencent cloud
The following keys are available as mentioned [here](https://cloud.tencent.com/document/sdk/Go#dc4aa78b-5240-403f-b68d-da41afef2a14):
* TENCENTCLOUD_SECRET_ID
* TENCENTCLOUD_SECRET_KEY
* TENCENTCLOUD_REGION

### Aliyun cloud
The following keys are available as mentioned [here](https://next.api.aliyun.com/api-tools/sdk/Ecs?version=2014-05-26&language=go-tea&tab=primer-doc#doc-full-code-demo):
* ALIBABA_CLOUD_ACCESS_KEY_ID
* ALIBABA_CLOUD_ACCESS_KEY_SECRET
* ALIBABA_CLOUD_REGION

### Azure
The following keys are available as mentioned [here](https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication):
* AZURE_CLIENT_ID
* AZURE_TENANT_ID
* AZURE_CLIENT_SECRET
* AZURE_SUBSCRIPTION_ID
