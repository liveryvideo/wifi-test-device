import TableBuilder from "./modules/TableBuilder.js";

const store = {
    page: 0,
    range: 20,
    remaining: 0
}

const PreviousButton = document.getElementById("button-prev");
const NextButton = document.getElementById("button-next");

PreviousButton.addEventListener("click", ()=>{prevPage()}, false);
NextButton.addEventListener("click", ()=>{nextPage()}, false);

updateLogs();

const logsContainer = document.getElementById("logs-container");

async function updateLogs(){
    const start = store.page*store.range;
    const end = (store.page*store.range)+store.range;

    const response = await fetchLogs(start, end);

    const json = JSON.parse(response);
    store.remaining = json.Remaining;

    if(store.remaining <= 0) {
        NextButton.disabled = true;
    }else{
        NextButton.disabled = false;
    }

    if(store.page <= 0){
        PreviousButton.disabled = true;
    }else{
        PreviousButton.disabled = false;
    }

    buildLogTable(json.Logs);
}

function fetchLogs(start, end) {
    const request = new XMLHttpRequest();
    request.open("GET", "/api/logs?start=" + start + "&end=" + end);
    console.log("GET", "/api/logs?start=" + start + "&end=" + end)
    return new Promise(resolve => {
        request.onload = ()=>{resolve(request.response);}
        request.send();
    });
}

function getLogsTable(){    
    let table = document.getElementById("logs-table")
    if(!table){
        table = document.createElement("table");
        table.id = "logs-table";
    }
    return table;
}

function buildLogTable(logs){
    const logsContainer = document.getElementById("logs-container");

    const table = getLogsTable();
    table.innerHTML = "";

    TableBuilder.addTableRow(table, ["Log", "time"], "th")
    for(let log of logs){
        TableBuilder.addTableRow(table, [log.Value, log.Time])
    }

    logsContainer.appendChild(table);
}

function prevPage(){
    if(store.page > 0){
        store.page--;
    }
    updateLogs();
}

function nextPage(){
    console.log(store)
    if(store.remaining > 0){
        store.page++
    }
    updateLogs();
}