import Common from './common'
import Config from 'config'

export const LOG_REQUEST = 'LOG_REQUEST';
export const LOG_SUCCESS = 'LOG_SUCCESS';
export const LOG_FAILURE = 'LOG_FAILURE';

export const getLogs = () => (dispatch, getState) => {
    let state = getState();
    if (state.logs.fetchLogs.loading) {
        return
    }
    dispatch({type: LOG_REQUEST});
    let reader = null;
    let loop = (response) => {
        setTimeout(() => {
            reader.read().then((b) => {
                let text = new TextDecoder("utf-8").decode(b.value);
                dispatch({
                    type: LOG_SUCCESS,
                    status: response.status,
                    text,
                    receivedAt: Date.now()
                });
                if (b.done) {
                    init = false;
                }
                loop(response);
            }).catch((error) => {
                dispatch({
                    type: LOG_FAILURE,
                    status: 500,
                    error: error,
                    receivedAt: Date.now()
                });
            })
        }, 0);
    };
    fetch(Config.api + "/logs", {
        method: 'get',
        headers: {'connection': 'keep-alive'},
    }).then((response) => {
        reader = response.body.getReader();
        loop(response);
    }).catch((error) => {
        dispatch({
            type: LOG_FAILURE,
            status: 500,
            error: error,
            receivedAt: Date.now()
        });
    })
};

const logs = (state = {
    fetchLogs: Common.initRequest,
    list: [],                   //数据
}, action) => {
    switch (action.type) {
        case LOG_REQUEST:
            return {
                ...state,
                fetchLogs: {
                    ...state.fetchLogs,
                    loading: true,
                    status: 0,
                    error: null,
                }
            };
        case LOG_SUCCESS:
            return {
                ...state,
                fetchLogs: {
                    ...state.fetchLogs,
                    status: 200,
                    error: null,
                    lastUpdated: action.receivedAt
                },
                list: [...state.list, action.text]
            };
        case LOG_FAILURE:
            return {
                ...state,
                fetchLogs: {
                    ...state.fetchLogs,
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

export default logs