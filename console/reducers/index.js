import {combineReducers} from 'redux'
import nodes from './nodes'
import services from './services'
import logs from './logs'

const reducers = combineReducers({
    nodes,
    services,
    logs,
});

export default reducers