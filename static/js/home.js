import TableBuilder from "./modules/TableBuilder.js";

updateStatus();

const homeContainer = document.getElementById("home-container");

async function updateStatus(){
    const response = await fetchDeviceStatus();

    networks = JSON.parse(response);

    for(network of networks) {
        buildNetworkTable(network);
    }
}

function fetchDeviceStatus() {
    const request = new XMLHttpRequest();
    request.open("GET", "/api/status");
    return new Promise(resolve => {
        request.onload = ()=>{resolve(request.response);}
        request.send();
    });
}

function buildNetworkTable(network) {
    const table = document.createElement("table");
    const header = TableBuilder.addTableRow(table, "th", [network.Name.substring(0,network.Name.length-1), ""]);

    table.innerHTML = header;

    const keys = Object.keys(network)
    for(key of keys) {
        if(key == "Name"){continue;}
        const tableRow = TableBuilder.addTableRow(table, [key, network[key]]);
    }
    
    homeContainer.appendChild(table);
}