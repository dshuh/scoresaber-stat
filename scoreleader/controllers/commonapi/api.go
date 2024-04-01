package commonapi

import (
	"net/http"
	"strconv"

	"stove-gitlab.sginfra.net/backend/template/config"
	"stove-gitlab.sginfra.net/backend/template/constant"
	"stove-gitlab.sginfra.net/backend/template/framework/util"
	"stove-gitlab.sginfra.net/backend/template/framework/web/base"
	"stove-gitlab.sginfra.net/backend/template/models"

	"github.com/labstack/echo/v4"
)

const (
	headerContentType = "Content-Type"
	headerOctetStream = "application/octet-stream"
)

// GetHealthCheck 서버 상태 체크
func GetHealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

// GetStatus 서비스 상태 체크
func GetStatus(c echo.Context) error {
	conf := config.GetInstance()
	resp := base.Response{}

	getServerDetail := func() map[string]interface{} {
		server := util.GetProcessInfo()
		result := map[string]interface{}{
			"hostname":          server.HostInfo.Hostname,
			"create_time":       server.ProcessInfo.CreateTime,
			"cmd_line":          server.ProcessInfo.Cmdline,
			"cpu_count":         server.CPUCount,
			"cpu_count_logical": server.CPUCountLogical,
			"cpu_percent":       server.CPUPercent,
			"disk_percent":      server.DiskPercent,
			"memory_percent":    server.MemoryPercent,
		}
		return result
	}

	getConfigDetail := func() map[string]interface{} {
		result := map[string]interface{}{
			"environment": conf.Template.EnvName,
			"revision":    constant.RevVersion,
			"release":     constant.TagVersion,
		}
		return result
	}

	getDBDetail := func() (bool, map[string]interface{}) {
		result := map[string]interface{}{
			"type":      conf.DBType,
			"host":      conf.Database.MysqlDB.Host,
			"status_ok": false,
		}
		err := models.GetDBStatus()
		if err != nil {
			return false, result
		}
		result["status_ok"] = true
		return true, result
	}

	status := true
	dbStatus, dbDetail := getDBDetail()
	if !dbStatus {
		status = false
	}

	redisCheck, err := models.PingRedis()
	if err != nil {
		redisCheck = false
	}

	statusRedis := struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		DB       string `json:"db"`
		StatusOK bool   `json:"status_ok"`
	}{
		Host:     conf.Database.MysqlDB.Host,
		Port:     strconv.Itoa(conf.Database.RedisDB.Port),
		DB:       strconv.Itoa(conf.Database.RedisDB.DefaultDB),
		StatusOK: redisCheck,
	}

	resp.Value = map[string]interface{}{
		"status_ok": status,
		"status_detail": map[string]interface{}{
			"server": getServerDetail(),
			"config": getConfigDetail(),
			"db":     dbDetail,
			"redis":  statusRedis,
		},
	}

	resp.OK()

	return c.JSON(http.StatusOK, resp)
}

// GetVersion API 버전 조회
func GetVersion(c echo.Context, maxVersion string) error {
	resp := base.Response{}
	resp.Value = map[string]interface{}{
		"version":  maxVersion,
		"revision": constant.RevVersion,
		"release":  constant.TagVersion,
	}
	resp.OK()

	return c.JSON(http.StatusOK, resp)
}

// GetRealIP Real IP 조회
func GetRealIP(c echo.Context) error {
	resp := base.Response{}
	resp.Value = map[string]interface{}{
		"real_ip": c.Request().Header.Get(echo.HeaderXRealIP)}
	resp.OK()

	return c.JSON(http.StatusOK, resp)
}
