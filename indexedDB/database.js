var Database = function(dbname, tablename, datakey, version) {
  //private members
  var db = null,
        dbname = dbname,
        tablename = tablename,
        datakey = datakey,
        version = version,
      trace = function(msg) {
          //Traces
          console.log(msg);
      },
      init = function() {

          //Make indexedDB compatible
          if (compatibility()) {
              open();
          }
      },
      compatibility = function() {

          trace("window.indexedDB: " + window.indexedDB);
          trace("window.mozIndexedDB: " + window.mozIndexedDB);
          trace("window.webkitIndexedDB: " + window.webkitIndexedDB);
          trace("window.msIndexedDB: " + window.msIndexedDB);

          window.indexedDB = window.indexedDB || window.webkitIndexedDB || window.mozIndexedDB;

          trace("window.IDBTransaction: " + window.IDBTransaction);
          trace("window.webkitIDBTransaction: " + window.webkitIDBTransaction);
          trace("window.msIDBTransaction: " + window.msIDBTransaction);

          window.IDBTransaction = window.IDBTransaction || window.webkitIDBTransaction || window.msIDBTransaction || window.OIDBTransaction;

          trace("window.IDBKeyRange: " + window.IDBKeyRange);
          trace("window.webkitIDBKeyRange: " + window.webkitIDBKeyRange);
          trace("window.msIDBKeyRange: " + window.msIDBKeyRange);

          window.IDBKeyRange = window.IDBKeyRange || window.webkitIDBKeyRange || window.msIDBKeyRange;

          if (window.indexedDB) {
              var span = document.querySelector("header h1 span");
              span.textContent = "Yes";
              span.style.color = "green";
              return true;
          }

          trace("Your browser does not support a stable version of IndexedDB.");
          return false;

      },
      deletedb = function(dbname) {
          trace("Delete " + dbname);

          var request = window.indexedDB.deleteDatabase(dbname);

          request.onsuccess = function() {
              trace("Database " + dbname + " deleted!");
          };

          request.onerror = function(event) {
              trace("deletedb(); error: " + event);
          };
      },
      open = function() {
          //3.1. Open a database async
          var request = window.indexedDB.open(tablename, version);

          //3.2 The database has changed its version (For IE 10 and Firefox)
          request.onupgradeneeded = function(event) {

              trace("Upgrade needed!");

              db = event.target.result;

              modifydb(); //Here we can modify the database
          };

          request.onsuccess = function(event) {
              trace("Database opened");

              db = event.target.result;

              //3.2 The database has changed its version (For Chrome)
              if (version != db.version && window.webkitIndexedDB) {

                  trace("version is different");

                  var setVersionreq = db.setVersion(version);

                  setVersionreq.onsuccess = modifydb; //Here we can modify the database
              }

              trace("Let's paint");
              items(); //4. Read our previous objects in the store (If there are any)
          };

          request.onerror = function(event) {
              trace("Database error: " + event);
          };
      },
      modifydb = function() {
          //3.3 Create / Modify object stores in our database 
          //2.Delete previous object store
          if (db.objectStoreNames.contains(tablename)) {
              db.deleteObjectStore(tablename);
              trace("db.deleteObjectStore(tablename);");
          }

          //3.Create object store
          var store = db.createObjectStore(tablename, {
              keyPath: "leaderboardId"
          });


      },        
      add = function(data) {
          //4. Add objects
          trace(tablename + " add(" + data[datakey] + "," + data + ")");

          // var trans = db.transaction(JSON.parse("[" + tablename + "]"), "readwrite");
          var trans = db.transaction([tablename], "readwrite");

          var store = trans.objectStore(tablename);
              
          var request = store.put(data);

          request.onsuccess = function(event) {
              trace("wish added!");
              items(); //5 Read items after adding
          };
      },
      items = function() {
          //5. Read
          trace("items(); called");

          var list = document.getElementById("list"),
              trans = db.transaction([tablename], "readonly"),
              store = trans.objectStore(tablename);

          list.innerHTML = "";

          var keyRange = IDBKeyRange.lowerBound(0);
          var cursorRequest = store.openCursor(keyRange);

          cursorRequest.onsuccess = function(event) {
              trace("Cursor opened!");

              var result = event.target.result;

              if (result === false || result === null){
                  return;
              }
              
              render(result.value); //4.1 Create HTML elements for this object
              result.continue ();

          };
      },
      render = function(item) {
          //5.1 Create DOM elements
          trace("Render items");

          var list = document.getElementById("list"),
              li = document.createElement("li"),
              a = document.createElement("a"),
              text = document.createTextNode(item.text);

          a.textContent = " X";

          //6. Delete elements
          a.addEventListener("click", function() {

              del(item[datakey]);

          });

          li.appendChild(text);
          li.appendChild(a);
          list.appendChild(li);
      },
      del = function(datakey) {
          //6. Delete items
          var transaction = db.transaction([tablename], "readwrite");
          var store = transaction.objectStore(tablename);

          var request = store.delete(datakey);

          request.onsuccess = function(event) {
              trace("Item deleted!");
              items(); //5.1 Read items after deleting
          };

          request.onerror = function(event) {
              trace("Error deleting: " + e);
          };
      };

  //public members
  return {
      init: init,
      add: add
  };
};

window.onload = function() {

  var database = new Database("scoresaber", "usr76561198830502286", "leaderboardId", 1);
  database.init(); //database name and database version
//   var data = {
//         "rank": 395,
//         "scoreId": 59576761,
//         "score": 1442332,
//         "unmodififiedScore": 1442332,
//         "mods": "",
//         "pp": 363.267,
//         "weight": 0.4406849377227,
//         "timeSet": "2022-01-26T12:26:33.000Z",
//         "leaderboardId": 367651,
//         "songHash": "27860892DD00DF00DA30A3DB12C23A6999EB853F",
//         "songName": "Night Raid With A Dragon",
//         "songSubName": "(Fvrwvrd Remix)",
//         "songAuthorName": "Camellia",
//         "levelAuthorName": "Jabob vs. Narwhal",
//         "difficulty": 7,
//         "difficultyRaw": "_Expert_SoloStandard",
//         "maxScore": 1529155
//     };
//   database.add(367651, data);
  $("#btnSave").click(function () {
      alert("aa");
    var data = {
        "rank": 395,
        "scoreId": 59576761,
        "score": 1442332,
        "unmodififiedScore": 1442332,
        "mods": "",
        "pp": Number($("#wish").val()),
        "weight": 0.4406849377227,
        "timeSet": "2022-01-26T12:26:33.000Z",
        "leaderboardId": 367651,
        "songHash": "27860892DD00DF00DA30A3DB12C23A6999EB853F",
        "songName": "Night Raid With A Dragon",
        "songSubName": "(Fvrwvrd Remix)",
        "songAuthorName": "Camellia",
        "levelAuthorName": "Jabob vs. Narwhal",
        "difficulty": 7,
        "difficultyRaw": "_Expert_SoloStandard",
        "maxScore": 1529155
    };
    database.add(data);
  });
};