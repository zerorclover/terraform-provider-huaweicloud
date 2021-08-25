module github.com/huaweicloud/terraform-provider-huaweicloud

go 1.14

require (
	github.com/hashicorp/errwrap v1.0.0
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/huaweicloud/golangsdk v0.0.0-20210406114125-0f07278d722c
	github.com/jen20/awspolicyequivalence v1.1.0
	github.com/smartystreets/goconvey v0.0.0-20190222223459-a17d461953aa // indirect
	github.com/stretchr/testify v1.4.0
	github.com/unknwon/com v1.0.1
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/huaweicloud/golangsdk v0.0.0-20210406114125-0f07278d722c => github.com/zerorclover/golangsdk v0.0.0-20210825072916-d7719366a85a
