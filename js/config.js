function displayEnvLayer() {
    root.isLocalhost = true;
    root.isInternal = true;

    var html = '<select id="selEnv" style="width:100%;">';

    for(var i=0;i<config.environments.length;i++){
        html += '<option value="' + config.environments[i].prefix_internal_url + '">' + config.environments[i].name + '</option>';
    }
    html += '</select>'

    return html;
}

function displayContentsLayer(menuConfig) {
    return '\
        <ul data-role="listview">\
            <li data-role="controlgroup">\
                ' + menuConfig.search_area.title + '\
            </li>\
            <li data-role="fieldcontain">\
                ' + displaySearchArea(menuConfig.search_area.controls) + '\
            </li>\
            ' + displayContentsArea(menuConfig) + '\
        </ul>\
    ';
};

function displayMultiSearchAreaLayer(menuConfig) {

    var multiSearchAreaHtml = "";    

    for(var i=0; i<menuConfig.search_area.length; i++) {
        multiSearchAreaHtml = multiSearchAreaHtml + 
                                '<li data-role="controlgroup">\
                                    ' + menuConfig.search_area[i].title + '\
                                </li>\
                                    <li data-role="fieldcontain">\
                                        ' + displaySearchArea(menuConfig.search_area[i].controls) + '\
                                </li>';
    }
    return '\
        <ul data-role="listview">\
            '+multiSearchAreaHtml+'\
        </ul>\
    ';
};

function displaySearchArea(controls) {
    
    var html = '';

    var controls_index = 0;
    var widthCount = 0;
    do {
        var grid_index = 95, block_index =97; // ascii 97 == 'a', grid는 1개일 경우 solo이므로 95로 설정
        widthCount = ($.inArray("{newline}", controls.slice(controls_index, controls.length)) == -1) ? controls.slice(controls_index, controls.length).length : $.inArray("{newline}", controls.slice(controls_index, controls.length));
        if(widthCount == 0) break;
        var widthPercentage = (widthCount != 0) ? parseFloat(100 / widthCount).toFixed(2) : 100;
        var grid_value = (grid_index + widthCount == 96) ? 
            "solo" : 
            String.fromCharCode(grid_index + widthCount);
        
        html += '<div class="ui-grid-' + grid_value + '" style="font-size:medium;width:100%;">';
        for(var i=0; i<widthCount; i++) {
            var control = config.controls.find(x => x.id === controls[i+controls_index]);
            html += '\
                <div class="ui-block-' + String.fromCharCode(block_index + i) + '" style="width:' + widthPercentage + '%">\
                    ' + displayControl(control) + '\
                </div>\
            ';
        }
        html += '</div>';
        controls_index += widthCount + 1;
    } while (true) 
    return html;
}

function displayContentsArea(menuConfig) {
    if (menuConfig.contents_area.grid){
        var gridConfig = menuConfig.contents_area.grid;
    
        var html = '\
            <li data-role="controlgroup">\
                ' + menuConfig.contents_area.title + '\
                <div id="dvtitleArea"></div>\
            </li>\
            <li data-role="controlgroup">\
                ' + displayColumnOptions(gridConfig.columns, gridConfig.checkbox_name) + '\
            </li>\
            <li data-role="controlgroup">\
                <div class="ui-grid-solo" style="width:100%;overflow:auto;">\
                    <table id="' +  gridConfig.id + '" style="width:100%;"></table>\
                    <div id="' + gridConfig.toolbar.id + '"></div>\
                </div>\
            </li>\
        ';
    } else {
        var html = '\
            <li data-role="controlgroup">\
                ' + menuConfig.contents_area.title + '\
            </li>\
        ';
    }
    

    return html;
}

function displayControl(control) {
    var html = "";
    if(control.tag == "select") {
        html = '<select name="' + control.name + '" id="' + control.id + '" title="' + control.title + '" data-mini="true">';
        for(var i=0; i<control.options.length; i++) {
            var selected = (control.options[i].selected == "selected") ? "selected=\"selected\"" : "";
            var name = (control.options[i].name != undefined && control.options[i].name != "") ? " name=\"" + control.options[i].name + "\"" : "";
            html += '<option ' + name + 'value="' + control.options[i].value + '" ' + selected + '>' + control.options[i].text + '</option>';
        }
        html += '</select>';
    } else if(control.tag == "text") {
        var type = (control.type != undefined && control.type != "") ? control.type : "text";
        html = '<input id="' + control.id + '" name="' + control.name + '" type="' + type + '" value="' + control.value + '" placeholder="' + control.placeholder + '" title="' + control.title + '" data-mini="true" />';
    } else if(control.tag == "button") {
        html = '<input id="' + control.id + '" type="button" value="' + control.value + '" title="' + control.title + '" data-mini="true"/>';
    } else if(control.tag == "password") {
        var type = (control.type != undefined && control.type != "") ? control.type : "password";
        html = '<input id="' + control.id + '" name="' + control.name + '" type="' + type + '" value="' + control.value + '" placeholder="' + control.placeholder + '" title="' + control.title + '" data-mini="true" />';
    } else {
        html = '<' + control.tag + ' id="' + control.id + '" title="' + control.title + '"></' + control.tag + '>';
    }
    return html;
}

function displayColumnOptions(gridColumns, checkbox_name) {
    if(gridColumns == undefined || gridColumns.length <= 0) {
        console.log("config error");
        return;
    }
    var innerHtml = '<table style="width:100%"><tr style="width:100%">';
    for(var i=0;i<gridColumns.length; i++) {
        var checked = "";
        if (gridColumns[i].hidden != undefined) {
            if(!gridColumns[i].hidden) {
                checked = "checked=\"checked\"";
            }
        } else {
            checked = "checked=\"checked\"";
        }
        
        if(i != 0 && i%8 == 0) {
            innerHtml += '</tr><tr>';
        }

        innerHtml += '<td style="width:8%">';
        innerHtml += '<input type="checkbox" id="chk_' + gridColumns[i].name + '" value="' + gridColumns[i].text +'" name="' + checkbox_name + '" ' + checked + '>';
        innerHtml += '<label for="chk_' + gridColumns[i].name + '">' + gridColumns[i].text + '</label>';
        innerHtml += "</td>";
        
        if(i == gridColumns.length-1) {
            innerHtml += '</tr>';
        }
    }
    innerHtml += '</table>';
    return innerHtml;
};

function setColumns(gridColumns, checkedColNames) {
    var result = [];
    var colNames = []
    var colModels = [];
    for(var i=0;i<gridColumns.length; i++) {
        if (jQuery.inArray(gridColumns[i].text, checkedColNames) == -1) {
            gridColumns[i].hidden = true;
            if(gridColumns[i].editrules == undefined) {
                gridColumns[i].editrules = {};
            }
            gridColumns[i].editrules.edithidden = true;
            gridColumns[i].editrules.required = true;
        } else {
            gridColumns[i].hidden = false;
            gridColumns[i].editrules = null;
        }
        colNames.push(gridColumns[i].text);
        colModels.push(gridColumns[i]);
    }
    result.push(colNames);
    result.push(colModels);
    return result;
}

function setAPIParameters(apiObj, obj) {
    var requestData = {};
    
    for(var i=0; i<apiObj.parameters.length; i++) {
        var parameter = (typeof apiObj.parameters[i] == "string") 
        ? config.apis.parameters.find(x => x.id === apiObj.parameters[i])
        : apiObj.parameters[i];
        
        parameter.in = (parameter.in == undefined || parameter.in == "") ? "query" : parameter.in;
        parameter.value_type = (parameter.value_type == undefined || parameter.value_type == "") ? "value" : parameter.value_type;
        parameter.type = (parameter.type == undefined || parameter.type == "") ? "string" : parameter.type;
        parameter.value = (parameter.value == undefined || parameter.value == "") ? "" : parameter.value;
        
        var dynamicValue = (parameter.value_type == "control") ? $("#" + parameter.value).val() // string
        : (parameter.value_type == "variable") ? eval(parameter.value) // variable
        : (parameter.value_type == "value") ? parameter.value // string
        : (parameter.value_type == "object") ? parameter.value // object
        : parameter.value;

        //dataType is "number" "boolean" "string" "object" "undefined" "function"
        var dataType = typeof dynamicValue;

        var parameterValue = (parameter.type == "string" && dataType == "object") ? JSON.stringify(dynamicValue) // string
        : ((parameter.type == "integer" || parameter.type == "int" || parameter.type == "number") && !$.isNumeric(dynamicValue)) ? null // int
        : (parameter.type == "integer" || parameter.type == "int" || parameter.type == "number") ? parseInt(dynamicValue) // int
        : (parameter.type == "bool" || parameter.type == "boolean") ? (dynamicValue == "true") // boolean
        : (parameter.type == "array" && (dynamicValue == null || dynamicValue == "")) ? [] //$.parseJSON("[]") 
        : (parameter.type == "array" && dataType !== "object") ? $.parseJSON(dynamicValue) 
        : (parameter.type == "object" && (dynamicValue == null || dynamicValue == "")) ? {} 
        : (parameter.type == "object" && dataType !== "object") ? $.parseJSON(dynamicValue) 
        : (parameter.type == "file" && parameter.value_type == "control") ? $("#" + parameter.value).prop('files')[0] 
        : dynamicValue;

        if(parameter.value_type == "object") {
            var jsonObj = {};
            for(var j=0; j<dynamicValue.length; j++) {
                
                var tmpParameter = (typeof dynamicValue[j] == "string") 
                ? config.apis.parameters.find(x => x.id === dynamicValue[j])
                : dynamicValue[j];

                var tmpDynamicValue = (tmpParameter.value_type == "control") ? $("#" + tmpParameter.value).val() // string
                : (tmpParameter.value_type == "variable") ? eval(dynamicValue[ij].value) // variable
                : (tmpParameter.value_type == "value") ? tmpParameter.value // string
                : tmpParameter.value;
                
                dataType = typeof tmpDynamicValue;

                var tmpParameterValue = (tmpParameter.type == "string" && dataType == "object") ? JSON.stringify(tmpDynamicValue) // string
                : (tmpParameter.type == "integer" || tmpParameter.type == "int" || tmpParameter.type == "number") ? parseInt(tmpDynamicValue) // int
                : (tmpParameter.type == "bool" || tmpParameter.type == "boolean") ? (tmpDynamicValue == "true") // boolean
                : (tmpParameter.type == "array" && (tmpDynamicValue == null || tmpDynamicValue == "")) ? []
                : (tmpParameter.type == "array" && tmpDynamicValue != "") ? $.parseJSON(tmpDynamicValue) 
                : (tmpParameter.type == "object" && (tmpDynamicValue == null || tmpDynamicValue == "")) ? {}
                : (tmpParameter.type == "object" && tmpDynamicValue != "") ? $.parseJSON(tmpDynamicValue) 
                : tmpDynamicValue;
                
                jsonObj[tmpParameter.name] = tmpParameterValue;
                if(tmpParameter.required) {
                    if(
                        (tmpParameterValue == undefined || tmpParameterValue == null) ||
                        (tmpParameter.type == "string" && tmpParameterValue == "") ||
                        (tmpParameter.type == "string" && dataType == "object" && tmpParameterValue == "[]") ||
                        ((tmpParameter.type == "integer" || tmpParameter.type == "int" || tmpParameter.type == "number") && !$.isNumeric(tmpParameterValue)) ||
                        (tmpParameter.type == "array" && tmpParameterValue == []) ||
                        (tmpParameter.type == "object" && tmpParameterValue == {})
                    ) {
                        panelAlertLayer("request body parameter invalid : " + tmpParameter.name, "info");
                        return null;
                    }
                }
            }
            if (requestData[parameter.in] == null) requestData[parameter.in] = {};
            
            requestData[parameter.in][parameter.name] = (parameter.type == "object") ? jsonObj : JSON.stringify(jsonObj);
            if(parameter.required) {
                if( requestData[parameter.in][parameter.name] == undefined || 
                    requestData[parameter.in][parameter.name] == null ||
                    requestData[parameter.in][parameter.name] == []
                ) {
                    panelAlertLayer("request parameter invalid : " + parameter.name, "info");
                    return null;
                }
            }
        } else {
            if (requestData[parameter.in] == null) requestData[parameter.in] = {};
            requestData[parameter.in][parameter.name] = parameterValue;
            if(parameter.required) {
                if(
                    (parameterValue == undefined || parameterValue == null) ||
                    (parameter.type == "string" && parameterValue == "") ||
                    (parameter.type == "string" && dataType == "object" && parameterValue == "[]") ||
                    ((parameter.type == "integer" || parameter.type == "int" || parameter.type == "number") && !$.isNumeric(parameterValue)) ||
                    (parameter.type == "array" && parameterValue == []) ||
                    (parameter.type == "object" && parameterValue == {})
                ) {
                    panelAlertLayer("request parameter invalid : " + parameter.name, "info");
                    return null;
                }
            }
        }
    }
    
    if(apiObj.content_type == "application/json") {
        requestData["query"] = JSON.stringify(requestData["query"]);
    }
    
    if(apiObj.content_type == "multipart/form-data") {
        var result = new FormData();
        for ( var key in requestData["query"] ) {
            result.append(key, requestData["query"][key]);
        }
        requestData["query"] = result;
    }
    
    return requestData;
}