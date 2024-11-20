package sqlparse

import (
	"github.com/violet-eva-01/starrocks/util"
	"regexp"
	"sort"
	"strings"
)

const (
	extractTime = 96 + iota
	with
	from
	join
	Select
	Alter
	Insert
	Create
	Drop
	Update
	Delete
	Truncate
)

func (p *Parse) assignment(code int, tbl []string, isExists bool) {

	tbl = addSpecialCharacters(util.RemoveRepeatElementAndToLower(tbl))

	if isExists {
		tbl = util.RemoveMatchElement(tbl, p.excludeSign)
	}

	switch code {
	case extractTime:
		p.extractTime = tbl
	case with:
		p.withTablaName = tbl
	case from:
		p.fromTableName = tbl
	case join:
		p.joinTableName = tbl
	case Select:
		otherArr := append(p.extractTime, p.excludeTables...)
		otherArr = append(otherArr, p.withTablaName...)
		otherArr = append(otherArr, p.DeleteTableName...)
		sort.Strings(otherArr)
		tbl = util.RemoveCoincideElement(append(p.fromTableName, p.joinTableName...), otherArr, false)
		p.SelectTableName = tbl
	case Insert:
		p.InsertTableName = tbl
	case Create:
		p.CreatTableName = tbl
	case Drop:
		p.DropTableName = tbl
	case Alter:
		p.AlterTableName = tbl
	case Delete:
		p.DeleteTableName = tbl
	case Update:
		p.UpdateTableName = tbl
	case Truncate:
		p.TruncateTableName = tbl
	default:
		return
	}

	p.assign(code, tbl)
}

func addSpecialCharacters(list []string) (result []string) {
	for _, v := range list {
		compile, _ := regexp.Compile("(`[^`]+`|[a-z0-9_]+)")
		allString := compile.FindAllString(v, -1)
		var tmpArr []string
		for _, str := range allString {
			if !strings.Contains(str, "`") {
				str = "`" + str + "`"
			}
			tmpArr = append(tmpArr, str)
		}
		v = strings.Join(tmpArr, ".")
		result = append(result, v)
	}
	return
}

func (p *Parse) assign(code int, tbl []string) {
	switch code {
	case Create, Drop, Insert, Alter, Delete, Update, Truncate, Select:
	default:
		return
	}
	i := 0
	for _, table := range tbl {
		var t Table
		compile, _ := regexp.Compile("(`[^`]+`|[a-z0-9_]+)")
		allString := compile.FindAllString(table, -1)
		length := len(allString)
		switch length {
		case 1:
			if p.DbName == "" {
				p.ErrorTables = append(p.ErrorTables, table)
				continue
			}
			t.TableName = strings.ReplaceAll(allString[0], "`", "")
			t.DbName = strings.ReplaceAll(p.DbName, "`", "")
			t.Catalog = strings.ReplaceAll(p.Catalog, "`", "")
		case 2:
			t.TableName = strings.ReplaceAll(allString[1], "`", "")
			t.DbName = strings.ReplaceAll(allString[0], "`", "")
			t.Catalog = strings.ReplaceAll(p.Catalog, "`", "")
		case 3:
			t.TableName = strings.ReplaceAll(allString[2], "`", "")
			t.DbName = strings.ReplaceAll(allString[1], "`", "")
			t.Catalog = strings.ReplaceAll(allString[0], "`", "")
		default:
			p.ErrorTables = append(p.ErrorTables, table)
			continue
		}
		t.Action = code
		t.Index = i
		p.ParseTables = append(p.ParseTables, t)
		i += 1
	}
}
