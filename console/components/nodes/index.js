import React from 'react'
import { connect } from 'react-redux'

import { getList } from '../../reducers/nodes'

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
        return nodes.list.map((node, index) => {
            return <div key={node.name}>
                {node.name}
                {node.ip}
                {node.port}
            </div>
        })
    }
    render() {

        return <div>{this.renderList()}</div>
    }
}

export default connect(mapStateToProps,mapDispatchToProps)(Nodes)