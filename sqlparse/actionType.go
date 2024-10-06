package sqlparse

import (
	"github.com/violet-eva-01/starrocks/util"
	"regexp"
	"strings"
)

const (
	ExtractTime = 96
	With        = 97
	From        = 98
	Join        = 99
	Select      = 100
	Insert      = 102
	Create      = 103
	Drop        = 104
)

func (p *Parse) assignment(code int, tbl []string) {

	tbl = util.RemoveRepeatElementAndToLower(tbl)

	switch code {
	case ExtractTime:
		p.ExtractTime = tbl
	case With:
		p.WithTablaName = tbl
	case From:
		p.FromTableName = tbl
	case Join:
		p.JoinTableName = tbl
	case Select:
		otherArr := append(p.ExtractTime, p.DirtyData...)
		tbl = util.RemoveCoincideElement(append(p.FromTableName, p.JoinTableName...), append(p.WithTablaName, otherArr...))
		p.SelectTableName = tbl
	case Insert:
		p.InsertTableName = tbl
	case Create:
		p.CreatTableName = tbl
	case Drop:
		p.DropTableName = tbl
	default:
		return
	}

	p.assign(code, tbl)
}

func (p *Parse) assign(code int, tbl []string) {
	switch code {
	case Create, Drop, Insert, Select:
		i := 0
		for _, table := range tbl {
			var t Table
			table = strings.ReplaceAll(table, " ", "")
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
	default:
		return
	}
}
