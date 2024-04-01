package config

import (
	"sync"

	"stove-gitlab.sginfra.net/backend/template/framework/web/config"
)

// Template 공통 설정
type Template struct {
	ApplicationName   string `json:"application_name" yaml:"application_name"`
	ApplicationKey    string `json:"application_key" yaml:"application_key"`
	ApplicationSecret string `json:"application_secret" yaml:"application_secret"`
	Debug             bool   `json:"debug" yaml:"debug"`
	APIDocs           bool   `json:"api_docs" yaml:"api_docs"`
	APIDocsPath       string `json:"api_docs_path" yaml:"api_docs_path"`
	EnvName           string `json:"env_name" yaml:"env_name"`
	EnvOwner          string `json:"env_owner" yaml:"env_owner"`
	RKCredencialLocal string `json:"rk_credencial_local" yaml:"rk_credencial_local"`
	RKCredencialS3    string `json:"rk_credencial_s3" yaml:"rk_credencial_s3"`
	JWTSecret         string `json:"jwt_secret" yaml:"jwt_secret"`
	JWTTimeout        int    `json:"jwt_timeout" yaml:"jwt_timeout"`
	LdapUrlPrimary    string `json:"ldap_url_primary" yaml:"ldap_url_primary"`
	LdapUrlSecondary  string `json:"ldap_url_secondary" yaml:"ldap_url_secondary"`
	LdapBaseDN        string `json:"ldap_base_dn" yaml:"ldap_base_dn"`
	LdapUsername      string `json:"ldap_username" yaml:"ldap_username"`
	LdapPassword      string `json:"ldap_password" yaml:"ldap_password"`
}

// TemplateConfig 는 어플리케이션의 설정을 처리한다.
type TemplateConfig struct {
	config.Config `yaml:",inline"`
	Template      Template `json:"template" yaml:"template"`
	Version       string
}

// currentConfig 는 전역 Config 싱글턴 객체이다.
// GetInstance 함수에 의해 싱글턴 객체를 참조할 수 있다.
var currentConfig *TemplateConfig
var once sync.Once

// GetInstance 함수는 Config 싱글턴 객체를 참조한다.
func GetInstance(filepath ...string) *TemplateConfig {
	once.Do(func() {
		if len(filepath) <= 0 {
			panic(config.ErrInitConfigFailed)
		}
		currentConfig = &TemplateConfig{}
		if err := config.Load(filepath[0], currentConfig); err != nil {
			currentConfig = nil
		}
	})

	return currentConfig
}
