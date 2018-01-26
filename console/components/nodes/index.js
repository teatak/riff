import React from 'react'
import { connect } from 'react-redux'
import { Switch, Route } from 'react-router-dom'
import { getList } from '../../reducers/nodes'
import Node from './node'
import { NavLink } from 'react-router-dom'

import './index.css'

const mapStateToProps = (state, ownProps) => {
    return {
        nodes: state.nodes
    }
};

const mapDispatchToProps = (dispatch) => {
    return {
        getList: () => {
            dispatch(getList());
        },
    }
};

class Nodes extends React.Component {
    constructor(props) {
        super(props);
    }
    componentWillMount() {
        this.props.getList();
    }
    renderList() {
        const { nodes } = this.props;
        return <ul className="list">
            {nodes.list.map((node, index) => {
                let className = "node-link";
                if(node.state === "Suspect" || node.state === "Dead") {
                    className += " dead";
                }
                return <li className="item" key={node.name}>
                    <NavLink className={className} to={"/nodes/"+node.name}>
                        <span className="name">{node.name}</span>
                        <span className="ipport">{node.ip}:{node.port}</span>
                        </NavLink>
                </li>
            })}
            </ul>
    }
    render() {
        return <div className="nodes">
                {this.renderList()}
            <div className="detail">
                <Switch>
                    <Route path="/nodes/:nodeName" component={Node}  />
                </Switch>
            </div>
        </div>
    }
}

export default connect(mapStateToProps,mapDispatchToProps)(Nodes)