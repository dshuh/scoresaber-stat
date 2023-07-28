function Ranking(controlData) {

    var obj = this;
    obj.pageRankingLayer = $(controlData.pageRankingLayer);

    obj.menuConfig = config.menu.find(x => x.id === "ranking");
    obj.gridConfig = obj.menuConfig.contents_area.grid;
    obj.apiConfig = config.apis.ranking;
    obj.checkedColNames = [];

    // #region 개발자 영역 > 멤버변수 선언 및 초기화
    obj.selectedRankIDs = []; // Grid에서 bulk 처리 시 사용.
    obj.pageIndex = 1; // list 페이징 처리를 page_no로 처리할 경우 사용.
    obj.pageSize = 50; // list 페이징 처리를 page_no로 처리할 경우 사용.
    obj.dataCount = 0; // list 조회 시 전체 카운트 제공 시 사용.
    obj.total_count = 0;
    obj.statusSearch = "none";
    // #endregion

    // 서비스 활성화를 위한 초기화 설정을 담당한다.
    obj.init = function () {
        obj.createRankingLayer();
        obj.pageRankingLayer.show();
        //obj.selPageSize.val(getUrlParameter('page_size'));
        // obj.rankingApiCall("pages");
    };

    // 서비스를 비활성화를 위한 리셋 설정을 담당한다.
    obj.reset = function () {
        // #region 개발자 영역 > 멤버변수 초기화
        obj.selectedRankIDs = [];
        obj.pageIndex = 1;
        obj.dataCount = 0;
        obj.total_count = 0;
        obj.loopCnt = 1;
        obj.statusSearch = "none";
        // #endregion

        obj.checkedColNames = [];
        obj.pageRankingLayer.empty().hide();
    };

    // search 영역과 contents 영역의 UI를 생성하고, control을 생성하고 이벤트를 등록한다. 
    obj.createRankingLayer = function () {

        var html = displayContentsLayer(obj.menuConfig);

        obj.pageRankingLayer.empty().append(html).trigger("create");

        obj.defineElements();
        obj.defineElementsEvent();
    };

    // config 파일의 root > menu > {search_area | contents_area} > controls 항목들에 대한 변수를 정의한다.
    obj.defineElements = function() {
		obj.selCountryList = $("#selCountryList");
		obj.selRankSortBy = $("#selRankSortBy");
		obj.selRankOrder = $("#selRankOrder");
		obj.btnRankSearch = $("#btnRankSearch");
		obj.rankingContainer = $("#rankingContainer");
		obj.rankingPager = $("#rankingPager");
    };
    
    obj.defineElementsEvent = function() {
        // #region 개발자 영역 > Config 정의한 Controls 이벤트 등록
        // config 파일의 root > menu > {search_area | contents_area} > controls 항목들에 대한 이벤트 함수를 등록한다.
		
        obj.selCountryList.change(function(){
            obj.btnRankSearch.click();
		});
        obj.selRankSortBy.change(function(){
            obj.btnRankSearch.click();
		});
        obj.selRankOrder.change(function(){
            obj.btnRankSearch.click();
		});

		obj.btnRankSearch.click(function () {
			obj.pageIndex = 1;

            obj.setGrid(); //Grid 조회 시 주석 해제
            obj.rankingApiCall("get_ranking");
		});
		// obj.btnSearch.click(function () {
		// 	obj.setGrid(); //Grid 조회 시 주석 해제
            
        //     obj.pageIndex = 1;
        //     obj.rankingApiCall(config.controls.find(x => x.id === this.id).api);
		// });
		// obj.btnAllSearch.click(function () {
        //     if (obj.statusAllSearch == "none") {
        //         panelAlertLayer("Please search after clicking the Search Profile button.", "info");
        //     } else if (obj.statusAllSearch == "ready") {
        //         $(this).text("Stop").button("refresh");
        //         obj.setGrid(); //Grid 조회 시 주석 해제
        //         obj.pageIndex = 1;
		// 	    obj.rankingApiCall(config.controls.find(x => x.id === this.id).api);
        //     } else if (obj.statusAllSearch == "process") {
        //         $(this).text("Resume").button("refresh");
        //         obj.statusAllSearch = "stop";
        //         //obj.max_page = 0;
        //     } else if (obj.statusAllSearch == "resume") {
        //         $(this).text("Stop").button("refresh");
		// 	    obj.rankingApiCall(config.controls.find(x => x.id === this.id).api);
        //     } else if (obj.statusAllSearch == "complete") {
        //         $(this).text("All Search").button("refresh");
        //     } else {
                
        //     }
		// });

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
        obj.rankingContainer.jqGrid('clearGridData').jqGrid({
            datatype: "local",
            colNames: columnInfo[0],
            colModel: columnInfo[1],
            pager: "#rankingPager",
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
        $.fn.fmatter.convertGap = function (cellValue,rowObject,options) {
            return cellValue > 0 ? '▲'+ Math.abs(cellValue) : cellValue < 0 ? '▼'+ Math.abs(cellValue) : '-';
        };
        $.fn.fmatter.convertPP = function (cellValue,rowObject,options) {
			if(cellValue != 0) {
                return cellValue.toFixed(2) + "pp";
            } else {
                return "-";
            }
        };
        $.fn.fmatter.convertAccuracy = function (cellValue,rowObject,options) {
			return (cellValue == 0) ? "-" : cellValue.toFixed(2) +"%";
        };
        $.fn.fmatter.convertTime = function (cellValue,rowObject,options) {
			return ConvertToString("timestamp", cellValue);
        };
        $.fn.fmatter.convertPlayerName = function (cellValue,rowObject,options) {
			return '<a href="https://www.beatleader.xyz/u/' + options.id + '" target="_blank"><font color="green">' + cellValue + '</font></a>';
        };
        $.fn.fmatter.convertAvartar = function (cellValue,rowObject,options) {
			return "<img src='" + cellValue + "' border=0 width=24 height=24 />";
        };
        $.fn.fmatter.convertClans = function (cellValue,rowObject,options) {
			var clans = '';
            for(var i=0;i<cellValue.length;i++) {
                clans += '<a href="https://www.beatleader.xyz/clan/' + cellValue[i].tag + '" target="_blank" style="bgcolor:black;"><font color="' + cellValue[i].color + '">' + cellValue[i].tag + '</font></a> ';
            }
            return clans;
        };
        $.fn.fmatter.convertCountry = function (cellValue,rowObject,options) {
			return "<img src='https://www.beatleader.xyz/assets/flags/" + cellValue.toLowerCase() + ".png' title=" + cellValue + " border=0  width=16 height=11/>";
        };
        // #endregion
        
        obj.navGrid();
    }

    // Grid Data Binding
    obj.bindGrid = function(data) {
        console.log(data.length);
		for(var i=0; i<data.length; i++) {
            // var globalRankPage = parseInt((data[i].rank - 1) / 50) + 1;
            // var countryRankPage = parseInt((data[i].countryRank - 1) / 50) + 1;
            data[i].globalRankGap = data[i].lastWeekRank-data[i].rank;
            data[i].countryRankGap = data[i].lastWeekCountryRank-data[i].countryRank;
            data[i].globalPPGap = (data[i].pp-data[i].lastWeekPp).toFixed(2);
            // var clans = '';
            
            // data[i].clans = clans;
            data[i].averageRank = data[i].scoreStats.averageRank.toFixed(2);
            data[i].accuracy = (data[i].scoreStats.averageRankedAccuracy*100);
            data[i].weightedAccuracy = (data[i].scoreStats.averageWeightedRankedAccuracy*100);
            data[i].topPp = data[i].scoreStats.topPp;
            data[i].maxStreak = data[i].scoreStats.maxStreak;
            data[i].watchedReplays = data[i].scoreStats.anonimusReplayWatched;
            data[i].sspPlays = data[i].scoreStats.sspPlays;
            data[i].ssPlays = data[i].scoreStats.ssPlays;
            data[i].spPlays = data[i].scoreStats.spPlays;
            data[i].sPlays = data[i].scoreStats.sPlays;
            data[i].aPlays = data[i].scoreStats.aPlays;
            data[i].lastScoreTime = data[i].scoreStats.lastScoreTime.toString();
            data[i].totalPlayCount = data[i].scoreStats.totalPlayCount;
            data[i].rankedPlayCount = data[i].scoreStats.rankedPlayCount;
            
            // #region 개발자 영역 > paging(more) 처리를 위한 셋팅
			if(i == data.length - 1) {
				obj.pageIndex++;
			}
			// #endregion
			// #region 개발자 영역 > Grid에 UUID Mapping
			obj.rankingContainer.jqGrid("addRowData", data[i].id, data[i]);
			// #endregion
		}
		var gridCount = obj.rankingContainer.getGridParam("reccount");
		var percentage = (obj.dataCount != 0) ? parseFloat((gridCount*100) / obj.total_count).toFixed(2) : 0;
		var gridCaption = "Searched <font color='blue'><b>" + gridCount + "</b></font> out of <font color='red'><b>" + obj.total_count + "</b></font> - <font color='green'><b>" + percentage + "%</b></font> Completed.";
        obj.rankingContainer.jqGrid("setCaption", gridCaption);
    }

    // Grid 하단 Navigation Bar 버튼 및 이벤트 등록
    obj.navGrid = function() {
        var toolbar = obj.gridConfig.toolbar;
        var ctrlToolbarID = "#" + toolbar.id;
        obj.rankingContainer.jqGrid("navGrid", ctrlToolbarID, {
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
            obj.rankingContainer.jqGrid("navButtonAdd", ctrlToolbarID, toolbar.buttons[i]);
            $("#" + toolbar.buttons[i].id).off("click").click(function () {
                eval("javascirpt:" + toolbar.buttons.find(x => x.id === this.id).callback + "()");
            });
        }
        
        // #region 개발자 영역 > Grid navButton Callback 구현부
        
		var callback_navButton_search = function() {
            obj.rankingContainer[0].toggleToolbar();
		};
		var callback_navButton_download = function() {
            obj.selectedRankIDs = obj.rankingContainer.getGridParam("selarrrow");
            if (obj.selectedRankIDs.length == 0) {
                panelAlertLayer("Please select the data to download with the checkbox.", "info");
                return;
            }

            var gridData = obj.rankingContainer.getGridParam('data');
            var downloadRank = [];
            for( var i=0;i<obj.selectedRankIDs.length;i++) {            
                downloadRank.push(gridData.find(x => x.id === obj.selectedRankIDs[i]));
            }

            var str = JSON.stringify(downloadRank);
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
            obj.loopCnt = 1;
            obj.rankingApiCall("get_ranking");
		};

        // #endregion
    }

    // #region Call APIs
    obj.rankingApiCall = function (apiID) {
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
        
		var callback_get_ranking = function(data) {
            obj.dataCount = (data != undefined) ? data.data.length : 0;
            obj.max_page = parseInt((data.metadata.total - 1) / data.metadata.itemsPerPage) + 1;
            obj.total_count = data.metadata.total;//(data.metadata.total * data.metadata.itemsPerPage);
            if (obj.dataCount > 0) {
                obj.bindGrid(data.data);
            } else {
                return null
            }
		};

        // #endregion
        
        ajaxCall(apiObj.path, requestData, fnSuccess, fnError, apiObj.method, apiObj.content_type, apiObj.response_content_type);
    };
    // #endregion
}