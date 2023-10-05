
$(document).keydown(function (e) {
    if (e.target.nodeName != "INPUT" && e.target.nodeName != "TEXTAREA") {
        if (e.keyCode === 8) {
            return false;
        }
    }
});

$(window).resize(function () {
    resizeContent();
});

function resizeContent() {
    var header_obj = $("div[data-role='header']");
    var footer_obj = $("div[data-role='footer']");
    var browserHeight = document.documentElement.clientHeight;

    var headerHeight = parseInt(header_obj.height()) +
        parseInt(header_obj.css("padding-bottom")) +
        parseInt(header_obj.css("padding-top")) +
        parseInt(header_obj.css("boarder-top-width")) +
        parseInt(header_obj.css("padding-bottom-width"));

    var footerHeight = parseInt(footer_obj.height()) +
        parseInt(footer_obj.css("padding-bottom")) +
        parseInt(footer_obj.css("padding-top")) +
        parseInt(footer_obj.css("boarder-top-width")) +
        parseInt(footer_obj.css("padding-bottom-width"));

    var contentHeight = browserHeight - headerHeight - footerHeight;

    $("div[data-role='content']").css("height", contentHeight);
}

function startPreventParentScroll() {
    $('#home').addClass('stop-scrolling');
    $('#dvMainContent').addClass('stop-scrolling').bind('touchmove', function (e) { e.preventDefault() });
}

function endPreventParentScroll() {
    $('#home').removeClass('stop-scrolling');
    $('#dvMainContent').removeClass('stop-scrolling').unbind('touchmove');
}

function startPreventAllScroll() {
    $('body').addClass('stop-scrolling').bind('touchmove', function (e) { e.preventDefault() });
}

function endPreventAllScroll() {
    $('body').removeClass('stop-scrolling').unbind('touchmove');
}

// 뒤로가기 막기
$(document).keydown(function (e) {
    if (e.target.nodeName != "INPUT" && e.target.nodeName != "TEXTAREA") {
        if (e.keyCode === 8) {
            return false;
        }
    }
});

$.ajaxSetup({ 
    cache: true
});

$(document).ajaxStart(function () {
    //$.mobile.loading('show');
});

$(document).ajaxStop(function () {
    //$.mobile.loading('hide');
});

$(document).ajaxError(function (e, jqxhr, settings, exception) {
    console.log("ajaxError" + exception);
    if (jqxhr.status != 500 && jqxhr.status != 404) {
        panelAlertLayer(jqxhr.status + '//' + jqxhr.responseText);
        //window.location.reload(true);
    }
    else {
        execScript(jqxhr.responseText);
    }
});

/*비동기호출용ajax*/
function ajaxCall(url, data, success, error, method, contentType, responseContentType) {
    var selEnv = $("#selEnv").val();

    if(contentType == undefined || contentType == "") {
        contentType = "application/json; charset=UTF-8;";
    }
    
    if(responseContentType == undefined || responseContentType == "") {
        responseContentType = "application/json";
    }
    
    for ( var key in data["path"] ) {
        url = url.replace(key, data["path"][key]);
    }
    if(url.substring(0,1) == "/" && url.indexOf("http") < 0){
        url = selEnv + url;
    }

    if (contentType == "multipart/form-data") {
        $.ajax({
            url: url,
            type: method,
            data: data["query"],
            processData: false,
            contentType: false,
            dataType: "json",
            cache: false,
            beforeSend: function (xhr) {
                for ( var key in data["header"] ) {
                    if (root.isInternal) {
                        xhr.setRequestHeader(key, data["header"][key]);
                    } else {
                        if (key != "Application-Key" && key != "Application-Secret") {
                            if (key == "Access-Token" && root.tokenKey == "post_user_token") {
                                xhr.setRequestHeader("Authorization", data["header"][key]);
                            } else {
                                xhr.setRequestHeader(key, data["header"][key]);
                            }
                        }
                    }
                }
            },
            complete: function () {
            }
        }).done(success).fail(error);
    } else {
        if(responseContentType == "application/json") {
            $.ajax({
                url: url,
                type: method,
                data: data["query"],
                contentype: contentType,
                dataType: 'json',
                beforeSend: function (xhr) {
                    xhr.setRequestHeader("Content-Type", contentType);
                    for ( var key in data["header"] ) {
                        if (root.isInternal) {
                            xhr.setRequestHeader(key, data["header"][key]);
                        } else {
                            if (key != "Application-Key" && key != "Application-Secret") {
                                if (key == "Access-Token" && root.tokenKey == "post_user_token") {
                                    xhr.setRequestHeader("Authorization", data["header"][key]);
                                } else {
                                    xhr.setRequestHeader(key, data["header"][key]);
                                }
                            }
                        }
                    }
                },
                complete: function () {
                }
            }).done(success).fail(error);
        } else if(responseContentType == "application/octet-stream") {
            let urlParameters = Object.entries(data["query"]).map(e => e.join('=')).join('&');
            var oReq = new XMLHttpRequest();
            oReq.open(method, url + "?" + urlParameters);
            oReq.responseType = "arraybuffer";
            for ( var key in data["header"] ) {
                if (root.isInternal) {
                    oReq.setRequestHeader(key, data["header"][key]);
                } else {
                    if (key != "Application-Key" && key != "Application-Secret") {
                        if (key == "Access-Token" && root.tokenKey == "post_user_token") {
                            oReq.setRequestHeader("Authorization", data["header"][key]);
                        } else {
                            oReq.setRequestHeader(key, data["header"][key]);
                        }
                    }
                }
            }
    
            oReq.onload = function(oEvent) {
                if (this.status === 200) {
                    var blob = new Blob([oReq.response], {type: contentType});
                    var url = URL.createObjectURL( blob );
                    var link = document.createElement('a');
                    link.href = url;
                    link.download = data["path"]["{filename}"];
                    link.click();
                }
            };
            oReq.send();
        } else {
        }
    }
}

function validationAjaxSuccess(data) {
    if (data.return !== undefined) {
        if (data.message != "OK") {
            var msg = (data.error == undefined || data.error == null) ? "[" + data.return + "] " + data.message : "[" + data.return + "] " + data.message + " : " + data.error
            panelAlertLayer(msg, "error");
            return false;
        }
        else {
            if (data.return == 0) {
                return true;
            }
            else {
                return false;
            }
        }
    } else if (data.result !== undefined) {
        if (data.message != "OK") {
            var msg = (data.error == undefined || data.error == null) ? "[" + data.result + "] " + data.message : "[" + data.result + "] " + data.message + " : " + data.error
            panelAlertLayer(msg, "error");
            return false;
        }
        else {
            if (data.result == "000") {
                return true;
            }
            else {
                return false;
            }
        }
    }
    
};

function panelAlertLayer(msg, type) {

    var color = (type == undefined || type == null || type == "") ? "red" : //디폴트 실패로..
        (type == "error") ? "red" : // 실패
        (type == "info") ? "orange" : //정보(유효성체크)
        (type == "success") ? "green" : //성공
        "black"; // etc

    var ctrl = $("#panelAlertLayer");

    var html = '\
        <div id="alertHeader" data-role="header" style="background-color:' + color + ';">\
            <div style="text-align:center; color:white; vertical-align:middle;"><h4>' + msg + '</h4></div>\
        </div>\
    ';

    ctrl.empty().append(html).trigger('create');

    ctrl.css("height", $("#alertHeader").height());

    ctrl.animate({ 'top': '0' }, 200);

    setTimeout("$(\"#panelAlertLayer\").animate({ 'top': '-100%' }, 1000);", 2000);
}

function ConvertToCSV(header, objArray) {
    var array = typeof objArray != 'object' ? JSON.parse(objArray) : objArray;
    var str = header + '\r\n';

    for (var i = 0; i < array.length; i++) {
        var line = '';
        for (var index in array[i]) {
            if (line != '') line += ','

            line += array[i][index];
        }

        str += line + '\r\n';
    }

    return str;
}

function getUrlParameter(sParam) {
    var sPageURL = window.location.search.substring(1),
        sURLVariables = sPageURL.split('&'),
        sParameterName,
        i;

    for (i = 0; i < sURLVariables.length; i++) {
        sParameterName = sURLVariables[i].split('=');

        if (sParameterName[0] === sParam) {
            return sParameterName[1] === undefined ? true : decodeURIComponent(sParameterName[1]);
        }
    }
}

function toggleCheckbox(ctrlName, chkName) {
    var flag = $('input:checkbox[name="' + ctrlName + '"]').is(":checked");
    $('input:checkbox[name="' + chkName + '"]').each(function() {
        this.checked = flag; //checked 처리
   });
}

function createLeftMenuLayer(obj) {

    var html = '\
        <ul data-role="listview">\
            <li><a data-role="button">Game Tools</a></li>\
            <li><a href="../scoresaber/index.html" target="_blank">ScoreSaber</a></li>\
            <li><a href="../scoresaber2/index.html" target="_blank">ScoreSaber2</a></li>\
            <li><a href="../beatleader/index.html" target="_blank">BeatLeader</a></li>\
            <li></li>\
        </ul>\
    ';

    obj.navpanel_left.empty().append(html).trigger('create');
};