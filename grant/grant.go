package grant

import (
	"github.com/violet-eva-01/starrocks/util"
	"strings"
)

type StarRocksPermission int

const (
	Grant StarRocksPermission = iota
	Node
	CreateResourceGroup
	CreateResource
	CreateExternalCatalog
	Plugin
	Repository
	BlackList
	File
	Operate
	CreateGlobalFunction
	CreateStorageVolume
	Usage
	All
	Impersonate
	Apply
	CreateMaskingPolicy
	CreateRowAccessPolicy
	CreatePipe
	CreateWarehouse
	Security
	CreateDatabase
	CreateTable
	CreateView
	CreateFunction
	CreateMaterializedView
	Refresh
	Select
	Alter
	Insert
	Create
	Drop
	Update
	Delete
)

var starRocksPermissionNames = []string{
	"GRANT",
	"NODE",
	"CREATE RESOURCE GROUP",
	"CREATE RESOURCE",
	"CREATE EXTERNAL CATALOG",
	"PLUGIN",
	"REPOSITORY",
	"BLACKLIST",
	"FILE",
	"OPERATE",
	"CREATE GLOBAL FUNCTION",
	"CREATE STORAGE VOLUME",
	"USAGE",
	"ALL",
	"IMPERSONATE",
	"APPLY",
	"CREATE MASKING POLICY",
	"CREATE ROW ACCESS POLICY",
	"CREATE PIPE",
	"CREATE WAREHOUSE",
	"SECURITY",
	"CREATE DATABASE",
	"CREATE TABLE",
	"CREATE VIEW",
	"CREATE FUNCTION",
	"CREATE MATERIALIZED VIEW",
	"REFRESH",
	"SELECT",
	"ALTER",
	"INSERT",
	"CREATE",
	"DROP",
	"UPDATE",
	"DELETE",
}

func ParsePermissionName(str string) StarRocksPermission {

	index := util.FindIndex(strings.ToUpper(str), starRocksPermissionNames)
	if index == -1 {
		return -1
	} else {
		return StarRocksPermission(index)
	}

}

func (sp StarRocksPermission) String() string {

	if sp >= Grant && sp <= Delete {
		return starRocksPermissionNames[sp]
	}

	return "nil"
}

func (sp StarRocksPermission) RegexpString() string {
	if sp >= Grant && sp <= Delete {
		return strings.ReplaceAll(starRocksPermissionNames[sp], " ", "\\s+")
	}
	return "nil"
}

type StarRocksPermissionType int

const (
	System StarRocksPermissionType = iota
	Warehouse
	ResourceGroup
	Resource
	User
	GlobalFunction
	Function
	Catalog
	StorageVolume
	MaskingPolicy
	RowAccessPolicy
	Database
	Table
	MaterializedView
	View
)

var starRocksPermissionTypeNames = []string{
	"SYSTEM",
	"WAREHOUSE",
	"RESOURCE GROUP",
	"RESOURCE",
	"USER",
	"GLOBAL FUNCTION",
	"FUNCTION",
	"CATALOG",
	"STORAGE VOLUME",
	"MASKING POLICY",
	"ROW ACCESS POLICY",
	"DATABASE",
	"TABLE",
	"MATERIALIZED VIEW",
	"VIEW",
}

func (spt StarRocksPermissionType) String() string {
	if spt >= System && spt <= View {
		return starRocksPermissionTypeNames[spt]
	}
	return "nil"
}

func (spt StarRocksPermissionType) RegexpString() string {
	if spt >= System && spt <= View {
		return strings.ReplaceAll(starRocksPermissionTypeNames[spt], " ", "\\s+")
	}
	return "nil"
}

func ParsePermissionTypeName(str string) StarRocksPermissionType {
	index := util.FindIndex(strings.ToUpper(str), starRocksPermissionTypeNames)
	if index == -1 {
		return -1
	} else {
		return StarRocksPermissionType(index)
	}
}
