
Date.prototype.yyyyMMdd = function() {
    var yyyy = this.getFullYear().toString();
    var mm = (this.getMonth()+1).toString(); // getMonth() is zero-based
    var dd  = this.getDate().toString();
    return yyyy + (mm[1]?mm:"0"+mm[0]) + (dd[1]?dd:"0"+dd[0]); // padding
};

Date.prototype.yyyyMM = function() {
    var yyyy = this.getFullYear().toString();
    var mm = (this.getMonth()+1).toString(); // getMonth() is zero-based
    return yyyy + (mm[1]?mm:"0"+mm[0]); // padding
};

Date.prototype.getWeek = function() {
    var date = new Date(this.getTime());
    date.setHours(0, 0, 0, 0);
    // Thursday in current week decides the year.
    date.setDate(date.getDate() + 3 - (date.getDay() + 6) % 7);
    // January 4 is always in week 1.
    var week1 = new Date(date.getFullYear(), 0, 4);
    // Adjust to Thursday in week 1 and count number of weeks from date to week1.
    return 1 + Math.round(((date.getTime() - week1.getTime()) / 86400000 - 3 + (week1.getDay() + 6) % 7) / 7);
}

Date.prototype.getWeekYear = function() {
    var date = new Date(this.getTime());
    date.setDate(date.getDate() + 3 - (date.getDay() + 6) % 7);
    return date.getFullYear();
}

Date.prototype.yyyyWeekNo = function () { 
    // var yyyy = this.getFullYear().toString();
    // var weekNo = $.datepicker.iso8601Week(this).toString();
    // return yyyy + (weekNo[1]?weekNo:"0"+weekNo[0]); // padding
    var date = new Date(this.getTime());
    var weekYear = date.getWeekYear().toString();
    var weekNo = date.getWeek().toString();
    return weekYear + (weekNo[1]?weekNo:"0"+weekNo[0]); // padding
}

Date.prototype.yyyyMMddHHmmss = function() {    
    var yyyy = this.getFullYear().toString();
    var mm = (this.getMonth()+1).toString(); // getMonth() is zero-based
    var dd  = this.getDate().toString();
    var hhmmss = this.getHours().toString() + this.getMinutes().toString() + this.getSeconds().toString();
    return yyyy + (mm[1]?mm:"0"+mm[0]) + (dd[1]?dd:"0"+dd[0]) + hhmmss; // padding
};

function ConvertToString(key, value) {
    if (key == "timestamp") {
        if (value.length == 10) {
            value = value + "000";
        } else if (value.length == 19) {
            value = value.substring(0, 13);
        }
        date = new Date(parseInt(value));
        var year = date.getFullYear();

        var month = (1 + date.getMonth()).toString();
        month = month.length > 1 ? month : '0' + month;

        var day = date.getDate().toString();
        day = day.length > 1 ? day : '0' + day;

        var hour = date.getHours().toString();
        hour = hour.length > 1 ? hour : '0' + hour;
        
        var minutes = date.getMinutes().toString();
        minutes = minutes.length > 1 ? minutes : '0' + minutes;

        return year + "-" + month + "-" + day + " " + hour  + ':' + minutes;
    } else if (key == "time") {
        var date = new Date(value.timeSet).toISOString().
        replace(/T/, ' ').      // replace T with a space
        replace(/\..+/, '');
        return date;
    }
    else if (key == "pp") {
        return value.pp.toFixed(2);
    }
    else if (key == "stars") {
        return value.stars + "★";
    }
    else if (key == "mypp") {
        return (value.pp * value.weight).toFixed(2);
    }
    else if (key == "maxPP") {
        return value.maxPP.toFixed(2);
    }
    else if (key == "weight") {
        return value.weight.toFixed(5);
    }
    else if (key == "accuracy") {
        return ((value.score * 100) / value.maxScore).toFixed(2);
    }
    else if (key == "rating") {
        return value.rating.toFixed(2);
    }
    else if (key == "durationSeconds") {
        var min = value.durationSeconds / 60;
        var sec = value.durationSeconds % 60;
        return Math.floor(min) + "m " + sec + "s";
    }
    else if (key == "difficulty") {
        return (value.difficulty == 1) ? "Easy" : 
            (value.difficulty == 3) ? "Normal" : 
            (value.difficulty == 5) ? "Hard" : 
            (value.difficulty == 7) ? "Expert" : 
            (value.difficulty == 9) ? "Expert+" : 
            value.difficulty;
    }
    else {

    }
    return null
}

function popupInfoLayer(name) {
    var ctrl = "#popupInfo";
    var html = '<p>';
    if (name == "service_type") {
        html += '    서비스 타입<br />';
        html += '    - application_no > App 번호<br />';
        html += '    - game_no > 게임 번호<br />';
        html += '    - game_id > 게임 ID<br />';
    }
    if (name == "service_id") {
        html += '    서비스 ID<br />';
        html += '    - 서비스 타입에 대응하는 서비스 ID<br />';
    }
    if (name == "user_type") {
        html += '    사용자 식별타입<br />';
        html += '    - member_no > 회원번호 기반으로 랭킹을 구성할 경우.<br />';
        html += '    - profile_id > 프로필 ID 기반으로 랭킹을 구성할 경우.<br />';
        html += '    - character_no > Character 번호 기반으로 랭킹을 구성할 경우.<br />';
        html += '    - game_no > 게임번호 기반으로 랭킹을 구성할 경우.<br />';
        html += '    - game_id > 게임ID 기반으로 랭킹을 구성할 경우.<br />';
        html += '    - application_no > App 번호 기반으로 랭킹을 구성할 경우.<br />';
    }
    html += '</p>';
    $(ctrl).empty().append(html);
    $(ctrl).trigger('create');
    $(ctrl).popup('open');
}

function removeMapProperties() {
    for(var i=0;i<mapList.length;i++){
        delete mapList[i].artist;
        delete mapList[i].beatSaverKey;
        delete mapList[i].downloads;
        delete mapList[i].duration;
        delete mapList[i].id;
        //delete mapList[i].mapper;
        delete mapList[i].downvotes;
        delete mapList[i].recentScores;
        delete mapList[i].scores;
        delete mapList[i].upvotes;
        delete mapList[i]._id;
    }
}

function sortJSONArray(list, prop, asc) {
    list.sort(function(a, b) {
        if (asc) {
            return (a[prop] > b[prop]) ? 1 : ((a[prop] < b[prop]) ? -1 : 0);
        } else {
            return (b[prop] > a[prop]) ? 1 : ((b[prop] < a[prop]) ? -1 : 0);
        }
    });
}

function removeBeatSaberKeyString(str) {
    console.log(str);
    const regex = /x{2,}/;
    const match = str.match(regex);
    
    if (match) {
      const index = match.index;
      return str.substring(0, index);
    }
    return str;
  }

  function getPlaylistImage() {
    return "data:image/png;base64,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";
  }