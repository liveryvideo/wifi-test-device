import TableBuilder from "./modules/TableBuilder.js";

updateLogs();

const logsContainer = document.getElementById("logs-container");

async function updateLogs(){
    const response = await fetchLogs();

    logs = JSON.parse(response);
    buildLogTable(logs);
}

function fetchLogs() {
    const request = new XMLHttpRequest();
    request.open("GET", "/api/logs");
    return new Promise(resolve => {
        request.onload = ()=>{resolve(request.response);}
        request.send();
    });
}

function buildLogTable(logs){
    const logsContainer = document.getElementById("logs-container");
    logsContainer.innerHTML = "";

    const table = document.createElement("table");

    TableBuilder.addTableRow(table, ["Log", "time"], "th")
    for(let log of logs){
        TableBuilder.addTableRow(table, [log.Value, log.Time])
    }

    logsContainer.appendChild(table);
}