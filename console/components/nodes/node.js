import React from 'react'
import { withRouter } from 'react-router-dom'

class Node extends React.Component {
    constructor(props) {
        super(props);
    }
    componentWillMount() {
        let nodeName = this.props.match.params.nodeName;
    }
    render() {
        return <div>{this.props.match.params.nodeName}</div>
    }
}

export default withRouter(Node)