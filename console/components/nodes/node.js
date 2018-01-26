import React from 'react'
import { withRouter } from 'react-router-dom'
import { connect } from 'react-redux'
import { getNode} from "../../reducers/nodes";

const mapStateToProps = (state, ownProps) => {
    return {
        nodes: state.nodes
    }
};

const mapDispatchToProps = (dispatch) => {
    return {
        getNode: (nodeName) => {
            dispatch(getNode(nodeName));
        }
    }
};

class Node extends React.Component {
    constructor(props) {
        super(props);
    }
    componentWillMount() {
        let nodeName = this.props.match.params.nodeName;
        this.props.getNode(nodeName)
    }
    componentWillReceiveProps(nextProps) {
        if(this.props.match.params.nodeName !== nextProps.match.params.nodeName) {
            this.props.getNode(nextProps.match.params.nodeName)
        }
    }
    renderList() {
        const { nodes } = this.props;
        if (nodes.data.services) {
            return <ul className="services">
                <li className="title">Services</li>
                {nodes.data.services.map((service, index) => {
                    let className = "item "+service.state.toLowerCase();
                    return <li className={className} key={service.name}>
                        <span className="name">{service.name}</span>
                        <span className="ipport">{service.port !== 0?":"+service.port:""}</span>
                    </li>
                })}
            </ul>
        }
    }
    render() {
        const { nodes } = this.props;
        return <div>
            <div className="title">
                <span className="name">{nodes.data.name}</span>
                <span className="ipport">{nodes.data.ip}:{nodes.data.port}</span>
            </div>
            {this.renderList()}
            </div>
    }
}

export default withRouter(connect(mapStateToProps,mapDispatchToProps)(Node))