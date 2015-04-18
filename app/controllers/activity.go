package controllers

import (
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"server.findfriend.com/app/models"
	"server.findfriend.com/utils"
	"strconv"
	"time"
)

type Activity struct {
	*revel.Controller
}

//发布活动
func (c Activity) Publish(id, title, category, content string) revel.Result {
	c.Validation.Required(id)
	c.Validation.Required(title)
	c.Validation.Required(category)
	c.Validation.Required(content)
	if c.Validation.HasErrors() {
		revel.WARN.Println("All params is required!")
		return c.RenderText("fail")
	}

	var activity models.Activity
	activity.AccountId, _ = strconv.Atoi(c.Params.Get("id"))
	activity.Title = c.Params.Get("title")
	activity.Content = c.Params.Get("content")
	activity.Category, _ = strconv.Atoi(c.Params.Get("category"))

	//接收上传文件
	file, header, err := c.Request.FormFile("imgFile")
	if err != nil {
		revel.ERROR.Println(err)
		return c.RenderText("fail")
	}

	//保存文件
	activity.Image, err = utils.Upload(revel.BasePath+"/upload/activityimg/", file, header)
	if err != nil {
		return c.RenderText("fail")
	}

	//产生缩略图
	err = utils.Thumbnail(revel.BasePath+"\\upload\\activityimg\\", activity.Image, "thumb_"+activity.Image, 2, 300, 0)
	if err != nil {
		return c.RenderText("fail")
	}

	err = models.ActivityModel.Insert(activity)
	if err != nil {
		return c.RenderText("fail")
	}

	return c.RenderText("ok")
}

//获取活动分页列表
func (c Activity) GetPagingList(userid, actiminid, category string) revel.Result {
	var nodata = make([]string, 0, 0)
	c.Validation.Required(userid)
	c.Validation.Required(actiminid)
	c.Validation.Required(category)
	if c.Validation.HasErrors() {
		revel.WARN.Println("All params is required!")
		return c.RenderJson(nodata)
	}

	account_id, _ := strconv.Atoi(userid)
	min_activity_id, _ := strconv.Atoi(actiminid)
	_category, _ := strconv.Atoi(category)
	page_num := 10

	//先从缓存中取数据，取不到数据再查询数据库并入缓存
	var list []models.ActivityList
	var err error
	if min_activity_id == -1 {
		list, err = models.ActivityModel.GetPagingListByCategory(min_activity_id, _category, page_num)
		if err != nil {
			return c.RenderJson(nodata)
		}
	} else {
		c_key := "activitylist_" + actiminid + "_" + category + "_" + strconv.Itoa(page_num)
		if err := cache.Get(c_key, &list); err != nil {
			list, err = models.ActivityModel.GetPagingListByCategory(min_activity_id, _category, page_num)
			if err != nil {
				return c.RenderJson(nodata)
			}
			go cache.Set(c_key, list, 24*time.Hour)
		}
	}

	if len(list) == 0 {
		return c.RenderJson(nodata)
	}

	LoopActivityList(list, account_id)

	return c.RenderJson(list)
}

//下拉更新活动列表
func (c Activity) GetNewList(userid, actimaxid, category string) revel.Result {
	var nodata = make([]string, 0, 0)
	c.Validation.Required(userid)
	c.Validation.Required(actimaxid)
	c.Validation.Required(category)
	if c.Validation.HasErrors() {
		revel.WARN.Println("All params is required!")
		return c.RenderJson(nodata)
	}

	account_id, _ := strconv.Atoi(userid)
	max_activity_id, _ := strconv.Atoi(actimaxid)
	_category, _ := strconv.Atoi(category)

	list, err := models.ActivityModel.GetNewListByCategory(max_activity_id, _category)
	if err != nil {
		return c.RenderJson(nodata)
	}

	if len(list) == 0 {
		return c.RenderJson(nodata)
	}

	LoopActivityList(list, account_id)

	return c.RenderJson(list)
}

func LoopActivityList(list []models.ActivityList, account_id int) {
	for k, v := range list {
		//获取用户名和头像，先取缓存
		var account models.Account
		if err := cache.Get("account_"+strconv.Itoa(v.AccountId), &account); err != nil {
			account, err = models.AccountModel.GetAccountInfoById(v.AccountId)
			if err != nil {
				revel.WARN.Println("Get account info fail:", v.AccountId)
			} else {
				go cache.Set("account_"+strconv.Itoa(v.AccountId), account, 24*time.Hour)
			}
		}
		list[k].AccountName = account.Name
		list[k].AccountAvator = account.Avator

		var err error
		//获取活动的参与数、喜欢数、评论数
		list[k].JoinNum, err = models.ActivityModel.GetJoinNumByActivityId(v.Id)
		if err != nil {
			revel.WARN.Println("Get activity join number fail:", v.Id)
		}
		list[k].LikeNum, err = models.ActivityModel.GetLikeNumByActivityId(v.Id)
		if err != nil {
			revel.WARN.Println("Get activity like number fail:", v.Id)
		}
		list[k].CommentNum, err = models.ActivityModel.GetCommentNumByActivityId(v.Id)
		if err != nil {
			revel.WARN.Println("Get activity comment number fail:", v.Id)
		}

		//判断当前用户是否点击过参与和喜欢
		list[k].IsJoin, err = models.ActivityModel.IsJoin(v.Id, account_id)
		if err != nil {
			revel.WARN.Printf("Get user(id=%d) is join activity(id=%d) fail!", account_id, v.Id)
		}
		list[k].IsLike, err = models.ActivityModel.IsLike(v.Id, account_id)
		if err != nil {
			revel.WARN.Printf("Get user(id=%d) is like activity(id=%d) fail!", account_id, v.Id)
		}
	}

	return
}
