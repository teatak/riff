import Common from './common'

export const NODES_REQUEST = 'NODES_REQUEST';
export const NODES_SUCCESS = 'NODES_SUCCESS';
export const NODES_FAILURE = 'NODES_FAILURE';

export const NODE_REQUEST = 'NODE_REQUEST';
export const NODE_SUCCESS = 'NODE_SUCCESS';
export const NODE_FAILURE = 'NODE_FAILURE';

//获取Product
export const getList = () => (dispatch, getState) => {
    const query = `{
    nodes {
        name
        ip
        port
        dataCenter
        state
        version
    }
}`;
    dispatch({ type: NODES_REQUEST });
    Common.fetch({query},(json,error) => {
        if(error==null) {
            dispatch({
                type: NODES_SUCCESS,
                status: 200,
                json,
                receivedAt: Date.now()
            });
        } else {
            dispatch({
                type: NODES_FAILURE,
                status: 500,
                error: error,
                receivedAt: Date.now()
            });
        }
    })
};

export const getNode = (nodeName) => (dispatch, getState) => {
    let node = (nodeName === undefined)?"node:server":"node(name:\""+nodeName+"\")";
    let query = `{
    `+node+` {
        name
        ip
        port
        dataCenter
        snapShot
        state
        version
        services {
            name
            ip
            port
            state
        } 
    }
}`;
    dispatch({ type: NODE_REQUEST });
    Common.fetch({query},(json,error) => {
        if(error==null) {
            dispatch({
                type: NODE_SUCCESS,
                status: 200,
                json,
                receivedAt: Date.now()
            });
        } else {
            dispatch({
                type: NODE_FAILURE,
                status: 500,
                error: error,
                receivedAt: Date.now()
            });
        }
    })
};

const nodes = (
    state = {
        fetchNodes:     Common.initRequest,
        fetchNode:      Common.initRequest,
        list:           [],                   //数据
        data:           {}
    }, action) => {
    switch (action.type) {
        case NODES_REQUEST:
            return { ...state,
                fetchNodes: {
                    ...state.fetchNodes,
                    loading: true,
                    status: 0,
                }
            };
        case NODES_SUCCESS:
            return {
                ...state,
                fetchNodes: { ...state.fetchNodes,
                    loading: false,
                    status: 200,
                    lastUpdated: action.receivedAt
                },
                list: action.json.data.nodes
            };
        case NODES_FAILURE:
            return {
                ...state,
                fetchNodes: {
                    ...state.fetchNodes,
                    loading: false,
                    status: action.status,
                    error: action.error,
                    lastUpdated: action.receivedAt
                },
            };
        case NODE_REQUEST:
            return { ...state,
                fetchNode: {
                    ...state.fetchNode,
                    loading: true,
                    status: 0
                }
            };
        case NODE_SUCCESS:
            return {
                ...state,
                fetchNode: { ...state.fetchNode,
                    loading: false,
                    status: 200,
                    lastUpdated: action.receivedAt
                },
                data: action.json.data.node
            };
        case NODE_FAILURE:
            return {
                ...state,
                fetchNode: {
                    ...state.fetchNode,
                    loading: false,
                    status: action.status,
                    error: action.error,
                    lastUpdated: action.receivedAt
                },
            };
        default:
            return state
    }
};

export default nodes