function Compare(controlData) {

    var obj = this;
    obj.pageCompareLayer = $(controlData.pageCompareLayer);

    obj.menuConfig = config.menu.find(x => x.id === "compare");
    obj.gridConfig = obj.menuConfig.contents_area.grid;
    obj.apiConfig = config.apis.compare;
    obj.checkedColNames = [];

    // #region 개발자 영역 > 멤버변수 선언 및 초기화
    obj.selectedRankIDs = []; // Grid에서 bulk 처리 시 사용.
    obj.pageIndex1 = 1; // list 페이징 처리를 page_no로 처리할 경우 사용.
    obj.pageIndex2 = 1; // list 페이징 처리를 page_no로 처리할 경우 사용.
    obj.dataCount1 = 0; // list 조회 시 전체 카운트 제공 시 사용.
    obj.dataCount2 = 0; // list 조회 시 전체 카운트 제공 시 사용.
    obj.total_count = 0;
    obj.statusAllSearch1 = "none";
    obj.statusAllSearch2 = "none";
    obj.country = "KR";

    obj.p1Data = [];
    obj.p2Data = [];
    // #endregion

    // 서비스 활성화를 위한 초기화 설정을 담당한다.
    obj.init = function () {
        obj.createCompareLayer();
        obj.pageCompareLayer.show();
        // obj.txtUserID1.val(getUrlParameter('user1'));
        // obj.txtUserID2.val(getUrlParameter('user2'));
    };

    // 서비스를 비활성화를 위한 리셋 설정을 담당한다.
    obj.reset = function () {
        // #region 개발자 영역 > 멤버변수 초기화
        obj.selectedRankIDs = [];
        obj.pageIndex1 = 1;
        obj.pageIndex2 = 1;
        obj.dataCount1 = 0;
        obj.dataCount2 = 0;
        obj.total_count = 0;
        obj.statusAllSearch1 = "none";
        obj.statusAllSearch2 = "none";
        // #endregion

        obj.checkedColNames = [];
        obj.pageCompareLayer.empty().hide();
    };

    // search 영역과 contents 영역의 UI를 생성하고, control을 생성하고 이벤트를 등록한다. 
    obj.createCompareLayer = function () {

        var html = displayContentsLayer(obj.menuConfig);

        obj.pageCompareLayer.empty().append(html).trigger("create");

        obj.defineElements();
        obj.defineElementsEvent();
    };

    // config 파일의 root > menu > {search_area | contents_area} > controls 항목들에 대한 변수를 정의한다.
    obj.defineElements = function() {
		obj.selUserList1 = $("#selUserList1");
		obj.txtUserID1 = $("#txtUserID1");
		obj.selUserList2 = $("#selUserList2");
		obj.txtUserID2 = $("#txtUserID2");
        obj.txtGetPages = $("#txtGetPages");
		obj.btnCompareSearch = $("#btnCompareSearch");
		obj.compareContainer = $("#compareContainer");
		obj.comparePager = $("#comparePager");
    };
    
    obj.defineElementsEvent = function() {
        // #region 개발자 영역 > Config 정의한 Controls 이벤트 등록
        // config 파일의 root > menu > {search_area | contents_area} > controls 항목들에 대한 이벤트 함수를 등록한다.
		obj.selUserList1.change(function(){
            obj.txtUserID1.val(this.value);
		});
        obj.selUserList1.change();
        obj.selUserList2.change(function(){
            obj.txtUserID2.val(this.value);
		});
        obj.selUserList2.change();
        obj.btnCompareSearch.click(function () {
            obj.setGrid(); //Grid 조회 시 주석 해제
            obj.pageIndex1 = 1;
            obj.pageIndex2 = 1;
            obj.dataCount1 = 0;
            obj.dataCount2 = 0;
            obj.statusAllSearch1 = "none";
            obj.statusAllSearch2 = "none";
            obj.p1Data = [];
            obj.p2Data = [];
            $("#dvtitleArea").empty().trigger("create");
            obj.compareApiCall("get_full1");
            obj.compareApiCall("get_full2");
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
        obj.compareContainer.jqGrid('clearGridData').jqGrid({
            datatype: "local",
            colNames: columnInfo[0],
            colModel: columnInfo[1],
            pager: "#comparePager",
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
        $.fn.fmatter.convertWeight = function (cellValue,rowObject,options) {
			return ConvertToString("weight", options);
        };
        // $.fn.fmatter.convertDifficulty = function (cellValue,rowObject,options) {
		// 	return ConvertToString("diff", options);
        // };
        $.fn.fmatter.convertAccuracy = function (cellValue,rowObject,options) {
			return ConvertToString("accuracy", options);
        };
        $.fn.fmatter.convertTime = function (cellValue,rowObject,options) {
			return ConvertToString("time", options);
        };
        // #endregion
        
        obj.navGrid();
    }

    // Grid Data Binding
    obj.bindGrid = function() {
        var resultData = {};
        var p1WinCount = 0;
        var p1LoseCount = 0;
        var p1RankGap = 0;
        var p1AccuracyGap = 0;
        var p1PPGap = 0;
        var p1OnlyPlayCount = 0;
        var p2OnlyPlayCount = 0;
        for(var i=0; i<obj.p1Data.length; i++) {
            var isDuplicate = false;
            for(var j=0; j<obj.p2Data.length; j++) {
                if(obj.p1Data[i].leaderboard.difficulty.leaderboardId == obj.p2Data[j].leaderboard.difficulty.leaderboardId) {
                    isDuplicate = true;

                    var mapInfo = mapList.find(x => x.uid === obj.p1Data[i].leaderboard.difficulty.leaderboardId);
                    if(mapInfo != undefined && mapInfo != null) {
                        resultData["stars"] = mapInfo.stars + "★";
                    } else {
                        resultData["stars"] = 0 + "★";
                    }

                    resultData["leaderboardId"] = obj.p1Data[i].leaderboard.difficulty.leaderboardId;
                    resultData["songName"] = obj.p1Data[i].leaderboard.songName;
                    resultData["songAuthorName"] = obj.p1Data[i].leaderboard.songAuthorName;
                    resultData["diff"] = obj.p1Data[i].diff;
                    if(obj.p1Data[i].score.rank < obj.p2Data[j].score.rank) {
                        resultData["winner"] = obj.p1Nickname;
                        p1WinCount++;
                    } else {
                        resultData["winner"] = obj.p2Nickname;
                        p1LoseCount++;
                    }
                    resultData["p1_rank"] = obj.p1Data[i].score.rank;
                    resultData["p2_rank"] = obj.p2Data[j].score.rank;
                    resultData["rank_gap"] = obj.p2Data[j].score.rank - obj.p1Data[i].score.rank;
                    p1RankGap += obj.p2Data[j].score.rank - obj.p1Data[i].score.rank;

                    resultData["p1_accuracy"] = obj.p1Data[i].accuracy;
                    resultData["p2_accuracy"] = obj.p2Data[j].accuracy;
                    resultData["accuracy_gap"] = (obj.p1Data[i].accuracy - obj.p2Data[j].accuracy).toFixed(2);
                    p1AccuracyGap += parseFloat(obj.p1Data[i].accuracy - obj.p2Data[j].accuracy);
                    
                    resultData["p1_pp"] = obj.p1Data[i].score.pp;
                    resultData["p2_pp"] = obj.p2Data[j].score.pp;
                    resultData["pp_gap"] = (obj.p1Data[i].score.pp - obj.p2Data[j].score.pp).toFixed(2);
                    p1PPGap += parseFloat(obj.p1Data[i].score.pp - obj.p2Data[j].score.pp);

                    var p1Date = new Date(obj.p1Data[i].score.timeSet).toISOString().
                        replace(/T/, ' ').      // replace T with a space
                        replace(/\..+/, '');
                    var p2Date = new Date(obj.p2Data[j].score.timeSet).toISOString().
                        replace(/T/, ' ').      // replace T with a space
                        replace(/\..+/, '');
                    resultData["p1_timeSet"] = p1Date;
                    resultData["p2_timeSet"] = p2Date;
                    
                    obj.compareContainer.jqGrid("addRowData", resultData["leaderboardId"], resultData);
                }
            }
            if (!isDuplicate) {
                
                var mapInfo = mapList.find(x => x.uid === obj.p1Data[i].leaderboard.difficulty.leaderboardId);
                if(mapInfo != undefined && mapInfo != null) {
                    resultData["stars"] = mapInfo.stars + "★";
                } else {
                    resultData["stars"] = 0 + "★";
                }
                
                resultData["leaderboardId"] = obj.p1Data[i].leaderboard.difficulty.leaderboardId;
                resultData["songName"] = obj.p1Data[i].leaderboard.songName;
                resultData["songAuthorName"] = obj.p1Data[i].leaderboard.songAuthorName;
                resultData["diff"] = obj.p1Data[i].diff;
                resultData["winner"] = "-";
                p1OnlyPlayCount++;
                resultData["p1_rank"] = obj.p1Data[i].score.rank;
                resultData["p2_rank"] = 0;
                resultData["rank_gap"] = 0;

                resultData["p1_accuracy"] = obj.p1Data[i].accuracy;
                resultData["p2_accuracy"] = 0;
                resultData["accuracy_gap"] = 0;
                
                resultData["p1_pp"] = obj.p1Data[i].score.pp;
                resultData["p2_pp"] = 0;
                resultData["pp_gap"] = 0;

                var p1Date = new Date(obj.p1Data[i].score.timeSet).toISOString().
                    replace(/T/, ' ').      // replace T with a space
                    replace(/\..+/, '');
                resultData["p1_timeSet"] = p1Date;
                // var p2Date = new Date(obj.p2Data[i].score.timeSet).toISOString().
                //     replace(/T/, ' ').      // replace T with a space
                //     replace(/\..+/, '');
                resultData["p2_timeSet"] = "-";
                
                obj.compareContainer.jqGrid("addRowData", resultData["leaderboardId"], resultData);
            }
        }

        for(var i=0; i<obj.p2Data.length; i++) {
            var isDuplicate = false;
            for(var j=0; j<obj.p1Data.length; j++) {
                if(obj.p2Data[i].leaderboard.difficulty.leaderboardId == obj.p1Data[j].leaderboard.difficulty.leaderboardId) {
                    isDuplicate = true;
                }
            }
            if (!isDuplicate) {
                                
                var mapInfo = mapList.find(x => x.uid === obj.p2Data[i].leaderboard.difficulty.leaderboardId);
                if(mapInfo != undefined && mapInfo != null) {
                    resultData["stars"] = mapInfo.stars + "★";
                } else {
                    resultData["stars"] = 0 + "★";
                }

                resultData["leaderboardId"] = obj.p2Data[i].leaderboard.difficulty.leaderboardId;
                resultData["songName"] = obj.p2Data[i].leaderboard.songName;
                resultData["songAuthorName"] = obj.p2Data[i].leaderboard.songAuthorName;
                resultData["diff"] = obj.p2Data[i].diff;
                resultData["winner"] = "-";
                p2OnlyPlayCount++;
                resultData["p1_rank"] = 0;
                resultData["p2_rank"] = obj.p2Data[i].score.rank;
                resultData["rank_gap"] = 0;
                resultData["p1_accuracy"] = 0;
                resultData["p2_accuracy"] = obj.p2Data[i].accuracy;
                resultData["accuracy_gap"] = 0;
                resultData["p1_pp"] = 0;
                resultData["p2_pp"] = obj.p2Data[i].score.pp;
                resultData["pp_gap"] = 0;
                resultData["p1_timeSet"] = "-";
                var p2Date = new Date(obj.p2Data[i].score.timeSet).toISOString().
                    replace(/T/, ' ').      // replace T with a space
                    replace(/\..+/, '');
                resultData["p2_timeSet"] = p2Date;
                
                obj.compareContainer.jqGrid("addRowData", resultData["leaderboardId"], resultData);
            }
        }
        
		var gridCount = obj.compareContainer.getGridParam("reccount");
		var gridCaption = "Searched <font color='blue'><b>" + gridCount + "</b></font> Completed.";
        obj.compareContainer.jqGrid("setCaption", gridCaption);

        var html = '<hr>';
        html += '<h4><font color="green"><b>' + obj.p1Nickname + '</b></font> Summary Info(Versus : <font color="blue">' + obj.p2Nickname + '</font>)</h4>';
        html += "<h4><b><font color='green'>" + obj.p1Nickname + "</font> is <font color='blue'>" + p1WinCount + "</font> win <font color='red'>" + p1LoseCount + "</font> lose.</b></h4>";
        html += '<h4><b>- Rank Gap : ' + p1RankGap + '<br/>- Accuracy Gap : ' + p1AccuracyGap.toFixed(2) + '<br/>- PP Gap : ' + p1PPGap.toFixed(2) + '</b></h4>';
        html += "<h4>Only <b><font color='green'>" + obj.p1Nickname + "</font> Play Count : <font color='red'>" + p1OnlyPlayCount + "</font></b></h4>";
        html += "<h4>Only <b><font color='blue'>" + obj.p2Nickname + "</font> Play Count : <font color='red'>" + p2OnlyPlayCount + "</font></b></h4>";
        $("#dvtitleArea").append(html).trigger("create");
    }

    // Grid 하단 Navigation Bar 버튼 및 이벤트 등록
    obj.navGrid = function() {
        var toolbar = obj.gridConfig.toolbar;
        var ctrlToolbarID = "#" + toolbar.id;
        obj.compareContainer.jqGrid("navGrid", ctrlToolbarID, {
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
            obj.compareContainer.jqGrid("navButtonAdd", ctrlToolbarID, toolbar.buttons[i]);
            $("#" + toolbar.buttons[i].id).off("click").click(function () {
                eval("javascirpt:" + toolbar.buttons.find(x => x.id === this.id).callback + "()");
            });
        }
        
        // #region 개발자 영역 > Grid navButton Callback 구현부
        
		var callback_navButton_search = function() {
            obj.compareContainer[0].toggleToolbar();
		};
		var callback_navButton_download = function() {
            obj.selectedRankIDs = obj.compareContainer.getGridParam("selarrrow");
            if (obj.selectedRankIDs.length == 0) {
                panelAlertLayer("Please select the data to download with the checkbox.", "info");
                return;
            }

            var gridData = obj.compareContainer.getGridParam('data');
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
            obj.compareApiCall("get_player");
		};

        // #endregion
    }

    obj.displayData = function (index, data) {
        //"<a href='https://scoresaber.com/u/76561198830502286" + page + "' target='_blank'>" + cellValue + "</a>"
        //https://scoresaber.com/global/2&country=kr
        var globalRankPage = parseInt((data.rank - 1) / 50) + 1;
        var countryRankPage = parseInt((data.countryRank - 1) / 50) + 1;
        html = '<hr>';
        //html += '<h4><b><font color="red">P' + index + '</font> : <a href="https://scoresaber.com/u/' + data.playerId + '" target="_blank"><font color="green">' + data.name + '</font></a> <a href="https://scoresaber.com/rankings?page=1&countries=' + data.country + '" target="_blank"><font color="blue">(' + data.country + ')</font></a></b></h4>';
        html += '<h4><b><a href="https://scoresaber.com/u/' + data.id + '" target="_blank">'+ "<img src='" + data.profilePicture + "' width=50 height=50 boarder=0 />" + '<font color="green"> ' + data.name + '</font></a> <a href="https://scoresaber.com/rankings?page=1&countries=' + data.country + '" target="_blank"><font color="blue">(' + data.country + ')</font></a></b></h4>';
        html += '<h4><b>Global Rank : <a href="https://scoresaber.com/rankings?page=' + globalRankPage + '" target="_blank"><font color="red">' + data.rank + '</font></a>' + ' (<a href="https://scoresaber.com/rankings?page=' + countryRankPage + '&countries=' + data.country + '" target="_blank"><font color="orange">' + data.countryRank + '</font></a>)</b></h4>';
        html += '<h4><b>PP : ' + data.pp + ', Avg Accuracy : ' + data.scoreStats.averageRankedAccuracy.toFixed(2) + '</b></h4>';
        html += '<h4><b>Play Count(Rank Count) : ' + data.scoreStats.totalPlayCount + '(' + data.scoreStats.rankedPlayCount + ')</b></h4>';
        $("#dvtitleArea").append(html).trigger("create");
    }

    // #region Call APIs
    obj.compareApiCall = function (apiID) {
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
		var callback_get_full1 = function(data) {
            var pages = Number(obj.txtGetPages.val());
            obj.statusAllSearch1 = "ready";
            obj.p1Nickname = data.name;
            obj.country = data.country;
            obj.max_page1 = parseInt((data.scoreStats.totalPlayCount - 1) / 8) + 1;
            obj.max_page1 = (obj.max_page1 > pages) ? pages : obj.max_page1;
            obj.total_count = data.scoreStats.totalPlayCount
            obj.displayData(1, data);
            
            if (obj.statusAllSearch1 == "ready") {
                // obj.setGrid(); //Grid 조회 시 주석 해제
                obj.pageIndex1 = 1;
			    obj.compareApiCall("get_all_player1");
            } else if (obj.statusAllSearch1 == "process") {
                obj.statusAllSearch1 = "stop";
                //obj.max_page1 = 0;
            } else if (obj.statusAllSearch1 == "resume") {
			    obj.compareApiCall("get_all_player1");
            } else if (obj.statusAllSearch1 == "complete") {
            } else {
                
            }
        };
        
		var callback_get_full2 = function(data) {
            var pages = Number(obj.txtGetPages.val());
            obj.statusAllSearch2 = "ready";
            obj.p2Nickname = data.name;
            obj.max_page2 = parseInt((data.scoreStats.totalPlayCount - 1) / 8) + 1;
            obj.max_page2 = (obj.max_page2 > pages) ? pages : obj.max_page2;
            obj.total_count = data.scoreStats.totalPlayCount
            obj.displayData(2, data);
            
            if (obj.statusAllSearch2 == "ready") {
                //obj.setGrid(); //Grid 조회 시 주석 해제
                obj.pageIndex2 = 1;
			    obj.compareApiCall("get_all_player2");
            } else if (obj.statusAllSearch2 == "process") {
                obj.statusAllSearch2 = "stop";
                //obj.max_page2 = 0;
            } else if (obj.statusAllSearch2 == "resume") {
			    obj.compareApiCall("get_all_player2");
            } else if (obj.statusAllSearch2 == "complete") {
            } else {
                
            }
		};

        // #region 개발자 영역 > API callback 함수
		var callback_get_all_player1 = function(data) {
            obj.dataCount1 = (data.playerScores != undefined) ? data.playerScores.length : 0;
            if (obj.dataCount1 > 0) {
                for(var i=0; i<data.playerScores.length; i++) {
                    data.playerScores[i].accuracy = (data.playerScores[i].leaderboard.maxScore > 0) ? 
                        ((data.playerScores[i].score.baseScore * 100) / data.playerScores[i].leaderboard.maxScore).toFixed(2) : 
                        0;
                    data.playerScores[i].mypp = (data.playerScores[i].score.pp * data.playerScores[i].score.weight).toFixed(2);
                    console.log("data.playerScores[i].leaderboard.difficulty.difficulty", data.playerScores[i].leaderboard.difficulty.difficulty)
                    data.playerScores[i].diff = (data.playerScores[i].leaderboard.difficulty.difficulty == 1) ? "Easy" : 
                        (data.playerScores[i].leaderboard.difficulty.difficulty == 3) ? "Normal" : 
                        (data.playerScores[i].leaderboard.difficulty.difficulty == 5) ? "Hard" : 
                        (data.playerScores[i].leaderboard.difficulty.difficulty == 7) ? "Expert" : 
                        (data.playerScores[i].leaderboard.difficulty.difficulty == 9) ? "Expert+" : 
                        data.playerScores[i].leaderboard.difficulty.difficulty;
                    // #region 개발자 영역 > paging(more) 처리를 위한 셋팅
                    obj.p1Data.push(data.playerScores[i]);
                    if(i == data.playerScores.length - 1) {
                        obj.pageIndex1++;
                    }
                    // #endregion
                }
                //obj.bindGrid(data.scores);
            } else {
                return null
            }
            var pages = Number(obj.txtGetPages.val());
            var processCount = (obj.pageIndex1 > pages) ? pages : obj.pageIndex1;
            var gridCaption = "Searched <font color='red'><b>" + processCount + " / " + pages + "</b></font> Completed.";
            obj.compareContainer.jqGrid("setCaption", gridCaption);

            if (obj.statusAllSearch1 == "ready" || obj.statusAllSearch1 == "process" || obj.statusAllSearch1 == "resume") {
                if (obj.pageIndex1 <= obj.max_page1) {
                    obj.statusAllSearch1 = "process";
                    obj.compareApiCall("get_all_player1");
                } else {
                    obj.statusAllSearch1 = "complete";
                    if (obj.statusAllSearch1 == "complete" && obj.statusAllSearch2 == "complete") {
                        obj.bindGrid();
                    }
                    return null
                }
            } else if (obj.statusAllSearch1 == "stop") {
                obj.statusAllSearch1 = "resume";
                return null;
            } else {

            }

            // if (obj.pageIndex2 <= obj.max_page2) {
            //     obj.statusAllSearch2 = "process";
            //     obj.compareApiCall("get_all_player");
            // }
        };
        
        // #region 개발자 영역 > API callback 함수
		var callback_get_all_player2 = function(data) {
            obj.dataCount2 = (data.playerScores != undefined) ? data.playerScores.length : 0;
            if (obj.dataCount2 > 0) {
                for(var i=0; i<data.playerScores.length; i++) {
                    data.playerScores[i].accuracy = (data.playerScores[i].leaderboard.maxScore > 0) ? 
                        ((data.playerScores[i].score.baseScore * 100) / data.playerScores[i].leaderboard.maxScore).toFixed(2) : 
                        0;
                    data.playerScores[i].mypp = (data.playerScores[i].score.pp * data.playerScores[i].score.weight).toFixed(2);
                    data.playerScores[i].diff = (data.playerScores[i].leaderboard.difficulty.difficulty == 1) ? "Easy" : 
                        (data.playerScores[i].leaderboard.difficulty.difficulty == 3) ? "Normal" : 
                        (data.playerScores[i].leaderboard.difficulty.difficulty == 5) ? "Hard" : 
                        (data.playerScores[i].leaderboard.difficulty.difficulty == 7) ? "Expert" : 
                        (data.playerScores[i].leaderboard.difficulty.difficulty == 9) ? "Expert+" : 
                        data.playerScores[i].leaderboard.difficulty.difficulty;
                    // #region 개발자 영역 > paging(more) 처리를 위한 셋팅
                    obj.p2Data.push(data.playerScores[i]);
                    if(i == data.playerScores.length - 1) {
                        obj.pageIndex2++;
                    }
                    // #endregion
                }
                //obj.bindGrid(data.scores);
            } else {
                return null
            }
            var pages = Number(obj.txtGetPages.val());
            var processCount = (obj.pageIndex2 > pages) ? pages : obj.pageIndex2;
            var gridCaption = "Searched <font color='red'><b>" + processCount + " / " + pages + "</b></font> Completed.";
            obj.compareContainer.jqGrid("setCaption", gridCaption);

            if (obj.statusAllSearch2 == "ready" || obj.statusAllSearch2 == "process" || obj.statusAllSearch2 == "resume") {
                if (obj.pageIndex2 <= obj.max_page2) {
                    obj.statusAllSearch2 = "process";
                    obj.compareApiCall("get_all_player2");
                } else {
                    obj.statusAllSearch2 = "complete";
                    if (obj.statusAllSearch1 == "complete" && obj.statusAllSearch2 == "complete") {
                        obj.bindGrid();
                    }
                    return null
                }
            } else if (obj.statusAllSearch2 == "stop") {
                obj.statusAllSearch2 = "resume";
                return null;
            } else {

            }

            // if (obj.pageIndex2 <= obj.max_page2) {
            //     obj.statusAllSearch2 = "process";
            //     obj.compareApiCall("get_all_player");
            // }
		};

        // #endregion
        
        ajaxCall(apiObj.path, requestData, fnSuccess, fnError, apiObj.method, apiObj.content_type, apiObj.response_content_type);
    };
    // #endregion
}