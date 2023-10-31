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
    obj.loopCnt = 1;
    obj.statusSearch = "none";
    // #endregion

    // 서비스 활성화를 위한 초기화 설정을 담당한다.
    obj.init = function () {
        obj.createRankingLayer();
        obj.pageRankingLayer.show();
        //obj.selPageSize.val(getUrlParameter('page_size'));
        obj.selPageIndex.val(getUrlParameter('page_index'));
        obj.rankingApiCall("pages");
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
		obj.selPageIndex = $("#selPageIndex");
		obj.selPageSize = $("#selPageSize");
		obj.btnRankSearch = $("#btnRankSearch");
		obj.rankingContainer = $("#rankingContainer");
		obj.rankingPager = $("#rankingPager");
    };
    
    obj.defineElementsEvent = function() {
        // #region 개발자 영역 > Config 정의한 Controls 이벤트 등록
        // config 파일의 root > menu > {search_area | contents_area} > controls 항목들에 대한 이벤트 함수를 등록한다.
		obj.selPageSize.change(function(){
            //obj.btnRankSearch.click();
        });

        obj.selPageIndex.change(function(){
            obj.btnRankSearch.click();
        });
        
		obj.btnRankSearch.click(function () {
			obj.setGrid(); //Grid 조회 시 주석 해제
            obj.loopCnt = 1;
            obj.pageSize = parseInt(obj.selPageSize.val());
            if (obj.pageSize == 1) {
                obj.pageIndex = (obj.pageSize * parseInt(obj.selPageIndex.val()));
            } else {
                obj.pageIndex = (obj.pageSize * (parseInt(obj.selPageIndex.val())-1)) + 1;
            }
            
            obj.statusSearch = "none";
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
        $.fn.fmatter.convertPlayerName = function (cellValue,rowObject,options) {
			return '<a href="https://scoresaber.com/u/' + options.id + '" target="_blank"><font color="green">' + options.name + '</font></a>';
        };
        $.fn.fmatter.convertAvartar = function (cellValue,rowObject,options) {
			return "<img src='" + options.profilePicture + "' border=0 width=24 height=24 />";
        };
        $.fn.fmatter.convertCountry = function (cellValue,rowObject,options) {
			return "<img src='https://new.scoresaber.com/api/static/flags/" + options.country.toLowerCase() + ".png' title=" + options.country + " border=0  width=16 height=11/>";
        };
        $.fn.fmatter.convertPeriodRankGap = function (cellValue,rowObject,options) {
            var result = "";
            if (cellValue < 0) {
                result = "▼" + Math.abs(cellValue);
            } else if (cellValue > 0) {
                result = "▲" + cellValue;
            } else if (cellValue == 0) {
                result = "-";
            }
            return result;
        };
        $.fn.fmatter.convertMonthRankGap = function (cellValue,rowObject,options) {
            var result = "";
            var splitHistory = options.histories.split(",");
            if (splitHistory[splitHistory.length - 30] > 900000) {
                result = "NEW";
            } else {
                var monthGap = splitHistory[splitHistory.length - 30] - options.rank;
                if (monthGap < 0) {
                    result = "▼" + Math.abs(monthGap);
                } else if (monthGap > 0) {
                    result = "▲" + monthGap;
                } else if (monthGap == 0) {
                    result = "-";
                }
            }
            return result;
        };
        $.fn.fmatter.convertDailyRankGap = function (cellValue,rowObject,options) {
            var result = "";
            var splitHistory = options.histories.split(",");
            if (splitHistory[splitHistory.length - 2] > 900000) {
                result = "NEW";
            } else {
                var dailyGap = splitHistory[splitHistory.length-2] - options.rank;
                if (dailyGap < 0) {
                    result = "▼" + Math.abs(dailyGap);
                } else if (dailyGap > 0) {
                    result = "▲" + dailyGap;
                } else if (dailyGap == 0) {
                    result = "-";
                }
            }
            return result;
        };
        // #endregion
        
        obj.navGrid();
    }

    // Grid Data Binding
    obj.bindGrid = function(data) {
		for(var i=0; i<data.length; i++) {
            var splitHistory = data[i].histories.split(",");
            //data[i].weekly_gap = data[i].difference;
            if (splitHistory[splitHistory.length - 7] > 900000) {
                data[i].weekly_gap = "NEW";
            } else {
                data[i].weekly_gap = splitHistory[splitHistory.length-7] - data[i].rank;
            }
            if (splitHistory[splitHistory.length - 1] > 900000) {
                data[i].daily_gap = "NEW";
            } else {
                data[i].daily_gap = splitHistory[splitHistory.length-1] - data[i].rank;
            }
            if (splitHistory.length < 30) {
                data[i].monthly_gap = "NEW";
            } else if (splitHistory[splitHistory.length - 30] > 900000) {
                data[i].monthly_gap = "COMEBACK";
            } else {
                data[i].monthly_gap = splitHistory[splitHistory.length - 30] - data[i].rank;
            }
        
            // #region 개발자 영역 > paging(more) 처리를 위한 셋팅
			if(i == data.length - 1) {
				obj.pageIndex++;
			}
			// #endregion
			// #region 개발자 영역 > Grid에 UUID Mapping
			obj.rankingContainer.jqGrid("addRowData", data[i].rank, data[i]);
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
		var callback_get_pages = function(data) {
            obj.dataCount = (data.metadata != undefined) ? data.metadata.total : 0;
            obj.total_count = data.pages*50;
            var optionHtml = "";
            for (var i=2;i<=obj.dataCount;i++) {
                if(i <= 10 || i%5 == 0) {
                    optionHtml += '<option value="' + i + '">' + i + ' page</option>';
                }
            }
            obj.selPageIndex.append(optionHtml).trigger("create");;
        };
        
		var callback_get_ranking = function(data) {
            if (obj.dataCount > 0) {
                obj.bindGrid(data.players);
                
                if(obj.loopCnt < parseInt(obj.selPageSize.val())){
                    obj.rankingApiCall("get_ranking");
                    obj.loopCnt++;
                }
            } else {
                return null
            }
		};

        // #endregion
        
        ajaxCall(apiObj.path, requestData, fnSuccess, fnError, apiObj.method, apiObj.content_type, apiObj.response_content_type);
    };
    // #endregion
}