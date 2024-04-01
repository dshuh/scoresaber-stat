package models

// CreateTemplate Template/01. Template Meta 생성
func CreateTemplate(template *Template) error {
	db := getDB()
	return db.CreateTemplate(template)
}

// UpdateTemplate Template/02. Template Meta 수정
func UpdateTemplate(template *Template) error {
	db := getDB()
	return db.UpdateTemplate(template)
}

// DeleteTemplate Template/03. Template Meta 삭제
func DeleteTemplate(template *Template) error {
	db := getDB()
	return db.DeleteTemplate(template)
}

// GetTemplate Template/04. Template Meta 조회
func GetTemplate(template *Template) error {
	db := getDB()
	return db.GetTemplate(template)
}

// GetTemplates Template/05. Template List 조회
func GetTemplates(templateList *TemplateList) error {
	db := getDB()
	return db.GetTemplates(templateList)
}

// SetTemplateRedis Set Template to Redis
func SetTemplateRedis(template *Template) error {
	db := getDB()
	return db.SetTemplateRedis(template)
}

// GetTemplateRedis Get Template to Redis
func GetTemplateRedis(template *Template) error {
	db := getDB()
	return db.GetTemplateRedis(template)
}

// RemoveTemplateRedis Remove Template to Redis
func RemoveTemplateRedis(template *Template) error {
	db := getDB()
	return db.RemoveTemplateRedis(template)
}

// RefreshRedis Redis Refresh
func RefreshRedis(templateList *TemplateList) error {
	db := getDB()
	return db.RefreshRedis(templateList)
}
