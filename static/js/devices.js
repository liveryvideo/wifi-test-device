getDevices();

setInterval(getDevices, 1000);

function setDevices(request) {
    const deviceContainer = document.getElementById("table-container");
    deviceContainer.innerHTML = "";

    const table = document.createElement("table");
    
    const headers = document.createElement("tr");
    addTableData(headers, "Host Name", "th");
    addTableData(headers, "Client Identifier", "th");
    addTableData(headers, "Link Address", "th");
    addTableData(headers, "Expiration Time", "th");
    addTableData(headers, "Connected", "th");

    table.appendChild(headers);
    deviceContainer.append(table);

    const response = JSON.parse(request.response)
    for(let i = 0; i < response.length; i++){
        const deviceRow = createDeviceRow(response[i])
        table.appendChild(deviceRow)
    }
}

function addTableData(tableRow, contents, tag="td"){
    const dataElement = document.createElement(tag);
    dataElement.innerHTML = contents;
    tableRow.appendChild(dataElement);
}

function createDeviceRow(device){
    const deviceRow = document.createElement("tr");
    addTableData(deviceRow, device.HostName);
    addTableData(deviceRow, device.ClientIdentifier);
    addTableData(deviceRow, device.LinkAddress);
    addTableData(deviceRow, device.ExpirationTime);
    addTableData(deviceRow, device.Connected);
    return deviceRow;
}

function getDevices() {
    const request = new XMLHttpRequest();
    request.onload = function(){
        setDevices(request)
    };
    request.open("GET", "/api/devices");
    request.send();
}