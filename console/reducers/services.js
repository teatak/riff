import Common from '../common'
import Config from 'config'

export const SERVICES_REQUEST = 'SERVICES_REQUEST';
export const SERVICES_SUCCESS = 'SERVICES_SUCCESS';
export const SERVICES_FAILURE = 'SERVICES_FAILURE';

export const SERVICE_WATCH = 'SERVICE_WATCH';
export const SERVICE_RESET = 'SERVICE_RESET';
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

export const isWatch = (serviceName) => (dispatch, getState) => {
    let state = getState();
    dispatch({
        type: SERVICE_WATCH,
        isWatch: !state.services.isWatch,
    });
    dispatch(getService(serviceName));
};

let initReader = null;

export const cancelWatch = () => (dispatch, getState) => {
    let state = getState();
    if (state.services.fetchService.loading) {
        if (initReader !== null) {
            initReader.cancel();
            initReader = null;
            dispatch({type: SERVICE_RESET});
        }
    }
};

const buildQuery = (serviceName) => {
    let query = `{
    service(name:"` + serviceName + `",state:All) {
        name
        nodes {
            name
            ip
            port
            rpcPort
            state
            isSelf
            config
        } 
    }
}`;
    return query;
};

export const getService = (serviceName, state) => (dispatch, getState) => {
    if (!Common.isIe()) {
        let state = getState();
        if (state.services.isWatch) {
            dispatch(watchService(serviceName, state));
            return;
        }
    }
    let query = buildQuery(serviceName);

    dispatch(cancelWatch());
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

let retryCount = 0;
const watchService = (serviceName, state) => (dispatch, getState) => {
    let query = buildQuery(serviceName);
    dispatch(cancelWatch());
    dispatch({type: SERVICE_REQUEST});

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
                                type: SERVICE_SUCCESS,
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
    let param = serviceName === undefined ? "type=service" : "type=service&name=" + serviceName;
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
                type: SERVICE_FAILURE,
                status: 500,
                error: "Server connect closed",
                receivedAt: Date.now()
            });
        }
    }).catch((error) => {
        //throw error
        if (retryCount < 3) {
            retryCount++;
            setTimeout(() => {
                dispatch(watchService(serviceName, state));
            }, retryCount * 1000);
        } else {
            dispatch({
                type: SERVICE_FAILURE,
                status: 500,
                error: error.message,
                receivedAt: Date.now()
            });
        }
    })
};

const services = (state = {
    fetchServices: Common.initRequest,
    fetchService: Common.initRequest,
    isWatch: true,
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
        case SERVICE_RESET:
            return {
                ...state,
                fetchService: {
                    ...state.fetchService,
                    loading: false,
                    status: 0,
                    error: null,
                }
            };
        case SERVICE_WATCH:
            return {
                ...state,
                isWatch: action.isWatch
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
                    loading: true,
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