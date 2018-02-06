import React from 'react'
import {connect} from 'react-redux'
import {NavLink, Route, Switch} from 'react-router-dom'
import {getList} from '../../reducers/services'
import Service from './service'
import Search from '../icons/search'

import './index.pcss'

const mapStateToProps = (state, ownProps) => {
    return {
        services: state.services
    }
};

const mapDispatchToProps = (dispatch) => {
    return {
        getList: () => {
            dispatch(getList());
        }
    }
};

class Services extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            filter: ""
        }
    }

    componentWillMount() {
        this.props.getList();
    }

    componentWillReceiveProps(nextProps) {
        const locationChanged = nextProps.location !== this.props.location;
        if (locationChanged) {
            this.props.getList()
        }
    }

    onChange = (e) => {
        this.setState({filter: e.target.value});
    };

    renderList() {
        const {services} = this.props;
        return <ul className="list">
            <li className="filter">
                <Search/>
                <input placeholder="Filter by name" onChange={this.onChange} value={this.state.filter}/>
            </li>
            {services.list.map((service, index) => {
                if (service.name.toLowerCase().indexOf(this.state.filter.toLowerCase()) > -1) {
                    let className = "service-link";
                    return <li className="item" key={service.name}>
                        <NavLink className={className} to={"/services/" + service.name}>
                            <span className="name">{service.name}</span>
                        </NavLink>
                    </li>
                }
            })}
        </ul>
    }

    render() {
        return <div className="services">
            {this.renderList()}
            <div className="detail">
                <Switch>
                    <Route path="/services/:serviceName" component={Service}/>
                </Switch>
            </div>
        </div>
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(Services)