// @author: violet-eva @date  : 2024/9/30 @notes :
package sqlparse

import (
	"fmt"
	"testing"
)

func TestGetTableNames(t *testing.T) {
	join := ",x as ("
	parse := NewParse(join, "", "violet", "")
	parse.getTableNames(With)
	for _, i := range parse.WithTablaName {
		fmt.Println(i)
	}
	for _, i := range parse.ErrorTables {
		fmt.Println(i)
	}
}

func TestParse_StmtClearAnnotation(t *testing.T) {
	//	querySQL := `, /*
	//力度:xxx
	//门店.门店.xxx
	//匹配率:啥啥啥
	//*/ dt as ( xxx )
	//`
	querySQL := `\nwith \nat_20_item as (select xxxxx)`
	//querySQL := "from aa_bb_cc.aa_bb_cc a , aa_bb_cc.aa_bb_cc b"
	//querySQL := "select distinct `ab_from` as `ac_from` from (select xxx,aaaa,bbbb) "
	//querySQL := "select distinct ab_from as ac_from from `zz`.`bb`.`aa`"
	newParse := NewParse(querySQL, "", "violet", "")
	newParse.StmtClearAnnotation()
	fmt.Println(newParse.Query)
	newParse.getTableNames(From)
	fmt.Println(newParse.FromTableName)

}
