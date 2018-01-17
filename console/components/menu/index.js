import React from 'react'

import './index.css'

class Menu extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return <div className="nav-menu">
            <img src="/static/images/logo.svg" />
            <span className="title">Console</span>
        </div>
    }
}

export default Menu