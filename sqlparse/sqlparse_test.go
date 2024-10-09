// @author: violet-eva @date  : 2024/9/30 @notes :
package sqlparse

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetTableNames(t *testing.T) {
	joinSQL := `'690799282442" AND 6259=(SELECT COUNT(*) FROM SYSMASTER:SYSPAGHDR)--'`
	parse := NewParse(joinSQL, "", "violet", "")
	parse.QueryClearAnnotation()
	parse.DebugGetSelectTables()
}

func TestGetDBName(t *testing.T) {
	query := "use `aa`. `bb`;select * from cc"
	parse := NewParse(query, "", "violet", "test")
	parse.getUse(query)
	fmt.Println("1" + parse.Catalog + "2")
	fmt.Println("1" + parse.DbName + "2")
}

func TestParse_StmtClearAnnotation(t *testing.T) {
	//	querySQL := `, /*
	//力度:xxx
	//门店.门店.xxx
	//匹配率:啥啥啥
	//*/ dt as ( xxx )
	//`
	//querySQL := `\nwith \nat_20_item as (select xxxxx)`
	//querySQL := "from aa_bb_cc.aa_bb_cc a , aa_bb_cc.aa_bb_cc b"
	//querySQL := "select distinct `ab_from` as `ac_from` from (select xxx,aaaa,bbbb) "
	//querySQL := "select distinct ab_from as ac_from from `zz`.`bb`.`aa`"
	//newParse := NewParse(querySQL, "", "violet", "")
	//newParse.StmtClearAnnotation()
	//fmt.Println(newParse.Query)
	//newParse.getTableNames(from)
	//fmt.Println(newParse.fromTableName)
	//querySQL := `select * from aa where a like 'from aa' and b like 'a\'b' and c like '\\aaa'`
	querySQL := `select * from aa where '--' d like " a\" b -- " `
	//querySQL := "with ddl_总计  as () ,ttl_总计 as () select * from ddl_总计 "
	//querySQL := "use hive.alden;set catalog a;with `TMP` as (select * from `aa`) select * from tmp , a.`tmp1` , `b`.tmp2 join `c`.`tmp3`"
	parse := NewParse(querySQL, "", "violet", "")
	parse.InitAllUseTable()
	parse.DebugGetSelectTables()
	//parse.InitAllUseTable(true)
	//parse.DebugGetSelectTables()
	str := "aaa"
	sp := strings.Split(str, ";")
	fmt.Println(len(sp))
	for _, i := range parse.ParseTables {
		fmt.Println(i)
	}
}
