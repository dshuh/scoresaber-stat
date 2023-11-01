function Map(controlData) {

    var obj = this;
    obj.pageMapLayer = $(controlData.pageMapLayer);

    obj.menuConfig = config.menu.find(x => x.id === "map");
    obj.gridConfig = obj.menuConfig.contents_area.grid;
    obj.apiConfig = config.apis.map;
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
		
        obj.selUserList = $("#selUserList");
        obj.selCountryList = $("#selCountryList");
		obj.txtUserID = $("#txtUserID");
		obj.btnUserMoreSearch = $("#btnUserMoreSearch");

		obj.selStartStar = $("#selStartStar");
		obj.selEndStar = $("#selEndStar");
        obj.selMapSortValue = $("#selMapSortValue");
        obj.txtCountry = $("#txtCountry");
		obj.btnMapSearch = $("#btnMapSearch");
		obj.mapContainer = $("#mapContainer");
		obj.mapPager = $("#mapPager");
    };
    
    obj.defineElementsEvent = function() {
        // #region 개발자 영역 > Config 정의한 Controls 이벤트 등록
        // config 파일의 root > menu > {search_area | contents_area} > controls 항목들에 대한 이벤트 함수를 등록한다.
		obj.selCountryList.change(function(){
            obj.pageUserIndex = 1;
            obj.selUserList.empty();
            obj.btnUserMoreSearch.click();
		});
        obj.selUserList.change(function(){
            obj.txtUserID.val(this.value);
            obj.mapApiCall("get_player_accgraph");
		});
        obj.selUserList.change();
        obj.btnUserMoreSearch.click(function () {
            // obj.setGrid(); //Grid 조회 시 주석 해제
            obj.mapApiCall(config.controls.find(x => x.id === this.id).api);
		});

        var optionHtml = "";
        for (var i=0;i<=20;i++) {
            optionHtml += '<option value="' + i + '">' + i + ' stars</option>';
        }
        obj.selStartStar.append(optionHtml).trigger("create");
        obj.selEndStar.append(optionHtml).trigger("create");

        obj.selStartStar.val("0").prop("selected", true).change();
        obj.selEndStar.val("20").prop("selected", true).change();

		obj.btnMapSearch.click(function () {
            obj.pageIndex = 1;
			obj.setGrid(); //Grid 조회 시 주석 해제
            obj.mapApiCall(config.controls.find(x => x.id === this.id).api);
            // var dataList = mapList.filter(function (el) {
            //     return el.stars >= obj.selStartStar.val() &&
            //            el.stars < obj.selEndStar.val()
            //   });
            // obj.bindGrid(dataList); //Grid 조회 시 주석 해제
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
        // #region 개발자 영역 > Grid Formatter callback 함수 구현부
        $.fn.fmatter.convertDownload = function (cellValue,rowObject,options) {
			if(cellValue != "") {
                return "<a href='" + cellValue + "'><img src='https://w7.pngwing.com/pngs/596/75/png-transparent-download-now-download-icon-download-button-download-logo-flat-icon-flat-logo-flat-image-button-flat-round-thumbnail.png' width='20px' height='20px' /></a>";
            } else {
                return "";
            }
        };
        $.fn.fmatter.convertRank = function (cellValue,rowObject,options) {
            //var page = parseInt((cellValue / 12) + 1);
            var page = parseInt((cellValue - 1) / 10) + 1;
			return "<a href='https://www.beatleader.xyz/leaderboard/global/" + options.leaderboard.id + "/" + page + "' target='_blank'>" + cellValue + "</a>";
        };
        $.fn.fmatter.convertCoverImage = function (cellValue,rowObject,options) {
			return "<img src='" + options.coverImage + "' width=30 height=30 boarder=0 />";
        };
        $.fn.fmatter.convertSongTitle = function (cellValue,rowObject,options) {
			return "<a href='https://www.beatleader.xyz/leaderboard/global/" + options.id + "?page=1' target='_blank'>" + cellValue + "</a>";
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
			if (cellValue == "") return cellValue;
            return (cellValue*100).toFixed(2) + "%";
            // return ConvertToString("accuracy", options) + "%";
        };
        $.fn.fmatter.convertRating = function (cellValue,rowObject,options) {
			return ConvertToString("rating", options) + "%";
        };
        $.fn.fmatter.convertSongTime = function (cellValue,rowObject,options) {
            var min = cellValue / 60;
            var sec = cellValue % 60;
            return Math.floor(min) + "m " + sec + "s";
        };
        $.fn.fmatter.convertTime = function (cellValue,rowObject,options) {
			if (cellValue == "") return cellValue;
            return ConvertToString("timestamp", cellValue);
        };
        
        $.fn.fmatter.convertStars = function (cellValue,rowObject,options) {
			return ConvertToString("stars", options);
        };
        $.fn.fmatter.convertSSStars = function (cellValue,rowObject,options) {
			return cellValue + "★";
        };
        $.fn.fmatter.convertGap = function (cellValue,rowObject,options) {
            return cellValue > 0 ? '▲'+ Math.abs(cellValue) : cellValue < 0 ? '▼'+ Math.abs(cellValue) : '-';
        };
        $.fn.fmatter.convertAccRating = function (cellValue,rowObject,options) {
			return options.accRating.toFixed(2) + "★";
        };
        $.fn.fmatter.convertTechRating = function (cellValue,rowObject,options) {
			return options.techRating.toFixed(2) + "★";
        };
        $.fn.fmatter.convertPassRating = function (cellValue,rowObject,options) {
			return options.passRating.toFixed(2) + "★";
        };
        
        
        // #endregion
        
        obj.navGrid();
    }

    // Grid Data Binding
    obj.bindGrid = function(data) {
		for(var i=0; i<data.length; i++) {
            var scoreSaberData = mapList.find(x => x.beatSaverKey === data[i].song.id && x.maxScore === data[i].difficulty.maxScore);
            if (scoreSaberData != undefined) {
                data[i].ss_stars = scoreSaberData.stars.toFixed(2);
                data[i].stars_gap = (data[i].difficulty.stars - scoreSaberData.stars).toFixed(2);
            } else {
                data[i].ss_stars = "0";
                data[i].stars_gap = "0";
            }
            var userData = obj.pData.find(x => x.leaderboardId === data[i].id);
            if (userData == undefined) {
                data[i].modifiers = "";
                data[i].modded_stars = "";
                data[i].accuracy = "";
                data[i].timeset = "";
            } else {
                if (userData.modifiers != "") {
                    data[i].modifiers = userData.modifiers + "(" + userData.stars.toFixed(2) + "★)";
                } else {
                    data[i].modifiers = "";
                }
                data[i].accuracy = userData.acc;
                data[i].timeset = userData.timeset;
            }
            data[i].songName = data[i].song.name;
            data[i].songAuthorName = data[i].song.author;
            data[i].stars = (data[i].difficulty.stars != null) ? data[i].difficulty.stars.toFixed(2) : "0";
            data[i].modifiers = data[i].modifiers;
            
            data[i].diff = (data[i].difficulty.value == 1) ? "Easy" : 
                (data[i].difficulty.value == 3) ? "Normal" : 
                (data[i].difficulty.value == 5) ? "Hard" : 
                (data[i].difficulty.value == 7) ? "Expert" : 
                (data[i].difficulty.value == 9) ? "Expert+" : 
                data[i].difficulty.difficultyName;
            data[i].accRating  = data[i].difficulty.accRating;
            data[i].techRating = data[i].difficulty.techRating;
            data[i].passRating = data[i].difficulty.passRating;
            data[i].diffValue = data[i].difficulty.value;

            data[i].coverImage = data[i].song.coverImage;
            data[i].levelAuthorName = data[i].song.mapper;
            data[i].download = data[i].song.downloadUrl;
            data[i].durationSeconds = data[i].song.duration;
            data[i].upvotes = data[i].positiveVotes;
            data[i].downvotes = data[i].negativeVotes;
            data[i].votes = data[i].positiveVotes - data[i].negativeVotes;
            
            data[i].noteCount = data[i].difficulty.notes;
            data[i].njs = data[i].difficulty.njs;
            data[i].nps = data[i].difficulty.nps;
            data[i].rankedTime = data[i].difficulty.rankedTime.toString();
            data[i].diffValue = data[i].difficulty.value;

            data[i].beatSaverKey = data[i].song.id;
            data[i].id = data[i].song.id;

            data[i].hash = data[i].song.hash;
            // data[i].difficulties = data[i].song.difficulties;
            data[i].difficulties = [];
            var difficult = {};
            difficult.name = data[i].difficulty.difficultyName;
            difficult.characteristic = data[i].difficulty.modeName;
            data[i].difficulties.push(difficult);
            // for(var j=0;j<data[i].difficulty.length;j++){
            //     var difficult = {};
            //     difficult.name = data[i].difficulty[j].difficultyName;
            //     difficult.characteristic = data[i].difficulty[j].modeName;
            //     data[i].difficulties.push(difficult);
            // }

			// #region 개발자 영역 > paging(more) 처리를 위한 셋팅
			if(i == data.length - 1) {
				obj.pageIndex++;
			}

			obj.mapContainer.jqGrid("addRowData", data[i].id, data[i]);
		}
		var gridCount = obj.mapContainer.getGridParam("reccount");
		var percentage = (obj.dataCount != 0) ? parseFloat((gridCount*100) / obj.total_count).toFixed(2) : 0;
		var gridCaption = "Searched <font color='blue'><b>" + gridCount + "</b></font> out of <font color='red'><b>" + obj.total_count + "</b></font> - <font color='green'><b>" + percentage + "%</b></font> Completed.";
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
		var callback_navButton_playlist = function() {
            obj.selectedMapIDs = obj.mapContainer.getGridParam("selarrrow");
            if (obj.selectedMapIDs.length == 0) {
                panelAlertLayer("Please select the data to playlist with the checkbox.", "info");
                return;
            }

            var gridData = obj.mapContainer.getGridParam('data');
            // var playlistMap = [];
            var playListData = {};
            playListData.playlistTitle = "beatleader_maps";
            playListData.playlistAuthor = "hudson";
            playListData.image = getPlaylistImage();
            playListData.songs = [];
            for( var i=0;i<obj.selectedMapIDs.length;i++) {            
                // playlistMap.push(gridData.find(x => x.id === obj.selectedMapIDs[i]));
                var mapData = gridData.find(x => x.id === obj.selectedMapIDs[i]);
                var song = {};
                song.hash = mapData.hash;
                song.songName = mapData.songName;
                song.levelAuthorName = mapData.levelAuthorName;
                song.difficulties = mapData.difficulties;
                // for (var j=0;j<mapData.difficulties.length;j++){
                //     song.difficulties.push(mapData.difficulties[j]);
                // }
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
            obj.mapApiCall("get_leaderboards");
		};

        // #endregion
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
		var callback_get_leaderboards = function(data) {
            obj.dataCount = (data != undefined) ? data.data.length : 0;
            obj.total_count = data.metadata.total;
            if (obj.dataCount > 0) {
                obj.bindGrid(data.data);
            } else {
                return null
            }
		};

        
        var callback_get_ranking = function(data) {
            obj.userDataCount = (data != undefined) ? data.data.length : 0;
            obj.user_max_page = parseInt((data.metadata.total - 1) / data.metadata.itemsPerPage) + 1;
            obj.user_total_count = data.metadata.total;//(data.metadata.total * data.metadata.itemsPerPage);
            if (obj.userDataCount > 0) {
                var iterHtml = "";
                iterHtml += "<option value='' selected='selected'>select</option>";
                for (var i = 0; i < data.data.length; i++) {
                    iterHtml += "<option value='" + data.data[i].name + "'>(" + data.data[i].rank + ") " + data.data[i].name + "</option>";
                }
                obj.selUserList.append(iterHtml).trigger('create');
                obj.pageUserIndex++;
            } else {
                return null
            }
		};

        // var callback_get_full = function(data) {
        //     obj.mapApiCall("get_player_accgraph");
        // };
        
        // #region 개발자 영역 > API callback 함수
		var callback_get_player_accgraph = function(data) {
            obj.dataCount = (data != undefined) ? data.length : 0;
            if (obj.dataCount > 0) {
                obj.pData = [];
                for(var i=0; i<data.length; i++) {
                    obj.pData.push(data[i]);
                    // #endregion
                }
                //obj.bindGrid(data.scores);
            } else {
                return null
            }
            // var gridCaption = "<font color='blue'>Processing Player1 Records : </font><font color='red'>" + obj.dataCount1 + "</font>";
            // obj.compareContainer.jqGrid("setCaption", gridCaption);
        };

        // #endregion
        
        ajaxCall(apiObj.path, requestData, fnSuccess, fnError, apiObj.method, apiObj.content_type, apiObj.response_content_type);
    };
    // #endregion
}