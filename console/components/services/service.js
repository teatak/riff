import React from 'react'
import {NavLink, withRouter} from 'react-router-dom'
import {connect} from 'react-redux'
import {getService} from "../../reducers/services";
import ArrowDown from '../icons/arrowDown'
import ArrowUp from '../icons/arrowUp'

const mapStateToProps = (state, ownProps) => {
    return {
        services: state.services
    }
};

const mapDispatchToProps = (dispatch) => {
    return {
        getService: (nodeName) => {
            dispatch(getService(nodeName));
        },
    }
};

class Service extends React.Component {
    constructor(props) {
        super(props);
        this.state = {toggle: {}};
    }

    componentWillMount() {
        let serviceName = this.props.match.params.serviceName;
        this.props.getService(serviceName)
    }

    componentWillReceiveProps(nextProps) {
        const locationChanged = nextProps.location !== this.props.location;
        if (locationChanged) {
            this.setState({toggle: {}});
            this.props.getService(nextProps.match.params.serviceName)
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
        const {services} = this.props;
        if (services.data.nodes) {
            return <ul className="nestnodes">
                <li className="title">Nodes</li>
                {services.data.nodes.map((node, index) => {
                    let className = "item " + node.state.toLowerCase();
                    return <li className={className} key={node.name}>
                        <div className="basic">
                            <div className="toggle" onClick={() => {
                                this.toggle(node.name)
                            }}>{this.state.toggle[node.name] ? <ArrowUp/> : <ArrowDown/>}</div>
                            <span className="name">
                                <NavLink to={"/nodes/" + node.name}>
                                    {node.name}
                                </NavLink>
                            </span>
                            <span className="ipport">{node.ip+(node.port !== 0 ? ":" + node.port : "")}</span>
                        </div>

                        {this.state.toggle[node.name] ? <pre>
                            {node.config}
                        </pre> : null}
                    </li>
                })}
            </ul>
        }
    }

    render() {
        const {services} = this.props;
        if (services.fetchService.status === 404) {
            return <div className="error">Not Found</div>
        }
        if (services.fetchService.status === 500) {
            return <div className="error">{services.fetchService.error}</div>
        }
        if (services.fetchService.status === 200) {
            return <div>
                <div className="title">
                    <span className="name">{services.data.name}</span>
                </div>
                {this.renderList()}
            </div>
        } else {
            return null
        }
    }
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Service))