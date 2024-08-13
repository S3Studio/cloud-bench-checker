// Package definition:
// Definition of conf file in yaml format
//
// See github.com/s3studio/cloud-bench-checker/doc/Baseline.md for details
package definition

// CloudType: Cloud type, aka connector type
type CloudType string

const (
	TENCENT_CLOUD CloudType = "tencent_cloud"
	TENCENT_COS   CloudType = "tencent_cos"
	ALIYUN_CLOUD  CloudType = "aliyun"
	ALIYUN_OSS    CloudType = "aliyun_oss"
	K8S           CloudType = "k8s"
	AZURE         CloudType = "azure"
)

type ParamType string

const (
	PARAM_INT         ParamType = "int"
	PARAM_STRING      ParamType = "string"
	PARAM_STRING_LIST ParamType = "string_list"
)

type PaginationType int

const (
	PAGEINATION_DEFAULT PaginationType = 0
	PAGE_OFFSET_LIMIT   PaginationType = 1
	PAGE_CURPAGE_SIZE   PaginationType = 2
	PAGE_NOPAGEINATION  PaginationType = 3
	PAGE_MARKER         PaginationType = 4
)

type OutputFormat string

const (
	OUTPUT_FORMAT_CSV  OutputFormat = "csv"
	OUTPUT_FORMAT_JSON OutputFormat = "json"
)

type ConfOption struct {
	PageSize       int          `yaml:"page_size"`
	OutputFormat   OutputFormat `yaml:"output_format"`
	OutputFilename string       `yaml:"output_filename"`
	OutputMetadata []string     `yaml:"output_metadata"`
	OutputRiskOnly bool         `yaml:"output_risk_only"`
	ServerHideYaml bool         `yaml:"server_hide_yaml"`
}

type ConfProfile map[string]string

const (
	PROFILE_ENV = "$ENV"
)

type ConfJsonPathCmd struct {
	Path string `yaml:"path"`
}

type ConfTencentCloudCmd struct {
	Service    string         `yaml:"service"`
	Version    string         `yaml:"version"`
	Action     string         `yaml:"action"`
	ExtraParam map[string]any `yaml:"extra_param"`
}

type ConfTencentCOSCmd struct {
	Service string `yaml:"service"`
	Action  string `yaml:"action"`
	// Use of ExtraParam for Tencent COS is not required at this time
}

type ConfAliyunCloudCmd struct {
	Endpoint           string         `yaml:"endpoint"`
	EndpointWithRegion bool           `yaml:"endpoint_with_region"`
	Version            string         `yaml:"version"`
	Action             string         `yaml:"action"`
	ExtraParam         map[string]any `yaml:"extra_param"`
}

type ConfAliyunOSSCmd struct {
	Action string `yaml:"action"`
}

type ConfK8sListCmd struct {
	Namespace   string         `yaml:"namespace"`
	Group       string         `yaml:"group"`
	Version     string         `yaml:"version"`
	Resource    string         `yaml:"resource"`
	ListOptions map[string]any `yaml:"list_options"`
}

type ConfAzureCmd struct {
	Provider string `yaml:"provider"`
	Version  string `yaml:"version"`
	RsType   string `yaml:"rs_type"` // Resource type
	Action   string `yaml:"action"`
}

type ConfListCmd struct {
	TencentCloud ConfTencentCloudCmd `yaml:"tencent_cloud"`
	TencentCOS   ConfTencentCOSCmd   `yaml:"tencent_cos"`
	Aliyun       ConfAliyunCloudCmd  `yaml:"aliyun"`
	AliyunOSS    ConfAliyunOSSCmd    `yaml:"aliyun_oss"`
	K8sList      ConfK8sListCmd      `yaml:"k8s_list"`
	Azure        ConfAzureCmd        `yaml:"azure"`

	DataListJsonPath    string `yaml:"data_list_json_path"`
	ConvertObjectToList bool   `yaml:"convert_object_to_list"`
}

type ConfPaginator struct {
	PaginationType PaginationType `yaml:"pagination_type"`

	OffsetType    ParamType `yaml:"offset_type"`
	OffsetName    string    `yaml:"offset_name"`
	LimitType     ParamType `yaml:"limit_type"`
	LimitName     string    `yaml:"limit_name"`
	RespTotalName string    `yaml:"resp_total_name"`

	MarkerName     string `yaml:"marker_name"`
	NextMarkerName string `yaml:"next_marker_name"`
	TruncatedName  string `yaml:"truncated_name"`
}

type ConfConstraintK8s struct {
	Version string `yaml:"version"`
}

type ConfConstraint struct {
	ConstraintK8s ConfConstraintK8s `yaml:"k8s"`
}

type ConfListor struct {
	Id         int            `yaml:"id"`
	CloudType  CloudType      `yaml:"cloud_type"`
	RsType     string         `yaml:"rs_type"` // Human readable resource type
	ListCmd    ConfListCmd    `yaml:"list_cmd"`
	Paginator  ConfPaginator  `yaml:"paginator"`
	Constraint ConfConstraint `yaml:"constraint"`
}

type ConfExtractCmd struct {
	IdJsonPath   string `yaml:"id_jsonpath"`
	NameJsonPath string `yaml:"name_jsonpath"`
	IdConst      string `yaml:"id_const"`
	NormalizeId  bool   `yaml:"normalize_id"`

	// Way to use id when extracting from cloud
	IdParamName string    `yaml:"id_param_name"`
	IdParamType ParamType `yaml:"id_param_type"`

	// Way to extract prop
	ExtractJsonPath ConfJsonPathCmd     `yaml:"extract_jsonpath"`
	TencentCloud    ConfTencentCloudCmd `yaml:"tencent_cloud"`
	TencentCOS      ConfTencentCOSCmd   `yaml:"tencent_cos"`
	Aliyun          ConfAliyunCloudCmd  `yaml:"aliyun"`
	AliyunOSS       ConfAliyunOSSCmd    `yaml:"aliyun_oss"`
	Azure           ConfAzureCmd        `yaml:"azure"`
	// Extract command for k8s is not required at this time

	// Way to extract prop using a list of commands as chain,
	// need to be checked before decommenting
	//
	// CmdChain []ConfExtractCmd `yaml:"cmd_chain"`
}

type ConfValidator struct {
	ValidateSchema   string            `yaml:"validate_schema"`
	DynValidateValue map[string]string `yaml:"dyn_validate_value"` // Modify value in validate_schema
	ValueJsonPath    string            `yaml:"value_jsonpath"`     // JsonPath to extract the actual value to be displayed
}

type ConfChecker struct {
	CloudType  CloudType      `yaml:"cloud_type"`
	Listor     []int          `yaml:"listor"`
	ExtractCmd ConfExtractCmd `yaml:"extract_cmd"`
	Validator  ConfValidator  `yaml:"validator"`
}

type ConfBaseline struct {
	Tag      []string          `yaml:"tag"`
	Metadata map[string]string `yaml:"metadata"`
	Checker  []ConfChecker     `yaml:"checker"`
}

type ConfFile struct {
	Option   ConfOption     `yaml:"option"`
	Profile  ConfProfile    `yaml:"profile"`
	Listor   []ConfListor   `yaml:"listor"`
	Baseline []ConfBaseline `yaml:"baseline"`
}
