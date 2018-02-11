class Ws {
    ws = null;
    watchMsg = null;
    onWatch = null;
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
            if (this.watchMsg) {
                this.send(this.watchMsg)
            }
        };
    };
    send = (msg) => {
        this.ws.send(
            JSON.stringify(msg)
        )
    };
    watch = (msg, onWatch) => {
        this.watchMsg = msg;
        this.onWatch = onWatch;
        this.send(this.watchMsg);
    };
}

export default new Ws()