package server

import (
	_type2 "DY-DanMu/dbConn/_type"
	"DY-DanMu/dbConn/config"
	"DY-DanMu/web/server/_type"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/olivere/elastic/v7"
	Log "github.com/sirupsen/logrus"
)

type ResultItem struct {
	Rid     string `json:"Rid"`
	Id      string `json:"Id"`
	Payload struct {
		Bl    int    `json:"bl"`
		Bnn   string `json:"bnn"`
		Brid  string `json:"brid"`
		Cid   string `json:"cid"`
		Cst   int64  `json:"cst"`
		Ct    string `json:"ct"`
		Dms   string `json:"dms"`
		El    string `json:"el"`
		Hc    string `json:"hc"`
		Ic    string `json:"ic"`
		Level int    `json:"level"`
		Lk    string `json:"lk"`
		Nn    string `json:"nn"`
		Rid   string `json:"rid"`
		Sahf  string `json:"sahf"`
		Txt   string `json:"txt"`
		Type  string `json:"type"`
		Uid   string `json:"uid"`
		Urlev int    `json:"urlev"`
	} `json:"Payload"`
}

type UserBarrageResult struct {
	ResultList []ResultItem `json:"result_list"`
	Hits       int64        `json:"hits"`
}

type SelectMiddlerWare struct {
	RedisClient *redis.Client
	Client      *elastic.Client
	Conn        *sql.DB
	Index       string
	Count       int
}

// HandlerEsResutl:处理Es返回数据转换为[]ResultItem
func HandlerEsResutl(hits []*elastic.SearchHit) []ResultItem {
	var typ ResultItem
	var resultTyp []ResultItem
	for _, res := range hits {
		//fmt.Println(string(res.Source))
		err := json.Unmarshal(res.Source, &typ)
		if err != nil {
			Log.Error(err)
			continue
		}
		resultTyp = append(resultTyp, typ)
	}
	return resultTyp
}

// UserQuery 查询用户弹幕
func (e *SelectMiddlerWare) UserQuery(data _type.UserSearchStruct, result *UserBarrageResult) error {
	if e.Client == nil {
		Log.Warn("ElasticSearch client is nil, skipping query")
		return nil
	}
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewTermQuery("Payload.nn.keyword", data.UserName))
	boolQuery.Filter(elastic.NewRangeQuery("Payload.cst").Gte(data.StartTime).Lte(data.EndTime))
	pretty := e.Client.Search(data.EsIndex).Query(boolQuery).From(data.From).Size(10).Sort("_score",
		false).Sort("Payload.cst", false).Pretty(true)
	searchResult, err := pretty.Do(context.Background())
	if err != nil {
		Log.Error(err)
		return err
	}
	resultTyp := HandlerEsResutl(searchResult.Hits.Hits)
	if len(resultTyp) > 0 {
		*result = UserBarrageResult{
			ResultList: resultTyp,
			Hits:       searchResult.Hits.TotalHits.Value,
		}
	} else {
		*result = UserBarrageResult{
			[]ResultItem{},
			0,
		}
	}
	return nil
}

// BarrageAll 查询所有弹幕 分页返回
func (e *SelectMiddlerWare) BarrageAll(data _type.BarrageAllStruct, result *UserBarrageResult) error {
	if e.Client == nil {
		// 降级：从 MySQL 查询
		var rows *sql.Rows
		var err error
		if data.Rid != "" {
			rows, err = e.Conn.Query(fmt.Sprintf("SELECT nn, txt, cst FROM %s WHERE rid = ? ORDER BY cst DESC LIMIT 10 OFFSET ?", config.MysqlDBName+"."+config.MysqlTableName), data.Rid, data.From)
		} else {
			rows, err = e.Conn.Query(fmt.Sprintf("SELECT nn, txt, cst FROM %s ORDER BY cst DESC LIMIT 10 OFFSET ?", config.MysqlDBName+"."+config.MysqlTableName), data.From)
		}

		if err != nil {
			Log.Error(err)
			return err
		}
		defer rows.Close()

		var resultList []ResultItem
		for rows.Next() {
			var nn, txt string
			var cst int64
			if err := rows.Scan(&nn, &txt, &cst); err != nil {
				Log.Error(err)
				continue
			}
			item := ResultItem{}
			item.Payload.Nn = nn
			item.Payload.Txt = txt
			item.Payload.Cst = cst
			resultList = append(resultList, item)
		}
		*result = UserBarrageResult{
			ResultList: resultList,
			Hits:       int64(len(resultList)), // MySQL 简单查询暂不返回总数
		}
		return nil
	}
	query := elastic.NewMatchAllQuery()
	pretty := e.Client.Search(data.EsIndex).From(data.From).Query(query).Sort("Payload.cst", false).From(data.From).Size(10).Pretty(true)
	searchResult, err := pretty.Do(context.Background())
	if err != nil {
		Log.Error(err)
		return err
	}
	resultTyp := HandlerEsResutl(searchResult.Hits.Hits)
	if len(resultTyp) > 0 {
		*result = UserBarrageResult{
			ResultList: resultTyp,
			Hits:       searchResult.Hits.TotalHits.Value,
		}
	} else {
		*result = UserBarrageResult{
			[]ResultItem{},
			0,
		}
	}
	return nil
}

// SearchFieldAll:查询所有弹幕 分页返回每页10条
func (e *SelectMiddlerWare) SearchFieldAll(data _type.QueryAllFieldStruct, result *UserBarrageResult) error {
	if e.Client == nil {
		Log.Warn("ElasticSearch client is nil, skipping query")
		return nil
	}
	allQuery := elastic.NewBoolQuery() //
	query := elastic.NewMatchQuery("Payload.txt", data.Query)
	allQuery.Must(query)
	pretty := e.Client.Search(data.EsIndex).Query(allQuery).Sort("_score", false).Sort("Payload.cst",
		false).From(data.From).Size(10).Pretty(true)
	searchResult, err := pretty.Do(context.Background())
	if err != nil {
		Log.Error(err)
		return err
	}
	resultTyp := HandlerEsResutl(searchResult.Hits.Hits)
	if len(resultTyp) > 0 {
		*result = UserBarrageResult{
			ResultList: resultTyp,
			Hits:       searchResult.Hits.TotalHits.Value,
		}
	} else {
		*result = UserBarrageResult{
			[]ResultItem{},
			0,
		}
	}
	return nil
}

// BarrageCount:查询一段时间内的弹幕总数
func (e *SelectMiddlerWare) BarrageCount(data _type.BarrageCountStruct, result *_type2.BarrageCount) error {
	res, err := e.RedisClient.Get("BarrageCount").Result()
	if err != nil {
		return err
	}
	count, err := strconv.Atoi(res)
	if err != nil {
		return err
	}
	result.Count = count
	return nil
}

// StatisticsBarrageForTime:查询弹幕频率
func (e *SelectMiddlerWare) StatisticsBarrageForTime(data _type.StatisticsBarrageStruct, result *[]_type2.BarrageStatisticsCountResult) error {
	bytes, err := e.RedisClient.Get("BarrageStatisticsCountResult").Bytes()
	if err != nil {
		prepare, err := e.Conn.Prepare(fmt.Sprintf("select COUNT(*) As a,txt from %s WHERE ?>cst<? and txt != ? GROUP BY txt ORDER BY a Desc LIMIT ?", config.MysqlDBName+"."+config.MysqlTableName))
		if err != nil {
			Log.Error(err)
			return err
		}
		query, err := prepare.Query(data.StartTime, data.EndTime, "为主播点了个赞", data.From)
		if err != nil {
			Log.Error(err)
			return err
		}
		var barrageStatistic _type2.BarrageStatisticsCountResult
		for query.Next() {
			err := query.Scan(&barrageStatistic.Count, &barrageStatistic.Txt)
			if err != nil {
				return err
			}
			*result = append(*result, barrageStatistic)
		}
		redisData, err := json.Marshal(*result)
		if err != nil {
			Log.Error(err)
		} else {
			err = e.RedisClient.Set("BarrageStatisticsCountResult", redisData, time.Hour).Err()
			if err != nil {
				Log.Error(err)
			}
		}
		return nil
	} else {
		err := json.Unmarshal(bytes, result)
		if err != nil {
			panic(err)
		}
		return nil
	}
}

// StatisticsUserBarrageForTime:查询用户弹幕频率
func (e *SelectMiddlerWare) StatisticsUserBarrageForTime(data _type.StatisticsBarrageStruct, result *[]_type2.BarrageStatisticsUserCountResult) error {
	bytes, err := e.RedisClient.Get("StatisticsUserBarrageForTime").Bytes()
	if err != nil {
		prepare, err := e.Conn.Prepare(fmt.Sprintf("select COUNT(*) As a, nn,uid from %s WHERE ?>cst<? GROUP BY nn,uid ORDER BY a Desc LIMIT ?", config.MysqlDBName+"."+config.MysqlTableName))
		if err != nil {
			Log.Error(err)
			return err
		}
		query, err := prepare.Query(data.StartTime, data.EndTime, data.From)
		if err != nil {
			Log.Error(err)
			return err
		}
		var barrageStatistic _type2.BarrageStatisticsUserCountResult
		for query.Next() {
			err := query.Scan(&barrageStatistic.Count, &barrageStatistic.UserName, &barrageStatistic.UserId)
			if err != nil {
				return err
			}
			*result = append(*result, barrageStatistic)
		}
		redisData, err := json.Marshal(*result)
		if err != nil {
			Log.Error(err)
		} else {
			err = e.RedisClient.Set("StatisticsUserBarrageForTime", redisData, time.Hour*24).Err()
			if err != nil {
				Log.Error(err)
			}
		}
		return nil
	} else {
		err := json.Unmarshal(bytes, result)
		if err != nil {
			panic(err)
		}
		return nil
	}
}

// GetDanmuFrequency: 获取弹幕词频统计供 AI 分析
func (e *SelectMiddlerWare) GetDanmuFrequency(req _type.DanmuFrequencyRequest, result *[]_type2.BarrageStatisticsCountResult) error {
	limit := req.Limit
	// 如果 limit 为 0，默认取前 1000 条
	if limit <= 0 {
		limit = 1000
	}

	var queryStr string
	var prepare *sql.Stmt
	var err error

	if req.Rid != "" {
		queryStr = fmt.Sprintf("SELECT COUNT(*) as count, txt FROM %s WHERE rid = ? GROUP BY txt ORDER BY count DESC LIMIT ?", config.MysqlDBName+"."+config.MysqlTableName)
		prepare, err = e.Conn.Prepare(queryStr)
		if err != nil {
			Log.Error(err)
			return err
		}
		defer prepare.Close()
		rows, err := prepare.Query(req.Rid, limit)
		if err != nil {
			Log.Error(err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item _type2.BarrageStatisticsCountResult
			if err := rows.Scan(&item.Count, &item.Txt); err != nil {
				Log.Error(err)
				continue
			}
			*result = append(*result, item)
		}
	} else {
		queryStr = fmt.Sprintf("SELECT COUNT(*) as count, txt FROM %s GROUP BY txt ORDER BY count DESC LIMIT ?", config.MysqlDBName+"."+config.MysqlTableName)
		prepare, err = e.Conn.Prepare(queryStr)
		if err != nil {
			Log.Error(err)
			return err
		}
		defer prepare.Close()
		rows, err := prepare.Query(limit)
		if err != nil {
			Log.Error(err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item _type2.BarrageStatisticsCountResult
			if err := rows.Scan(&item.Count, &item.Txt); err != nil {
				Log.Error(err)
				continue
			}
			*result = append(*result, item)
		}
	}

	return nil
}
