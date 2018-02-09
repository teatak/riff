import React from 'react'
import {Redirect, Route, Switch} from 'react-router-dom'
import Menu from './components/menu'
import Nodes from './components/nodes'
import Services from './components/services'
import Explorer from './components/explorer'
import Logs from './components/logs'

class App extends React.Component {
    constructor(props) {
        super(props);
        this.state = {logs: false};

        let loc = window.location, uri = "";
        if (loc.protocol === "https:") {
            uri = "wss:";
        } else {
            uri = "ws:";
        }
        uri += "//" + loc.host+"/ws";
        this.start(uri)
    }

     start = (url) => {
        let ws = new WebSocket(url);
        ws.onmessage = (evt) => { console.log(evt) };
        ws.onclose = () => {
            setTimeout(() => {this.start(url)}, 5000);
        };
        ws.onopen = (evt) => {
            console.log(evt);
            setTimeout(() => {ws.send(
                JSON.stringify({
                    text:"create"
                })
            )});
        };
    };

    toggleLogs = () => {
        this.setState({logs: !this.state.logs});
    };

    render() {
        return <div>
            <Menu/>
            <div className="container">
                <Switch>
                    <Redirect exact from="/" to='/nodes'/>
                    <Route strict path="/nodes" component={Nodes}/>
                    <Route strict path="/services" component={Services}/>
                    <Route strict path="/explorer" component={Explorer}/>
                </Switch>
            </div>
            <div className="logs-container">
                {!this.state.logs ? <div className="handle" onClick={this.toggleLogs}>Event Log</div> : null}
                {this.state.logs ? <Logs onClose={this.toggleLogs}/> : null}
            </div>
        </div>
    }
}

export default App