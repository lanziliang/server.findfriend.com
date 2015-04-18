package controllers

import (
	//"fmt"
	"github.com/revel/revel"
	"server.findfriend.com/app/models"
	"server.findfriend.com/utils"
	"strconv"
)

type Account struct {
	*revel.Controller
}

//获取用户信息，用户不存在返回{}
func (c Account) GetInfo() revel.Result {
	name := c.Params.Get("username")
	user, err := models.AccountModel.GetAccountInfoByName(name)
	if err != nil {
		data := make(map[string]interface{})
		return c.RenderJson(data)
	}

	return c.RenderJson(user)
}

//用户注册
func (c Account) Signup() revel.Result {
	var user models.Account
	user.Email = c.Params.Get("email")
	user.Pwd = c.Params.Get("password")
	user.Phone = c.Params.Get("phone")
	user.Name = c.Params.Get("username")
	user.Sex, _ = strconv.Atoi(c.Params.Get("usersex"))

	var res string
	err := models.AccountModel.Insert(user)
	if err != nil {
		res = "fail"
	} else {
		res = "ok"
	}
	return c.RenderText(res)
}

//用户上传头像
func (c Account) UploadAvator(id string) revel.Result {

	//获取用户id
	c.Validation.Required(id)
	if c.Validation.HasErrors() {
		revel.WARN.Println("Id is required!")
		return c.RenderText("fail")
	}

	//接收上传文件
	file, header, err := c.Request.FormFile("imgFile")
	if err != nil {
		revel.ERROR.Println(err)
		return c.RenderText("fail")
	}

	filename, err := utils.Upload(revel.BasePath+"/upload/avator/"+id+"/", file, header)
	if err != nil {
		return c.RenderText("fail")
	}

	//产生缩略图
	err = utils.Thumbnail(revel.BasePath+"\\upload\\avator\\"+id+"\\", filename, "thumb_"+filename, 2, 300, 0)
	if err != nil {
		return c.RenderText("fail")
	}

	userId, _ := strconv.Atoi(id)
	err = models.AccountModel.UpdateAvator(userId, filename)
	if err != nil {
		return c.RenderText("fail")
	}

	return c.RenderText("ok")
}

//获取个人中心 我发布和我参与的活动
func (c Account) GetPersonalActivity(id string) revel.Result {
	var nodata = make([]string, 0, 0)
	c.Validation.Required(id)
	if c.Validation.HasErrors() {
		revel.WARN.Println("Id is required!")
		return c.RenderJson(nodata)
	}

	_id, _ := strconv.Atoi(id)

	activities, err := models.AccountModel.GetActivityByAccountId(_id)
	if err != nil {
		return c.RenderJson(nodata)
	}
	if len(activities) == 0 {
		return c.RenderJson(nodata)
	}
	return c.RenderJson(activities)
}

//用户点击参与活动
func (c Account) JoinActivity(userid, activityid string) revel.Result {
	c.Validation.Required(userid)
	c.Validation.Required(activityid)
	if c.Validation.HasErrors() {
		revel.WARN.Println("Id is required!")
		return c.RenderText("fail")
	}

	uid, _ := strconv.Atoi(userid)
	aid, _ := strconv.Atoi(activityid)
	err := models.ActivityModel.AddJoin(aid, uid)
	if err != nil {
		return c.RenderText("fail")
	}

	return c.RenderText("ok")
}

//用户取消参与活动
func (c Account) CancelJoinActivity(userid, activityid string) revel.Result {
	c.Validation.Required(userid)
	c.Validation.Required(activityid)
	if c.Validation.HasErrors() {
		revel.WARN.Println("Id is required!")
		return c.RenderText("fail")
	}

	uid, _ := strconv.Atoi(userid)
	aid, _ := strconv.Atoi(activityid)
	affect, err := models.ActivityModel.DeleteJoin(aid, uid)
	if err != nil || affect == 0 {
		return c.RenderText("fail")
	}

	return c.RenderText("ok")
}

//用户点击喜欢活动
func (c Account) LikeActivity(userid, activityid string) revel.Result {
	c.Validation.Required(userid)
	c.Validation.Required(activityid)
	if c.Validation.HasErrors() {
		revel.WARN.Println("Id is required!")
		return c.RenderText("fail")
	}

	uid, _ := strconv.Atoi(userid)
	aid, _ := strconv.Atoi(activityid)
	err := models.ActivityModel.AddLike(aid, uid)
	if err != nil {
		return c.RenderText("fail")
	}

	return c.RenderText("ok")
}

//用户取消喜欢活动
func (c Account) CancelLikeActivity(userid, activityid string) revel.Result {
	c.Validation.Required(userid)
	c.Validation.Required(activityid)
	if c.Validation.HasErrors() {
		revel.WARN.Println("Id is required!")
		return c.RenderText("fail")
	}

	uid, _ := strconv.Atoi(userid)
	aid, _ := strconv.Atoi(activityid)
	affect, err := models.ActivityModel.DeleteLike(aid, uid)
	if err != nil || affect == 0 {
		return c.RenderText("fail")
	}

	return c.RenderText("ok")
}
