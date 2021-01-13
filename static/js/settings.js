getSettings();

function setSettings(request){
    const response = JSON.parse(request.response);
    
    setPacketSettings(response);
    setLatencySettings(response);
    setBandwidthSettings(response);
}

function setPacketSettings(response){
    document.getElementById("packet-loss-number").value = response.Loss;
    document.getElementById("packet-corruption-number").value = response.Corruption;
    document.getElementById("packet-duplication-number").value = response.Duplication;
}

function setLatencySettings(response){
    document.getElementById("latency-base-number").value = response.LatencyRule.BaseLatency;
    document.getElementById("latency-variation-number").value = response.LatencyRule.Variation;
    document.getElementById("latency-correlation-number").value = response.LatencyRule.Correlation;
}

function setBandwidthSettings(response){
    document.getElementById("bandwidth-rate-number").value = response.BandwidthRule.Rate;
    document.getElementById("bandwidth-burst-number").value = response.BandwidthRule.Burst;
    document.getElementById("bandwidth-max-number").value = response.BandwidthRule.MaxLatency;
}

function getSettings(){
    const request = new XMLHttpRequest();
    request.onload = function(){
        setSettings(request)
    };
    request.open("GET", "/api/settings");
    request.send();
}

function savePacketManipulation() {
    const loss = parseInt(document.getElementById("packet-loss-number").value);
    const corruption = parseInt(document.getElementById("packet-corruption-number").value);
    const duplication = parseInt(document.getElementById("packet-duplication-number").value);

    const packetManipulation = {
        Loss: loss,
        Corruption: corruption,
        Duplication: duplication
    };

    // TO-DO: Think about this being in global scoping.
    sendSettings(packetManipulation);
}

function saveLatency() {
    const baseLatency = parseInt(document.getElementById("latency-base-number").value);
    const variation = parseInt(document.getElementById("latency-variation-number").value);
    const correlation = parseInt(document.getElementById("latency-correlation-number").value);

    const latencyRule = {
        BaseLatency: baseLatency,
        Variation: variation,
        Correlation: correlation
    };

    sendSettings({LatencyRule: latencyRule})
}

function saveBandwidth(e) {
    const rate = document.getElementById("bandwidth-rate-number").value;
    const burst = document.getElementById("bandwidth-burst-number").value;
    const maxLatency = parseInt(document.getElementById("bandwidth-max-number").value);

    const bandwidthRule = {
        Rate: rate,
        Burst: burst,
        MaxLatency: maxLatency 
    };

    sendSettings({BandwidthRule: bandwidthRule})
}

function sendSettings(settings){
    const json = JSON.stringify(settings);

    const request = new XMLHttpRequest();
    request.onload = function() {
        if(request.status >= 200 && request.status < 300){
            showMessagePopup("Changes successfully saved.")
        }else if(request.status >= 400 && request.status < 500){
            showMessagePopup("A client-side error occured.", "rgb(187,33,36)")
        }else if(request.status >= 500 && request.status < 600){
            showMessagePopup("A server-side error occured.", "rgb(187,33,36)")
        }
    };
    request.open("POST", "/api/settings");
    request.send(json);
}

let messagePopupIntervalId = -1;
let messagePopupTime = 3;

function showMessagePopup(message, color="rgb(34,187,51)") {
    const messagePopupElement = document.getElementById("message-popup");
    messagePopupElement.innerHTML = message;
    messagePopupElement.style.backgroundColor = color;
    messagePopupElement.style.display = "inline-block";
    messagePopupTime = 3;

    if(messagePopupIntervalId < 0) {
        messagePopupIntervalId = setInterval(() => {
            if(messagePopupTime <= 0){
                messagePopupElement.style.display = "none";
                clearInterval(messagePopupIntervalId);
                messagePopupIntervalId = 0;
            }
            messagePopupTime--
        }, 1000);
    }

}