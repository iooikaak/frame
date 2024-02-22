# 项目简介

1.基于xormplus库，做啦封装，增加啦tracing

# 快速开始

* 创建引擎

```go

    db, err := xorm.New(&config.Config{
                        		Type:     "mysql",
                        		Server:   "127.0.0.1",
                        		Port:     3306,
                        		Database: "db",
                        		User:     "db",
                        		Password: "123456",
                        	})

        if err != nil {
                panic(err)
        }

	//step 2: 把db保存到全局，方便项目中使用


```

* db查询方式 - Exec方法，纯sql写法

```go
 sql :="select * from  user  where name = ?"
 sess := db.NewSession(context.Background())
 defer sess.Close()
 res, err = sess.Exec(sql,  "xorm")

```

* db查询方式 - QueryString方法

```go
   sess := db.NewSession(context.Background())
   defer sess.Close()
  //当调用QueryString时，第一个返回值results为[]map[string]string的形式。
  sql_2_1 := "select * from user"
  results, err := sess.QueryString(sql_2_1)

  //SqlMapClient和SqlTemplateClient传参方式类同如下2种，具体参见第6种方式和第7种方式
  sql_2_2 := "select id,userid,title,createdatetime,content from Article where id=?"
  results, err := sess.SQL(sql_2_2, 2).QueryString()

  sql_2_3 := "select id,userid,title,createdatetime,content from Article where id=?id"
  paramMap_2_3 := map[string]interface{}{"id": 2}
  results, err := sess.SQL(sql_2_3, &paramMap_2_3).QueryString()


```

* db查询方式 - get方法

```go
    sess := db.NewSession(context.Background())
    defer sess.Close()
  //使用Sql函数与Get函数组合可以查询单条数据 以Sql与Get函数组合为例：
  //获得单条数据的值，并存为结构体
  var article Article
  has, err := sess.Sql("select * from article where id=?", 2).Get(&article)

  //获得单条数据的值并存为map
  var valuesMap1 = make(map[string]string)
  has, err := sess.Sql("select * from article where id=?", 2).Get(&valuesMap1)

  var valuesMap2 = make(map[string]interface{})
  has, err := sess.Sql("select * from article where id=?", 2).Get(&valuesMap2)

  var valuesMap3 = make(map[string]xorm.Value)
  has, err := sess.Sql("select * from article where id=?", 2).Get(&valuesMap3)

  //获得单条数据的值并存为xorm.Record
  record := make(xorm.Record)
  has, err = sess.SQL("select * from article where id=?", 2).Get(&record)
  id := record["id"].Int64()
  content := record["content"].NullString()

  //获得单条数据某个字段的值
  var title string
  has, err := sess.Sql("select title from article where id=?", 2).Get(&title)

  var id int
  has, err := sess.Sql("select id from article where id=?", 2).Get(&id)

```

* db查询方式 - find方法

```go

sess := db.NewSession(context.Background())
defer sess.Close()

//find 返回多条记录
var categories []Category
err := sess.Sql("select * from category where id =?", 16).Find(&categories)

paramMap_6 := map[string]interface{}{"id": 2}
err := sess.Sql("select * from category where id =?id", &paramMap_6).Find(&categoriesdb)

```

* db更新操作 - update方法

```go

sess := db.NewSession(context.Background())
defer sess.Close()

//更新数据使用Update方法，Update方法的第一个参数为需要更新的内容，可以为一个结构体指针或者一个Map[string]interface{}类型。当传入的为结构体指针时，只有非空和0的field才会被作为更新的字段。当传入的为Map类型时，key为数据库Column的名字，value为要更新的内容。

user := new(User)
user.Name = "myname"
affected, err := sess.Id(id).Update(user)
/*
这里需要注意，Update会自动从user结构体中提取非0和非nil得值作为需要更新的内容，因此，如果需要更新一个值为0，则此种方法将无法实现，因此有两种选择：

1.通过添加Cols函数指定需要更新结构体中的哪些值，未指定的将不更新，指定了的即使为0也会更新。
*/
affected, err := sess.Id(id).Cols("age").Update(&user)
/*
2.通过传入map[string]interface{}来进行更新，但这时需要额外指定更新到哪个表，因为通过map是无法自动检测更新哪个表的。
*/
affected, err := sess.Table(new(User)).Id(id).Update(map[string]interface{}{"age":0})

```

* db事务操作例子

```go
session := engine.NewSession(context.Background())
defer session.Close()
// add Begin() before any action
err := session.Begin()
user1 := Userinfo{Username: "xiaoxiao", Departname: "dev", Alias: "lunny", Created: time.Now()}
_, err = session.Insert(&user1)
if err != nil {
    session.Rollback()
    return
}
user2 := Userinfo{Username: "yyy"}
_, err = session.Where("id = ?", 2).Update(&user2)
if err != nil {
    session.Rollback()
    return
}

_, err = session.Exec("delete from userinfo where username = ?", user2.Username)
if err != nil {
    session.Rollback()
    return
}

// add Commit() after all actions
err = session.Commit()
if err != nil {
    return
}
```
