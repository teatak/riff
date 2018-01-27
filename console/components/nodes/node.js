import React from 'react'
import {withRouter} from 'react-router-dom'
import {connect} from 'react-redux'
import {getNode} from "../../reducers/nodes";
import ArrowDown from '../icons/arrowDown'
import ArrowUp from '../icons/arrowUp'

const mapStateToProps = (state, ownProps) => {
    return {
        nodes: state.nodes
    }
};

const mapDispatchToProps = (dispatch) => {
    return {
        getNode: (nodeName) => {
            dispatch(getNode(nodeName));
        },
    }
};

class Node extends React.Component {
    constructor(props) {
        super(props);
        this.state = {toggle: {}};
    }

    componentWillMount() {
        let nodeName = this.props.match.params.nodeName;
        this.props.getNode(nodeName)
    }

    componentWillReceiveProps(nextProps) {
        if (this.props.match.params.nodeName !== nextProps.match.params.nodeName) {
            this.setState({toggle: {}});
            this.props.getNode(nextProps.match.params.nodeName)
        }
    }

    toggle = (name) => {
        this.setState({
            toggle: {
                ...this.state.toggle,
                [name]: !this.state.toggle[name]
            }
        });
    };

    renderList() {
        const {nodes} = this.props;
        if (nodes.data.services) {
            return <ul className="services">
                <li className="title">Services</li>
                {nodes.data.services.map((service, index) => {
                    let className = "item " + service.state.toLowerCase();
                    return <li className={className} key={service.name}>
                        <div className="basic" onClick={() => {
                            this.toggle(service.name)
                        }}>
                            <div className="toggle">{this.state.toggle[service.name] ? <ArrowUp/> : <ArrowDown/>}</div>
                            <span className="name">{service.name}</span>
                            <span className="ipport">{service.port !== 0 ? ":" + service.port : ""}</span>
                        </div>

                        {this.state.toggle[service.name] ? <pre>
                            {service.config}
                        </pre> : null}
                    </li>
                })}
            </ul>
        }
    }

    render() {
        const {nodes} = this.props;
        if (nodes.fetchNode.status === 404) {
            return <div className="error">Not Found</div>
        }
        if (nodes.fetchNode.status === 500) {
            return <div className="error">{nodes.fetchNode.error}</div>
        }
        if (nodes.fetchNode.status === 200) {
            return <div>
                <div className="title">
                    <span className="name">{nodes.data.name}</span>
                    <span className="ipport">{nodes.data.ip}:{nodes.data.port}</span>
                </div>
                {this.renderList()}
            </div>
        } else {
            return null
        }
    }
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Node))