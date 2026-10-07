package fixtures

import "time"

type decodeAndValidateRequest struct {
	// BEWARE : the flag of URLParam should match the const string URLParam
	URLParam    string `json:"-" path:"url_param" validate:"numeric"`
	Text        string `json:"text" validate:"max=10"`
	DefaultInt  int    `json:"defaultInt" default:"10.0"` // MATCH /type mismatch between field type and default value type in default tag/
	DefaultInt2 int    `json:"defaultInt2" default:"10"`
	// MATCH:12 /unknown option "inline" in json tag/
	DefaultInt3      int             `json:"defaultInt2,inline" default:"11"` // MATCH /duplicated tag name "defaultInt2" in json tag/
	DefaultString    string          `json:"defaultString" default:"foo"`
	DefaultBool      bool            `json:"defaultBool" default:"trues"` // MATCH /type mismatch between field type and default value type in default tag/
	DefaultBool2     bool            `json:"defaultBool2" default:"true"`
	DefaultBool3     bool            `json:"defaultBool3" default:"false"`
	DefaultFloat     float64         `json:"defaultFloat" default:"f10.0"` // MATCH /type mismatch between field type and default value type in default tag/
	DefaultFloat2    float64         `json:"defaultFloat2" default:"10.0"`
	MandatoryStruct  mandatoryStruct `json:"mandatoryStruct" required:"trues"` // MATCH /required should be "true" or "false" in required tag/
	MandatoryStruct2 mandatoryStruct `json:"mandatoryStruct2" required:"true"`
	MandatoryStruct4 mandatoryStruct `json:"mandatoryStruct4" required:"false"`
	OptionalStruct   *optionalStruct `json:"optionalStruct,omitempty"`
	OptionalQuery    string          `json:"-" querystring:"queryfoo"`
	optionalQuery    string          `json:"-" querystring:"queryfoo"` // MATCH /tag on not-exported field optionalQuery/
	// No-reg test for bug https://github.com/revive-lint/revive/issues/208
	Tiret       string `json:"-,"`
	BadTiret    string `json:"other,"`               // MATCH /option can not be empty in json tag/
	ForOmitzero string `json:"forOmitZero,omitzero"` // MATCH /prior Go 1.24, option "omitzero" is unsupported in json tag/
	// MATCH:30 /option can not be empty in json tag/
	BadTiret string `json:"other,"` // MATCH /duplicated tag name "other" in json tag/
}

type RangeAllocation struct {
	metav1.TypeMeta   `json:",inline"` // MATCH /unknown option "inline" in json tag/
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Range             string `json:"range,flow"`  // MATCH /unknown option "flow" in json tag/
	Data              []byte `json:"data,inline"` // MATCH /unknown option "inline" in json tag/
}

type RangeAllocation struct {
	metav1.TypeMeta   `bson:",minsize"`
	metav1.ObjectMeta `bson:"metadata,omitempty"`
	Range             string `bson:"range,flow"` // MATCH /unknown option "flow" in bson tag/
	Data              []byte `bson:"data,inline"`
}

type TestContextSpecificTags2 struct {
	A       int       `asn1:"explicit,tag:1"`
	B       int       `asn1:"tag:2"`
	S       string    `asn1:"tag:0,utf8"`
	Ints    []int     `asn1:"set"`
	Version int       `asn1:"optional,explicit,default:0,tag:000"` // MATCH /duplicated tag number 0 in asn1 tag/
	Time    time.Time `asn1:"explicit,tag:4,other"`                // MATCH /unknown option "other" in asn1 tag/
	X       int       `asn1:"explicit,tag:invalid"`                // MATCH /tag must be a number but is "invalid" in asn1 tag/
}

type VirtualMachineRelocateSpecDiskLocator struct {
	DynamicData

	DiskId          int32                           `xml:"diskId,attr,cdata"`
	Datastore       ManagedObjectReference          `xml:"datastore,chardata,innerxml"`
	DiskMoveType    string                          `xml:"diskMoveType,omitempty,comment"`
	DiskBackingInfo BaseVirtualDeviceBackingInfo    `xml:"diskBackingInfo,omitempty,any"`
	Profile         []BaseVirtualMachineProfileSpec `xml:"profile,omitempty,other"` // MATCH /unknown option "other" in xml tag/
}

type TestDuplicatedxmlTags struct {
	A int `xml:"a"`
	B int `xml:"a"` // MATCH /duplicated tag name "a" in xml tag/
	C int `xml:"c"`
}

type TestDuplicatedbsonTags struct {
	A int `bson:"b"`
	B int `bson:"b"` // MATCH /duplicated tag name "b" in bson tag/
	C int `bson:"c"`
}

type TestDuplicatedYAMLTags struct {
	A int `yaml:"b"`
	B int `yaml:"c"`
	C int `yaml:"c"` // MATCH /duplicated tag name "c" in yaml tag/
}

type TestDuplicatedProtobufTags struct {
	A int `protobuf:"varint,name=b"`
	B int `protobuf:"varint,name=c"`
	C int `protobuf:"varint,name=c"` // MATCH /duplicated tag name "c" in protobuf tag/
}

type SpannerDuplicatedTags struct {
	A int `spanner:"field_a"`
	B int `spanner:"field_a"` // MATCH /duplicated tag name "field_a" in spanner tag/
	C int `spanner:"field_c"`
}

// test case from
// sigs.k8s.io/kustomize/api/types/helmchartargs.go

type HelmChartArgs struct {
	ChartName        string                 `json:"chartName,omitempty" yaml:"chartName,omitempty"`
	ChartVersion     string                 `json:"chartVersion,omitempty" yaml:"chartVersion,omitempty"`
	ChartRepoURL     string                 `json:"chartRepoUrl,omitempty" yaml:"chartRepoUrl,omitempty"`
	ChartHome        string                 `json:"chartHome,omitempty" yaml:"chartHome,omitempty"`
	ChartRepoName    string                 `json:"chartRepoName,omitempty" yaml:"chartRepoName,omitempty"`
	HelmBin          string                 `json:"helmBin,omitempty" yaml:"helmBin,omitempty"`
	HelmHome         string                 `json:"helmHome,omitempty" yaml:"helmHome,omitempty"`
	Values           string                 `json:"values,omitempty" yaml:"values,omitempty"`
	ValuesLocal      map[string]interface{} `json:"valuesLocal,omitempty" yaml:"valuesLocal,omitempty"`
	ValuesMerge      string                 `json:"valuesMerge,omitempty" yaml:"valuesMerge,omitempty"`
	ReleaseName      string                 `json:"releaseName,omitempty" yaml:"releaseName,omitempty"`
	ReleaseNamespace string                 `json:"releaseNamespace,omitempty" yaml:"releaseNamespace,omitempty"`
	ExtraArgs        []string               `json:"extraArgs,omitempty" yaml:"extraArgs,omitempty"`
}

// Test message for holding primitive types.
type Simple struct {
	OBool                *bool    `protobuf:"varint,1,req,json=oBool"`                           // MATCH /mandatory option "name" not found in protobuf tag/
	OInt32               *int32   `protobuf:"varint,2,opt,name=o_int32,jsonx=oInt32"`            // MATCH /unknown option "jsonx" in protobuf tag/
	OInt32Str            *int32   `protobuf:"varint,3,rep,name=o_int32_str,name=oInt32Str"`      // MATCH /duplicated option "name" in protobuf tag/
	OInt64               *int64   `protobuf:"varint,4,opt,json=oInt64,name=o_int64,json=oInt64"` // MATCH /duplicated option "json" in protobuf tag/
	OSint32Str           *int32   `protobuf:"zigzag32,11,opt,name=o_sint32_str,json=oSint32Str"`
	OSint64Str           *int64   `protobuf:"zigzag64,13,opt,name=o_sint32_str,json=oSint64Str"` // MATCH /duplicated tag name "o_sint32_str" in protobuf tag/
	OFloat               *float32 `protobuf:"fixed32,14,opt,name=o_float,json=oFloat"`
	ODouble              *float64 `protobuf:"fixed64,014,opt,name=o_double,json=oDouble"`      // MATCH /duplicated tag number 14 in protobuf tag/
	ODoubleStr           *float64 `protobuf:"fixed6,17,opt,name=o_double_str,json=oDoubleStr"` // MATCH /invalid tag name "fixed6" in protobuf tag/
	OString              *string  `protobuf:"bytes,18,opt,name=o_string,json=oString"`
	OString2             *string  `protobuf:"bytes,name=ameno"`
	OString3             *string  `protobuf:"bytes,name=ameno"` // MATCH /duplicated tag name "ameno" in protobuf tag/
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type RequestQueryOption struct {
	Properties           []string `url:"properties,comma,omitempty"`
	CustomProperties     []string `url:"-"`
	Associations         []string `url:"associations,brackets,omitempty"`
	Associations2        []string `url:"associations2,semicolon,omitempty"`
	Associations3        []string `url:"associations3,space,brackets,omitempty"`
	Associations4        []string `url:"associations4,numbered,omitempty"`
	Associations5        []string `url:"associations5,space,semicolon,omitempty"` // MATCH /can not set both "semicolon" and "space" as delimiters in url tag/
	PaginateAssociations bool     `url:"paginateAssociations,int,omitempty"`
	Archived             bool     `url:"archived,myURLOption"` // MATCH /unknown option "myURLOption" in url tag/
	IDProperty           string   `url:"idProperty,omitempty"`
}

type Fields struct {
	Field      string `datastore:",noindex,flatten,omitempty"`
	OtherField string `datastore:",unknownOption"` // MATCH /unknown option "unknownOption" in datastore tag/
}

type MapStruct struct {
	Field1     string `mapstructure:",squash,reminder,omitempty"`
	OtherField string `mapstructure:",unknownOption"` // MATCH /unknown option "unknownOption" in mapstructure tag/
}

type ValidateUser struct {
	Id          string  `validate:"omitempty,min=3,max=32"`
	Username    string  `validate:"required,min=3,max=32"`
	Email       string  `validate:"required,email"`
	Password    string  `validate:"required,min=8,max=32"`
	Biography   string  `validate:"min=0,max=1000"`
	DisplayName string  `validate:"displayName,min=3,max=32"` // MATCH /unknown option "displayName" in validate tag/
	Complex     string  `validate:"gt=0,dive,keys,eq=1|eq=2,endkeys,required"`
	BadComplex  string  `validate:"gt=0,keys,eq=1|eq=2,endkeys,required"`              // MATCH /option "keys" must follow a "dive" option in validate tag/
	BadComplex2 string  `validate:"gt=0,dive,eq=1|eq=2,endkeys,required"`              // MATCH /option "endkeys" without a previous "keys" option in validate tag/
	BadComplex3 string  `validate:"gt=0,dive,keys,eq=1|eq=2,endkeys,endkeys,required"` // MATCH /option "endkeys" without a previous "keys" option in validate tag/
	Issue1367   string  `validate:"required_without=ExternalValue,excluded_with=ExternalValue"`
	Handle      *string `validate:"omitnil,omitempty"`
	Score       int     `validate:"omitzero"`
}

type TomlUser struct {
	Username string `toml:"username,omitempty"`
	Location string `toml:"location,unknown"` // MATCH /unknown option "unknown" in toml tag/
}

type PropertiesTags struct {
	Field int               `properties:"-"`
	Field int               `properties:"myName"`
	Field int               `properties:"myName,default=15"`
	Field int               `properties:"myName,default=sString"` // MATCH /type mismatch between field type and default value type in properties tag/
	Field int               `properties:",default:15"`            // MATCH /unknown or malformed option "default:15" in properties tag/
	Field int               `properties:",default=15,default=2"`  // MATCH /duplicated option "default" in properties tag/
	Field time.Time         `properties:"date,layout=2006-01-02"`
	Field time.Time         `properties:",layout=2006-01-02"`
	Field time.Time         `properties:"date,layout"`            // MATCH /unknown or malformed option "layout" in properties tag/
	Field time.Time         `properties:"date,layout=  "`         // MATCH /option "layout" not of the form layout=value in properties tag/
	Field string            `properties:"date,layout=2006-01-02"` // MATCH /layout option is only applicable to fields of type time.Time in properties tag/
	Field []string          `properties:",default=a;b;c"`
	Field map[string]string `properties:"myName,omitempty"` // MATCH /unknown or malformed option "omitempty" in properties tag/
}

type SpannerUser struct {
	ID        int       `spanner:"user_id"`
	Name      string    `spanner:"full_name"`
	Email     string    `spanner:"-"` // Valid: ignore field
	CreatedAt time.Time `spanner:"created_at"`
	UpdatedAt time.Time `spanner:"updated_at,unknown"` // MATCH /unknown option "unknown" in spanner tag/
}

type Codec struct {
	_something struct{} `codec:",omitempty,int"` // MATCH /tag on not-exported field _something/
	_struct    struct{} `codec:",omitempty,int"` // do not match, _struct has special meaning for codec tag
	Field1     string   `codec:"-"`
	Field2     int      `codec:"myName"`
	Field3     int32    `codec:",omitempty"` // MATCH /redundant option "omitempty", already set for all fields in codec tag/
	Field4     bool     `codec:"f4,int"`     // MATCH /redundant option "int", already set for all fields in codec tag/
	field5     bool     // unexported, so skipped
	Anon
}

type TestDuplicatedCodecTags struct {
	_struct struct{} `json:",omitempty"` // MATCH /tag on not-exported field _struct/
	A       int      `codec:"field_a"`
	B       int      `codec:"field_a"` // MATCH /duplicated tag name "field_a" in codec tag/
	C       int      `codec:"field_c"`
}

type Cbor struct {
	RepeatedStr string `cbor:"errors"`
	Repeated    string `cbor:"1,keyasint"`
	Useless     string `cbor:"-,omitempty"`              // MATCH /useless option omitempty for ignored field in cbor tag/
	Inputs      string `cbor:",keyasint"`                // MATCH /tag name for option "keyasint" should be an integer in cbor tag/
	Outputs     string `cbor:"inputs,keyasint"`          // MATCH /tag name for option "keyasint" should be an integer in cbor tag/
	Errors      string `cbor:"errors,optempty,keyasint"` // MATCH /unknown option "optempty" in cbor tag/
	Inputs2     string `cbor:"-10,omitempty"`            // MATCH /integer tag names are only allowed in presence of "keyasint" option in cbor tag/
	Outputs2    string `cbor:"-12,toarray"`              // MATCH /tag name for option "toarray" should be empty in cbor tag/
	RepeatedInt string `cbor:"1,keyasint"`               // MATCH /duplicated integer key 1 in cbor tag/
	RepeatedStr string `cbor:"errors,omitempty"`         // MATCH /duplicated tag name errors in cbor tag/
	Useless     string `cbor:",toarray,omitempty"`       // MATCH /options "omitempty" and "omitzero" are ignored in presence of "toarray" option in cbor tag/
	Useless2    string `cbor:",omitzero,toarray"`        // MATCH /options "omitempty" and "omitzero" are ignored in presence of "toarray" option in cbor tag/
	// OK
	InputsOk   string `cbor:"8,keyasint"`
	OutputsOk  string `cbor:"-100,keyasint"`
	ErrorsOk   string `cbor:"-1,keyasint"`
	InputsOk2  string `cbor:"inputs,omitempty"`
	OutputsOk2 string `cbor:",toarray"`
}

// Tag numbers of asn1, cbor and protobuf live in separate namespaces,
// so reusing the same number across them is not a duplicate.
type CrossFormatTagNumbers struct {
	A int    `asn1:"tag:1"`
	B string `cbor:"1,keyasint"`
	C *int32 `protobuf:"varint,1,opt,name=c"`
}
