class Ws {
    ws = null;
    watchParam = null;
    onWatch = null;
    logsParam = null;
    onLogs = null;
    start = () => {
        let loc = window.location, uri = "";
        if (loc.protocol === "https:") {
            uri = "wss:";
        } else {
            uri = "ws:";
        }
        uri += "//" + loc.host + "/ws";
        this.ws = new WebSocket(uri);
        this.ws.onmessage = (evt) => {
            if (evt.data) {
                let response = JSON.parse(evt.data);
                if (response.event) {
                    switch (response.event) {
                        case "NodeChange" :
                        case "ServiceChange" :
                            if (this.onWatch) {
                                this.onWatch(response.body);
                            }
                            break;
                        case "Logs":
                            if (this.onLogs) {
                                this.onLogs(response.body);
                            }
                    }
                }

            }
        };
        this.ws.onclose = () => {
            setTimeout(() => {
                this.start();
            }, 5000);
        };
        this.ws.onopen = (evt) => {
            //send watch and logs param when reconnect
            if (this.watchParam) {
                this.send(this.watchParam);
            }
            if (this.logsParam) {
                this.send(this.logsParam);
            }
        };
    };
    send = (msg) => {
        this.ws.send(
            JSON.stringify(msg)
        )
    };
    watch = (param, onWatch) => {
        this.watchParam = param;
        this.onWatch = onWatch;
        this.send(this.watchParam);
    };
    logs = (param, onLogs) => {
        this.logsParam = param;
        this.onLogs = onLogs;
        this.send(this.logsParam);
    }
}

export default new Ws()