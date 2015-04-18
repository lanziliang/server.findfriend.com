package models

import (
	"github.com/revel/revel"
	"strconv"
	"text/template"
)

var (
	ActivityModel = new(Activity)
)

type Activity struct {
	AccountId int    `json:"account_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Image     string `json:"img"`
	Category  int    `json:"category"`
}

//添加新活动
func (m *Activity) Insert(activity Activity) error {
	sql := "INSERT INTO activity (account_id, title, content, img, category) VALUES (?, ?, ?, ?, ?)"

	stmt, err := DbLocal.Prepare(sql)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}

	_, err = stmt.Exec(activity.AccountId, activity.Title, activity.Content, activity.Image, activity.Category)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}

	return nil
}

type ActivityList struct {
	Id            int    `json:"id"`
	AccountId     int    `json:"account_id"`
	AccountName   string `json:"account_name"`
	AccountAvator string `json:"account_avator"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Img           string `json:"img"`
	Category      int    `json:"category"`
	CreateDate    string `json:"create_date"`
	JoinNum       int    `json:"join_num"`
	LikeNum       int    `json:"like_num"`
	CommentNum    int    `json:"comment_num"`
	IsJoin        int    `json:"is_join"`
	IsLike        int    `json:"is_like"`
}

//获取活动分页列表
//参数说明: minId上一页最后一条活动id(minId=-1表示获取第一页数据)
func (m *Activity) GetPagingListByCategory(minId, category, pageNum int) ([]ActivityList, error) {
	var activityLists []ActivityList

	//拼接sql，对字符串进行转义处理防止SQL注入
	_minId := template.HTMLEscapeString(strconv.Itoa(minId))
	where := " WHERE "
	switch category {
	case 0:
		where += "category > 0"
	case 1:
		where += "category = 1"
	case 2:
		where += "category = 2"
	case 3:
		where += "category = 3"
	}
	if minId != -1 {
		where += " AND id < " + _minId
	}
	sql := `SELECT id,account_id,title,content,img,category,create_date 
			FROM activity` + where + `
			ORDER BY id DESC
			LIMIT 0, ?`

	rows, err := DbLocal.Query(sql, pageNum)
	if err != nil {
		revel.ERROR.Println(err)
		return activityLists, err
	}

	for rows.Next() {
		var activityList ActivityList
		err = rows.Scan(&activityList.Id, &activityList.AccountId, &activityList.Title, &activityList.Content, &activityList.Img, &activityList.Category, &activityList.CreateDate)
		if err != nil {
			revel.ERROR.Println(err)
			return activityLists, err
		}

		activityLists = append(activityLists, activityList)
	}

	return activityLists, nil
}

//获取最新活动列表
//参数说明: maxId第一条活动id
func (m *Activity) GetNewListByCategory(maxId, category int) ([]ActivityList, error) {
	var activityLists []ActivityList

	//拼接sql，对字符串进行转义处理防止SQL注入
	where := " WHERE "
	switch category {
	case 0:
		where += "category > 0 AND id > ?"
	case 1:
		where += "category = 1 AND id > ?"
	case 2:
		where += "category = 2 AND id > ?"
	case 3:
		where += "category = 3 AND id > ?"
	}
	sql := `SELECT id,account_id,title,content,img,category,create_date 
			FROM activity` + where + `
			ORDER BY id DESC`

	rows, err := DbLocal.Query(sql, maxId)
	if err != nil {
		revel.ERROR.Println(err)
		return activityLists, err
	}

	for rows.Next() {
		var activityList ActivityList
		err = rows.Scan(&activityList.Id, &activityList.AccountId, &activityList.Title, &activityList.Content, &activityList.Img, &activityList.Category, &activityList.CreateDate)
		if err != nil {
			revel.ERROR.Println(err)
			return activityLists, err
		}

		activityLists = append(activityLists, activityList)
	}

	return activityLists, nil
}

func (m *Activity) GetJoinNumByActivityId(id int) (int, error) {
	var num int

	sql := "SELECT count(1) AS num FROM activity_join WHERE activity_id = ?"
	err := DbLocal.QueryRow(sql, id).Scan(&num)
	if err != nil {
		revel.ERROR.Println(err)
		return num, err
	}

	return num, nil
}

func (m *Activity) GetLikeNumByActivityId(id int) (int, error) {
	var num int

	sql := "SELECT count(1) AS num FROM activity_like WHERE activity_id = ?"
	err := DbLocal.QueryRow(sql, id).Scan(&num)
	if err != nil {
		revel.ERROR.Println(err)
		return num, err
	}

	return num, nil
}

func (m *Activity) GetCommentNumByActivityId(id int) (int, error) {
	var num int

	sql := "SELECT count(1) AS num FROM activity_comments WHERE activity_id = ?"
	err := DbLocal.QueryRow(sql, id).Scan(&num)
	if err != nil {
		revel.ERROR.Println(err)
		return num, err
	}

	return num, nil
}

//判断用户是否点击参与某个活动
func (m *Activity) IsJoin(activityId, accountId int) (int, error) {
	var num int

	sql := "SELECT count(1) AS num FROM activity_join WHERE activity_id = ? AND account_id = ?"
	err := DbLocal.QueryRow(sql, activityId, accountId).Scan(&num)
	if err != nil {
		revel.ERROR.Println(err)
		return num, err
	}

	return num, nil
}

//判断用户是否点击喜欢某个活动
func (m *Activity) IsLike(activityId, accountId int) (int, error) {
	var num int

	sql := "SELECT count(1) AS num FROM activity_like WHERE activity_id = ? AND account_id = ?"
	err := DbLocal.QueryRow(sql, activityId, accountId).Scan(&num)
	if err != nil {
		revel.ERROR.Println(err)
		return num, err
	}

	return num, nil
}

func (m *Activity) AddJoin(activityId, accountId int) error {
	sql := "INSERT INTO activity_join (activity_id, account_id) VALUES (?, ?)"
	stmt, err := DbLocal.Prepare(sql)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}

	_, err = stmt.Exec(activityId, accountId)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}

	return nil
}

func (m *Activity) DeleteJoin(activityId, accountId int) (int64, error) {
	sql := "DELETE FROM activity_join WHERE activity_id = ? AND account_id = ?"
	stmt, err := DbLocal.Prepare(sql)
	if err != nil {
		revel.ERROR.Println(err)
		return 0, err
	}

	res, err := stmt.Exec(activityId, accountId)
	if err != nil {
		revel.ERROR.Println(err)
		return 0, err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		revel.ERROR.Println(err)
		return 0, err
	}

	if affect == 0 {
		revel.ERROR.Println("No rows exist!")
		return affect, nil
	}
	return affect, nil
}

func (m *Activity) AddLike(activityId, accountId int) error {
	sql := "INSERT INTO activity_Like (activity_id, account_id) VALUES (?, ?)"
	stmt, err := DbLocal.Prepare(sql)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}

	_, err = stmt.Exec(activityId, accountId)
	if err != nil {
		revel.ERROR.Println(err)
		return err
	}

	return nil
}

func (m *Activity) DeleteLike(activityId, accountId int) (int64, error) {
	sql := "DELETE FROM activity_like WHERE activity_id = ? AND account_id = ?"
	stmt, err := DbLocal.Prepare(sql)
	if err != nil {
		revel.ERROR.Println(err)
		return 0, err
	}

	res, err := stmt.Exec(activityId, accountId)
	if err != nil {
		revel.ERROR.Println(err)
		return 0, err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		revel.ERROR.Println(err)
		return 0, err
	}

	if affect == 0 {
		revel.ERROR.Println("No rows exist!")
		return affect, nil
	}
	return affect, nil
}
