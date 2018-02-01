import Common from './common'

export const SERVICES_REQUEST = 'SERVICES_REQUEST';
export const SERVICES_SUCCESS = 'SERVICES_SUCCESS';
export const SERVICES_FAILURE = 'SERVICES_FAILURE';

export const SERVICE_REQUEST = 'SERVICE_REQUEST';
export const SERVICE_SUCCESS = 'SERVICE_SUCCESS';
export const SERVICE_FAILURE = 'SERVICE_FAILURE';

//获取Product
export const getList = () => (dispatch, getState) => {
    const query = `{
    services {
        name
    }
}`;
    dispatch({type: SERVICES_REQUEST});
    Common.fetch({query}, (json, error, status) => {
        if (status === 200) {
            dispatch({
                type: SERVICES_SUCCESS,
                status: status,
                json,
                receivedAt: Date.now()
            });
        } else {
            dispatch({
                type: SERVICES_FAILURE,
                status: status,
                error: error,
                receivedAt: Date.now()
            });
        }
    })
};

export const getService = (serviceName,state) => (dispatch, getState) => {
    let query = `{
    service(name:"`+serviceName+`",state:All) {
        name
        nodes {
            name
            ip
            port
            state
            isSelf
            config
        } 
    }
}`;
    dispatch({type: SERVICE_REQUEST});
    Common.fetch({query}, (json, error, status) => {
        if (status === 200) {
            dispatch({
                type: SERVICE_SUCCESS,
                status: status,
                json,
                receivedAt: Date.now()
            });
        } else {
            dispatch({
                type: SERVICE_FAILURE,
                status: status,
                error: error,
                receivedAt: Date.now()
            });
        }
    })
};

const services = (state = {
    fetchServices: Common.initRequest,
    fetchService: Common.initRequest,
    list: [],                   //数据
    data: {}
}, action) => {
    switch (action.type) {
        case SERVICES_REQUEST:
            return {
                ...state,
                fetchServices: {
                    ...state.fetchServices,
                    loading: true,
                    status: 0,
                    error: null,
                }
            };
        case SERVICES_SUCCESS:
            return {
                ...state,
                fetchServices: {
                    ...state.fetchServices,
                    loading: false,
                    status: 200,
                    error: null,
                    lastUpdated: action.receivedAt
                },
                list: action.json.data.services
            };
        case SERVICES_FAILURE:
            return {
                ...state,
                fetchServices: {
                    ...state.fetchServices,
                    loading: false,
                    status: action.status,
                    error: action.error,
                    lastUpdated: action.receivedAt
                },
            };
        case SERVICE_REQUEST:
            return {
                ...state,
                fetchService: {
                    ...state.fetchService,
                    loading: true,
                    status: 0,
                    error: null,
                }
            };
        case SERVICE_SUCCESS:
            return {
                ...state,
                fetchService: {
                    ...state.fetchService,
                    loading: false,
                    status: 200,
                    error: null,
                    lastUpdated: action.receivedAt
                },
                data: action.json.data.service
            };
        case SERVICE_FAILURE:
            return {
                ...state,
                fetchService: {
                    ...state.fetchService,
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

export default services