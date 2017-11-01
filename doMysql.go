package main 
import (
    "fmt"
    "strings"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
    "time"
)

var db *sql.DB
var dataSource string

func initDB() {
    db, _ = sql.Open("mysql", dataSource)
}

func writeDB(table_name string, dic map[string] interface{}){
    if db == nil {
        fmt.Println("this is writedb init")
        initDB()
    }
    str_sql := "insert into " + table_name + " ("
    var slice_col []string
    var slice_val []interface{}
    var slice_x []string
    for k, v := range dic {
        slice_col = append(slice_col, k)
        slice_x = append(slice_x, "?")
        slice_val = append(slice_val, v)
    }
    slice_col = append(slice_col, "create_time")
    slice_x = append(slice_x, "?")
    slice_val = append(slice_val, time.Now().Unix())
    slice_col = append(slice_col, "update_time")
    slice_x = append(slice_x, "?")
    slice_val = append(slice_val, time.Now().Unix())

    str_sql = str_sql + strings.Join(slice_col, ", ") + ") values (" + strings.Join(slice_x, ", ")  + ")"
    fmt.Println(str_sql)

    //stmtIns, err := db.Prepare(str_sql)
    //if err != nil {
    //    panic(err.Error())
    //}
    _, err := db.Exec(str_sql, slice_val...)
    if err != nil {
        fmt.Println("insert error; dic: ", dic, err.Error())
    }
}

func updateDB(table_name string, dic map[string] interface{}, cond map[string] interface{}){
    if db == nil {
        fmt.Println("this is updatedb init")
        initDB()
    }
    str_sql := "update " + table_name + " set "
    var slice_col []string
    var slice_val []interface{}
    for k, v := range dic {
        slice_col = append(slice_col, k + "=?")
        slice_val = append(slice_val, v)
    }
    slice_col = append(slice_col, "update_time=?")
    slice_val = append(slice_val, time.Now().Unix())
    str_sql = str_sql + strings.Join(slice_col, ", ") + " where 1=1 "
    for k, v := range cond {
        str_sql = str_sql + " and " + k + "=? "
        slice_val = append(slice_val, v)
    }
    _, err := db.Exec(str_sql, slice_val...)
    if err != nil {
        panic(err.Error())
    }
}
func deleteDB(table_name string, cond string){
    if db == nil {
        fmt.Println("this is deletedb init")
        initDB()
    }
    str_sql := "delete from " + table_name + " where " + cond
    fmt.Println(str_sql)
    result, err := db.Exec(str_sql)
    if err != nil {
        panic(err.Error())
    }
    fmt.Println(result.RowsAffected())
}

func selectDB(str_sql string) []map[string]interface{}{
    if db == nil {
        fmt.Println("this is query init")
        initDB()
    }
    rows, err := db.Query(str_sql)
    if err != nil {
        panic(err.Error())
    }
    columns, _ := rows.Columns()
    slice_ret := make([]map[string]interface{}, 0)

    scanArgs := make([]interface{}, len(columns))
    values := make([]interface{}, len(columns))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    for rows.Next() {
        err = rows.Scan(scanArgs...)
        ret := make(map[string]interface{})
        for i, col := range values {
            if col == nil{
                ret[columns[i]] = col
            }else{
                ret[columns[i]] = string(col.([]byte))
            }
        }
        slice_ret = append(slice_ret, ret)
    }
    return slice_ret
}

func artistDB(dic map[string] interface{}){
    if len(dic) == 10 {
        executeSql := "select max(artist_id) as max_id from music_app_artist;"
        ret := selectDB(executeSql)
        fmt.Println(ret)
        var max_id int
        if ret[0]["max_id"] == nil {
            max_id = 100001
        }else{
            max_id, _ = strconv.Atoi(ret[0]["max_id"].(string))
        }
        max_id += 1

        dic["artist_id"] = max_id
        dic["related_spider_id"] = dic["id"]
        delete(dic, "id")
        writeDB("music_app_artist", dic)
    }
}

func albumDB(albumList []map[string] interface{}, artist_id interface{}, la interface{}){
    for _, album := range albumList{
        album["artist_id"] = artist_id
        album["la"] = la
        album["status"] = 4
        executeSql := "select max(album_id) as max_id from music_app_album;"
        ret := selectDB(executeSql)
        var max_id int
        if ret[0]["max_id"] == nil {
            max_id = 400001
        }else{
            max_id, _ = strconv.Atoi(ret[0]["max_id"].(string))
        }
        max_id += 1

        album["album_id"] = max_id
        writeDB("music_app_album", album)
    }
}
//func batchWriteDB(table_name string, arr_dic []map[string] interface{})
//    if db == nil {
//        fmt.Println("this is writedb init")
//        initDB()
//    }
//    str_sql := "insert into " + table_name + " ("
//    dic := arr_dic[0]
//    var slice_col []string
//    var slice_val []interface{}
//    var slice_x []string
//    for k, v := range dic {
//        slice_col = append(slice_col, k)
//        slice_x = append(slice_x, "?")
//        slice_val = append(slice_val, v)
//    }
//    str_sql = str_sql + strings.Join(slice_col, ", ") + ") values (" + strings.Join(slice_x, ", ")  + ")"
//    fmt.Println(str_sql)
//
//    //stmtIns, err := db.Prepare(str_sql)
//    //if err != nil {
//    //    panic(err.Error())
//    //}
//    _, err := db.Exec(str_sql, slice_val...)
//    if err != nil {
//        panic(err.Error())
//    }
//}
