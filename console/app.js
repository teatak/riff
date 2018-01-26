import React from 'react'
import { Redirect, Switch, Route } from 'react-router-dom'
import Menu from './components/menu'
import Nodes from './components/nodes'
import Services from './components/services'
import Explorer from './components/explorer'

class App extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return <div>
            <Menu/>
            <div className="container">
                <Switch>
                    <Redirect exact from="/" to='/nodes'/>
                    <Route strict path="/nodes" component={Nodes} />
                    <Route strict path="/services" component={Services} />
                    <Route strict path="/explorer" component={Explorer} />
                </Switch>
            </div>
        </div>
    }
}

export default App