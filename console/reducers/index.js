import {combineReducers} from 'redux'
import nodes from './nodes'
import services from './services'
import logs from './logs'
import mutation from './mutation'

const reducers = combineReducers({
    nodes,
    services,
    logs,
    mutation,
});

export default reducers