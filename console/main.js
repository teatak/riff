import 'whatwg-fetch'
import React from 'react'
import { render } from 'react-dom'
import { Provider } from 'react-redux'
import thunk from 'redux-thunk'
import { createStore, applyMiddleware } from 'redux'
import { Router, Switch, Route } from 'react-router-dom'
import BrowserHistory from './history/browserhistory'
import reducer from './reducers'

import App from './app'

import './style/main.css'

const store = createStore(reducer,applyMiddleware(thunk));
const router = <Provider store={store}>
    <div>
        <Router history={BrowserHistory}>
            <Switch>
                <Route strict path="/" component={App} />
            </Switch>
        </Router>
    </div>
</Provider>;

render(router, document.getElementById('root'));
