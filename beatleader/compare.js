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
        obj.userPageIndex1 = 1;
        obj.userPageIndex2 = 1;
        obj.dataCount1 = 0;
        obj.dataCount2 = 0;
        obj.total_count = 0;
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
        obj.selCountryList1 = $("#selCountryList1");
		obj.txtUserID1 = $("#txtUserID1");
		obj.btnUserMoreSearch1 = $("#btnUserMoreSearch1");

		obj.selUserList2 = $("#selUserList2");
        obj.selCountryList2 = $("#selCountryList2");
		obj.txtUserID2 = $("#txtUserID2");
		obj.btnUserMoreSearch2 = $("#btnUserMoreSearch2");

		obj.btnCompareSearch = $("#btnCompareSearch");
		obj.compareContainer = $("#compareContainer");
		obj.comparePager = $("#comparePager");
    };
    
    obj.defineElementsEvent = function() {
        // #region 개발자 영역 > Config 정의한 Controls 이벤트 등록
        // config 파일의 root > menu > {search_area | contents_area} > controls 항목들에 대한 이벤트 함수를 등록한다.
		obj.selCountryList1.change(function(){
            obj.pageUserIndex1 = 1;
            obj.selUserList1.empty();
            obj.btnUserMoreSearch1.click();
		});
        obj.selUserList1.change(function(){
            obj.txtUserID1.val(this.value);
		});
        obj.selUserList1.change();

        obj.btnUserMoreSearch1.click(function () {
            // obj.setGrid(); //Grid 조회 시 주석 해제
            obj.compareApiCall(config.controls.find(x => x.id === this.id).api);
		});

        obj.selCountryList2.change(function(){
            obj.pageUserIndex2 = 1;
            obj.selUserList2.empty();
            obj.btnUserMoreSearch2.click();
		});
        obj.selUserList2.change(function(){
            obj.txtUserID2.val(this.value);
		});
        obj.selUserList2.change();

        obj.btnUserMoreSearch2.click(function () {
            // obj.setGrid(); //Grid 조회 시 주석 해제
            obj.compareApiCall(config.controls.find(x => x.id === this.id).api);
		});

        obj.btnCompareSearch.click(function () {
            obj.setGrid(); //Grid 조회 시 주석 해제
            obj.pageIndex1 = 1;
            obj.pageIndex2 = 1;
            obj.dataCount1 = 0;
            obj.dataCount2 = 0;
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
        $.fn.fmatter.convertDownload = function (cellValue,rowObject,options) {
            if(cellValue != "") {
                return "<a href='" + cellValue + "'><img src='https://w7.pngwing.com/pngs/596/75/png-transparent-download-now-download-icon-download-button-download-logo-flat-icon-flat-logo-flat-image-button-flat-round-thumbnail.png' width='25px' height='25px' /></a>";
            } else {
                return "";
            }
        };
        $.fn.fmatter.convertSongTitle = function (cellValue,rowObject,options) {
			return "<a href='https://www.beatleader.xyz/leaderboard/global/" + options.leaderboard.id + "?page=1&countries=" + obj.country + "' target='_blank'>" + cellValue + "</a>";
        };
        $.fn.fmatter.convertDifficulty = function (cellValue,rowObject,options) {
			return ConvertToString("difficulty", options);
        };
        $.fn.fmatter.convertAccuracy = function (cellValue,rowObject,options) {
            return (cellValue*100).toFixed(2) + "%";
			// return (cellValue <= 0) ? "-" : (cellValue*100).toFixed(2) + "%";
        };
        $.fn.fmatter.convertStar = function (cellValue,rowObject,options) {
			if(cellValue != null) {
                return cellValue.toFixed(2) + "★";
            } else {
                return "-";
            }
        };
        $.fn.fmatter.convertTime = function (cellValue,rowObject,options) {
			return ConvertToString("timestamp", cellValue);
        };
        // #endregion
        
        obj.navGrid();
    }

    // Grid Data Binding
    obj.bindGrid = function() {
        console.log("bindGrid start");

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

                if(obj.p1Data[i].leaderboardId == obj.p2Data[j].leaderboardId) {
                    isDuplicate = true;
                    
                    resultData["id"] = obj.p2Data[j].leaderboardId;
                    resultData["stars"] = obj.p2Data[j].stars;
                    resultData["leaderboardId"] = obj.p2Data[j].leaderboardId;
                    resultData["songName"] = obj.p2Data[j].songName;
                    resultData["mapper"] = obj.p2Data[j].mapper;
                    resultData["diff"] = obj.p2Data[j].diff;
                    resultData["modifiers_1"] = obj.p1Data[i].modifiers;
                    resultData["modifiers_2"] = obj.p2Data[j].modifiers;
                    resultData["accRating"] = obj.p2Data[j].accRating;
                    resultData["techRating"] = obj.p2Data[j].techRating;
                    resultData["passRating"] = obj.p2Data[j].passRating;
                    resultData["timeset_1"] = obj.p1Data[i].timeset;
                    resultData["timeset_2"] = obj.p2Data[j].timeset;
                    
                    if(obj.p1Data[i].acc < obj.p2Data[j].acc) {
                        resultData["winner"] = obj.p2Nickname;
                        p1LoseCount++;
                    } else {
                        resultData["winner"] = obj.p1Nickname;
                        p1WinCount++;
                    }

                    resultData["accuracy_1"] = obj.p1Data[i].acc;
                    resultData["accuracy_2"] = obj.p2Data[j].acc;
                    resultData["accuracy_gap"] = obj.p1Data[i].acc - obj.p2Data[j].acc;
                    p1AccuracyGap += parseFloat(obj.p1Data[i].acc - obj.p2Data[j].acc);
                    
                    obj.compareContainer.jqGrid("addRowData", resultData["id"], resultData);
                }
            }
            if (!isDuplicate) {
                
                resultData["id"] = obj.p1Data[i].leaderboardId;
                resultData["stars"] = obj.p1Data[i].stars;
                resultData["leaderboardId"] = obj.p1Data[i].leaderboardId;
                resultData["songName"] = obj.p1Data[i].songName;
                resultData["mapper"] = obj.p1Data[i].mapper;
                resultData["diff"] = obj.p1Data[i].diff;
                resultData["modifiers_1"] = obj.p1Data[i].modifiers;
                resultData["modifiers_2"] = "";
                resultData["winner"] = "-";
                p1OnlyPlayCount++;
                resultData["accuracy_1"] = obj.p1Data[i].acc;
                resultData["accuracy_2"] = 0;
                resultData["accuracy_gap"] = 0;
                resultData["accRating"] = obj.p1Data[i].accRating;
                resultData["techRating"] = obj.p1Data[i].techRating;
                resultData["passRating"] = obj.p1Data[i].passRating;
                resultData["timeset_1"] = obj.p1Data[i].timeset;
                resultData["timeset_2"] = "";
                
                obj.compareContainer.jqGrid("addRowData", resultData["leaderboardId"], resultData);
            }
        }

        for(var i=0; i<obj.p2Data.length; i++) {
            var isDuplicate = false;
            for(var j=0; j<obj.p1Data.length; j++) {
                if(obj.p2Data[i].leaderboardId == obj.p1Data[j].leaderboardId) {
                    isDuplicate = true;
                }
            }
            if (!isDuplicate) {

                resultData["id"] = obj.p2Data[i].leaderboardId;
                resultData["stars"] = obj.p2Data[i].stars;
                resultData["leaderboardId"] = obj.p2Data[i].leaderboardId;
                resultData["songName"] = obj.p2Data[i].songName;
                resultData["mapper"] = obj.p2Data[i].mapper;
                resultData["diff"] = obj.p2Data[i].diff;
                resultData["modifiers_1"] = "";
                resultData["modifiers_2"] = obj.p2Data[i].modifiers;
                resultData["winner"] = "-";
                p2OnlyPlayCount++;
                resultData["accuracy_1"] = 0;
                resultData["accuracy_2"] = obj.p2Data[i].acc;
                resultData["accuracy_gap"] = 0;
                resultData["accRating"] = obj.p2Data[i].accRating;
                resultData["techRating"] = obj.p2Data[i].techRating;
                resultData["passRating"] = obj.p2Data[i].passRating;
                resultData["timeset_1"] = "";
                resultData["timeset_2"] = obj.p2Data[i].timeset;

                obj.compareContainer.jqGrid("addRowData", resultData["leaderboardId"], resultData);
            }
        }
        
		var gridCount = obj.compareContainer.getGridParam("reccount");
		var gridCaption = "Searched <font color='blue'><b>" + gridCount + "</b></font> Completed.";
        obj.compareContainer.jqGrid("setCaption", gridCaption);

        var html = '<hr>';
        html += '<h4><font color="green"><b>' + obj.p1Nickname + '</b></font> Summary Info(Versus : <font color="blue">' + obj.p2Nickname + '</font>)</h4>';
        html += "<h4><b><font color='green'>" + obj.p1Nickname + "</font> is <font color='blue'>" + p1WinCount + "</font> win <font color='red'>" + p1LoseCount + "</font> lose.</b></h4>";
        html += '<h4><b>- Accuracy Gap : ' + p1AccuracyGap.toFixed(2) + '</b></h4>';
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
        //"<a href='https://www.beatleader.xyz/u/76561198830502286" + page + "' target='_blank'>" + cellValue + "</a>"
        //https://www.beatleader.xyz/global/2&country=kr
        var globalRankPage = parseInt((data.rank - 1) / 50) + 1;
        var countryRankPage = parseInt((data.countryRank - 1) / 50) + 1;
        var globalRankGap = (data.lastWeekRank-data.rank) > 0 ? '↑'+(data.lastWeekRank-data.rank) : (data.lastWeekRank-data.rank) < 0 ? '↓'+(data.rank-data.lastWeekRank) : '-';
        var countryRankGap = (data.lastWeekCountryRank-data.countryRank) > 0 ? '↑'+(data.lastWeekCountryRank-data.countryRank) : (data.lastWeekCountryRank-data.countryRank) < 0 ? '↓'+(data.countryRank-data.lastWeekCountryRank) : '-';
        var globalPPGap = (data.pp-data.lastWeekPp) > 0 ? '↑'+(data.pp-data.lastWeekPp).toFixed(2)+"pp" : (data.pp-data.lastWeekPp) < 0 ? '↓'+(data.lastWeekPp-data.pp).toFixed(2)+"pp" : '-';
        var clansHtml = '';
        for(var i=0;i<data.clans.length;i++) {
            clansHtml += '<a href="https://www.beatleader.xyz/clan/' + data.clans[i].tag + '" target="_blank"><font color="' + data.clans[i].color + '">' + data.clans[i].tag + '</font></a> ';
        }
        html = '<hr>';
        html += '<h4><b>P' + index + ' : <a href="https://www.beatleader.xyz/u/' + data.id + '" target="_blank">'+ "<img src='" + data.avatar + "' width=50 height=50 boarder=0 />" + '<font color="green"> ' + data.name + '</font></a> <a href="https://www.beatleader.xyz/rankings?page=1&countries=' + data.country + '" target="_blank"><font color="blue">(' + data.country + ')</font></a> ' + clansHtml + '</b></h4>';
        html += '<h4><b><font color="darkgray">Global / Country Rank</font> : <a href="https://www.beatleader.xyz/rankings?page=' + globalRankPage + '" target="_blank"><font color="red">' + data.rank + '(' + globalRankGap + ')</font></a>' + ' / <a href="https://www.beatleader.xyz/rankings?page=' + countryRankPage + '&countries=' + data.country + '" target="_blank"><font color="orange">' + data.countryRank + '(' + countryRankGap + ')</font></a></b></h4>';
        html += '<h4><b><font color="darkgray">PP</font> : <font color="red">' + data.pp + 'pp</font>(<font color="green">' + globalPPGap + '</font>) (Acc : <font color="#ff7300">' + data.accPp.toFixed(2) + '</font> , Tech : <font color="orange">' + data.techPp.toFixed(2) + '</font> , Pass : <font color="blue">' + data.passPp.toFixed(2) + '</font>)</b></h4>';
        html += '<h4><b><font color="darkgray">Avg Accuracy</font> : <font color="red">' + (data.scoreStats.averageRankedAccuracy*100).toFixed(2) + '%</font> (SS+:<font color="#ff7300">' + data.scoreStats.sspPlays + '</font>, SS:<font color="orange">' + data.scoreStats.ssPlays + '</font>, S+:<font color="blue">' + data.scoreStats.spPlays + '</font>, S:<font color="green">' + data.scoreStats.sPlays + '</font>, A:<font color="gray">' + data.scoreStats.aPlays + '</font>)</b></h4>';
        html += '<h4><b><font color="darkgray">Play Count(Rank Count)</font> : <font color="red">' + data.scoreStats.totalPlayCount + '</font>(<font color="green">' + data.scoreStats.rankedPlayCount + '</font>)</b></h4>';
        $("#dvtitleArea").append(html).trigger("create");

        // //"<a href='https://scoresaber.com/u/76561198830502286" + page + "' target='_blank'>" + cellValue + "</a>"
        // //https://scoresaber.com/global/2&country=kr
        // var globalRankPage = parseInt((data.rank - 1) / 50) + 1;
        // var countryRankPage = parseInt((data.countryRank - 1) / 50) + 1;
        // html = '<hr>';
        // //html += '<h4><b><font color="red">P' + index + '</font> : <a href="https://scoresaber.com/u/' + data.playerId + '" target="_blank"><font color="green">' + data.name + '</font></a> <a href="https://scoresaber.com/rankings?page=1&countries=' + data.country + '" target="_blank"><font color="blue">(' + data.country + ')</font></a></b></h4>';
        // html += '<h4><b><a href="https://scoresaber.com/u/' + data.id + '" target="_blank">'+ "<img src='" + data.profilePicture + "' width=50 height=50 boarder=0 />" + '<font color="green"> ' + data.name + '</font></a> <a href="https://scoresaber.com/rankings?page=1&countries=' + data.country + '" target="_blank"><font color="blue">(' + data.country + ')</font></a></b></h4>';
        // html += '<h4><b>Global Rank : <a href="https://scoresaber.com/rankings?page=' + globalRankPage + '" target="_blank"><font color="red">' + data.rank + '</font></a>' + ' (<a href="https://scoresaber.com/rankings?page=' + countryRankPage + '&countries=' + data.country + '" target="_blank"><font color="orange">' + data.countryRank + '</font></a>)</b></h4>';
        // html += '<h4><b>PP : ' + data.pp + ', Avg Accuracy : ' + data.scoreStats.averageRankedAccuracy.toFixed(2) + '</b></h4>';
        // html += '<h4><b>Play Count(Rank Count) : ' + data.scoreStats.totalPlayCount + '(' + data.scoreStats.rankedPlayCount + ')</b></h4>';
        // $("#dvtitleArea").append(html).trigger("create");
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

        
        var callback_get_ranking1 = function(data) {
            obj.userDataCount1 = (data != undefined) ? data.data.length : 0;
            obj.user1_max_page = parseInt((data.metadata.total - 1) / data.metadata.itemsPerPage) + 1;
            obj.user_total_count1 = data.metadata.total;//(data.metadata.total * data.metadata.itemsPerPage);
            if (obj.userDataCount1 > 0) {
                var iterHtml = "";
                iterHtml += "<option value='' selected='selected'>select</option>";
                for (var i = 0; i < data.data.length; i++) {
                    iterHtml += "<option value='" + data.data[i].id + "'>(" + data.data[i].rank + ") " + data.data[i].name + "</option>";
                }
                obj.selUserList1.append(iterHtml).trigger('create');
                obj.pageUserIndex1++;
            } else {
                return null
            }
		};
        
        var callback_get_ranking2 = function(data) {
            obj.userDataCount2 = (data != undefined) ? data.data.length : 0;
            obj.user2_max_page = parseInt((data.metadata.total - 1) / data.metadata.itemsPerPage) + 1;
            obj.user_total_count2 = data.metadata.total;//(data.metadata.total * data.metadata.itemsPerPage);
            if (obj.userDataCount2 > 0) {
                var iterHtml = "";
                iterHtml += "<option value='' selected='selected'>select</option>";
                for (var i = 0; i < data.data.length; i++) {
                    iterHtml += "<option value='" + data.data[i].id + "'>(" + data.data[i].rank + ") " + data.data[i].name + "</option>";
                }
                obj.selUserList2.append(iterHtml).trigger('create');
                obj.pageUserIndex2++;
            } else {
                return null
            }
		};

		var callback_get_full1 = function(data) {
            obj.p1Nickname = data.name;
            obj.country = data.country;
            // obj.max_page1 = parseInt((data.scoreStats.totalPlayCount - 1) / 8) + 1;
            // obj.max_page1 = (obj.max_page1 > pages) ? pages : obj.max_page1;
            obj.total_count = data.scoreStats.totalPlayCount
            obj.displayData(1, data);
            
            obj.pageIndex1 = 1;
            obj.compareApiCall("get_player_accgraph1");
            // if (obj.statusAllSearch1 == "ready") {
            //     // obj.setGrid(); //Grid 조회 시 주석 해제
            //     obj.pageIndex1 = 1;
			//     obj.compareApiCall("get_all_player1");
            // } else if (obj.statusAllSearch1 == "process") {
            //     obj.statusAllSearch1 = "stop";
            //     //obj.max_page1 = 0;
            // } else if (obj.statusAllSearch1 == "resume") {
			//     obj.compareApiCall("get_all_player1");
            // } else if (obj.statusAllSearch1 == "complete") {
            // } else {
                
            // }
        };
        
		var callback_get_full2 = function(data) {
            obj.p2Nickname = data.name;
            // obj.max_page2 = parseInt((data.scoreStats.totalPlayCount - 1) / 8) + 1;
            // obj.max_page2 = (obj.max_page2 > pages) ? pages : obj.max_page2;
            obj.total_count = data.scoreStats.totalPlayCount
            obj.displayData(2, data);
            
            obj.pageIndex2 = 1;
            obj.compareApiCall("get_player_accgraph2");
            // if (obj.statusAllSearch2 == "ready") {
            //     //obj.setGrid(); //Grid 조회 시 주석 해제
            //     obj.pageIndex2 = 1;
			//     obj.compareApiCall("get_all_player2");
            // } else if (obj.statusAllSearch2 == "process") {
            //     obj.statusAllSearch2 = "stop";
            //     //obj.max_page2 = 0;
            // } else if (obj.statusAllSearch2 == "resume") {
			//     obj.compareApiCall("get_all_player2");
            // } else if (obj.statusAllSearch2 == "complete") {
            // } else {
                
            // }
		};

        // #region 개발자 영역 > API callback 함수
		var callback_get_player_accgraph1 = function(data) {
            obj.dataCount1 = (data != undefined) ? data.length : 0;
            if (obj.dataCount1 > 0) {
                for(var i=0; i<data.length; i++) {
                    obj.p1Data.push(data[i]);
                    // #endregion
                }
                //obj.bindGrid(data.scores);
            } else {
                return null
            }
        };
        
        // #region 개발자 영역 > API callback 함수
		var callback_get_player_accgraph2 = function(data) {
            obj.dataCount2 = (data != undefined) ? data.length : 0;
            if (obj.dataCount2 > 0) {
                for(var i=0; i<data.length; i++) {
                    obj.p2Data.push(data[i]);
                    // #endregion
                }
                //obj.bindGrid(data.scores);
            } else {
                return null
            }
            var gridCaption = "<font color='red'>" + obj.dataCount2 + "</font> Completed.";
            obj.compareContainer.jqGrid("setCaption", gridCaption);
            setTimeout(function() {
                obj.bindGrid();
            }, 3000);
		};

        // #endregion
        
        ajaxCall(apiObj.path, requestData, fnSuccess, fnError, apiObj.method, apiObj.content_type, apiObj.response_content_type);
    };
    // #endregion
}