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
			return "<a href='https://www.beatleader.xyz/leaderboard/global/" + options.leaderboardId + "' target='_blank'>" + cellValue + "</a>";
        };
        $.fn.fmatter.convertDifficulty = function (cellValue,rowObject,options) {
			if (cellValue == "ExpertPlus") {
                return "Expert+"; 
            }
            else {
                return cellValue;
            }
            // return ConvertToString("difficulty", options);
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
			if (cellValue != "") {
                return ConvertToString("timestamp", cellValue);
            }
            return cellValue;
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

                if(obj.p1Data[i].leaderboardId == obj.p2Data[j].leaderboardId) {
                    isDuplicate = true;
                    
                    resultData["id"] = obj.p1Data[i].leaderboardId;
                    resultData["stars_1"] = obj.p1Data[i].stars;
                    resultData["stars_2"] = obj.p2Data[j].stars;
                    resultData["leaderboardId"] = obj.p1Data[i].leaderboardId;
                    resultData["songName"] = obj.p1Data[i].songName;
                    resultData["mapper"] = obj.p1Data[i].mapper;
                    resultData["diff"] = obj.p1Data[i].diff;
                    resultData["modifiers_1"] = obj.p1Data[i].modifiers;
                    resultData["modifiers_2"] = obj.p2Data[j].modifiers;
                    resultData["accRating"] = obj.p1Data[i].accRating;
                    resultData["techRating"] = obj.p1Data[i].techRating;
                    resultData["passRating"] = obj.p1Data[i].passRating;
                    resultData["timeset_1"] = obj.p1Data[i].timeset;
                    resultData["timeset_2"] = obj.p2Data[j].timeset;
                    
                    if (obj.p1Data[i].stars == obj.p2Data[j].stars) {
                        if(obj.p1Data[i].acc < obj.p2Data[j].acc) {
                            resultData["winner"] = obj.p2Nickname;
                            p1LoseCount++;
                        } else {
                            resultData["winner"] = obj.p1Nickname;
                            p1WinCount++;
                        }
                    } 
                    // else {
                    //     if(obj.p1Data[i].stars < obj.p2Data[j].stars && obj.p1Data[i].acc < obj.p2Data[j].acc) {
                    //         resultData["winner"] = obj.p2Nickname;
                    //         p1LoseCount++;
                    //     } else if(obj.p1Data[i].stars > obj.p2Data[j].stars && obj.p1Data[i].acc > obj.p2Data[j].acc) {
                    //         resultData["winner"] = obj.p1Nickname;
                    //         p1WinCount++;
                    //     }
                    // }

                    resultData["accuracy_1"] = obj.p1Data[i].acc;
                    resultData["accuracy_2"] = obj.p2Data[j].acc;
                    resultData["accuracy_gap"] = obj.p1Data[i].acc - obj.p2Data[j].acc;
                    p1AccuracyGap += parseFloat(obj.p1Data[i].acc - obj.p2Data[j].acc);

                    resultData["hash"] = obj.p1Data[i].hash;
                    resultData["levelAuthorName"] = obj.p1Data[i].mapper;
                    
                    resultData["difficulties"] = [];
                    var difficult = {};
                    difficult.name = obj.p1Data[i].diff;
                    difficult.characteristic = obj.p1Data[i].mode;
                    resultData["difficulties"].push(difficult);
                    
                    obj.compareContainer.jqGrid("addRowData", resultData["id"], resultData);
                }
            }
            if (!isDuplicate) {
                
                resultData["id"] = obj.p1Data[i].leaderboardId;
                resultData["stars_1"] = obj.p1Data[i].stars;
                resultData["stars_2"] = "";
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

                resultData["hash"] = obj.p1Data[i].hash;
                resultData["levelAuthorName"] = obj.p1Data[i].mapper;
                
                resultData["difficulties"] = [];
                var difficult = {};
                difficult.name = obj.p1Data[i].diff;
                difficult.characteristic = obj.p1Data[i].mode;
                resultData["difficulties"].push(difficult);
                
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
                resultData["stars_1"] = "";
                resultData["stars_2"] = obj.p2Data[i].stars;
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

                resultData["hash"] = obj.p2Data[i].hash;
                resultData["levelAuthorName"] = obj.p2Data[i].mapper;
                
                resultData["difficulties"] = [];
                var difficult = {};
                difficult.name = obj.p2Data[i].diff;
                difficult.characteristic = obj.p2Data[i].mode;
                resultData["difficulties"].push(difficult);

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
		
		var callback_navButton_playlist = function() {
            obj.selectedMapIDs = obj.compareContainer.getGridParam("selarrrow");
            if (obj.selectedMapIDs.length == 0) {
                panelAlertLayer("Please select the data to playlist with the checkbox.", "info");
                return;
            }

            var gridData = obj.compareContainer.getGridParam('data');
            // var playlistMap = [];
            var playListData = {};
            playListData.playlistTitle = "beatleader_compares";
            playListData.playlistAuthor = "hudson";
            playListData.image = "data:image/png;base64,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";
            playListData.songs = [];
            for( var i=0;i<obj.selectedMapIDs.length;i++) {            
                // playlistMap.push(gridData.find(x => x.id === obj.selectedMapIDs[i]));
                var compareData = gridData.find(x => x.id === obj.selectedMapIDs[i]);
                console.log(compareData);
                var song = {};
                song.hash = compareData.hash;
                song.songName = compareData.songName;
                song.levelAuthorName = compareData.levelAuthorName;
                song.difficulties = compareData.difficulties;
                playListData.songs.push(song);
            }

            var str = JSON.stringify(playListData);
            var blob = new Blob( [str], {
                type: "application/octet-stream"
            });

            var url = URL.createObjectURL( blob );
            var link = document.createElement("a");
            link.href = url;
            link.download = new Date().yyyyMMddHHmmss() + ".bplist";
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
        };
        
		var callback_get_full2 = function(data) {
            obj.p2Nickname = data.name;
            // obj.max_page2 = parseInt((data.scoreStats.totalPlayCount - 1) / 8) + 1;
            // obj.max_page2 = (obj.max_page2 > pages) ? pages : obj.max_page2;
            obj.total_count = data.scoreStats.totalPlayCount
            obj.displayData(2, data);
            
            obj.pageIndex2 = 1;
            obj.compareApiCall("get_player_accgraph2");
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
            var gridCaption = "<font color='blue'>Processing Player1 Records : </font><font color='red'>" + obj.dataCount1 + "</font>";
            obj.compareContainer.jqGrid("setCaption", gridCaption);
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
            var gridCaption = "<font color='blue'>Processing Player2 Records : </font><font color='red'>" + obj.dataCount2 + "</font>";
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