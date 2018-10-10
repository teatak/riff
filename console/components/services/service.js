import React from 'react'
import {NavLink, withRouter} from 'react-router-dom'
import {connect} from 'react-redux'
import {getService,} from '../../reducers/services'
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
        services: state.services,
        mutation: state.mutation,
    }
};

const mapDispatchToProps = (dispatch) => {
    return {
        getService: (nodeName) => {
            dispatch(getService(nodeName));
        },

        mutationService: (services, cmd) => {
            dispatch(mutationService(services, cmd));
        },
    }
};

class Service extends React.Component {
    constructor(props) {
        super(props);
        this.state = {toggle: {}, check: {}};
        this.serviceName = "";
    }

    componentWillMount() {
        this.serviceName = this.props.match.params.serviceName;
        if (this.serviceName !== undefined) {
            this.props.getService(this.serviceName)
        }
    }

    componentWillReceiveProps(nextProps) {
        const locationChanged = nextProps.location !== this.props.location;
        if (locationChanged) {
            this.setState({toggle: {}, check: {}});
            this.serviceName = nextProps.match.params.serviceName;
            this.props.getService(this.serviceName)
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

    check = (node) => {
        const {services} = this.props;
        if (this.state.check[node.name]) {
            //remove
            let check = this.state.check;
            delete check[node.name];
            this.setState({check: check});
        } else {
            let check = this.state.check;
            check[node.name] = {
                name: services.data.name,
                ip: node.ip,
                port: node.rpcPort,
            };
            this.setState({check: check});
        }
    };

    mutationService = (cmd) => {
        this.props.mutationService(this.state.check, cmd);
    };

    checkAll = () => {
        const {services} = this.props;
        if (Object.keys(this.state.check).length === services.data.nodes.length) {
            //uncheck all
            this.setState({check: {}});
        } else {
            let check = this.state.check;
            services.data.nodes.map((node, index) => {
                check[node.name] = {
                    name: services.data.name,
                    ip: node.ip,
                    port: node.rpcPort,
                };
            });
            this.setState({check: check});
        }
    };

    renderList() {
        const {services} = this.props;
        if (services.data.nodes) {

            return <ul className="nestnodes">
                <li className="nesttitle">Nodes</li>
                {services.data.nodes.map((node, index) => {
                    let className = "item " + node.state.toLowerCase();
                    return <li className={className} key={node.name}>
                        <div className="basic">
                            <CheckCircle className={this.state.check[node.name] ? "checked" : ""} onClick={() => {
                                this.check(node);
                            }}
                            />
                            <span className="name">
                                <NavLink to={"/nodes/" + node.name}>
                                    {node.name}
                                </NavLink>
                            </span>
                            <span className="ipport">{node.ip + (node.port !== 0 ? ":" + node.port : "")}</span>
                            <div className="toggle" onClick={() => {
                                this.toggle(node.name)
                            }}>{this.state.toggle[node.name] ? <ArrowUp/> : <ArrowDown/>}</div>
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
        const {services, mutation} = this.props;

        return <div>
            <div className="title">
                {services.data.nodes && services.data.nodes.length > 0 ? <CheckCircle
                    className={Object.keys(this.state.check).length === services.data.nodes.length ? "checked" : ""}
                    onClick={() => {
                        this.checkAll();
                    }}
                /> : null}
                <span className="name">{services.data.name}</span>
                {mutation.mutationService.loading ? <span className="tools">
                    <Spinner/>
                    </span> :
                    (Object.keys(this.state.check).length > 0 ? <span className="tools">
                    <Play className="start"
                          onClick={() => {
                              this.mutationService("START");
                          }}
                    />
                    <Stop className="stop"
                          onClick={() => {
                              this.mutationService("STOP");
                          }}
                    />
                    <Replay className="restart"
                            onClick={() => {
                                this.mutationService("RESTART");
                            }}
                    /></span> : null)
                }
            </div>
            {this.renderList()}
        </div>

    }
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Service))