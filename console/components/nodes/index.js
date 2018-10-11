import React from 'react'
import {connect} from 'react-redux'
import {NavLink, Route, Switch} from 'react-router-dom'
import {getList} from '../../reducers/nodes'
import Node from './node'
import Search from '../icons/search'

import './index.pcss'

const mapStateToProps = (state, ownProps) => {
    return {
        nodes: state.nodes
    }
};

const mapDispatchToProps = (dispatch) => {
    return {
        getList: () => {
            dispatch(getList());
        }
    }
};

class Nodes extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            filter: ""
        }
    }

    componentWillMount() {
        this.props.getList();
    }

    componentWillReceiveProps(nextProps) {
        const locationChanged = nextProps.location !== this.props.location;
        if (locationChanged) {
            this.props.getList()
        }
    }

    onChange = (e) => {
        this.setState({filter: e.target.value});
    };

    renderList() {
        const {nodes} = this.props;
        return <ul className="list">
            <li className="filter">
                <Search/>
                <input placeholder="Filter by name" onChange={this.onChange} value={this.state.filter}/>
            </li>
            {nodes.list.map((node, index) => {
                if (node.name.toLowerCase().indexOf(this.state.filter.toLowerCase()) > -1) {
                    let className = "node-link " + node.state.toLowerCase();
                    return <li className="item" key={node.name}>
                        <NavLink className={className} to={"/nodes/" + node.name}>
                            <span className="name">{node.name}</span>
                            <span className="ipport">{node.ip}:{node.rpcPort}</span>
                        </NavLink>
                    </li>
                }
            })}
        </ul>
    }

    render() {
        const {nodes} = this.props;
        return <div className="nodes">
            {this.renderList()}
            <div className="detail">
                <Switch>
                    <Route path="/nodes/:nodeName?" component={Node}/>
                </Switch>
            </div>
        </div>
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(Nodes)