# Baseline Config File Format
The baseline configuration file defines and describes benchmark baselines using the YAML format.

*Last updated: 2024-08-05*

---

## Option
Defines global options for running benchmark check.

### page_size
Defines the size of items on single page. Type: Integer

The value of this option is sent to the cloud via API, 
but may be accepted or ignored depending on the strategy of the cloud.

Default value: 50

> See [paginator](#paginator) for more details.

### output_format
Defines the format of the file to output the result.

Avaliable values:
* csv
* json

### output_filename
Defines the filename of the output containing the result.

### output_metadata
Defines fields to be extracted from the `metadata` properties of each baseline and output to the file.

If the specified field doesn't exist in the baseline, the output value for the field will be empty.

### output_risk_only
Defines how to filter and output the result. Type: Boolean

Avaliable values:
* true: Only output results with cloud resources in risk (failing the benchmark check)
* false: Output all cloud resources that have been checked by baseline, with no filter

---

## profile
Defines the name of profile for the cloud used in the benchmark check.

`$ENV` (to use environment values) and filename (to use a file under ".auth" directory) are valid values.

> See [Cloud authorization reference](./Auth.md) for more details.

The configuration of `profile` is a mapping with some of the following keys:
| Key | Cloud | Connector |
| - | - | - |
| aliyun | Aliyun | aliyun, aliyun_oss |
| azure | Azure | azure |
| k8s | Kubernetes | k8s |
| tencent | Tencent Cloud | tencent_cloud, tencent_cos |

The value of connector described above is used as the available value of
Listor (`listor.cloud_type`) and Checker (`baseline.checker.cloud_type`), as a reference to the profile.
For example, a Listor with `cloud_type` of `aliyun_oss` will load profile of `aliyun`
if it is set as the key of the `profile`.

For those clouds that require "region" in the API,
the value is defined in the profile definition binding to a specific region.
Multiple regions or all regions in a single profile are not supported currently.

---

## listor
Defines Listor that retrieve a list of resources from the cloud.

The configuration is a sequence that contains elements of mapping with the following properties:

### id
Identification of the Listor. Type: Integer

Needs to be unique for the Checker to reference, if not, the latter one with the same id will be omitted.

### cloud_type
Defines which cloud connector to use to retrieve data, and which profile to use for authorization.

> See [profile](#profile) for more details, and for avaliable values in the column of "Connector".

### rs_type
Defines human-friendly names of the resources retrieved by the Listor.

Does not affect the progress of the benchmark check.

### list_cmd
Defines how to list resource from the cloud.

Avaliable properties:

#### tencent_cloud
Defines how to list resource from Tencent cloud.

Avaliable properties:
| Key | Type | Description |
| - | - | - |
| service | string | "service" of API of Tencent cloud |
| version | string | "version" of API of Tencent cloud |
| action | string | "action" of API of Tencent cloud |
| extra_param | mapping | "actionParameters" of API of Tencent cloud |

#### tencent_cos
Defines how to list resource from Tencent COS.

The following properties are defined but **IGNORED in the configuration of Listor**:
| Key | Type | Description |
| - | - | - |
| service | string | "service" of API of Tencent COS |
| action | string | "action" of API of Tencent COS |

#### aliyun
Defines how to list resource from Aliyun.

Avaliable properties:
| Key | Type | Description |
| - | - | - |
| endpoint | string | "endpoint" of API of Aliyun |
| endpoint_with_region | boolean | Indicates whether to add the region string in the endpoint |
| version | string | "version" of API of Aliyun |
| action | string | "action" of API of Aliyun |
| extra_param | mapping | Data in query string of API of Aliyun |

*Example:*
1. Endpoint with region: `ecs.[region_id].aliyuncs.com`
1. Endpoint without region: `ims.aliyuncs.com`

*Note:*
1. `region_id` is defined in the profile.
1. Only "string" and "integer" values in extra_param are supported currently.

#### aliyun_oss
Defines how to list resource from Aliyun OSS.

The following properties are defined but **IGNORED in the configuration of Listor**:
| Key | Type | Description |
| - | - | - |
| action | string | "action" of API of Aliyun OSS |

#### k8s_list
Defines how to list resource from K8s.

Avaliable properties:
| Key | Type | Description |
| - | - | - |
| namespace | string | "namespace" of API of K8s |
| group | string | "group" of API of K8s |
| version | string | "version" of API of K8s |
| resource | string | "resource" of API of K8s |
| list_options | mapping | "ListOptions" of API of K8s |

If `group` and `version` are both empty,
the full qualified name of the resource will be searched automatically.
Therefore, it is easy to simply define `resource` of "pod" to list all the resources of "/apps/v1/pods".

`list_options` is useful for filtering resources on demand.

> See listor with id "1" (`rs_type`: pod_kube-apiserver) in
> [CIS_Kubernetes_Benchmark_v1.9.0.tmpl.conf](/template/CIS_Kubernetes_Benchmark_v1.9.0.tmpl.conf)
> to get an example using "component=kube-apiserver" as "labelSelector" with `list_options`.

#### azure
Defines how to list resource from Azure.

Avaliable properties:
| Key | Type | Description |
| - | - | - |
| provider | string | "namespace" of API of Azure |
| version | string | "version" of API of Azure |
| rs_type | string | "resource_type" of API of Azure |
| action | string | "action" of API of Azure, **IGNORED in the configuration of Listor** |

`rs_type` defined here is different from `rs_type` defined in `listor`.

A typical fully qualified endpoint for resources listing is:
```
/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.{provider}/{rs_type}?api-version={version}
```

"subscriptionId" is defined in the profile, and "resourceGroupName" can be omitted in many APIs.
`action` is used in `extract_cmd` described below.

#### data_list_json_path
Defines the JsonPath to get the list of resources from the result of API call.

Sometimes the value is not set and the default value for the corresponding cloud is used,
which is not listed in the document.

#### convert_object_to_list
Indicates whether the result of `data_list_json_path` is an object and should be put into a list for further use.
Type: Boolean

> See listor with id "2" (`rs_type`: password_policy) in 
> [CIS_Alibaba_Cloud_Foundation_Benchmark_v1.0.0.tmpl.conf](/template/CIS_Alibaba_Cloud_Foundation_Benchmark_v1.0.0.tmpl.conf)
> to get an example.

### paginator
Defines how to merge data list from multiple pages.

It is often the case that the data of resouces retrieved is too large to returned in a single response of API call,
so there are many different pagination types designed by the API provider, as follows:

> [i, j) below means starts with index i (inclusive) and ends with index j (exclusive)

* PAGE_OFFSET_LIMIT:
Returns items of [offset, offset + limit) in a single page, and "offset" starts with 0.
Total count of resources may or may not be returned in the response.

* PAGE_CURPAGE_SIZE:
Returns items of [(curpage - 1) * pagesize, curpage * pagesize) in a single page, and "curpage" starts with 1.
Total count of resources may or may not be returned in the response.

> See listor with id "7" (`rs_type`: rds) in
> [CIS_Alibaba_Cloud_Foundation_Benchmark_v1.0.0.tmpl.conf](/template/CIS_Alibaba_Cloud_Foundation_Benchmark_v1.0.0.tmpl.conf)
> to get an example.

* PAGE_MARKER:
Returns items with marker of empty string on 1st page.
If the result contains a valid value of next marker,
returns items on the next page if the next API call provides value of the marker.
Repeat until there is no more item to return and the value of next marker is set to empty. 

> See listor with id "2" (`rs_type`: disk) in
> [baseline.tmpl.conf](/template/baseline.tmpl.conf)
> to get an example.

* PAGE_NOPAGEINATION
No pagination control, and returns full list of resources in a single API call.

Avaliable properties for `paginator`:

#### pagination_type
Defines type of `pagination`. Type: Integer

Avaliable values as described above:
* 1: PAGE_OFFSET_LIMIT
* 2: PAGE_CURPAGE_SIZE
* 3: PAGE_NOPAGEINATION
* 4: PAGE_MARKER

There are some common uses of pagination of clouds.
If the value of `pagination_type` is not specific or is defined as "0",
the following default type and other related properties (not listed below) will be used:
| cloud_type | Default pagination type |
| - | - |
| tencent_cloud | PAGE_OFFSET_LIMIT |
| tencent_cos | PAGE_NOPAGEINATION |
| aliyun_oss | PAGE_MARKER |
| k8s | PAGE_NOPAGEINATION |
| azure | PAGE_MARKER |

#### offset_type
Defines type of "offset" parameter of API call.

"offset" stands for "offset" in PAGE_OFFSET_LIMIT, or "curpage" in PAGE_CURPAGE_SIZE.

Avaliable if: PAGE_OFFSET_LIMIT, PAGE_CURPAGE_SIZE

Avaliable values:
* int: Parse "offset" to integer
* string: Parse "offset" to string

#### offset_name
Defines name of "offset" parameter of API call.

"offset" stands for "offset" in PAGE_OFFSET_LIMIT, or "curpage" in PAGE_CURPAGE_SIZE.

Avaliable if: PAGE_OFFSET_LIMIT, PAGE_CURPAGE_SIZE

#### limit_type
Defines type of "limit" parameter of API call.

"limit" stands for list limit or list size of single page.

Avaliable if: PAGE_OFFSET_LIMIT, PAGE_CURPAGE_SIZE

Avaliable values:
* int: Parse "limit" to integer
* string: Parse "limit" to string

#### limit_name
Defines name of "limit" parameter of API call.

"limit" stands for list limit or list size of single page.

Avaliable if: PAGE_OFFSET_LIMIT, PAGE_CURPAGE_SIZE

#### resp_total_name
Defines name of dict key of total count of resources in the result returned by the API.

Avaliable if: PAGE_OFFSET_LIMIT, PAGE_CURPAGE_SIZE

#### marker_name
Defines name of "marker" parameter of API call.

Avaliable if: PAGE_MARKER

#### next_marker_name
Defines name of dict key of "next marker" returned by the API.

Avaliable if: PAGE_MARKER

#### truncated_name
Defines name of dict key of "truncated" returned by the API.

Some APIs return "next marker" as well as a value of "truncated" which is of type of boolean.
If "truncated" is set, "next marker" is available and can be used in the API call of the next page.
If set to "false" or missing, the result returned is the last page.

Avaliable if: PAGE_MARKER

---

## baseline
Defines Baseline that manage process of benchmark checking.

It is *recommended* that each Baseline corresponds to a single benchmark recommendation.

The configuration is a sequence that contains elements of mapping with the following properties:

### tag
Defines the tags of Baseline. Type: Sequence of string

The command tool uses it to filter Baseline on demond with the value provided by the "-t" argument,
making it easy for different customers, departments, etc. to use parts of one conf file that interest them.

### metadata
Defines the matadata of Baseline. Type: Mapping of string

Used to provide accordant value for the output. See [output_metadata](#output_metadata).

### checker
Defines Checker that extract required properties and validate that they meet the requirements of benchmark guidelines.

The configuration is a sequence that contains elements of mapping.

*Example:*
> 1. See Section of "4.2" in
>    [baseline.tmpl.conf](/template/baseline.tmpl.conf)
>    to get an example of one Baseline with multiple Checkers.
> 1. See Section of "7.2" in
>    ["CIS_Microsoft_Azure_Foundations_Benchmark_v2.1.0.tmpl.conf](/template/CIS_Microsoft_Azure_Foundations_Benchmark_v2.1.0.tmpl.conf)
>    to get an example of one Baseline with multiple Checkers of same `cloud_type` of `azure`.

Avaliable properties of Checker:

#### cloud_type
Defines which Listor and cloud connector to use to retrieve data, and which profile to use for authorization.

> See [profile](#profile) for more details, and for avaliable values in the column of "Connector".

#### listor
Defines the ids of Listor to get raw data of resources from. Type: Sequence of integer

The `cloud_type` of Listor must match the value of the `cloud_type` of Checker. See [cloud_type](#cloud_type)

> See Section of "5.7.4" in
> [CIS_Kubernetes_Benchmark_v1.9.0.tmpl.conf](/template/CIS_Kubernetes_Benchmark_v1.9.0.tmpl.conf)
> to get an example of one Checker with data merged together from multiple Listors.

#### extract_cmd
Defines how to extract required properties to be validated.

Avaliable properties:

* id_jsonpath

Defines JsonPath to extract id from previous data.

"id" will be used in the API call of cloud, and also be outputed to the result.

* name_jsonpath

Defines JsonPath to extract name from previous data.

"name" is a human-friendly string to identify the resource, and will be outputed to the result.

* id_const

Defines a static string for id of single resource.

`id_jsonpath` is omitted if `id_const` is defined.

> See Section of "1.7" in
> [CIS_Alibaba_Cloud_Foundation_Benchmark_v1.0.0.tmpl.conf](/template/CIS_Alibaba_Cloud_Foundation_Benchmark_v1.0.0.tmpl.conf)
> to get an example.

* normalize_id

Indicates whether the id should to be normalized
by outputing the last part of slices splitting by "/" to the result.
Type: Boolean

Useful to parse the full quartified id of Azure to a human-friendly value in the "id" column of the result,
other than the "name" column.
See `baseline.checker.extract_cmd.azure` for more details.

* id_param_name

Defines name of "id" parameter of API call.

Avaliable in `tencent_cloud` and `aliyun`.

* id_param_type

Defines type of "id" parameter of API call.

Avaliable in `tencent_cloud` and `aliyun`.

Avaliable values:
> * int: Parse "id" to integer
> * string: Parse "id" to string
> * string_list: Provide a list of string where the only element is "id"

* extract_jsonpath

Defines how to get data directly from the resource of Listor or previouse data.

The following `extract_cmd` methods (tencent_cloud, tencent_cos, etc.)
are all omitted if `extract_jsonpath` is defined.

Avaliable properties:
| Key | Type | Description |
| - | - | - |
| path | string | JsonPath to extract a sub element from previous data |

> See Section of "1.2.1" in
> [baseline.tmpl.conf](/template/baseline.tmpl.conf)
> to get an example.

* tencent_cloud

Defines how to get data from Tencent cloud.

Avaliable properties:
| Key | Type | Description |
| - | - | - |
| service | string | "service" of API of Tencent cloud |
| version | string | "version" of API of Tencent cloud |
| action | string | "action" of API of Tencent cloud |
| extra_param | mapping | "actionParameters" of API of Tencent cloud |

* tencent_cos

Defines how to get data from Tencent COS.

Avaliable properties:
| Key | Type | Description |
| - | - | - |
| service | string | "service" of API of Tencent COS |
| action | string | "action" of API of Tencent COS |

* aliyun

Defines how to get data from Aliyun.

Avaliable properties:
| Key | Type | Description |
| - | - | - |
| endpoint | string | "endpoint" of API of Aliyun |
| endpoint_with_region | boolean | Indicates whether to add the region string in the endpoint |
| version | string | "version" of API of Aliyun |
| action | string | "action" of API of Aliyun |
| extra_param | mapping | Data in query string of API of Aliyun |

*Example:*
1. Endpoint with region: `ecs.[region_id].aliyuncs.com`
1. Endpoint without region: `ims.aliyuncs.com`

*Note:*
1. `region_id` is defined in the profile.
1. Only "string" and "integer" values in extra_param are supported currently.

* aliyun_oss

Defines how to get data from Aliyun OSS.

Avaliable properties:
| Key | Type | Description |
| - | - | - |
| action | string | "action" of API of Aliyun OSS |

* azure

Defines how to get data from Azure.

The following properties are defined but **IGNORED in the configuration of Checker**
except **`action`**:
| Key | Type | Description |
| - | - | - |
| provider | string | "namespace" of API of Azure |
| version | string | "version" of API of Azure |
| rs_type | string | "resource_type" of API of Azure |
| action | string | "action" of API of Azure |

A typical fully qualified endpoint is:
```
/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.{provider}/{rs_type}/{rs_name}/{action}?api-version={version}
```

Normally the string before "action" is returned in full in the resource with dict key of "id".
So only the "id" and `action` (if defined) are needed to be combined to the endpoint by Checker,
and others are omitted.

#### validator
Defines how to validate the resource against the benchmark.

Avaliable properties:

* validate_schema

Defines JsonSchema to validate the property of resource.

To reduce false positives, the resource is only considered as "InRisk"
if the JsonSchema matches the property.
That is to say, if it is not sure whether the result would differ from the API documenation,
it is *recommended* to add a "required" schema to the target of the "object" type.

* dyn_validate_value

Defines dynamic value used in the `validate_schema`. Type: Mapping of string

The string with key of `dyn_validate_value` surrounded by "%"(which is "%key%") in the `validate_schema`
will be replaced to the value of the ralated key.
The replacement occurs before the validation of JsonSchema,
so the `validate_schema` required to match the format of JSON after the replacement
rather than before it (in the conf file).

It is useful if the acturl threshold is different for vary environments, and it is easier 
to modify the value of `dyn_validate_value` rather than to replace it manually in `validate_schema`.
More flexible uses of dynamic replacement are under consideration.

> See Section of "1.11" in
> [CIS_Alibaba_Cloud_Foundation_Benchmark_v1.0.0.tmpl.conf](/template/CIS_Alibaba_Cloud_Foundation_Benchmark_v1.0.0.tmpl.conf)
> to get an example.

* value_jsonpath

Defines how to get the actual value from the resource to output for further review.

It is useful to re-check the result with the actual value afterwards.
However, it is *NOT GRANTED* that the result of JsonPath will match the logic of JsonSchema.
It is the provider of conf file who is *responsible* for reducing misunderstanding of them.
