package dao

import (
	"github.com/druidcaesa/gotool"
	"github.com/go-xorm/xorm"
	"blog-server/models"
	"blog-server/models/req"
	"blog-server/pkg/page"
)

type ConfigDao struct {
}

func (d ConfigDao) sql(session *xorm.Session) *xorm.Session {
	return session.Table("sys_config")
}

// GetConfigKey 根据键名查询参数配置信息
func (d ConfigDao) GetConfigKey(param string) *models.SysConfig {
	config := models.SysConfig{}
	_, err := d.sql(SqlDB.NewSession()).Where("config_key = ?", param).Get(&config)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	return &config
}

// List 分页查询数据
func (d ConfigDao) List(query req.ConfigQuery) (*[]models.SysConfig, int64) {
	configs := make([]models.SysConfig, 0)
	session := d.sql(SqlDB.NewSession())
	if gotool.StrUtils.HasNotEmpty(query.ConfigName) {
		session.And("config_name like concat('%', ?, '%')", query.ConfigName)
	}
	if gotool.StrUtils.HasNotEmpty(query.ConfigType) {
		session.And("config_type = ?", query.ConfigType)
	}
	if gotool.StrUtils.HasNotEmpty(query.ConfigKey) {
		session.And("config_key like concat('%', ?, '%')", query.ConfigKey)
	}
	if gotool.StrUtils.HasNotEmpty(query.BeginTime) {
		session.And("date_format(create_time,'%y%m%d') >= date_format(?,'%y%m%d')", query.BeginTime)
	}
	if gotool.StrUtils.HasNotEmpty(query.EndTime) {
		session.And("date_format(create_time,'%y%m%d') <= date_format(?,'%y%m%d')", query.EndTime)
	}
	total, _ := page.GetTotal(session.Clone())
	err := session.Limit(query.PageSize, page.StartSize(query.PageNum, query.PageSize)).Find(&configs)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil, 0
	}
	return &configs, total
}

// CheckConfigKeyUnique 校验是否存在
func (d ConfigDao) CheckConfigKeyUnique(config models.SysConfig) int64 {
	session := d.sql(SqlDB.NewSession())
	if config.ConfigId > 0 {
		session.Where("config_id != ?", config.ConfigId)
	}
	count, err := session.And("config_key = ?", config.ConfigKey).Cols("config_id").Count()
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return 0
	}
	return count
}

// Insert 添加数据
func (d ConfigDao) Insert(config models.SysConfig) int64 {
	session := SqlDB.NewSession()
	session.Begin()
	insert, err := session.Insert(&config)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		session.Rollback()
		return 0
	}
	session.Commit()
	return insert
}

// GetById 查询数据
func (d ConfigDao) GetById(id int64) *models.SysConfig {
	config := models.SysConfig{}
	session := d.sql(SqlDB.NewSession())
	_, err := session.Where("config_id = ?", id).Get(&config)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	return &config
}

// Update 修改数据
func (d ConfigDao) Update(config models.SysConfig) int64 {
	session := SqlDB.NewSession()
	session.Begin()
	update, err := session.Where("config_id = ?", config.ConfigId).Update(&config)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		session.Rollback()
		return 0
	}
	session.Commit()
	return update
}

// CheckConfigByIds 根据id集合查询
func (d ConfigDao) CheckConfigByIds(list []int64) *[]models.SysConfig {
	configs := make([]models.SysConfig, 0)
	err := d.sql(SqlDB.NewSession()).In("config_id", list).Find(&configs)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	return &configs
}

// Remove 删除数据
func (d ConfigDao) Remove(list []int64) bool {
	session := SqlDB.NewSession()
	session.Begin()
	_, err := session.In("config_id", list).Delete(&models.SysConfig{})
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		session.Rollback()
		return false
	}
	session.Commit()
	return true
}

// FindAll 查询所有数据
func (d ConfigDao) FindAll() *[]models.SysConfig {
	configs := make([]models.SysConfig, 0)
	session := SqlDB.NewSession()
	err := session.Find(&configs)
	if err != nil {
		gotool.Logs.ErrorLog().Println(err)
		return nil
	}
	return &configs
}

