import React from 'react'
import Icon from './icon'
import './spinner.pcss'

let Spinner = (props) => (
    <Icon {...props} className="spinner" viewBox="0 0 600 600">
        <circle cx="300" cy="300" r="250" fill="none" stroke="#e6443b" strokeDasharray="392.699081698724155 3000"
                strokeWidth="60" transform="rotate(0 300 300)"/>
        <circle cx="300" cy="300" r="250" fill="none" stroke="#408cea" strokeDasharray="392.699081698724155 3000"
                strokeWidth="60" transform="rotate(90 300 300)"/>
        <circle cx="300" cy="300" r="250" fill="none" stroke="#fbc12c" strokeDasharray="392.699081698724155 3000"
                strokeWidth="60" transform="rotate(180 300 300)"/>
        <circle cx="300" cy="300" r="250" fill="none" stroke="#3fba66" strokeDasharray="392.699081698724155 3000"
                strokeWidth="60" transform="rotate(270 300 300)"/>
    </Icon>
);
export default Spinner