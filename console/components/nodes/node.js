import React from 'react'
import {NavLink, withRouter} from 'react-router-dom'
import {connect} from 'react-redux'
import {getNode} from '../../reducers/nodes'
import {mutationService} from '../../reducers/mutation'
import ArrowDown from '../icons/arrowDown'
import ArrowUp from '../icons/arrowUp'
import CheckCircle from '../icons/checkCircle'
import Play from '../icons/play'
import Stop from '../icons/stop'
import Replay from '../icons/replay';
import Spinner from "../icons/spinner";

const mapStateToProps = (state, ownProps) => {
    return {
        nodes: state.nodes,
        mutation: state.mutation,
    }
};

const mapDispatchToProps = (dispatch) => {
    return {
        getNode: (nodeName) => {
            dispatch(getNode(nodeName));
        },
        mutationService: (services, cmd) => {
            dispatch(mutationService(services, cmd));
        },
    }
};

class Node extends React.Component {
    constructor(props) {
        super(props);
        this.state = {toggle: {}, check: {}};
        this.nodeName = "";
    }

    componentWillMount() {
        this.nodeName = this.props.match.params.nodeName;
        this.props.getNode(this.nodeName)
    }

    componentWillReceiveProps(nextProps) {
        const locationChanged = nextProps.location !== this.props.location;
        if (locationChanged) {
            this.setState({toggle: {}, check: {}});
            this.nodeName = nextProps.match.params.nodeName;
            this.props.getNode(this.nodeName)
        }
    }

    componentWillUnmount() {
    }

    toggle = (name) => {
        if (this.state.toggle[name]) {
            //remove
            let toggle = this.state.toggle;
            delete toggle[name];
            this.setState({toggle: toggle});
        } else {
            let toggle = this.state.toggle;
            toggle[name] = true;
            this.setState({toggle: toggle});
        }
    };

    check = (service) => {
        const {nodes} = this.props;
        if (this.state.check[service.name]) {
            //remove
            let check = this.state.check;
            delete check[service.name];
            this.setState({check: check});
        } else {
            let check = this.state.check;
            check[service.name] = {
                name: service.name,
                ip: nodes.data.ip,
                port: nodes.data.port,
            };
            this.setState({check: check});
        }
    };

    mutationService = (cmd) => {
        this.props.mutationService(this.state.check, cmd);
    };

    checkAll = () => {
        const {nodes} = this.props;
        if (Object.keys(this.state.check).length === nodes.data.services.length) {
            //uncheck all
            this.setState({check: {}});
        } else {
            let check = this.state.check;
            nodes.data.services.map((service, index) => {
                check[service.name] = {
                    name: service.name,
                    ip: nodes.data.ip,
                    port: nodes.data.port,
                };
            });
            this.setState({check: check});
        }
    };

    renderList() {
        const {nodes} = this.props;
        if (nodes.data.services) {
            return <ul className="nestservices">
                <li className="nesttitle">Services</li>
                {nodes.data.services.map((service, index) => {
                    let className = "item " + service.state.toLowerCase();
                    return <li className={className} key={service.name}>
                        <div className="basic">
                            <CheckCircle className={this.state.check[service.name] ? "checked" : ""} onClick={() => {
                                this.check(service);
                            }}
                            />
                            <span className="name">
                                <NavLink to={"/services/" + service.name}>
                                    {service.name}
                                </NavLink>
                            </span>
                            <span className="ipport">{service.port !== 0 ? ":" + service.port : ""}</span>
                            <div className="toggle" onClick={() => {
                                this.toggle(service.name)
                            }}>{this.state.toggle[service.name] ? <ArrowUp/> : <ArrowDown/>}</div>
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
        const {nodes, mutation} = this.props;

        return <div>
            <div className="title">
                {nodes.data.services && nodes.data.services.length > 0 ? <CheckCircle
                    className={Object.keys(this.state.check).length === nodes.data.services.length ? "checked" : ""}
                    onClick={() => {
                        this.checkAll();
                    }}
                /> : null}
                <span className="name">{nodes.data.name}</span>
                <span className="ipport">{nodes.data.ip}</span>
                {mutation.mutationService.loading ? <span className="tools">
                    <Spinner/>
                    </span> :
                    (Object.keys(this.state.check).length > 0 ? <span className="tools">
                    <Play className="start"
                          onClick={() => {
                              this.mutationService("Start");
                          }}
                    />
                    <Stop className="stop"
                          onClick={() => {
                              this.mutationService("Stop");
                          }}
                    />
                    <Replay className="restart"
                            onClick={() => {
                                this.mutationService("Restart");
                            }}
                    /></span> : null)
                }
            </div>
            {this.renderList()}
        </div>
    }
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Node))