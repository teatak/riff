import React from 'react'
import { connect } from 'react-redux'
import { Switch, Route, Redirect } from 'react-router-dom'

class Home extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return <div>Home</div>
    }
}

export default Home