
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
    if (key == "time") {
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