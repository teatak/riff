import createHashHistory from 'history/createHashHistory'

const HashHistory = createHashHistory({
    hashType: 'slash' // Omit the leading slash
});

export default HashHistory