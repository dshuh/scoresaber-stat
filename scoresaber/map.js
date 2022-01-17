function Map(controlData) {

    var obj = this;
    obj.pageMapLayer = $(controlData.pageMapLayer);

    obj.menuConfig = config.menu.find(x => x.id === "map");
    obj.gridConfig = obj.menuConfig.contents_area.grid;
    obj.checkedColNames = [];

    // #region 개발자 영역 > 멤버변수 선언 및 초기화
    obj.selectedMapIDs = []; // Grid에서 bulk 처리 시 사용.
    obj.dataCount = 0; // list 조회 시 전체 카운트 제공 시 사용.
    obj.total_count = 0;
    // #endregion

    // 서비스 활성화를 위한 초기화 설정을 담당한다.
    obj.init = function () {
        obj.createMapLayer();
        obj.pageMapLayer.show();
    };

    // 서비스를 비활성화를 위한 리셋 설정을 담당한다.
    obj.reset = function () {
        // #region 개발자 영역 > 멤버변수 초기화
        obj.selectedMapIDs = [];
        obj.dataCount = 0;
        obj.total_count = 0;
        // #endregion

        obj.checkedColNames = [];
        obj.pageMapLayer.empty().hide();
    };

    // search 영역과 contents 영역의 UI를 생성하고, control을 생성하고 이벤트를 등록한다. 
    obj.createMapLayer = function () {

        var html = displayContentsLayer(obj.menuConfig);

        obj.pageMapLayer.empty().append(html).trigger("create");

        obj.defineElements();
        obj.defineElementsEvent();
    };

    // config 파일의 root > menu > {search_area | contents_area} > controls 항목들에 대한 변수를 정의한다.
    obj.defineElements = function() {
		obj.selStartStar = $("#selStartStar");
		obj.selEndStar = $("#selEndStar");
		obj.btnSearch = $("#btnSearch");
		obj.mapContainer = $("#mapContainer");
		obj.mapPager = $("#mapPager");
    };
    
    obj.defineElementsEvent = function() {
        // #region 개발자 영역 > Config 정의한 Controls 이벤트 등록
        // config 파일의 root > menu > {search_area | contents_area} > controls 항목들에 대한 이벤트 함수를 등록한다.
		
        var optionHtml = "";
        for (var i=0;i<=15;i++) {
            optionHtml += '<option value="' + i + '">' + i + ' stars</option>';
        }
        obj.selStartStar.append(optionHtml).trigger("create");
        obj.selEndStar.append(optionHtml).trigger("create");

        obj.selStartStar.val("0").prop("selected", true).change();
        obj.selEndStar.val("15").prop("selected", true).change();

		obj.btnSearch.click(function () {
			obj.setGrid(); //Grid 조회 시 주석 해제
            var dataList = mapList.filter(function (el) {
                return el.stars >= obj.selStartStar.val() &&
                       el.stars < obj.selEndStar.val()
              });
            obj.bindGrid(dataList); //Grid 조회 시 주석 해제
		});

        // #endregion

        // grid column checkbox 태그의 name 값은 root > menu > contents_area > grid > checkbox_name에 정의한다.
        if(obj.gridConfig != undefined && obj.gridConfig.checkbox_name != undefined ){
            // chk_columns 공통으로 사용 권장.
            obj.chk_columns = $('input[name=' + obj.gridConfig.checkbox_name + ']');
            
            // Columns 체크박스 change 이벤트
            obj.chk_columns.change(function(){
                if(this.checked) {
                    obj.checkedColNames.push(this.value);
                } else {
                    obj.checkedColNames.splice($.inArray(this.value, obj.checkedColNames), 1);
                }
            });
            // Config의 "columns" array의 hidden 값이 true가 아닌 값을 셋팅한다.
            obj.chk_columns.filter(":checked").each(function() {
                obj.checkedColNames.push(this.value);
            });
        }
    };
    
    obj.setGrid = function() {
        // 사용 여부 체크
        if(obj.gridConfig.columns == undefined || obj.gridConfig.columns.length <= 0) {
            console.log("config error");
            return;
        }
        // Columns 정보 셋팅
        var columnInfo = setColumns(obj.gridConfig.columns, obj.checkedColNames);

        // #region Grid 옵션 상세 설정
        obj.mapContainer.jqGrid('clearGridData').jqGrid({
            datatype: "local",
            colNames: columnInfo[0],
            colModel: columnInfo[1],
            pager: "#mapPager",
            rowList: [],
            rowNum: 10000, // Grid에 그릴 최대 갯수 지정
            bottompager:false,
            pgbuttons: false,
            pgtext: null,
            viewrecords: false,
            emptyrecords: "Empty data.",
            inlineEditing: {
				keys: true
            },
            //cellEdit:true,
            gridview: true,
            autoencode: true,
            ignoreCase: true,
            autowidth: true,
            scrollOffset: 0,
            shrinkToFit:false,
            forceFit:true,
            height: "auto",
            multiselect: true,
            beforeSelectRow: function(rowid, e){ 
                // #region 개발자 영역 > Grid loadComplete Event 구현부
                return ($.jgrid.getCellIndex(e.target) == 0)
                // #endregion
            },
            onSelectRow: function (rowid, e) {
                // #region 개발자 영역 > Grid onSelectRow Event 구현부
                // #endregion
            },
            afterEditCell: function (rowid, cellname, value, iRow, iCol) { 
                // #region 개발자 영역 > Grid afterEditCell Event 구현부
                // #endregion
            },
            afterSaveCell: function(rowid, cellname, value, iRow, iCol) {
                // #region 개발자 영역 > Grid onSelectRow Event 구현부
                // #endregion
            },
            ondblClickRow: function (rowid, iRow, iCol, e) {
                // #region 개발자 영역 > Grid ondblClickRow Event 구현부
                // #endregion
            },
            loadComplete: function (data) {
                // #region 개발자 영역 > Grid loadComplete Event 구현부
                // #endregion
            }
        });
        // #endregion

        // #region 개발자 영역 > Grid Formatter callback 함수 구현부
        $.fn.fmatter.convertDownload = function (cellValue,rowObject,options) {
			if(cellValue != "") {
                return "<a href='" + options.download + "'>" + cellValue + "</a>";
            }
            return "";
        };
        $.fn.fmatter.convertStars = function (cellValue,rowObject,options) {
			return ConvertToString("stars", options);
        };
        $.fn.fmatter.convertMapTitle = function (cellValue,rowObject,options) {
			return "<a href='https://scoresaber.com/leaderboard/" + options.uid + "' target='_blank'>" + cellValue + "</a>";
        };
        $.fn.fmatter.convertSongTitle = function (cellValue,rowObject,options) {
			return "<a href='https://scoresaber.com/leaderboard/" + options.leaderboardId + "' target='_blank'>" + cellValue + "</a>";
        };
        $.fn.fmatter.convertPP = function (cellValue,rowObject,options) {
			return ConvertToString("pp", options);
        };
        $.fn.fmatter.convertRating = function (cellValue,rowObject,options) {
			return ConvertToString("rating", options) + "%";
        };
        $.fn.fmatter.convertSongTime = function (cellValue,rowObject,options) {
			return ConvertToString("durationSeconds", options);
        };
        // #endregion
        
        obj.navGrid();
    }

    // Grid Data Binding
    obj.bindGrid = function(data) {
		for(var i=data.length-1; i>=0; i--) {
			obj.mapContainer.jqGrid("addRowData", data[i].uid, data[i]);
		}
		var gridCount = obj.mapContainer.getGridParam("reccount");
		var gridCaption = "Searched <font color='blue'><b>" + gridCount + "</b></font> Completed.";
		obj.mapContainer.jqGrid("setCaption", gridCaption);
    }

    // Grid 하단 Navigation Bar 버튼 및 이벤트 등록
    obj.navGrid = function() {
        var toolbar = obj.gridConfig.toolbar;
        var ctrlToolbarID = "#" + toolbar.id;
        obj.mapContainer.jqGrid("navGrid", ctrlToolbarID, {
            edit:toolbar.edit, 
            add: toolbar.add, 
            del: toolbar.del, 
            search: toolbar.search, 
            refresh: toolbar.refresh
        }).jqGrid("filterToolbar", {
            stringResult: true, 
            searchOnEnter: true,
            defaultSearch: "cn", 
            ignoreCase: true,
            searchOperators: true
        });

        for(var i=0;i<toolbar.buttons.length;i++) {
            obj.mapContainer.jqGrid("navButtonAdd", ctrlToolbarID, toolbar.buttons[i]);
            $("#" + toolbar.buttons[i].id).off("click").click(function () {
                eval("javascirpt:" + toolbar.buttons.find(x => x.id === this.id).callback + "()");
            });
        }
        
        // #region 개발자 영역 > Grid navButton Callback 구현부
        
		var callback_navButton_search = function() {
            obj.mapContainer[0].toggleToolbar();
		};
		var callback_navButton_download = function() {
            obj.selectedMapIDs = obj.mapContainer.getGridParam("selarrrow");
            if (obj.selectedMapIDs.length == 0) {
                panelAlertLayer("Please select the data to download with the checkbox.", "info");
                return;
            }

            var gridData = obj.mapContainer.getGridParam('data');
            var downloadMap = [];
            for( var i=0;i<obj.selectedMapIDs.length;i++) {            
                downloadMap.push(gridData.find(x => x.id === obj.selectedMapIDs[i]));
            }

            var str = JSON.stringify(downloadMap);
            var blob = new Blob( [str], {
                type: "application/octet-stream"
            });

            var url = URL.createObjectURL( blob );
            var link = document.createElement("a");
            link.href = url;
            link.download = new Date().yyyyMMddHHmmss() + ".json";
            link.click();
		};
		var callback_navButton_more = function() {
            obj.mapApiCall("get_map");
		};

        // #endregion
    }

    obj.displayData = function (data) {
        //"<a href='https://scoresaber.com/u/76561198830502286" + page + "' target='_blank'>" + cellValue + "</a>"
        //https://scoresaber.com/global/2&country=kr
        var globalRankPage = parseInt((data.mapInfo.rank - 1) / 50) + 1;
        var countryRankPage = parseInt((data.mapInfo.countryRank - 1) / 50) + 1;
        html = '<hr>';
        html += '<h4><b>Nickname : <a href="https://scoresaber.com/u/' + data.mapInfo.mapId + '" target="_blank"><font color="green">' + data.mapInfo.mapName + '</font></a> <a href="https://scoresaber.com/global?country=' + data.mapInfo.country + '" target="_blank"><font color="blue">(' + data.mapInfo.country + ')</font></a></b></h4>';
        html += '<h4><b>Global Rank : <a href="https://scoresaber.com/global/' + globalRankPage + '" target="_blank"><font color="red">' + data.mapInfo.rank + '</font></a>' + ' (<a href="https://scoresaber.com/global/' + countryRankPage + '&country=' + data.mapInfo.country + '" target="_blank"><font color="orange">' + data.mapInfo.countryRank + '</font></a>)</b></h4>';
        html += '<h4><b>PP : ' + data.mapInfo.pp + ', Avg Accuracy : ' + data.scoreStats.averageRankedAccuracy.toFixed(2) + '</b></h4>';
        html += '<h4><b>Play Count(Rank Count) : ' + data.scoreStats.totalPlayCount + '(' + data.scoreStats.rankedPlayCount + ')</b></h4>';
        $("#dvtitleArea").empty().append(html).trigger("create");
    }

    // #region Call APIs
    obj.mapApiCall = function (apiID) {
        var apiObj = obj.apiConfig.find(x => x.id === apiID);
        
        // #region api validation check
        // #endregion

        // set api parameters 
        var requestData = setAPIParameters(apiObj, obj);
        if(requestData == null) return;

        // #region request parameter 재정의
        // #endregion

        // confirm alert 설정
        if(apiObj.message.confirm != undefined && apiObj.message.confirm != ""){
            if (!confirm(apiObj.message.confirm)) {
                return;
            }
        }

        // api 성공 시 callback 함수
        var fnSuccess = function (data) {
            eval(apiObj.callback + "(data);");
        };

        // api 실패 시 callback 함수
        var fnError = function (err) {
            panelAlertLayer(apiObj.message.error + " > " + err, "error");
        };

        // #region 개발자 영역 > API callback 함수
		var callback_get_full = function(data) {
            obj.statusAllSearch = "ready";
            // obj.max_page = (data.scoreStats.totalPlayCount / 8).toFixed(0); 
            obj.max_page = parseInt((data.scoreStats.totalPlayCount - 1) / 8) + 1;
            obj.total_count = data.scoreStats.totalPlayCount
            obj.displayData(data);
		};

        // #region 개발자 영역 > API callback 함수
		var callback_get_map = function(data) {
            obj.dataCount = (data.scores != undefined) ? data.scores.length : 0;
            if (obj.dataCount > 0) {
                obj.bindGrid(data.scores);
            } else {
                return null
            }
		};


        // #region 개발자 영역 > API callback 함수
		var callback_get_all_map = function(data) {
            obj.dataCount = (data.scores != undefined) ? data.scores.length : 0;
            if (obj.dataCount > 0) {
                obj.bindGrid(data.scores);
            } else {
                return null
            }

            if (obj.statusAllSearch == "ready" || obj.statusAllSearch == "process" || obj.statusAllSearch == "resume") {
                if (obj.pageIndex <= obj.max_page) {
                    obj.statusAllSearch = "process";
                    obj.mapApiCall("get_all_map");
                } else {
                    obj.statusAllSearch = "complete";
                    return null
                }
            } else if (obj.statusAllSearch == "stop") {
                obj.statusAllSearch = "resume";
                return null;
            } else {

            }

            // if (obj.pageIndex <= obj.max_page) {
            //     obj.statusAllSearch = "process";
            //     obj.mapApiCall("get_all_map");
            // }
		};

        // #endregion
        
        ajaxCall(apiObj.path, requestData, fnSuccess, fnError, apiObj.method, apiObj.content_type, apiObj.response_content_type);
    };
    // #endregion
}