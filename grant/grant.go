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
	Truncate
)

var permissionNames = []string{
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

	index := util.FindIndex(strings.ToUpper(str), permissionNames)
	if index == -1 {
		return -1
	} else {
		return StarRocksPermission(index)
	}

}

func (sp StarRocksPermission) String() string {

	if sp >= CreateResourceGroup && sp <= Delete {
		return permissionNames[sp]
	}

	return "nil"
}

type StarRocksPermissionType int

const (
	System StarRocksPermissionType = iota
	ResourceGroup
	Resource
	User
	GlobalFunction
	Function
	Catalog
	StorageVolume
	Database
	Table
	MaterializedView
	View
)

var permissionTypeNames = []string{
	"SYSTEM",
	"RESOURCE GROUP",
	"RESOURCE",
	"USER",
	"GLOBAL FUNCTION",
	"FUNCTION",
	"CATALOG",
	"STORAGE VOLUME",
	"DATABASE",
	"TABLE",
	"MATERIALIZED VIEW",
	"VIEW",
}

func (spt StarRocksPermissionType) String() string {
	if spt >= System && spt <= View {
		return permissionTypeNames[spt]
	}
	return "nil"
}

func ParsePermissionTypeName(str string) StarRocksPermissionType {
	index := util.FindIndex(strings.ToUpper(str), permissionTypeNames)
	if index == -1 {
		return -1
	} else {
		return StarRocksPermissionType(index)
	}
}
