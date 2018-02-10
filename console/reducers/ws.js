import {changeNode} from './nodes'
import {changeService} from "./services";

class Ws {
    constructor() {
        this.ws = null;
    }

    start = () => (dispatch, getState) => {
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
                            dispatch(changeNode(response.body));
                            break;
                        case "ServiceChange" :
                            dispatch(changeService(response.body));
                            break;
                    }
                }

            }
        };
        this.ws.onclose = () => {
            setTimeout(() => {
                dispatch(this.start())
            }, 5000);
        };
        this.ws.onopen = (evt) => {

        };
    };
    send = (msg) => (dispatch, getState) => {
        this.ws.send(
            JSON.stringify(msg)
        )
    }
}

export default new Ws()