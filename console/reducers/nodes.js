import Common from '../common'
import Config from 'config'

export const NODES_REQUEST = 'NODES_REQUEST';
export const NODES_SUCCESS = 'NODES_SUCCESS';
export const NODES_FAILURE = 'NODES_FAILURE';

export const NODE_WATCH = 'NODE_WATCH';
export const NODE_RESET = 'NODE_RESET';
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

export const isWatch = (nodeName) => (dispatch, getState) => {
    let state = getState();
    dispatch({
        type: NODE_WATCH,
        isWatch: !state.nodes.isWatch,
    });
    dispatch(getNode(nodeName));
};

let initReader = null;

export const cancelWatch = () => (dispatch, getState) => {
    let state = getState();
    if (state.nodes.fetchNode.loading) {
        if (initReader !== null) {
            initReader.cancel();
            initReader = null;
            dispatch({type: NODE_RESET});
        }
    }
};

const buildQuery = (nodeName) => {
    let node = (nodeName === undefined) ? "node:server" : "node(name:\"" + nodeName + "\")";
    let query = `{
    ` + node + ` {
        name
        ip
        port
        dataCenter
        snapShot
        state
        isSelf
        version
        services {
            name
            port
            state
            config
        } 
    }
}`;
    return query;
};

export const getNode = (nodeName) => (dispatch, getState) => {
    if (!Common.isIe()) {
        let state = getState();
        if (state.nodes.isWatch) {
            dispatch(watchNode(nodeName));
            return;
        }
    }
    let query = buildQuery(nodeName);

    dispatch(cancelWatch());
    dispatch({type: NODE_REQUEST});

    Common.fetch({query}, (json, error, status) => {
        if (status === 200) {
            dispatch({
                type: NODE_SUCCESS,
                status: status,
                json,
                receivedAt: Date.now()
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

let retryCount = 0;
const watchNode = (nodeName) => (dispatch, getState) => {
    let query = buildQuery(nodeName);
    dispatch(cancelWatch());
    dispatch({type: NODE_REQUEST});

    let consume = (reader) => {
        initReader = reader;
        let total = 0;
        return new Promise((resolve, reject) => {
            function pump() {
                reader.read().then(({done, value}) => {
                    if (done) {
                        resolve();
                        return
                    }
                    total += value.byteLength;
                    let arr = Common.utf8ArrayToStr(value).split("\n");
                    arr.map((text) => {
                        if (text !== "") {
                            let json = JSON.parse(text);
                            dispatch({
                                type: NODE_SUCCESS,
                                status: 200,
                                json,
                                receivedAt: Date.now()
                            });
                        }
                    });
                    pump()
                }).catch(reject)
            }

            pump()
        })
    };
    let param = nodeName === undefined ? "type=node" : "type=node&name=" + nodeName;
    fetch(Config.api + "/watch?" + param, {
        method: 'post',
        headers: {'connection': 'keep-alive'},
        credentials: 'include',
        body: JSON.stringify({query})
    }).then((response) => {
        retryCount = 0;
        return consume(response.body.getReader())
    }).then(() => {
        if (initReader != null) {
            dispatch({
                type: NODE_FAILURE,
                status: 500,
                error: "Server connect closed",
                receivedAt: Date.now()
            });
        }
    }).catch((error) => {
        //throw error;
        if (retryCount < 3) {
            retryCount++;
            setTimeout(() => {
                dispatch(watchNode(nodeName));
            }, retryCount * 1000);
        } else {
            dispatch({
                type: NODE_FAILURE,
                status: 500,
                error: error.message,
                receivedAt: Date.now()
            });
        }
    })
};

const nodes = (state = {
    fetchNodes: Common.initRequest,
    fetchNode: Common.initRequest,
    isWatch: true,
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
        case NODE_RESET:
            return {
                ...state,
                fetchNode: {
                    ...state.fetchNode,
                    loading: false,
                    status: 0,
                    error: null,
                }
            };
        case NODE_WATCH:
            return {
                ...state,
                isWatch: action.isWatch
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