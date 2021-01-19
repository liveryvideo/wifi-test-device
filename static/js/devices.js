import TableBuilder from "./modules/TableBuilder.js";

getDevices();

setInterval(getDevices, 1000);

function setDevices(request) {
    const deviceContainer = document.getElementById("table-container");
    deviceContainer.innerHTML = "";

    let table = document.createElement("table");

    const headers = TableBuilder.addTableRow(table);
    TableBuilder.addTableRow(table, ["Host Name", "Client Identifier", "Link Address", "Expiration Time", "Connected"], "th")

    const response = JSON.parse(request.response)
    for(let i = 0; i < response.length; i++){
        const deviceRow = createDeviceRow(table, response[i])
    }

    deviceContainer.append(table);
}

function createDeviceRow(table, device){
    const deviceRow = TableBuilder.addTableRow(table);
    TableBuilder.addTableData(deviceRow, device.HostName);
    TableBuilder.addTableData(deviceRow, device.ClientIdentifier);
    TableBuilder.addTableData(deviceRow, device.LinkAddress);
    TableBuilder.addTableData(deviceRow, device.ExpirationTime);
    TableBuilder.addTableData(deviceRow, device.Connected);
}

function getDevices() {
    const request = new XMLHttpRequest();
    request.onload = function(){
        setDevices(request)
    };
    request.open("GET", "/api/devices");
    request.send();
}