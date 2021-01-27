import TableBuilder from "./modules/TableBuilder.js";

updateStatus();

const homeContainer = document.getElementById("home-container");

async function updateStatus(){
    const response = await fetchDeviceStatus();

    const status = JSON.parse(response);

    const deviceInformationHeader = document.createElement("h3");
    deviceInformationHeader.innerHTML = "Device Information";

    homeContainer.appendChild(deviceInformationHeader);

    buildHostDeviceTable(status);

    const networkStatusHeader = document.createElement("h3");
    networkStatusHeader.innerHTML = "Network Information";

    homeContainer.appendChild(networkStatusHeader);

    const hr = document.createElement("hr");
    homeContainer.appendChild(hr);

    for(let network of status.NetworkStatus) {
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

function buildHostDeviceTable(status) {
    const table = document.createElement("table");
    TableBuilder.addTableRow(table, [status.HostName, ""], "th");
    TableBuilder.addTableRow(table, ["Operating System", status.OperatingSystem]);
    TableBuilder.addTableRow(table, ["Software Version", status.TestDeviceSoftwareVersion]);
    homeContainer.appendChild(table);
}

function buildNetworkTable(network) {
    const table = document.createElement("table");
    const header = TableBuilder.addTableRow(table, [network.Name, ""], "th");

    const keys = Object.keys(network)
    for(let key of keys) {
        if(key == "Name"){continue;}
        if(key == "Addresses"){continue;}
        const tableRow = TableBuilder.addTableRow(table, [key, network[key]]);
    }

    for(let address of network.Addresses) {
        const tableRow = TableBuilder.addTableRow(table, [address.Name, address.Address]);
    }
    
    homeContainer.appendChild(table);
}
