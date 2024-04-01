package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"stove-gitlab.sginfra.net/backend/template/framework/util/log"
	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
	"stove-gitlab.sginfra.net/backend/template/utils"
)

// CreateTemplateTable Template Meta 테이블 생성
func (o *TemplateServiceDB) CreateTemplateTable() error {
	// Template Table 생성
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (`id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, `key` varchar(100) NOT NULL UNIQUE, `value` varchar(500), `reg_at` varchar(100), `upd_at` varchar(100))", getTemplateTable())
	log.Debugf("CreateTemplateTable, query: %s", query)

	_, err := o.mysqldb.Excute(query)
	if err != nil {
		log.Errorf("CreateTemplateTable, query error: %s", query)
		return err
	}

	return nil
}

// CreateTemplateTable Template Meta 테이블 생성
func (o *TemplateServiceDB) InitTemplateTable() error {
	// Template Table 생성
	query := fmt.Sprintf("DELETE FROM %s", getTemplateTable())
	log.Debugf("InitTemplateTable, query: %s", query)

	_, err := o.mysqldb.Excute(query)
	if err != nil {
		log.Errorf("InitTemplateTable, query error: %s", query)
		return err
	}

	return nil
}

// CreateTemplate Template/01. Template Meta 생성
func (o *TemplateServiceDB) CreateTemplate(template *Template) error {
	// 유효성 검사
	if err := validTemplateParams([]string{
		"key",
		"value",
	}, template); err != nil {
		return err
	}

	// Template Meta 등록
	template.RegAt = utils.GetMilliTimestamp()
	template.UpdAt = template.RegAt

	query := fmt.Sprintf("INSERT INTO %s "+
		"(`key`, `value`, `reg_at`, `upd_at`)", getTemplateTable())
	query = fmt.Sprintf("%v VALUES (%v)", query, strings.TrimSuffix(strings.Repeat("?,", strings.Count(query, ",")+1), ","))
	log.Debugf("CreateTemplate, query: %s", query)

	_, err := o.mysqldb.ExcuteWithArgs(query,
		template.Key,
		template.Value,
		template.RegAt,
		template.UpdAt,
	)

	if err != nil {
		log.Errorf("CreateTemplate, query error: %s", query)
		return err
	}

	go SetTemplateRedis(template)
	return nil
}

// UpdateTemplate Template/02. Template Meta 수정
func (o *TemplateServiceDB) UpdateTemplate(template *Template) error {
	// 유효성 검사
	if err := validTemplateParams([]string{
		"key",
		"value",
	}, template); err != nil {
		return err
	}

	// Template Meta 수정
	template.UpdAt = utils.GetMilliTimestamp()

	query := fmt.Sprintf("UPDATE %s SET "+
		"`value`=?, `upd_at`=? WHERE `key`=?", getTemplateTable())

	log.Debugf("UpdateTemplate, query: %s", query)

	ret, err := o.mysqldb.ExcuteWithArgs(query,
		template.Value,
		template.UpdAt,
		template.Key)

	if err != nil {
		log.Errorf("UpdateTemplate, query error: %s", query)
		return err
	}

	cnt, err := ret.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return base.MakeError(base.ResultCannotFindData, errors.New(errTemplateNotFound))
	}

	go SetTemplateRedis(template)

	return nil
}

// DeleteTemplate Template/03. Template Meta 삭제
func (o *TemplateServiceDB) DeleteTemplate(template *Template) error {
	// 유효성 검사
	if err := validTemplateParams([]string{
		"key",
	}, template); err != nil {
		return err
	}

	// Template Meta 삭제
	query := fmt.Sprintf("DELETE FROM %s WHERE `key`=?", getTemplateTable())

	log.Debugf("DeleteTemplate, query: %s", query)

	ret, err := o.mysqldb.ExcuteWithArgs(query, template.Key)

	if err != nil {
		return err
	}

	cnt, err := ret.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return base.MakeError(base.ResultCannotFindData, errors.New(errTemplateNotFound))
	}

	go RemoveTemplateRedis(template)

	return nil
}

// GetTemplate Template/04. Template Meta 조회
func (o *TemplateServiceDB) GetTemplate(template *Template) error {
	// 유효성 검사
	if err := validTemplateParams([]string{
		"key",
	}, template); err != nil {
		return err
	}

	// Redis 조회
	if err := GetTemplateRedis(template); err != nil {
		log.Warnf("redis hget error: %s, %s, %s", templateTable, template.Key, err)
	} else {
		return nil
	}

	// Template Meta 검색
	query := fmt.Sprintf("SELECT "+
		"`id`, `key`, `value`, `reg_at`, `upd_at` "+
		"FROM %s WHERE `key`='%s'", getTemplateTable(), template.Key)

	log.Debugf("GetTemplate, query: %s", query)

	if err := o.mysqldb.QueryRow(query,
		&template.ID,
		&template.Key,
		&template.Value,
		&template.RegAt,
		&template.UpdAt,
	); err != nil {
		if strings.Contains(err.Error(), errNotFound) {
			return base.MakeError(base.ResultCannotFindData, errors.New(errTemplateNotFound))
		}
		return err
	}

	return nil
}

// GetTemplates Template/05. Template List 조회
func (o *TemplateServiceDB) GetTemplates(templateList *TemplateList) error {
	// 검색 조건
	pageSize := templateList.PageSize
	pageOffset := (templateList.PageIndex - 1) * pageSize
	sort := sortRegAtDesc
	if templateList.SortOrder != "desc" {
		sort = sortRegAtAsc
	}

	var filter = ""
	StartRegAt, _ := strconv.ParseInt(templateList.StartRegAt, 10, 64)
	EndRegAt, _ := strconv.ParseInt(templateList.EndRegAt, 10, 64)
	filter = utils.AddQueryFilterInt64(filter, "`reg_at` >= ", StartRegAt)
	filter = utils.AddQueryFilterInt64(filter, "`reg_at` <= ", EndRegAt)
	filter = utils.AddQueryFilterString(filter, "`key` = ", templateList.Key)
	filter = utils.AddQueryFilterString(filter, "`value` = ", templateList.Value)

	// 클라이언트 리스트 카운트 검색
	queryCount := fmt.Sprintf("SELECT count(*) "+
		"FROM %s", getTemplateTable())
	if filter != "" {
		queryCount = fmt.Sprintf("%s %s", queryCount, filter)
	}

	log.Debugf("GetTemplates, query: %s", queryCount)

	if err := o.mysqldb.QueryRow(queryCount, &templateList.TotalCount); err != nil {
		if strings.Contains(err.Error(), errNotFound) {
			return base.MakeError(base.ResultCannotFindData, errors.New(errTemplateNotFound))
		}
		return err
	}

	// 클라이언트 리스트 검색
	query := fmt.Sprintf("SELECT * FROM %s", getTemplateTable())
	if filter != "" {
		query = fmt.Sprintf("%s %s", query, filter)
	}
	query = fmt.Sprintf("%s order by %s limit %d, %d", query, sort, pageOffset, pageSize)

	log.Debugf("GetTemplates, query: %s", query)

	rows, err := o.mysqldb.Query(query)
	if err != nil {
		if strings.Contains(err.Error(), errNotFound) {
			return base.MakeError(base.ResultCannotFindData, errors.New(errTemplateNotFound))
		}
		return err
	}
	defer rows.Close()

	template := NewTemplate()
	for rows.Next() {
		err := rows.Scan(
			&template.ID,
			&template.Key,
			&template.Value,
			&template.RegAt,
			&template.UpdAt,
		)
		if err != nil {
			return err
		}

		templateList.Templates = append(templateList.Templates, *template)
	}

	//templateList.TotalCount = len(templateList.Templates)

	return nil
}
