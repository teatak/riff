import React from 'react'
import Menu from './components/menu'

class App extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return <div>
            <Menu/>
        </div>
    }
}

export default App