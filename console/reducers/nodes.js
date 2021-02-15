import Common from "../common";
import Ws from "./ws";

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
        httpPort
        rpcPort
        dataCenter
        state
        isSelf
        version
    }
}`;
    dispatch({type: NODES_REQUEST});
    Common.fetch({query}, (json, error, status) => {
        if (status === 200) {
            dispatch({
                type: NODES_SUCCESS,
                status: status,
                json,
                receivedAt: Date.now()
            });
        } else {
            dispatch({
                type: NODES_FAILURE,
                status: status,
                error: error,
                receivedAt: Date.now()
            });
        }
    })
};

const buildQuery = (nodeName) => {
    let node = (nodeName === undefined) ? "node:server" : "node(name:\"" + nodeName + "\")";
    return `{
    ` + node + ` {
        name
        ip
        httpPort
        rpcPort
        dataCenter
        snapShot
        state
        isSelf
        version
        services {
            name
            ip
            port
            state
            config
            statusContent
            startTime
            progress {
                current
                total
                inProgress
            }
        } 
    }
}`;
};

export const changeNode = (json) => (dispatch, getState) => {
    dispatch({
        type: NODE_SUCCESS,
        status: status,
        json,
        receivedAt: Date.now()
    });
};

export const getNode = (nodeName) => (dispatch, getState) => {
    dispatch({type: NODE_REQUEST});
    let query = buildQuery(nodeName);
    Common.fetch({query}, (json, error, status) => {
        if (status === 200) {
            dispatch({
                type: NODE_SUCCESS,
                status: status,
                json,
                receivedAt: Date.now()
            });
            let name = nodeName === undefined ? "" : nodeName;
            Ws.watch({
                event: "Watch",
                body: {
                    name: name,
                    type: "node",
                    query: query
                }
            }, (json) => {
                dispatch({
                    type: NODE_SUCCESS,
                    status: status,
                    json,
                    receivedAt: Date.now()
                });
            });
        } else {
            dispatch({
                type: NODE_FAILURE,
                status: status,
                error: error,
                receivedAt: Date.now()
            });
        }
    })
};

const nodes = (state = {
    fetchNodes: Common.initRequest,
    fetchNode: Common.initRequest,
    list: [],                   //数据
    data: {}
}, action) => {
    switch (action.type) {
        case NODES_REQUEST:
            return {
                ...state,
                fetchNodes: {
                    ...state.fetchNodes,
                    loading: true,
                    status: 0,
                    error: null,
                }
            };
        case NODES_SUCCESS:
            return {
                ...state,
                fetchNodes: {
                    ...state.fetchNodes,
                    loading: false,
                    status: 200,
                    error: null,
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
            return {
                ...state,
                fetchNode: {
                    ...state.fetchNode,
                    loading: true,
                    status: 0,
                    error: null,
                }
            };
        case NODE_SUCCESS:
            return {
                ...state,
                fetchNode: {
                    ...state.fetchNode,
                    loading: true,
                    status: 200,
                    error: null,
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