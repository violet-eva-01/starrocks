package grant

import (
	"fmt"
	"testing"
)

func TestGrantParse(t *testing.T) {
	parse := NewParse("root", "default_catalog", "GRANT CREATE RESOURCE GROUP , CREATE RESOURCE , CREATE EXTERNAL CATALOG , REPOSITORY , BLACKLIST , FILE , OPERATE , CREATE STORAGE VOLUME  ON SYSTEM TO ROLE root")
	parse.Parse()
	fmt.Printf("%+v\n", parse.Authorize)
	/*	rs := fmt.Sprintf("(?i)%s", StarRocksPermission(2).RegexpString())
		fmt.Println(rs)
		compile := regexp.MustCompile(rs)
		fmt.Println(parse.AS)
		allString := compile.FindAllString(parse.AS, -1)
		fmt.Printf("%+v\n", allString)
		rs1 := "(?i)CREATE\\s+RESOURCE\\s+GROUP"
		compile1 := regexp.MustCompile(rs1)
		fmt.Println(parse.AS)
		allString1 := compile1.FindAllString(parse.AS, -1)
		fmt.Printf("%+v\n", allString1)*/
}

func TestRoleGrantParse(t *testing.T) {
	parse := NewParse("'violet-eva'@'%'", "default_catalog", "GRANT 'root', 'db_admin', 'aldentest', 'user_admin' TO 'violet-eva'@'%'")
	parse.Parse()
	fmt.Printf("%+v\n", parse.Authorize)
}

func TestFuncGrantParse(t *testing.T) {
	parse := NewParse("'violet-eva'@'%'", "default_catalog", "GRANT usage ON GLOBAL FUNCTION a(string,int) TO 'violet-eva'@'%'")
	parse.Parse()
	fmt.Printf("%+v\n", parse.Authorize)
}

func TestRegexpString(t *testing.T) {
	permissionType := StarRocksPermissionType(13).RegexpString()
	fmt.Printf("%+v\n", permissionType)
}

func Test11(t *testing.T) {
	for index, str := range starRocksPermissionTypeNames {
		switch StarRocksPermissionType(index) {
		case System:
			fmt.Println(str)
		default:
			fmt.Println(str)
		}
	}
}
