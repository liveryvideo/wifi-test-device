updateStatus();

const homeContainer = document.getElementById("home-container");

async function updateStatus(){
    const response = await fetchDeviceStatus();

    networks = JSON.parse(response);

    for(network of networks) {
        buildTable(network);
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

function buildTable(network) {
    const table = document.createElement("table");
    const header = "<tr><th>" + network.Name.substring(0,network.Name.length-1) + "</th><th></th></tr>";

    table.innerHTML = header;

    const keys = Object.keys(network)
    for(key of keys) {
        if(key == "Name"){continue;}
        const tableRow = document.createElement("tr");
        tableRow.innerHTML = "<tr><td>" + key + "</td><td>" + network[key] + "</td></tr>"
        table.appendChild(tableRow);
    }
    
    homeContainer.appendChild(table);
}