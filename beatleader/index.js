function TemplateLayer(controlData) {
    var obj = this;

    // Layer Template Element 세팅
    obj.navpanel_left = $(controlData.navpanel_left);
    obj.dvHeader = $(controlData.dvHeader);
    obj.dvFooter = $(controlData.dvFooter);
    obj.dvMenu = $(controlData.dvMenu);
    obj.dvMainContent = $(controlData.dvMainContent);

    obj.pageMapLayer = $(controlData.pageMapLayer);
	obj.map = new Map({ pageMapLayer: obj.pageMapLayer });

    obj.pagePlayerLayer = $(controlData.pagePlayerLayer);
	obj.player = new Player({ pagePlayerLayer: obj.pagePlayerLayer });
    
	obj.pageRankingLayer = $(controlData.pageRankingLayer);
	obj.ranking = new Ranking({ pageRankingLayer: obj.pageRankingLayer });
    
    obj.pageCompareLayer = $(controlData.pageCompareLayer);
	obj.compare = new Compare({ pageCompareLayer: obj.pageCompareLayer });

    obj.apiConfig = config.apis.auth;
    obj.isInternal = false;
    obj.isLocalhost = false;
    obj.appKey = "";
    obj.appSecret = "";
    obj.accessToken = "";
    obj.tokenKey = "";

    obj.init = function () {
        // Layer Template 생성
        obj.createHeaderLayer();
        obj.createMenuLayer();

		obj.hpMapLayer.addClass("ui-btn-active ui-state-persist");
		obj.hpMapLayer.trigger("click");
		// obj.hpPlayerLayer.addClass("ui-btn-active ui-state-persist");
		// obj.hpPlayerLayer.trigger("click");
		var contentHeight = document.documentElement.clientHeight * 0.85;
		obj.pageMapLayer.css("min-height", contentHeight);
		// obj.pagePlayerLayer.css("min-height", contentHeight);
    };

    obj.reset = function () {
    };

    obj.createHeaderLayer = function () {

        html = '\
            <div data-role="navbar">\
                <ul>\
                    <li>\
                        <div class="ui-grid-b" style="width:100%">\
                            <div class="ui-block-a" style="width:3%; float:left; padding-left:5px; padding-top:8px">\
                                <img id="openLeftMenuLayer" alt="목록" title="목록" src="../images/header_menu_list.png" width="25" height="26"/>\
                            </div>\
                            <div class="ui-block-b" style="width:90%; float:left;">\
                            ' + displayEnvLayer() + '\
                            </div>\
                            <div class="ui-block-c" style="width:4%; float:right; padding-left:5px; padding-top:5px">\
                                <img onclick="javascript:root.apiCall(\'version\');" src="../images/info.png" />\
                            </div>\
                        </div>\
                    </li>\
                </ul>\
            </div>\
        ';

        obj.dvHeader.empty().append(html).trigger('create');

        obj.selEnv = $("#selEnv");
        
        obj.openLeftMenuLayer = $("#openLeftMenuLayer");

        obj.openLeftMenuLayer.click(function () {
            obj.navpanel_left.panel('open');
        });
    };

    obj.createMenuLayer = function () {
        var html = '\
        <div data-role="navbar">\
            <ul>\
                <li><a id="hpMapLayer">Map</a></li>\
                <li><a id="hpPlayerLayer">Player</a></li>\
                <li><a id="hpRankingLayer">Ranking</a></li>\
                <li><a id="hpCompareLayer">Compare</a></li>\
            </ul>\
        </div>\
        ';

        obj.dvMenu.empty().append(html).trigger('create');

		obj.hpMapLayer = $("#hpMapLayer");
		obj.hpPlayerLayer = $("#hpPlayerLayer");
		obj.hpRankingLayer = $("#hpRankingLayer");
		obj.hpCompareLayer = $("#hpCompareLayer");

		obj.allMenuReset = function () {
			obj.map.reset();
			obj.player.reset();
            obj.ranking.reset();
			obj.compare.reset();
		};

		obj.hpMapLayer.click(function () {
			obj.allMenuReset();
			obj.map.init();
		});

		obj.hpPlayerLayer.click(function () {
			obj.allMenuReset();
			obj.player.init();
		});

		obj.hpRankingLayer.click(function () {
			obj.allMenuReset();
			obj.ranking.init();
		});

		obj.hpCompareLayer.click(function () {
			obj.allMenuReset();
			obj.compare.init();
		});

    };

    obj.apiCall = function (apiID) {
        if(obj.apiConfig == undefined) return;
        var apiObj = obj.apiConfig.find(x => x.id === apiID);

        if(apiID == "version") {
            apiObj = config.apis.system.find(x => x.id == "version");
        }

        // set api parameters
        var requestData = setAPIParameters(apiObj, obj);
        if(requestData == null) return;

        // api validation check

        // request parameter 재정의
        if(apiObj.id == obj.tokenKey) {
            var env = (obj.isInternal) ?
                config.environments.find(x => x.prefix_internal_url === obj.selEnv.val()) :
                config.environments.find(x => x.prefix_external_url === obj.selEnv.val());

            if(obj.isInternal) {
                apiObj.path = ( env.group == "global") ?
                    config.environments.find(x => x.prefix_internal_url === obj.selEnv.val()).user_token_url :
                    config.environments.find(x => x.prefix_internal_url === obj.selEnv.val()).infra_token_url;
            } else {
                apiObj.path =  ( env.group == "global") ?
                    config.environments.find(x => x.prefix_external_url === obj.selEnv.val()).user_token_url :
                    config.environments.find(x => x.prefix_external_url === obj.selEnv.val()).infra_token_url;
            }
            // apiObj.path = (obj.isInternal) ?
            //     config.environments.find(x => x.prefix_internal_url === obj.selEnv.val()).infra_token_url :
            //     config.environments.find(x => x.prefix_external_url === obj.selEnv.val()).infra_token_url
        }
        
        // confirm alert 설정
        if(apiObj.message.confirm != undefined && apiObj.message.confirm != ""){
            if (!confirm(apiObj.message.confirm)) {
                return;
            }
        }
        
        // api 성공 시 callback 함수
        var fnSuccess = function (data) {
            if(apiID == "version") {
                if(validationAjaxSuccess(data)) {
                    eval(apiObj.callback + "(data);");
                }
            } else if(obj.tokenKey == "post_user_token") {
                obj.accessToken = "bearer " + data.response_data.access_token;
                return
            } else  if(obj.tokenKey == "post_infra_token") {
                obj.accessToken = data.access_token;
            } else {
                if(validationAjaxSuccess(data)) {
                    eval(apiObj.callback + "(data);");
                }
            }
        };

        // api 실패 시 callback 함수
        var fnError = function (err) {
            panelAlertLayer(apiObj.message.error + " > " + err, "error");
        };

        // get method callback 함수
        var callback_get = function(data) {
            if(apiObj.id =="version") {
                var revision = data.data.revision;
                var version = data.data.version;
                var selEnv = $("#selEnv option:selected").text();
                alert(selEnv + "\nrevision : " + revision + "\n" + "version : " + version);
            }
        }
        
        // put method callback 함수
        var callback_put = function(data) {
        }

        // post method callback 함수
        var callback_post = function(data) {
        }
        
        // delete method callback 함수
        var callback_delete = function(data) {
        }

        ajaxCall(apiObj.path, requestData, fnSuccess, fnError, apiObj.method, apiObj.content_type, apiObj.response_content_type);
    };
}