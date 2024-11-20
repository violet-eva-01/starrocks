package grant

import (
	"fmt"
	"testing"
)

func TestGrantParse(t *testing.T) {
	parse := NewParse("root", "default_catalog", "GRANT CREATE RESOURCE GROUP , CREATE RESOURCE , CREATE EXTERNAL CATALOG , REPOSITORY , BLACKLIST , FILE , OPERATE , CREATE STORAGE VOLUME  ON SYSTEM TO ROLE root")
	parse.Parse()
	fmt.Printf("%+v\n", parse.Authorize)
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
