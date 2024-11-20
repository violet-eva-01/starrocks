package grant

import (
	"fmt"
	"regexp"
	"strings"
)

type AuthorizeParse struct {
	UserIdentity string
	Catalog      string
	// AS
	// @Description: Authorization Statement
	AS        string
	as        string
	Authorize Authorize
}

type Authorize struct {
	Catalog string
	//PermissionsType string
	Permissions   []string
	ObjectType    string
	ObjectName    string
	ObjectDBName  string
	ObjectTBLName string
	GranteeType   string
	GranteeName   string
	IP            string
}

func NewParse(userIdentity string, catalog string, as string) *AuthorizeParse {
	compile := regexp.MustCompile("(?i)^GRANT")
	return &AuthorizeParse{
		UserIdentity: userIdentity,
		Catalog:      catalog,
		AS:           as,
		as:           compile.ReplaceAllString(as, ""),
		Authorize:    Authorize{Catalog: catalog},
	}
}

func (ap *AuthorizeParse) Parse() {
	ap.getGrantee()
	ap.getObject()
	ap.getPermissions()
}

func (ap *AuthorizeParse) getGrantee() {

	compile := regexp.MustCompile("(?i) TO .*$")
	ap.as = compile.ReplaceAllString(ap.as, "")

	split := strings.Split(strings.TrimSpace(strings.ReplaceAll(ap.UserIdentity, "'", "")), "@")
	switch len(split) {
	case 1:
		ap.Authorize.GranteeName = split[0]
		ap.Authorize.GranteeType = "ROLE"
	case 2:
		ap.Authorize.GranteeName = split[0]
		ap.Authorize.IP = split[1]
		ap.Authorize.GranteeType = "USER"
	}
}

func (ap *AuthorizeParse) getObject() {
	compile := regexp.MustCompile("(?i) ON .*$")
	objectStr := compile.FindAllString(ap.as, -1)
	ap.as = compile.ReplaceAllString(ap.as, "")
	if len(objectStr) == 0 {
		ap.Authorize.ObjectType = "ROLE"
		return
	}

	for index, i := range permissionTypeNames {
		regexpStr := fmt.Sprintf("(?i)(ON %s|ON ALL %sS)", i, i)
		compile1 := regexp.MustCompile(regexpStr)
		if matchString := compile1.MatchString(objectStr[0]); matchString {
			ap.Authorize.ObjectType = i
			switch StarRocksPermissionType(index) {
			case Table, View, MaterializedView:
				if strings.Contains(objectStr[0], fmt.Sprintf("ON ALL %sS", i)) {
					ap.Authorize.ObjectTBLName = fmt.Sprintf("ALL %sS", i)
					all := strings.ReplaceAll(objectStr[0], fmt.Sprintf("ON ALL %sS", i), "")
					if strings.Contains(all, " ALL DATABASES ") {
						ap.Authorize.ObjectDBName = "ALL DATABASES"
					} else {
						ap.Authorize.ObjectDBName = strings.TrimSpace(strings.ReplaceAll(all, "IN DATABASE", ""))
					}
				} else {
					split := strings.Split(strings.TrimSpace(strings.ReplaceAll(objectStr[0], fmt.Sprintf("ON %s ", i), "")), ".")
					ap.Authorize.ObjectDBName = split[0]
					ap.Authorize.ObjectTBLName = split[1]
				}
			case Database:
				if strings.Contains(objectStr[0], fmt.Sprintf("ON ALL %sS", i)) {
					ap.Authorize.ObjectDBName = fmt.Sprintf("ALL %sS", i)
				} else {
					ap.Authorize.ObjectDBName = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(objectStr[0], fmt.Sprintf("ON %s ", i), ""), "'", ""))
				}
			case System:
				return
			default:
				if strings.Contains(objectStr[0], fmt.Sprintf("ON ALL %sS", i)) {
					ap.Authorize.ObjectName = fmt.Sprintf("ALL %sS", i)
				} else {
					ap.Authorize.ObjectName = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(objectStr[0], fmt.Sprintf("ON %s ", i), ""), "'", ""))
				}
			}
			return
		}
	}
}

func (ap *AuthorizeParse) getPermissions() {
	switch ap.Authorize.ObjectType {
	case "ROLE":
		split := strings.Split(strings.ReplaceAll(ap.as, "'", ""), ",")
		for _, i := range split {
			ap.Authorize.Permissions = append(ap.Authorize.Permissions, strings.TrimSpace(i))
		}
	default:
		for _, ptn := range permissionNames {
			regexpStr := fmt.Sprintf("(?i)\\s+%s(,)?", ptn)
			compile := regexp.MustCompile(regexpStr)
			if matchString := compile.MatchString(ap.as); matchString {
				ap.Authorize.Permissions = append(ap.Authorize.Permissions, ptn)
				ap.as = compile.ReplaceAllString(ap.as, " ")
			}
		}
	}
}
