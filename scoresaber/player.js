function Player(controlData) {

    var obj = this;
    obj.pagePlayerLayer = $(controlData.pagePlayerLayer);

    obj.menuConfig = config.menu.find(x => x.id === "player");
    obj.gridConfig = obj.menuConfig.contents_area.grid;
    obj.apiConfig = config.apis.player;
    obj.checkedColNames = [];

    // #region 개발자 영역 > 멤버변수 선언 및 초기화
    obj.selectedRankIDs = []; // Grid에서 bulk 처리 시 사용.
    obj.pageIndex = 1; // list 페이징 처리를 page_no로 처리할 경우 사용.
    obj.dataCount = 0; // list 조회 시 전체 카운트 제공 시 사용.
    obj.total_count = 0;
    obj.statusAllSearch = "none";
    obj.country = "KR";
    // #endregion

    // 서비스 활성화를 위한 초기화 설정을 담당한다.
    obj.init = function () {
        obj.createPlayerLayer();
        obj.pagePlayerLayer.show();
        obj.txtUserID.val(getUrlParameter('user'));
    };

    // 서비스를 비활성화를 위한 리셋 설정을 담당한다.
    obj.reset = function () {
        // #region 개발자 영역 > 멤버변수 초기화
        obj.selectedRankIDs = [];
        obj.pageIndex = 1;
        obj.dataCount = 0;
        obj.total_count = 0;
        obj.statusAllSearch = "none";
        // #endregion

        obj.checkedColNames = [];
        obj.pagePlayerLayer.empty().hide();
    };

    // search 영역과 contents 영역의 UI를 생성하고, control을 생성하고 이벤트를 등록한다. 
    obj.createPlayerLayer = function () {

        var html = displayContentsLayer(obj.menuConfig);

        obj.pagePlayerLayer.empty().append(html).trigger("create");

        obj.defineElements();
        obj.defineElementsEvent();
    };

    // config 파일의 root > menu > {search_area | contents_area} > controls 항목들에 대한 변수를 정의한다.
    obj.defineElements = function() {
		obj.txtUserID = $("#txtUserID");
		obj.selSearchType = $("#selSearchType");
		obj.btnSearch = $("#btnSearch");
		obj.btnAllSearch = $("#btnAllSearch");
		obj.playerContainer = $("#playerContainer");
		obj.playerPager = $("#playerPager");
    };
    
    obj.defineElementsEvent = function() {
        // #region 개발자 영역 > Config 정의한 Controls 이벤트 등록
        // config 파일의 root > menu > {search_area | contents_area} > controls 항목들에 대한 이벤트 함수를 등록한다.
		obj.selSearchType.change(function(){
		});
		obj.btnSearch.click(function () {
            obj.pageIndex = 1;

            // if (obj.statusAllSearch == "none") {
                obj.playerApiCall("get_full"); 
            // }
			
            obj.setGrid(); //Grid 조회 시 주석 해제
            obj.playerApiCall(config.controls.find(x => x.id === this.id).api);
		});
		obj.btnAllSearch.click(function () {
            
            if (obj.statusAllSearch == "none") {
                obj.playerApiCall("get_full");
                panelAlertLayer("Please click again.", "info");
            } else if (obj.statusAllSearch == "ready") {
                $(this).text("Stop").button("refresh");
                obj.setGrid(); //Grid 조회 시 주석 해제
                obj.pageIndex = 1;
			    obj.playerApiCall(config.controls.find(x => x.id === this.id).api);
            } else if (obj.statusAllSearch == "process") {
                $(this).text("Resume").button("refresh");
                obj.statusAllSearch = "stop";
                //obj.max_page = 0;
            } else if (obj.statusAllSearch == "resume" || obj.statusAllSearch == "stop") {
                $(this).text("Stop").button("refresh");
			    obj.playerApiCall(config.controls.find(x => x.id === this.id).api);
            } else if (obj.statusAllSearch == "complete") {
                $(this).text("All Search").button("refresh");
            } else {
                
            }
			// if (obj.total_count == 0) {
            //     panelAlertLayer("프로필 검색 버튼을 누른 후 검색해 주세요.", "info");
            // } else {
            //     obj.setGrid(); //Grid 조회 시 주석 해제
            //     obj.pageIndex = 1;
			//     obj.playerApiCall(config.controls.find(x => x.id === this.id).api);
            // }
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
        obj.playerContainer.jqGrid('clearGridData').jqGrid({
            datatype: "local",
            colNames: columnInfo[0],
            colModel: columnInfo[1],
            pager: "#playerPager",
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
			if(options.download != "") {
                return "<a href='https://api.beatsaver.com/download/key/" + cellValue + "'>" + cellValue + "</a>";
            } else {
                return "";
            }
        };
        $.fn.fmatter.convertRank = function (cellValue,rowObject,options) {
            //var page = parseInt((cellValue / 12) + 1);
            var page = parseInt((cellValue - 1) / 12) + 1;
			return "<a href='https://scoresaber.com/leaderboard/" + options.leaderboardId + "?page=" + page + "' target='_blank'>" + cellValue + "</a>";
        };
        $.fn.fmatter.convertSongTitle = function (cellValue,rowObject,options) {
			return "<a href='https://scoresaber.com/leaderboard/" + options.leaderboardId + "?page=1&countries=" + obj.country + "' target='_blank'>" + cellValue + "</a>";
        };
        $.fn.fmatter.convertPP = function (cellValue,rowObject,options) {
			return ConvertToString("pp", options);
        };
        $.fn.fmatter.convertMyPP = function (cellValue,rowObject,options) {
			return ConvertToString("mypp", options);
        };
        $.fn.fmatter.convertMaxPP = function (cellValue,rowObject,options) {
			return ConvertToString("maxPP", options);
        };
        $.fn.fmatter.convertWeight = function (cellValue,rowObject,options) {
			return ConvertToString("weight", options);
        };
        $.fn.fmatter.convertDifficulty = function (cellValue,rowObject,options) {
			return ConvertToString("difficulty", options);
        };
        $.fn.fmatter.convertAccuracy = function (cellValue,rowObject,options) {
			return ConvertToString("accuracy", options) + "%";
        };
        $.fn.fmatter.convertRating = function (cellValue,rowObject,options) {
			return ConvertToString("rating", options) + "%";
        };
        $.fn.fmatter.convertSongTime = function (cellValue,rowObject,options) {
			return ConvertToString("durationSeconds", options);
        };
        $.fn.fmatter.convertTime = function (cellValue,rowObject,options) {
			return ConvertToString("time", options);
        };
        $.fn.fmatter.convertYoutube = function (cellValue,rowObject,options) {
            var title = "Beat Saber " + options.songName + " - " + options.songAuthorName + " by " + options.levelAuthorName + " " 
                options.stars + " " + ConvertToString("difficulty", options) + " " + ConvertToString("accuracy", options) + "% #" + options.rank + " " + options.beatSaverKey;
			return title;
        };
        
        // #endregion
        
        obj.navGrid();
    }

    // Grid Data Binding
    obj.bindGrid = function(data) {
		for(var i=0; i<data.length; i++) {
            var mapInfo = mapList.find(x => x.uid === data[i].leaderboardId);
            if(mapInfo != undefined && mapInfo != null) {
                data[i].beatSaverKey = mapInfo.beatSaverKey;
                data[i].bpm = mapInfo.bpm;
                data[i].download = mapInfo.download;
                data[i].durationSeconds = mapInfo.durationSeconds;
                data[i].njs = mapInfo.njs;
                data[i].noteCount = mapInfo.noteCount;
                data[i].maxPP = mapInfo.pp;
                data[i].rating = mapInfo.rating.toFixed(2) + "%";
                data[i].stars = mapInfo.stars + "★";
            } else {
                data[i].beatSaverKey = "";
                data[i].bpm = 0;
                data[i].download = "";
                data[i].durationSeconds = 0;
                data[i].njs = 0;
                data[i].noteCount = 0;
                data[i].maxPP = 0;
                data[i].rating = 0 + "%";
                data[i].stars = 0 + "★";
            }

            data[i].accuracy = (data[i].maxScore > 0) ? 
                ((data[i].score * 100) / data[i].maxScore).toFixed(2) + "%" : 
                0;
            data[i].mypp = (data[i].pp * data[i].weight).toFixed(2);
            // data[i].diff = (data[i].difficulty == 1) ? "Easy" : 
            //     (data[i].difficulty == 3) ? "Normal" : 
            //     (data[i].difficulty == 5) ? "Hard" : 
            //     (data[i].difficulty == 7) ? "Expert" : 
            //     (data[i].difficulty == 9) ? "Expert+" : 
            //     data[i].difficulty;
			// #region 개발자 영역 > paging(more) 처리를 위한 셋팅
			if(i == data.length - 1) {
				obj.pageIndex++;
			}
			// #endregion
			// #region 개발자 영역 > Grid에 UUID Mapping
			obj.playerContainer.jqGrid("addRowData", data[i].id, data[i]);
			// #endregion
		}
		var gridCount = obj.playerContainer.getGridParam("reccount");
		var percentage = (obj.dataCount != 0) ? parseFloat((gridCount*100) / obj.total_count).toFixed(2) : 0;
		var gridCaption = "Searched <font color='blue'><b>" + gridCount + "</b></font> out of <font color='red'><b>" + obj.total_count + "</b></font> - <font color='green'><b>" + percentage + "%</b></font> Completed.";
		obj.playerContainer.jqGrid("setCaption", gridCaption);
    }

    // Grid 하단 Navigation Bar 버튼 및 이벤트 등록
    obj.navGrid = function() {
        var toolbar = obj.gridConfig.toolbar;
        var ctrlToolbarID = "#" + toolbar.id;
        obj.playerContainer.jqGrid("navGrid", ctrlToolbarID, {
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
            obj.playerContainer.jqGrid("navButtonAdd", ctrlToolbarID, toolbar.buttons[i]);
            $("#" + toolbar.buttons[i].id).off("click").click(function () {
                eval("javascirpt:" + toolbar.buttons.find(x => x.id === this.id).callback + "()");
            });
        }
        
        // #region 개발자 영역 > Grid navButton Callback 구현부
        
		var callback_navButton_search = function() {
            obj.playerContainer[0].toggleToolbar();
		};
		var callback_navButton_download = function() {
            obj.selectedRankIDs = obj.playerContainer.getGridParam("selarrrow");
            if (obj.selectedRankIDs.length == 0) {
                panelAlertLayer("Please select the data to download with the checkbox.", "info");
                return;
            }

            var gridData = obj.playerContainer.getGridParam('data');
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
            obj.playerApiCall("get_player");
		};

        // #endregion
    }

    obj.displayData = function (data) {
        //"<a href='https://scoresaber.com/u/76561198830502286" + page + "' target='_blank'>" + cellValue + "</a>"
        //https://scoresaber.com/global/2&country=kr
        var globalRankPage = parseInt((data.playerInfo.rank - 1) / 50) + 1;
        var countryRankPage = parseInt((data.playerInfo.countryRank - 1) / 50) + 1;
        html = '<hr>';
        html += '<h4><b>Nickname : <a href="https://scoresaber.com/u/' + data.playerInfo.playerId + '" target="_blank"><font color="green">' + data.playerInfo.playerName + '</font></a> <a href="https://scoresaber.com/rankings?page=1&countries=' + data.playerInfo.country + '" target="_blank"><font color="blue">(' + data.playerInfo.country + ')</font></a></b></h4>';
        html += '<h4><b>Global Rank : <a href="https://scoresaber.com/rankings?page=' + globalRankPage + '" target="_blank"><font color="red">' + data.playerInfo.rank + '</font></a>' + ' (<a href="https://scoresaber.com/rankings?page=' + countryRankPage + '&countries=' + data.playerInfo.country + '" target="_blank"><font color="orange">' + data.playerInfo.countryRank + '</font></a>)</b></h4>';
        html += '<h4><b>PP : ' + data.playerInfo.pp + ', Avg Accuracy : ' + data.scoreStats.averageRankedAccuracy.toFixed(2) + '</b></h4>';
        html += '<h4><b>Play Count(Rank Count) : ' + data.scoreStats.totalPlayCount + '(' + data.scoreStats.rankedPlayCount + ')</b></h4>';
        $("#dvtitleArea").empty().append(html).trigger("create");
    }

    // #region Call APIs
    obj.playerApiCall = function (apiID) {
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
            obj.country = data.playerInfo.country;
            obj.displayData(data);
		};

        // #region 개발자 영역 > API callback 함수
		var callback_get_player = function(data) {
            obj.dataCount = (data.scores != undefined) ? data.scores.length : 0;
            if (obj.dataCount > 0) {
                obj.bindGrid(data.scores);
            } else {
                return null
            }
		};


        // #region 개발자 영역 > API callback 함수
		var callback_get_all_player = function(data) {
            obj.dataCount = (data.scores != undefined) ? data.scores.length : 0;
            if (obj.dataCount > 0) {
                obj.bindGrid(data.scores);
            } else {
                return null
            }

            if (obj.statusAllSearch == "ready" || obj.statusAllSearch == "process" || obj.statusAllSearch == "resume") {
                if (obj.pageIndex <= obj.max_page) {
                    obj.statusAllSearch = "process";
                    obj.playerApiCall("get_all_player");
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
            //     obj.playerApiCall("get_all_player");
            // }
		};

        // #endregion
        
        ajaxCall(apiObj.path, requestData, fnSuccess, fnError, apiObj.method, apiObj.content_type, apiObj.response_content_type);
    };
    // #endregion
}