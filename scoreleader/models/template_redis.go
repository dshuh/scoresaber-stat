package models

import (
	"context"

	"stove-gitlab.sginfra.net/backend/template/framework/util/log"
)

// SetTemplateRedis SetTemplateRedis
func (o *TemplateServiceDB) SetTemplateRedis(template *Template) error {
	// 유효성 검사
	if err := validTemplateParams([]string{
		"key",
		"value",
	}, template); err != nil {
		return err
	}

	if err := o.redisdb.HSet(templateTable, template.Key, template.Value); err != nil {
		log.Errorf("redis hset error: %s, %s, %s, %s", templateTable, template.Key, template.Value, err)
		return err
	}
	return nil
}

// GetTemplateRedis GetTemplateRedis
func (o *TemplateServiceDB) GetTemplateRedis(template *Template) error {
	// 유효성 검사
	if err := validTemplateParams([]string{
		"key",
	}, template); err != nil {
		return err
	}

	var err error
	if template.Value, err = o.redisdb.HGet(templateTable, template.Key); err != nil {
		log.Warnf("redis hget error: %s, %s, %s", templateTable, template.Key, err)
		return err
	}
	return nil
}

// RemoveTemplateRedis RemoveTemplateRedis
func (o *TemplateServiceDB) RemoveTemplateRedis(template *Template) error {
	// 유효성 검사
	if err := validTemplateParams([]string{
		"key",
	}, template); err != nil {
		return err
	}

	if err := o.redisdb.HDel(templateTable, template.Key); err != nil {
		log.Errorf("redis hdel error: %s, %s, %s", templateTable, template.Key, err)
		return err
	}
	return nil
}

// RefreshRedis REDIS를 갱신한다.
func (o *TemplateServiceDB) RefreshRedis(templateList *TemplateList) error {

	if err := o.GetTemplates(templateList); err != nil {
		return err
	}

	if len(templateList.Templates) == 0 {
		return nil
	}

	ctx := context.Background()
	pipe := o.redisdb.Pipeline()

	// Redis에 갱신 대상인 자식 리소스들을 삭제한다.
	keys := make([]string, 0)
	for _, v := range templateList.Templates {
		keys = append(keys, v.Key)
	}
	pipe.HDel(ctx, templateTable, keys...)

	for _, v := range templateList.Templates {
		pipe.HSet(ctx, templateTable, v.Key, v.Value)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		if err != nil && err.Error() != errRedisPipeLineEmpty {
			log.Errorf("redis refresh error: %s", err)
			return err
		}
	}

	return nil
}
