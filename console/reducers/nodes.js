import Config from 'config'
import Common from './common'

export const NODES_REQUEST = 'NODES_REQUEST';
export const NODES_SUCCESS = 'NODES_SUCCESS';
export const NODES_FAILURE = 'NODES_FAILURE';

//获取Product
export const getList = () => (dispatch, getState) => {
    const query = `{
    nodes {
        name
        ip
        port
        dataCenter
    }
}`;
    dispatch({ type: NODES_REQUEST });
    fetch(Config.api, {
        method: 'post',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({query}),
    })
        .then(Common.checkStatus)
        .then(Common.parseJSON)
        .then(json => {
            dispatch({
                type: NODES_SUCCESS,
                status: 200,
                json,
                receivedAt: Date.now()
            });
        })
        .catch(error => {
            dispatch({
                type: NODES_FAILURE,
                status: 500,
                error: error,
                receivedAt: Date.now()
            });
        })
    //     .then(response => {
    //     return response.json()
    // }).then(json => {
    //     dispatch({
    //         type: NODES_SUCCESS,
    //         status: 200,
    //         json,
    //         receivedAt: Date.now()
    //     });
    // }).catch( ex => {
    //     dispatch({
    //         type: NODES_FAILURE,
    //         status: 500,
    //         error: ex,
    //         receivedAt: Date.now()
    //     });
    //     console.log('parsing failed', ex)
    // });
};

const nodes = (
    state = {
        fetch:          Common.initRequest,
        data:           [],                   //数据
    }, action) => {
    switch (action.type) {
        case NODES_REQUEST:
            return { ...state,
                fetch: {
                    ...state.fetch,
                    loading:true,
                    status: 0
                }
            };
        case NODES_SUCCESS:
            return {
                ...state,
                fetch: { ...state.fetch,
                    loading:false,
                    status: 200,
                    lastUpdated: action.receivedAt
                },
                data: action.json.data.nodes
            };
        case NODES_FAILURE:
            return {
                ...state,
                fetch: {
                    ...state.fetch,
                    loading:false,
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