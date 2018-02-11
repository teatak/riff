import Common from '../common'
import Ws from "./ws";

export const LOG_RESET = 'LOG_RESET';
export const LOG_REQUEST = 'LOG_REQUEST';
export const LOG_SUCCESS = 'LOG_SUCCESS';
export const LOG_FAILURE = 'LOG_FAILURE';

export const getLogs = () => (dispatch, getState) => {
    let state = getState();
    if (state.logs.fetchLogs.loading) {
        return
    }
    dispatch({type: LOG_REQUEST});
    Ws.logs({
        event: "Logs"
    }, (text) => {
        dispatch({
            type: LOG_SUCCESS,
            status: 200,
            text: text,
            receivedAt: Date.now()
        });
    });


};

const logs = (state = {
    fetchLogs: Common.initRequest,
    list: [],                   //数据
}, action) => {
    switch (action.type) {
        case LOG_RESET:
            return {
                ...state,
                fetchLogs: Common.initRequest,
                list: [],
            };
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