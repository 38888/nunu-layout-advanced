package casbin

import (
	"strconv"
	"sync"

	"github.com/pkg/errors"

	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

type CasbinInfo struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type RBAC interface {
	Init() *casbin.CachedEnforcer
	Clear(v int, p ...string) bool
	List(roleId int64) (pathMaps []CasbinInfo)
	Add(roleId int64, casbinInfos []CasbinInfo) error
	Update(oldPath, oldMethod, newPath, newMethod string) error
	Check(roleId int64, method, path string) (bool, error)
}

type rbac struct {
	db *gorm.DB
}

var (
	cachedEnforcer *casbin.CachedEnforcer
	once           sync.Once
)

func New(db *gorm.DB) RBAC {
	return &rbac{
		db: db,
	}
}

// Init 持久化到数据库  引入自定义规则
func (c *rbac) Init() *casbin.CachedEnforcer {
	once.Do(func() {
		a, err := gormadapter.NewAdapterByDB(c.db)
		text := `
		[request_definition]
		r = sub, obj, act 
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
		m, err := model.NewModelFromString(text)
		if err != nil {

			return
		}
		cachedEnforcer, err = casbin.NewCachedEnforcer(m, a)
		cachedEnforcer.SetExpireTime(60 * 60)
		_ = cachedEnforcer.LoadPolicy()
	})
	return cachedEnforcer
}

// Clear 清除匹配的权限
func (c *rbac) Clear(v int, p ...string) bool {
	e := c.Init()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}
func (c *rbac) Add(roleId int64, list []CasbinInfo) error {
	roleIdStr := strconv.Itoa(int(roleId))
	c.Clear(0, roleIdStr)
	var rules [][]string
	for _, v := range list {
		rules = append(rules, []string{roleIdStr, v.Path, v.Method})
	}
	e := c.Init()
	success, _ := e.AddPolicies(rules)
	if !success {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	err := e.InvalidateCache()
	if err != nil {
		return err
	}
	return nil
}
func (c *rbac) Update(oldPath, oldMethod, newPath, newMethod string) error {
	err := c.db.Model(&gormadapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	e := c.Init()
	err = e.InvalidateCache()
	if err != nil {
		return err
	}
	return err
}

// List 获取权限列表
func (c *rbac) List(roleId int64) (pathMaps []CasbinInfo) {
	roleIdStr := strconv.Itoa(int(roleId))

	e := c.Init()
	list, _ := e.GetFilteredPolicy(0, roleIdStr)
	for _, v := range list {
		pathMaps = append(pathMaps, CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

func (c *rbac) Check(roleId int64, method, path string) (bool, error) {
	e := c.Init()
	roleIdStr := strconv.Itoa(int(roleId))
	success, err := e.Enforce(roleIdStr, path, method)
	if err != nil {
		return false, err
	}
	return success, nil
}
