package admin

import (
	"blog-server/models"
	"blog-server/models/req"
	"blog-server/pkg/excels"
	"blog-server/pkg/library/user_util"
	"blog-server/pkg/page"
	"blog-server/pkg/resp"
	"blog-server/service"
	"github.com/druidcaesa/gotool"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type RoleApi struct {
	roleService service.RoleService
	userService service.UserService
}

// Find 分页查询角色数据
func (a RoleApi) Find(c *gin.Context) {
	query := req.RoleQuery{}
	if c.BindQuery(&query) == nil {
		list, i := a.roleService.FindList(query)
		p := page.Page{
			List:  list,
			Total: i,
			Size:  query.PageSize,
		}
		c.JSON(200, resp.Success(p))
	} else {
		c.JSON(500, resp.ErrorResp("参数绑定异常"))
	}
}

// GetRoleId 根据人roleId查询角色数据
func (a RoleApi) GetRoleId(c *gin.Context) {
	param := c.Param("roleId")
	roleId, err := strconv.ParseInt(param, 10, 64)
	if err == nil {
		role := a.roleService.SelectRoleByRoleId(roleId)
		c.JSON(200, resp.Success(role))
	} else {
		c.JSON(500, resp.ErrorResp("参数绑定异常"))
	}
}

// Add 添加角色业务操作
func (a RoleApi) Add(c *gin.Context) {
	role := models.SysRole{}
	if c.BindJSON(&role) == nil {
		if a.roleService.CheckRoleNameUnique(role) > 0 {
			c.JSON(500, resp.ErrorResp(500, "新增角色'"+role.RoleName+"'失败，角色名称已存在"))
		} else if a.roleService.CheckRoleKeyUnique(role) > 0 {
			c.JSON(500, resp.ErrorResp(500, "新增角色'"+role.RoleName+"'失败，角色权限已存在"))
		}
		role.CreateBy = user_util.GetUserInfo(c).UserName
		if a.roleService.Add(role) > 0 {
			c.JSON(200, resp.Success("保存成功"))
		} else {
			c.JSON(500, resp.Success("保存失败"))
		}
	} else {
		c.JSON(500, resp.ErrorResp("参数绑定异常"))
	}
}

// Edit 修改角色
func (a RoleApi) Edit(c *gin.Context) {
	role := models.SysRole{}
	if c.BindJSON(&role) == nil {
		if a.roleService.CheckRoleNameUnique(role) > 0 {
			c.JSON(500, resp.ErrorResp(500, "修改角色'"+role.RoleName+"'失败，角色名称已存在"))
		} else if a.roleService.CheckRoleKeyUnique(role) > 0 {
			c.JSON(500, resp.ErrorResp(500, "修改角色'"+role.RoleName+"'失败，角色权限已存在"))
		}
		role.CreateBy = user_util.GetUserInfo(c).UserName
		role.CreateTime = time.Now()
		if a.roleService.Update(role) > 0 {
			c.JSON(200, resp.Success("修改成功"))
		} else {
			c.JSON(500, resp.Success("修改失败"))
		}
	} else {
		c.JSON(500, resp.ErrorResp("参数绑定异常"))
	}
}

// Delete 删除角色
func (a RoleApi) Delete(c *gin.Context) {
	param := c.Param("roleId")
	roleId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		c.JSON(500, resp.ErrorResp("参数绑定异常"))
		return
	}
	if a.roleService.Delete(roleId) > 0 {
		c.JSON(200, resp.Success("删除成功"))
	} else {
		c.JSON(500, resp.Success("删除失败"))
	}
}

// ChangeStatus 状态修改
func (a RoleApi) ChangeStatus(c *gin.Context) {
	body := models.SysRole{}
	if c.BindJSON(&body) != nil {
		c.JSON(500, resp.ErrorResp("参数绑定异常"))
		return
	}
	allowed, s := a.roleService.CheckRoleAllowed(body.RoleId)
	if !allowed {
		c.JSON(500, resp.ErrorResp(s))
		return
	}
	body.UpdateTime = time.Now()
	body.UpdateBy = user_util.GetUserInfo(c).UserName
	if a.roleService.UpdateRoleStatus(&body) > 0 {
		resp.OK(c, "修改成功")
	} else {
		resp.Error(c)
	}
}

// AllocatedList 查询已分配用户角色列表
func (a RoleApi) AllocatedList(c *gin.Context) {
	query := req.UserQuery{}
	if c.BindQuery(&query) != nil {
		resp.Error(c)
		return
	}
	list, i := a.userService.GetAllocatedList(query)
	resp.OK(c, page.Page{
		List:  list,
		Total: i,
	})
}

// UnallocatedList 查询未分配用户角色列表
func (a RoleApi) UnallocatedList(c *gin.Context) {
	query := req.UserQuery{}
	if c.BindQuery(&query) != nil {
		resp.Error(c)
		return
	}
	list, i := a.userService.GetUnallocatedList(query)
	resp.OK(c, page.Page{
		List:  list,
		Total: i,
	})
}

// CancelAuthUser 取消授权用户
func (a RoleApi) CancelAuthUser(c *gin.Context) {
	roleUser := models.SysUserRole{}
	if c.BindJSON(&roleUser) != nil {
		resp.Error(c)
		return
	}
	if a.roleService.DeleteAuthUser(roleUser) > 0 {
		resp.OK(c, "操作成功")
	} else {
		resp.Error(c, "操作失败")
	}
}

// UpdateAuthUserAll 批量选择用户授权
func (a RoleApi) UpdateAuthUserAll(c *gin.Context) {
	body := req.UserRoleBody{}
	if c.Bind(&body) != nil {
		resp.Error(c)
		return
	}
	if a.roleService.InsertAuthUsers(body) > 0 {
		resp.OK(c, "操作成功")
	} else {
		resp.Error(c, "操作失败")
	}
}

// Export 导出Excel
func (a RoleApi) Export(c *gin.Context) {
	query := req.RoleQuery{}
	items := make([]interface{}, 0)
	if c.BindQuery(&query) == nil {
		list, _ := a.roleService.FindList(query)
		for _, role := range list {
			items = append(items, *role)
		}
		_, file := excels.ExportExcel(items, "角色表")
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+gotool.IdUtils.IdUUIDToRan(false)+".xlsx")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("FileName", gotool.IdUtils.IdUUIDToRan(false)+".xlsx")
		file.Write(c.Writer)
	} else {
		c.JSON(200, resp.ErrorResp(500, "参数错误"))
	}
}
