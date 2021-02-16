import React from "react";
import {NavLink, withRouter} from "react-router-dom";
import {connect} from "react-redux";
import {getNode} from "../../reducers/nodes";
import {mutationService,mutationAddService} from "../../reducers/mutation";
import ArrowDown from "../icons/arrowDown";
import ArrowUp from "../icons/arrowUp";
import CheckCircle from "../icons/checkCircle";
import Check from "../icons/check";
import Play from "../icons/play";
import Stop from "../icons/stop";
import Replay from "../icons/replay";
import Spinner from "../icons/spinner";
import Add from "../icons/add";
import toast, { Toaster } from 'react-hot-toast';

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
        mutationAddService: (ip, port, text, cb) => {
            dispatch(mutationAddService(ip, port, text, cb));
        }
    }
};


class Node extends React.Component {
    constructor(props) {
        super(props);
        this.state = {toggle: {}, check: {}, add:false, value: ""};
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
                port: nodes.data.rpcPort,
            };
            this.setState({check: check});
        }
    };
    handleChange = (event) => {
        this.setState({value: event.target.value});
    };
    addService = (ip, port) => {
        if (this.state.value ) {
            this.props.mutationAddService(ip, port, this.state.value, (success, error) => {
                if (success) {
                    this.props.getNode(this.nodeName);
                } else {
                    toast.error(error);
                }
            });
            this.setState({value:"", add:false})
        } else {
            toast.error("Config File Is Empty");
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
                    port: nodes.data.rpcPort,
                };
            });
            this.setState({check: check});
        }
    };
    getLastTime(time) {
        var sec =  Math.round(Date.now()/1000 - time)

        if (sec/(60*60*24) > 1) {
            return Math.round(sec/(60*60*24))+"D"
        }
        if (sec/(60*60) > 1) {
            return Math.round(sec/(60*60))+"H"
        }
        if (sec/(60) > 1) {
            return Math.round(sec/(60))+"M"
        }
        return sec+"S"
    }

    renderList() {
        const {nodes} = this.props;
        if (nodes.data.services) {
            return <ul className="nestservices">
                <li className="nesttitle">Services</li>
                {nodes.data.services.map((service, index) => {
                    let className = "item " + service.state.toLowerCase();
                    return <li className={className} key={service.name}>
                        <div className="basic">
                            {
                                this.state.check[service.name] ? <CheckCircle className="checked" onClick={() => {
                                    this.check(service);
                                }}
                                /> : <Check onClick={() => {
                                    this.check(service);
                                }}
                                />
                            }
                            <span className="name">
                                <NavLink to={"/services/" + service.name}>
                                    {service.name}
                                </NavLink>
                            </span>
                            {service.progress.inProgress?<React.Fragment>&nbsp;<Spinner/>&nbsp;{(service.progress.current/1024/1024).toFixed(2)}M</React.Fragment>:null}
                            <span className="ipport">
                                {service.port !== 0 ? service.ip + ":" + service.port : ""}
                            </span>
                            <div className="toggle" onClick={() => {
                                this.toggle(service.name)
                            }}>{this.state.toggle[service.name] ? <ArrowUp/> : <ArrowDown/>}</div>
                        </div>
                        {this.state.toggle[service.name] ? <div className="extend">CONFIG<pre>
                            {service.config}
                        </pre></div> : null}
                        {this.state.toggle[service.name] && service.statusContent !== "" ? <div className="extend">STATUS<pre>
                            {service.statusContent}
                        </pre></div> : null}
                        {/*{service.state.toLowerCase()==="alive"?<div className="footer"><Time/><span>{this.getLastTime(service.startTime)}</span></div>:null}*/}
                    </li>
                })}
            </ul>
        }
    }

    render() {
        const {nodes, mutation} = this.props;
        return <div>
            <Toaster
                toastOptions={{
                    className: 'toaster'
                }}
            />
            <div className="title">
                {nodes.data.services && nodes.data.services.length > 0 ? (
                    Object.keys(this.state.check).length === nodes.data.services.length ?
                        <CheckCircle
                            className="checked"
                            onClick={() => {
                                this.checkAll();
                            }}
                        />
                    : <Check
                            onClick={() => {
                                this.checkAll();
                            }}
                        />
                ) : null}
                <span className="name">{nodes.data.name}</span>
                <span className="ipport">{nodes.data.ip}</span>
                <span className="tools">
                    {mutation.mutationService.loading ?
                    <Spinner/>
                    :
                    <React.Fragment>
                        <Add className="add" title="add"
                              onClick={() => {
                                  this.setState({add:true})
                              }}/>
                        {Object.keys(this.state.check).length > 0 ? <React.Fragment>
                        <Play className="start" title="start"
                              onClick={() => {
                                  this.mutationService("START");
                              }}
                        />
                        <Stop className="stop" title="stop"
                              onClick={() => {
                                  this.mutationService("STOP");
                              }}
                        />
                        <Replay className="restart" title="restart"
                                onClick={() => {
                                    this.mutationService("RESTART");
                                }}
                        /></React.Fragment> : null}
                    </React.Fragment>
                }
                </span>
            </div>
            {this.renderList()}
            {this.state.add?<div className="addservice">
                <div className="text">
                    <textarea
                        value={this.state.value}
                        onChange={this.handleChange}
                    />
                </div>
                <div className="tools">
                    <input value="Close" type="button" onClick={() => {
                        this.setState({add:false});
                    }} />
                    <input value="Add" type="button" onClick={() => {
                        this.addService(nodes.data.ip,nodes.data.rpcPort)
                    }} />
                </div>
            </div>:null}
        </div>
    }
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(Node))