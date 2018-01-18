import React from 'react'

import './index.css'

class Menu extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return <div className="menu">
            <img src="/static/images/logo.svg" />
        </div>
    }
}

export default Menu